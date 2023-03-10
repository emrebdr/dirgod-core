package builder

import (
	"errors"

	"github.com/emrebdr/dirgod-core/interfaces"
	"github.com/emrebdr/dirgod-core/models"
	"github.com/emrebdr/dirgod-core/operations/create"
	"github.com/emrebdr/dirgod-core/utils"
)

type CreateFolderBuilder struct {
	Source      string `json:"source"`
	WorkingMode string `json:"workingMode"`
	Cache       bool   `json:"cache"`
	operation   interfaces.Operation
}

func (c *CreateFolderBuilder) Build() (interfaces.Operation, error) {
	if c.Source == "" {
		return nil, errors.New("source is empty")
	}

	workingMode, err := utils.SetWorkingMode(c.WorkingMode)
	if err != nil {
		return nil, err
	}

	c.operation = &create.CreateFolder{
		Source: c.Source,
		Options: models.OperationOptions{
			WorkingMode: workingMode,
			Cache:       c.Cache,
		},
	}

	return c.operation, nil
}

func (c *CreateFolderBuilder) GetName() string {
	return "CreateFolder"
}

func (c *CreateFolderBuilder) IsValid() bool {
	return c.operation != nil
}
