package builder

import (
	"errors"

	"github.com/emrebdr/dirgod-core/interfaces"
	"github.com/emrebdr/dirgod-core/models"
	"github.com/emrebdr/dirgod-core/operations/delete"
	"github.com/emrebdr/dirgod-core/utils"
)

type DeleteFileBuilder struct {
	Source      string `json:"source"`
	WorkingMode string `json:"workingMode"`
	Cache       bool   `json:"cache"`
	operation   interfaces.Operation
}

func (c *DeleteFileBuilder) Build() (interfaces.Operation, error) {
	if c.Source == "" {
		return nil, errors.New("source is empty")
	}

	workingMode, err := utils.SetWorkingMode(c.WorkingMode)
	if err != nil {
		return nil, err
	}

	c.operation = &delete.DeleteFile{
		Source: c.Source,
		Options: models.OperationOptions{
			WorkingMode: workingMode,
			Cache:       c.Cache,
		},
	}

	return c.operation, nil
}

func (c *DeleteFileBuilder) GetName() string {
	return "DeleteFile"
}

func (c *DeleteFileBuilder) IsValid() bool {
	return c.operation != nil
}
