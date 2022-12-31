package builder

import (
	"ena/dirgod/constants"
	"ena/dirgod/interfaces"
	"ena/dirgod/utils"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type NativeBuilder struct {
	base_path   string
	workingMode string
	cache       bool
	operations  []interfaces.Operation
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
	switch operationName {
	case "CreateFolder":
		var operationStruct CreateFolderBuilder
		return n.prepareOperation(&operationStruct, arguments)
	case "CreateFile":
		var operationStruct CreateFileBuilder
		return n.prepareOperation(&operationStruct, arguments)
	case "MoveFile":
        var operationStruct MoveFileBuilder
        return n.prepareOperation(&operationStruct, arguments)
	case "MoveFolder":
        var operationStruct MoveFolderBuilder
        return n.prepareOperation(&operationStruct, arguments)
	case "DeleteFile":
        var operationStruct DeleteFileBuilder
        return n.prepareOperation(&operationStruct, arguments)
	case "DeleteFolder":
        var operationStruct DeleteFolderBuilder
        return n.prepareOperation(&operationStruct, arguments)
	case "CopyFile":
        var operationStruct CopyFileBuilder
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
