package repository

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/emrebdr/dirgod-core/models"
)

func (r *Repository) Rollback(commitId string) error {
	var activeCommitId string = r.getActiveCommitId()
	if strings.HasPrefix(activeCommitId, commitId) {
		return errors.New("you cannot rollback to the current commit id")
	}

	fullCommitId := r.getCommitId(commitId)
	if fullCommitId == "" {
		return errors.New("commit id not found")
	}

	rootTree := r.getCommitRootTree(fullCommitId)
	if rootTree == nil {
		return errors.New("commit tree object couldn't find")
	}

	commitPath := r.findCommitPath(fullCommitId)
	if commitPath == "" {
		return errors.New("commit not found")
	}

	err := r.regenerateStructure(*rootTree, commitPath)
	if err != nil {
		return err
	}

	err = r.changeHeadPointer(fullCommitId)
	if err != nil {
		return err
	}

	err = r.changeRepositoryCommitPointer(fullCommitId)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) findCommitPath(commitId string) string {
	var commitPath string = ""
	err := filepath.Walk(r.Path+"/.dirgod/objects", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if len(commitId) > 10 {
				if info.Name() == commitId[:10] {
					commitPath = path
					return io.EOF
				}
			} else {
				if strings.HasPrefix(info.Name(), commitId) {
					commitPath = path
					return io.EOF
				}
			}
		}

		return nil
	})

	if err != nil && err != io.EOF {
		return ""
	}

	return commitPath
}

func (r *Repository) regenerateStructure(rootTree models.Tree, commitPath string) error {
	err := r.clearStructure()
	if err != nil {
		return err
	}

	err = r.createStructure(rootTree, commitPath)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) createStructure(tree models.Tree, commitPath string) error {
	err := os.MkdirAll(tree.Path, 0755)
	if err != nil {
		return err
	}

	var blobs []models.Blob
	for _, blobId := range tree.Blobs {
		var blobObject models.Blob
		blobContent, err := os.ReadFile(commitPath + "/" + blobId)
		if err != nil {
			return err
		}

		err = json.Unmarshal(blobContent, &blobObject)
		if err != nil {
			return err
		}

		blobs = append(blobs, blobObject)
	}

	for _, blob := range blobs {
		err := os.WriteFile(blob.Path, blob.Content, 0644)
		if err != nil {
			return err
		}
	}

	for _, treeId := range tree.Trees {
		var newTree models.Tree
		treeContent, err := os.ReadFile(commitPath + "/" + treeId)
		if err != nil {
			return err
		}

		err = json.Unmarshal(treeContent, &newTree)
		if err != nil {
			return err
		}

		err = r.createStructure(newTree, commitPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) clearStructure() error {
	var directoriesPath []string
	err := filepath.Walk(r.Path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasPrefix(path, r.Path+"/.dirgod") || path == r.Path {
			return nil
		}

		//! need to check .dirgodignore

		if info.IsDir() {
			directoriesPath = append(directoriesPath, path)
			return nil
		}

		err = os.Remove(path)
		if err != nil {
			return err
		}

		return nil
	})

	for _, dirPath := range directoriesPath {
		f, err := os.Open(dirPath)
		if err != nil {
			return err
		}

		defer f.Close()

		_, err = f.Readdirnames(1)
		if err == io.EOF {
			err = os.Remove(dirPath)
			if err != nil {
				return err
			}
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) changeHeadPointer(commitId string) error {
	HEADFile, err := os.ReadFile(r.Path + "/.dirgod/HEAD")
	if err != nil {
		return err
	}

	refPath := strings.Split(string(HEADFile), "ref: ")[1]
	err = os.WriteFile(r.Path+"/.dirgod/"+refPath, []byte(commitId), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) changeRepositoryCommitPointer(commitId string) error {
	r.CommitObject.CommitId = commitId
	err := r.saveReporefFile(r)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) PrintCommitLog() string {
	commits := r.getCommitHistory()
	var commitLogContent string = ""
	var activeCommitId string = r.getActiveCommitId()

	for _, commit := range commits {
		if commit.CommitId == activeCommitId {
			commitLogContent += "ACTIVE\n"
		}

		commitLogContent += fmt.Sprintf("CommitId: %s\nCommitter: %s\nMessage: %s\n===================\n", commit.CommitId, commit.Committer, commit.Message)
	}

	return commitLogContent
}

func (r *Repository) getCommitHistory() []models.CommitLog {
	branch, err := os.ReadFile(r.Path + "/.dirgod/HEAD")
	if err != nil {
		return nil
	}

	ref := strings.Split(string(branch), "ref: ")
	splitRef := strings.Split(ref[1], "/")
	branchName := splitRef[len(splitRef)-1]

	logFile, err := os.Open(r.Path + "/.dirgod/logs/" + branchName)
	if err != nil {
		return nil
	}

	defer logFile.Close()

	fileScanner := bufio.NewScanner(logFile)
	fileScanner.Split(bufio.ScanLines)

	var commits []models.CommitLog

	for fileScanner.Scan() {
		var newCommit models.CommitLog
		newCommit.Deserialize(fileScanner.Text())
		commits = append(commits, newCommit)
	}

	return commits
}

func (r *Repository) getActiveCommitId() string {
	HEADFileContent, err := os.ReadFile(r.Path + "/.dirgod/HEAD")
	if err != nil {
		return ""
	}

	reference := strings.Split(string(HEADFileContent), "ref: ")[1]
	commitId, err := os.ReadFile(r.Path + "/.dirgod/" + reference)
	if err != nil {
		return ""
	}

	return string(commitId)
}

func (r *Repository) getAllCommits() []models.Commit {
	var commits []models.Commit
	err := filepath.Walk(r.Path+"/.dirgod/refs/commits", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fileContent, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var commit models.Commit
		err = json.Unmarshal(fileContent, &commit)
		if err != nil {
			return err
		}

		commits = append(commits, commit)
		return nil
	})

	if err != nil {
		return nil
	}

	return commits
}

func (r *Repository) getCommitRootTree(commitId string) *models.Tree {
	commits := r.getAllCommits()
	var commitObject models.Commit
	for _, commit := range commits {
		if strings.HasPrefix(commit.CommitId, commitId) {
			commitObject = commit
			break
		}
	}

	var tree models.Tree
	fileContent, err := os.ReadFile(r.Path + "/.dirgod/objects/" + commitObject.CommitId[:10] + "/" + commitObject.Tree)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(fileContent, &tree)
	if err != nil {
		return nil
	}

	return &tree
}

func (r *Repository) getCommitId(commitId string) string {
	commits := r.getAllCommits()
	for _, commit := range commits {
		if strings.HasPrefix(commit.CommitId, commitId) {
			return commit.CommitId
		}
	}

	return ""
}
