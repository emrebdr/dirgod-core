package builder

import (
	"ena/dirgod/interfaces"
	"ena/dirgod/models"
	"ena/dirgod/operations/move"
	"ena/dirgod/utils"
	"errors"
)

type MoveFileBuilder struct {
	Source      string `json:"source"`
	WorkingMode string `json:"workingMode"`
	Cache       bool   `json:"cache"`
	operation   interfaces.Operation
}

func (m *MoveFileBuilder) Build() (interfaces.Operation, error) {
	if m.Source == "" {
		return nil, errors.New("source is empty")
	}

	workingMode, err := utils.SetWorkingMode(m.WorkingMode)
	if err != nil {
		return nil, err
	}

	m.operation = &move.MoveFile{
		Source: m.Source,
		Options: models.OperationOptions{
			WorkingMode: workingMode,
			Cache:       m.Cache,
		},
	}

	return m.operation, nil
}

func (m *MoveFileBuilder) GetName() string {
	return "MoveFile"
}

func (m *MoveFileBuilder) IsValid() bool {
	return m.operation != nil
}
