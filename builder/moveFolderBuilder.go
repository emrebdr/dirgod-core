package builder

import (
	"errors"

	"github.com/emrebdr/dirgod-core/interfaces"
	"github.com/emrebdr/dirgod-core/models"
	"github.com/emrebdr/dirgod-core/operations/move"
	"github.com/emrebdr/dirgod-core/utils"
)

type MoveFolderBuilder struct {
	Source      string `json:"source"`
	WorkingMode string `json:"workingMode"`
	Cache       bool   `json:"cache"`
	operation   interfaces.Operation
}

func (m *MoveFolderBuilder) Build() (interfaces.Operation, error) {
	if m.Source == "" {
		return nil, errors.New("source is empty")
	}

	workingMode, err := utils.SetWorkingMode(m.WorkingMode)
	if err != nil {
		return nil, err
	}

	m.operation = &move.MoveFolder{
		Source: m.Source,
		Options: models.OperationOptions{
			WorkingMode: workingMode,
			Cache:       m.Cache,
		},
	}

	return m.operation, nil
}

func (m *MoveFolderBuilder) GetName() string {
	return "MoveFolder"
}

func (m *MoveFolderBuilder) IsValid() bool {
	return m.operation != nil
}
