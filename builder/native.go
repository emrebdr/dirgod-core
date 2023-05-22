package builder

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/emrebdr/dirgod-core/constants"
	"github.com/emrebdr/dirgod-core/interfaces"
	"github.com/emrebdr/dirgod-core/repository"
	"github.com/emrebdr/dirgod-core/utils"
)

type NativeBuilder struct {
	base_path   string
	workingMode string
	cache       bool
	operations  []interfaces.Operation
}

func NewNativeBuilder() *NativeBuilder {
	return &NativeBuilder{}
}

func (n *NativeBuilder) addOperation(operation interfaces.Operation) {
	n.operations = append(n.operations, operation)
}

func (n *NativeBuilder) GetOperations() []interfaces.Operation {
	return n.operations
}

func (n *NativeBuilder) SetWorkingMode(workingMode string) *NativeBuilder {
	if workingMode != "" {
		if utils.Contains(constants.WorkingMode, strings.ToLower(workingMode)) {
			n.workingMode = workingMode
			return n
		} else {
			fmt.Println("Unknown working mode. Set default working mode")
		}
	}

	n.workingMode = "default"

	return n
}

func (n *NativeBuilder) SetCacheMode(cache string) *NativeBuilder {
	if cache != "" {
		if strings.ToLower(cache) == "true" {
			n.cache = true
			return n
		}
	}

	n.cache = false

	return n
}

func (n *NativeBuilder) SetBasePath(basePath string) *NativeBuilder {
	if filepath.IsAbs(basePath) {
		n.base_path = basePath
		return n
	}

	n.base_path = filepath.Join(os.Getenv("PWD"), basePath)
	return n
}

func (n *NativeBuilder) decodeOperations(operationName string, arguments []interface{}) (interfaces.Builder, error) {
	operationName = strings.ToLower(operationName)
	switch operationName {
	case "createfolder":
		var operationStruct CreateFolderBuilder
		return n.prepareOperation(&operationStruct, arguments)
	case "createfile":
		var operationStruct CreateFileBuilder
		return n.prepareOperation(&operationStruct, arguments)
	case "movefile":
		var operationStruct MoveFileBuilder
		return n.prepareOperation(&operationStruct, arguments)
	case "movefolder":
		var operationStruct MoveFolderBuilder
		return n.prepareOperation(&operationStruct, arguments)
	case "deletefile":
		var operationStruct DeleteFileBuilder
		return n.prepareOperation(&operationStruct, arguments)
	case "deletefolder":
		var operationStruct DeleteFolderBuilder
		return n.prepareOperation(&operationStruct, arguments)
	case "copyfile":
		var operationStruct CopyFileBuilder
		return n.prepareOperation(&operationStruct, arguments)
	case "copyfolder":
		var operationStruct CopyFolderBuilder
		return n.prepareOperation(&operationStruct, arguments)
	case "chmod":
		var operationStruct ChmodBuilder
		return n.prepareOperation(&operationStruct, arguments)
	case "chown":
		var operationStruct ChownBuilder
		return n.prepareOperation(&operationStruct, arguments)
	case "write":
		var operationStruct WriteBuilder
		return n.prepareOperation(&operationStruct, arguments)
	}

	return nil, errors.New("unknown operation")
}

func (n *NativeBuilder) setDefaultValues(operation interfaces.Builder) (interfaces.Builder, error) {
	defaultArguments := map[string]interface{}{
		"source":      n.base_path,
		"workingMode": n.workingMode,
		"cache":       n.cache,
	}

	encoded, err := json.Marshal(defaultArguments)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(encoded, operation)
	if err != nil {
		return nil, err
	}

	return operation, nil
}

func (n *NativeBuilder) decodeArguments(operation interfaces.Builder, arguments []interface{}) (interfaces.Builder, error) {
	for _, argument := range arguments {
		encoded, err := json.Marshal(argument)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(encoded, operation)
		if err != nil {
			return nil, err
		}
	}

	return operation, nil
}

func (n *NativeBuilder) prepareOperation(operation interfaces.Builder, arguments []interface{}) (interfaces.Builder, error) {
	_, err := n.setDefaultValues(operation)
	if err != nil {
		return nil, err
	}

	_, err = n.decodeArguments(operation, arguments)

	if err != nil {
		return nil, err
	}

	return operation, nil
}

func (n *NativeBuilder) CreateNewOperation(operationName string, arguments ...interface{}) error {
	if n.base_path == "" {
		path, err := os.Getwd()
		if err != nil {
			return err
		}
		n.base_path = path
	}

	operationStruct, err := n.decodeOperations(operationName, arguments)
	if err != nil {
		return err
	}

	operation, err := operationStruct.Build()
	if err != nil {
		return err
	}

	if isValid := operationStruct.IsValid(); !isValid {
		return errors.New("invalid operation")
	}

	n.addOperation(operation)
	return nil
}

func (n *NativeBuilder) Execute(commitMessage []byte) error {
	repo := repository.LoadRepository()
	if repo == nil {
		currDir, err := os.Getwd()
		if err != nil {
			return err
		}

		splitName := strings.Split(currDir, "/")
		folderName := splitName[len(splitName)-1]
		repo = repository.Init(folderName, "")
	}

	for _, operation := range n.operations {
		result := operation.Exec()
		if !result.Completed {
			return result.Err
		}
	}

	if commitMessage == nil {
		commitMessage = []byte(utils.GenerateId())
	}

	commitResult := repo.Commit(string(commitMessage))
	if commitResult != nil {
		return commitResult
	}

	return nil
}
