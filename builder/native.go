package builder

import (
	"ena/dirgod/interfaces"
	"encoding/json"
)

type NativeBuilder struct {
	operations []interfaces.Operation
}

func (n *NativeBuilder) AddOperation(operation interfaces.Operation) {
	n.operations = append(n.operations, operation)
}

func (n *NativeBuilder) GetOperations() []interfaces.Operation {
	return n.operations
}

func (n *NativeBuilder) CreateNewOperation(operationName string, arguments interface{}) error {
	switch operationName {
	case "CreateFolder":
		var operationStruct CreateFolderBuilder
		encoded, err := json.Marshal(arguments)

		if err != nil {
			return err
		}
		err = json.Unmarshal(encoded, &operationStruct)
		if err != nil {
			return err
		}

		operation, err := operationStruct.Build()
		if err != nil {
			return err
		}

		n.AddOperation(operation)
		return nil
	case "CreateFile":
		var operationStruct CreateFileBuilder
		encoded, err := json.Marshal(arguments)

		if err != nil {
			return err
		}

		err = json.Unmarshal(encoded, &operationStruct)
		if err != nil {
			return err
		}

		operation, err := operationStruct.Build()
		if err != nil {
			return err
		}

		n.AddOperation(operation)
		return nil
	}

	return nil
}
