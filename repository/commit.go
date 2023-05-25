package repository

import (
	"bufio"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/emrebdr/dirgod-core/interfaces"
	"github.com/emrebdr/dirgod-core/models"
	"github.com/emrebdr/dirgod-core/utils"
)

func (r *Repository) Commit(commitMessage string) error {
	err := r.createCommitObject(commitMessage)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) createCommitObject(message string) error {
	commitId := utils.GenerateId()

	commit := models.Commit{}
	commit.Author = fmt.Sprintf("%s %s", r.RootConfig.User.Email, r.RootConfig.User.Username)
	commit.Committer = fmt.Sprintf("%s %s", r.InternalConfig.User.Email, r.InternalConfig.User.Username) //! tbc
	commit.Message = message
	commit.CommitId = commitId

	tree, err := r.createTree(commitId[:10])
	if err != nil {
		fmt.Printf("tree: %v\n", err)
		return err
	}

	fmt.Printf("tree obj: %v\n", tree)

	commit.Tree = tree.TreeId

	r.CommitObject = commit

	err = r.saveReporefFile(r)
	if err != nil {
		fmt.Printf("repo ref: %v\n", err)
		return err
	}

	err = r.createCommitFiles(commitId, message)
	if err != nil {
		fmt.Printf("create commit file: %v\n", err)
		return err
	}

	err = r.createCommitLog(commitId, message)
	if err != nil {
		fmt.Printf("create commit log: %v\n", err)
		return err
	}

	return nil
}

func (r *Repository) createCommitFiles(commitId, message string) error {
	err := r.createCommitMsgFile(message)
	if err != nil {
		return err
	}

	err = r.createCommitHeadFile(commitId)
	if err != nil {
		return err
	}

	err = r.createCommitObjectFile()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) createCommitObjectFile() error {
	data, err := json.Marshal(r.CommitObject)
	if err != nil {
		return err
	}

	err = os.WriteFile(r.Path+"/.dirgod/refs/commits/"+r.CommitObject.CommitId, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) createCommitHeadFile(commitId string) error {
	//? If the branch system comes, this place will change
	file, err := os.Create(r.Path + "/.dirgod/refs/branches/main")
	if err != nil {
		return err
	}

	_, err = file.Write([]byte(commitId))
	if err != nil {
		return err
	}

	HEADFile, err := os.Create(r.Path + "/.dirgod/HEAD")
	if err != nil {
		return err
	}

	_, err = HEADFile.Write([]byte("ref: " + "refs/branches/main"))
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) createCommitMsgFile(message string) error {
	commitmsgFile, err := os.Create(r.Path + "/.dirgod/COMMIT_MSG")
	if err != nil {
		return err
	}

	writtenBytes, err := commitmsgFile.Write([]byte(message))
	if err != nil {
		return err
	}

	if len([]byte(message)) != writtenBytes {
		return errors.New("something went wrong while preparing commit message file")
	}

	return nil
}

func (r *Repository) createTree(ref string) (*models.Tree, error) {
	directoryMap := map[string]*models.Tree{}
	blobMap := map[string]*models.Blob{}

	err := filepath.Walk(r.Path, func(path string, info fs.FileInfo, err error) error {
		if strings.HasPrefix(path, r.Path+"/.dirgod") && !strings.HasPrefix(path, r.Path+"/.dirgodignore") { //! tbc
			return nil
		}

		if strings.HasPrefix(path, r.Path+"/.git") && !strings.HasPrefix(path, r.Path+"/.gitignore") {
			return nil
		}

		if err != nil {
			fmt.Printf("t walk: %v\n", err)
			return err
		}

		if info.IsDir() {
			directoryMap = r.addTreeObject(path, info, directoryMap)
			return nil
		}

		dirMap, newBlobMap, err := r.addBlobObject(path, info.Name(), directoryMap, blobMap)
		if err != nil {
			fmt.Printf("t walk for blob: %v\n", err)
			return err
		}

		directoryMap = dirMap
		blobMap = newBlobMap
		return nil
	})

	if err != nil {
		fmt.Printf("t walk err: %v\n", err)
		return nil, err
	}

	isCommitAvailable := false
	for _, object := range directoryMap {
		repoErr := r.checkRepositoryAlreadyExist()
		treeObject, err := r.findAndGetLastCommitObject(object.Name, object.Path)
		if err != nil && repoErr != nil {
			fmt.Printf("t is commitable: %v\n", err)
			return nil, err
		}

		if treeObject != nil {
			comparisonResult := r.compareTreeObject(treeObject.Path, *treeObject)
			if comparisonResult == "" {
				isCommitAvailable = true
				break
			}
		} else {
			isCommitAvailable = true
			break
		}

	}

	if !isCommitAvailable {
		return nil, errors.New("no changes have been made to any of the files")
	}

	for _, object := range directoryMap {
		err := r.createTreeObject(ref, object.TreeId, *object)
		if err != nil {
			fmt.Printf("t1: %v\n", err)
			return nil, err
		}
	}

	for _, object := range blobMap {
		err := r.createBlobObject(ref, object.BlobId, *object)
		if err != nil {
			fmt.Printf("t2: %v\n", err)
			return nil, err
		}
	}

	root := directoryMap[r.Path]

	return root, nil
}

func (r *Repository) createObject(ref, id string, object interface{}) error {
	if _, err := os.Stat(r.Path + "/.dirgod/objects/" + ref); os.IsNotExist(err) {
		err = os.MkdirAll(r.Path+"/.dirgod/objects/"+ref, 0755)
		if err != nil {
			return err
		}
	}

	objectFile, err := os.Create(r.Path + "/.dirgod/objects/" + ref + "/" + id)
	if err != nil {
		return err
	}

	defer objectFile.Close()

	serializedObject, err := json.Marshal(object)
	if err != nil {
		return err
	}

	err = binary.Write(objectFile, binary.BigEndian, serializedObject)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) createCommitLog(commitId, message string) error {
	_, err := os.Stat(r.Path + "/.dirgod/logs/main")
	if os.IsNotExist(err) {
		_, err = os.Create(r.Path + "/.dirgod/logs/main")
		if err != nil {
			return err
		}
	}

	readFile, err := os.OpenFile(r.Path+"/.dirgod/logs/main", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	var previousCommitId string = "0000000000000000000000000000000000000000"
	if len(fileLines) > 0 {
		lastCommit := fileLines[len(fileLines)-1]
		commitLog := &models.CommitLog{}
		commitLog.Deserialize(lastCommit)
		previousCommitId = commitLog.CommitId
	}

	var commitLog = models.CommitLog{
		PreviousCommitId: previousCommitId,
		CommitId:         commitId,
		Committer:        fmt.Sprintf("%s %s", r.InternalConfig.User.Username, r.InternalConfig.User.Email),
		Message:          message,
		Time:             time.Now().UnixMilli(),
	}

	_, err = readFile.WriteString(commitLog.Serialize())
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) loadObjectFile(content string) (interfaces.Object, error) {
	isTree := strings.Contains(content, "TreeId")
	if isTree {
		tree := &models.Tree{}
		err := json.Unmarshal([]byte(content), tree)
		if err != nil {
			return nil, err
		}

		return tree, nil
	}

	blob := &models.Blob{}
	err := json.Unmarshal([]byte(content), blob)
	if err != nil {
		return nil, err
	}

	return blob, nil
}

func (r *Repository) compareBlobObject(blob models.Blob) string {
	objectFilePath := ""
	err := filepath.Walk(r.Path+"/.dirgod/objects", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		blobFileContent, rErr := os.ReadFile(path)
		if rErr != nil {
			return rErr
		}

		blobObject, lErr := r.loadObjectFile(string(blobFileContent))
		if lErr != nil {
			return lErr
		}

		if blobObject.GetType() != "blob" {
			return nil
		}

		if blobObject.GetName() != blob.Name {
			return nil
		}

		if blobObject.GetChecksum() != blob.Checksum {
			return nil
		}

		objectFilePath, err = filepath.Abs(path)
		if err != nil {
			return err
		}

		return io.EOF
	})

	if err != nil && err != io.EOF {
		return ""
	}

	return objectFilePath
}

func (r *Repository) softLinkFile(path, target string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}

	filename := strings.Split(target, "/")
	symlink := filepath.Join(path, filename[len(filename)-1])
	err := os.Symlink(target, symlink)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) addBlobObject(path, filename string, directoryMap map[string]*models.Tree, blobMap map[string]*models.Blob) (map[string]*models.Tree, map[string]*models.Blob, error) {
	previousDir := r.getPreviousDirectory(path)

	newBlob := models.Blob{BlobId: utils.GenerateId(), Name: filename, Path: path}
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	newBlob.Content = content
	newBlob.Checksum = r.createBlobHash(path, filename, content)

	comparison := r.compareBlobObject(newBlob)
	if comparison != "" {
		blobFileContent, err := os.ReadFile(comparison)
		if err != nil {
			return nil, nil, err
		}

		err = json.Unmarshal(blobFileContent, &newBlob)
		if err != nil {
			return nil, nil, err
		}
	}

	_, ok := directoryMap[previousDir]
	if ok {
		directoryMap[previousDir].Blobs = append(directoryMap[previousDir].Blobs, newBlob.BlobId)
	}

	blobMap[newBlob.BlobId] = &newBlob

	return directoryMap, blobMap, nil
}

func (r *Repository) createBlobObject(ref, id string, newBlob models.Blob) error {
	objectPath := r.compareBlobObject(newBlob)
	if objectPath != "" {
		err := r.softLinkFile(r.Path+"/.dirgod/objects/"+ref, objectPath)
		if err != nil {
			return err
		}

		return nil
	}

	err := r.createObject(ref, newBlob.BlobId, newBlob)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) createTreeObject(ref, id string, newTree models.Tree) error {
	repoErr := r.checkRepositoryAlreadyExist()
	treeObject, err := r.findAndGetLastCommitObject(newTree.Name, newTree.Path)
	if err != nil && repoErr != nil {
		return err
	}

	if treeObject != nil {
		comparisonResult := r.compareTreeObject(treeObject.Path, *treeObject)
		if comparisonResult != "" {
			err = r.softLinkFile(r.Path+"/.dirgod/objects/"+ref, comparisonResult)
			if err != nil {
				return err
			}

			return nil
		}
	}

	err = r.createObject(ref, id, newTree)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) addTreeObject(directoryPath string, fileInfo fs.FileInfo, directoryMap map[string]*models.Tree) map[string]*models.Tree {
	previousDir := r.getPreviousDirectory(directoryPath)

	repoErr := r.checkRepositoryAlreadyExist()
	treeObject, err := r.findAndGetLastCommitObject(fileInfo.Name(), directoryPath)
	if err != nil && repoErr != nil {
		fmt.Printf("find err: %v\n", err)
		fmt.Printf("repo Err: %v\n", repoErr)
		return nil
	}

	newTree := models.Tree{TreeId: utils.GenerateId(), Name: fileInfo.Name(), Path: directoryPath}
	newTree.Checksum = r.createFolderHash(directoryPath)

	if treeObject != nil {
		comparisonResult := r.compareTreeObject(treeObject.Path, *treeObject)
		if comparisonResult != "" {
			treeContent, err := os.ReadFile(comparisonResult)
			if err != nil {
				fmt.Printf("read err: %v\n", err)
				return nil
			}

			err = json.Unmarshal(treeContent, &newTree)
			if err != nil {
				fmt.Printf("json err: %v\n", err)
				return nil
			}
		}
	}

	directoryMap[fileInfo.Name()] = &newTree
	if previousDir != "" {
		_, ok := directoryMap[previousDir]
		if ok {
			directoryMap[previousDir].Trees = append(directoryMap[previousDir].Trees, newTree.TreeId)
		}
	}

	return directoryMap
}

func (r *Repository) compareTreeObject(directoryPath string, tree models.Tree) string {
	treeHash := r.createFolderHash(directoryPath)
	if treeHash == tree.Checksum {
		path := r.findObjectPathFromObjectId(tree.TreeId)
		if path != "" {
			return path
		}
	}

	return ""
}

func (r *Repository) getPreviousDirectory(path string) string {
	splitPath := strings.Split(path, "/")
	previousDir := ""
	if len(splitPath) > 1 {
		previousDir = splitPath[len(splitPath)-2]
	}

	return previousDir
}

func (r *Repository) findAndGetLastCommitObject(name, path string) (*models.Tree, error) {
	snapshotPath := r.getCurrentSnapshotPath()
	if snapshotPath == "" {
		return nil, errors.New("couldn't find last commit object")
	}

	var objectFile *models.Tree
	err := filepath.Walk(snapshotPath, func(walkPath string, info fs.FileInfo, fErr error) error {
		if info.IsDir() {
			return nil
		}

		fileContent, err := os.ReadFile(walkPath)
		if err != nil {
			return err
		}

		object, err := r.loadObjectFile(string(fileContent))
		if err != nil {
			return err
		}

		if object.GetName() == name && object.GetChecksum() == r.createFolderHash(path) {
			commitObject := object.(*models.Tree)
			objectFile = commitObject

			return io.EOF
		}

		return nil
	})

	if err != nil && err != io.EOF {
		return nil, err
	}

	return objectFile, nil
}

func (r *Repository) getCurrentSnapshotPath() string {
	headFileContent, err := os.ReadFile(r.Path + "/.dirgod/HEAD")
	if err != nil {
		return ""
	}

	commitIdReference := strings.Split(string(headFileContent), "ref: ")
	if len(commitIdReference) < 1 {
		return ""
	}

	readCommitId, err := os.ReadFile(r.Path + "/.dirgod/" + commitIdReference[len(commitIdReference)-1])
	if err != nil {
		return ""
	}

	path := r.Path + "/.dirgod/objects/" + string(readCommitId)[:10]

	return path
}

func (r *Repository) createFolderHash(directoryPath string) string {
	var content string = ""
	err := filepath.Walk(directoryPath, func(path string, info fs.FileInfo, fErr error) error {
		if strings.HasPrefix(path, r.Path+"/.dirgod") && !strings.HasPrefix(path, r.Path+"/.dirgodignore") { //! tbc
			return nil
		}

		if strings.HasPrefix(path, r.Path+"/.git") && !strings.HasPrefix(path, r.Path+"/.gitignore") {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		fileContent, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		content += info.Name() + fmt.Sprintf("%x", md5.Sum([]byte(fileContent))) + path

		return nil
	})

	if err != nil {
		return ""
	}

	folderHash := fmt.Sprintf("%x", md5.Sum([]byte(content)))

	return folderHash
}

func (r *Repository) createBlobHash(blobPath, name string, fileContent []byte) string {
	content := name + blobPath + string(fileContent)
	hash := fmt.Sprintf("%x", md5.Sum([]byte(content)))

	return hash
}

func (r *Repository) findObjectPathFromObjectId(id string) string {
	var objectPath string = ""
	err := filepath.Walk(r.Path+"/.dirgod/objects", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if id == info.Name() {
			objectPath, err = filepath.Abs(path)
			if err != nil {
				return err
			}

			return io.EOF
		}

		return nil
	})

	if err != nil && err != io.EOF {
		return ""
	}

	return objectPath
}
