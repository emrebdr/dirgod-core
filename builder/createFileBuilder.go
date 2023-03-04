package builder

import (
	"errors"

	"github.com/emrebdr/dirgod-code/interfaces"
	"github.com/emrebdr/dirgod-code/models"
	"github.com/emrebdr/dirgod-code/operations/create"
	"github.com/emrebdr/dirgod-code/utils"
)

type CreateFileBuilder struct {
	Source      string `json:"source"`
	WorkingMode string `json:"workingMode"`
	Cache       bool   `json:"cache"`
	operation   interfaces.Operation
}

func (c *CreateFileBuilder) Build() (interfaces.Operation, error) {
	if c.Source == "" {
		return nil, errors.New("source is empty")
	}

	workingMode, err := utils.SetWorkingMode(c.WorkingMode)
	if err != nil {
		return nil, err
	}

	c.operation = &create.CreateFile{
		Source: c.Source,
		Options: models.OperationOptions{
			WorkingMode: workingMode,
			Cache:       c.Cache,
		},
	}

	return c.operation, nil
}

func (c *CreateFileBuilder) GetName() string {
	return "CreateFile"
}

func (c *CreateFileBuilder) IsValid() bool {
	return c.operation != nil
}
