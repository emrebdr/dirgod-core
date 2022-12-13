package builder

import (
	"encoding/json"
	"ena/dirgod/interfaces"
)

type NativeBuilder struct {
	Operation string
	Argumentes interface{}
	operations []interfaces.Operation
}

func (n *NativeBuilder) AddOperation(operation interfaces.Operation) {
	n.operations = append(n.operations, operation)
}

func (n *NativeBuilder) GetOperations() []interfaces.Operation {
	return n.operations
}

func (n *NativeBuilder) CreateNewOperation() (error) {
	switch n.Operation {
	case "CreateFolder":
		var operationStruct CreateFolderBuilder
		encoded, err := json.Marshal(n.Argumentes)

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
		encoded, err := json.Marshal(n.Argumentes)
		
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
