package builder

import (
	"errors"

	"github.com/emrebdr/dirgod-core/interfaces"
	"github.com/emrebdr/dirgod-core/models"
	"github.com/emrebdr/dirgod-core/operations/delete"
	"github.com/emrebdr/dirgod-core/utils"
)

type DeleteFolderBuilder struct {
	Source      string `json:"source"`
	WorkingMode string `json:"workingMode"`
	Cache       bool   `json:"cache"`
	operation   interfaces.Operation
}

func (d *DeleteFolderBuilder) Build() (interfaces.Operation, error) {
	if d.Source == "" {
		return nil, errors.New("source is empty")
	}

	workingMode, err := utils.SetWorkingMode(d.WorkingMode)
	if err != nil {
		return nil, err
	}

	d.operation = &delete.DeleteFolder{
		Source: d.Source,
		Options: models.OperationOptions{
			WorkingMode: workingMode,
			Cache:       d.Cache,
		},
	}

	return d.operation, nil
}

func (d *DeleteFolderBuilder) GetName() string {
	return "DeleteFolder"
}

func (d *DeleteFolderBuilder) IsValid() bool {
	return d.operation != nil
}
