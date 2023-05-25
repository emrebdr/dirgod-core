package repository

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/emrebdr/dirgod-core/models"
)

type Repository struct {
	Name           string                      `json:"name"`
	Description    string                      `json:"description"`
	CommitObject   models.Commit               `json:"commit"`
	RootConfig     models.RootDirgodConfig     `json:"rootConfig"`
	InternalConfig models.InternalDirgodConfig `json:"internalConfig"`
	Path           string                      `json:"path"`
}

func Init(name, description, path string) *Repository {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil
	}

	repo := &Repository{Name: name, Description: description, Path: absPath}
	repo.initializeRepository()
	return repo
}

func LoadRepository() *Repository {
	repo := &Repository{}
	err := filepath.Walk("./", func(path string, info fs.FileInfo, err error) error {
		if info.Name() == "reporef" {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			err = json.Unmarshal(content, &repo)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil || repo.Name == "" {
		return nil
	}

	return repo
}

func (r *Repository) checkRepositoryAlreadyExist() error {
	if _, err := os.Stat(r.Path + "/.dirgod/reporef"); !os.IsNotExist(err) {
		return errors.New("repository already initialized")
	}

	return nil
}

func (r *Repository) initializeRepository() error {
	err := r.checkRepositoryAlreadyExist()
	if err != nil {
		return err
	}

	err = os.MkdirAll(r.Path+"/.dirgod/objects", 0755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(r.Path+"/.dirgod/refs/branches", 0755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(r.Path+"/.dirgod/refs/commits", 0755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(r.Path+"/.dirgod/logs", 0755)
	if err != nil {
		return err
	}

	if _, err = os.Stat(r.Path + "/Dirgod"); os.IsNotExist(err) {
		_, err = os.Create(r.Path + "/Dirgod")
		if err != nil {
			return err
		}
	}

	_, err = os.Create(r.Path + "/.dirgodignore")
	if err != nil {
		return err
	}

	err = r.createREADMEFile()
	if err != nil {
		return err
	}

	err = r.createDescriptionFile()
	if err != nil {
		return err
	}

	err = r.loadRootConfigFile()
	if err != nil {
		return err
	}

	err = r.loadInternalConfigFile()
	if err != nil {
		return err
	}

	err = r.createCommitObject("Initial Commit")
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) loadRootConfigFile() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configFile, err := os.Open(homeDir + "/.dirgodconfig")
	if err != nil {
		return err
	}

	configFileInfo, err := configFile.Stat()
	if err != nil {
		return err
	}

	contentBytes := make([]byte, configFileInfo.Size())
	byteCount, err := configFile.Read(contentBytes)
	if err != nil {
		return err
	}

	if byteCount != int(configFileInfo.Size()) {
		return errors.New("something went wrong while reading root config file")
	}

	var config models.RootDirgodConfig
	err = json.Unmarshal(contentBytes, &config)
	if err != nil {
		return err
	}

	r.RootConfig = config
	return nil
}

func (r *Repository) loadInternalConfigFile() error {
	config := models.InternalDirgodConfig{}

	_, err := os.Stat(r.Path + "/.dirgod/config")
	if os.IsNotExist(err) {
		config.User = r.RootConfig.User

		file, err := os.Create(r.Path + "/.dirgod/config")
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return err
		}

		content, err := json.Marshal(config)
		if err != nil {
			return err
		}

		_, err = file.Write(content) //! error check for written byte count
		if err != nil {
			return err
		}
	} else {
		configContent, err := os.ReadFile(r.Path + "/.dirgod/config")
		if err != nil {
			return err
		}

		err = json.Unmarshal(configContent, &config)
		if err != nil {
			return err
		}
	}

	r.InternalConfig = config
	return nil
}

func (r *Repository) createREADMEFile() error {
	file, err := os.Create(r.Path + "/README.md")
	if err != nil {
		return err
	}

	text := "# " + r.Name + "\n\n" + "## " + r.Description

	writtenLen, err := file.Write([]byte(text))
	if err != nil {
		return err
	}

	if len([]byte(text)) != writtenLen {
		return errors.New("something went wrong while preparing README.md file")
	}

	return nil
}

func (r *Repository) createDescriptionFile() error {
	file, err := os.Create(r.Path + "/.dirgod/description")
	if err != nil {
		return err
	}

	context := r.Description
	if context == "" {
		context = "Unnamed repository; edit this file 'description' to name the repository."
	}

	byteCount, err := file.Write([]byte(context))
	if err != nil {
		return err
	}

	if len([]byte(context)) != byteCount {
		return errors.New("something went wrong while preparing description file")
	}

	return nil
}

func (r *Repository) saveReporefFile(repository interface{}) error {
	objectFile, err := os.Create(r.Path + "/.dirgod/reporef")
	if err != nil {
		return err
	}

	defer objectFile.Close()

	serializedObject, err := json.Marshal(repository)
	if err != nil {
		return err
	}

	err = binary.Write(objectFile, binary.BigEndian, serializedObject)
	if err != nil {
		return err
	}

	return nil
}
