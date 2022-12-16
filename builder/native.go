package builder

import (
	"ena/dirgod/interfaces"
	"encoding/json"
	"errors"
	"os"
)

type NativeBuilder struct {
	base_path   string
	workingMode string
	cache       string
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
		switch workingMode {
		case "force":
			n.workingMode = workingMode
			return n
		case "strict":
			n.workingMode = workingMode
			return n
		case "default":
			n.workingMode = workingMode
			return n
		default:
			n.workingMode = "default"
			return n
		}
	}

	n.workingMode = "default"

	return n
}

func (n *NativeBuilder) SetCacheMode(cache string) *NativeBuilder {
	if cache != "" {
		switch cache {
		case "enable":
			n.cache = cache
			return n
		case "disable":
			n.cache = cache
			return n
		default:
			n.cache = "disable"
			return n
		}
	}

	n.cache = "disable"

	return n
}

func (n *NativeBuilder) SetBasePath(basePath string) *NativeBuilder {
	n.base_path = basePath
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
	}

	return nil, errors.New("unknown operation")
}

func (n *NativeBuilder) setDefaultValues(operation interfaces.Builder) (interfaces.Builder, error) {
	defaultArguments := map[string]interface{}{
		"path":        n.base_path,
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

	n.addOperation(operation)
	return nil
}
