package builder

import (
	"errors"

	"github.com/emrebdr/dirgod-core/interfaces"
	"github.com/emrebdr/dirgod-core/models"
	"github.com/emrebdr/dirgod-core/operations/move"
	"github.com/emrebdr/dirgod-core/utils"
)

type MoveFileBuilder struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	WorkingMode string `json:"workingMode"`
	Cache       bool   `json:"cache"`
	operation   interfaces.Operation
}

func (m *MoveFileBuilder) Build() (interfaces.Operation, error) {
	if m.Source == "" {
		return nil, errors.New("source is empty")
	}

	if m.Destination == "" {
		return nil, errors.New("destination is empty")
	}

	workingMode, err := utils.SetWorkingMode(m.WorkingMode)
	if err != nil {
		return nil, err
	}

	m.operation = &move.MoveFile{
		Source:      m.Source,
		Destination: m.Destination,
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
