package builder

import (
	"ena/dirgod/interfaces"
	"ena/dirgod/models"
	"ena/dirgod/operations/create"
	"errors"
)

type CreateFolderBuilder struct {
	Source          string               `json:"source"`
	WorkingMode     string               `json:"workingMode"`
	Cache           bool                 `json:"cache"`
	createOperation interfaces.Operation `json:"-"`
}

func (c *CreateFolderBuilder) Build() (interfaces.Operation, error) {
	if c.Source == "" {
		return nil, errors.New("source is empty")
	}

	workingMode, err := c.setWorkingMode()
	if err != nil {
		return nil, err
	}

	c.createOperation = &create.CreateFolder{
		Source: c.Source,
		Options: models.OperationOptions{
			WorkingMode: workingMode,
			Cache:       c.Cache,
		},
	}

	return c.createOperation, nil
}

func (c *CreateFolderBuilder) GetName() string {
	return "CreateFolder"
}

func (c *CreateFolderBuilder) IsValid() bool {
	return c.createOperation != nil
}

func (c *CreateFolderBuilder) setWorkingMode() (models.Options, error) {
	if c.WorkingMode != "" {
		switch c.WorkingMode {
		case "force":
			return models.Force, nil
		case "strict":
			return models.Strict, nil
		case "default":
			return models.Default, nil
		default:
			return -1, errors.New("unknown working mode")
		}
	}

	return models.Default, nil
}
