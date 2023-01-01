package builder

import (
	"ena/dirgod/interfaces"
	"ena/dirgod/models"
	"ena/dirgod/operations/copy"
	"errors"
)

type CopyFolderBuilder struct {
	Source          string `json:"source"`
	Destination     string `json:"destination"`
	WorkingMode     string `json:"workingMode"`
	Cache           bool   `json:"cache"`
	createOperation interfaces.Operation
}

func (c *CopyFolderBuilder) Build() (interfaces.Operation, error) {
	if c.Source == "" {
		return nil, errors.New("source is empty")
	}

	if c.Destination == "" {
		return nil, errors.New("destination is empty")
	}

	workingMode, err := c.setWorkingMode()
	if err != nil {
		return nil, err
	}

	c.createOperation = &copy.CopyFolder{
		Source:      c.Source,
		Destination: c.Destination,
		Options: models.OperationOptions{
			WorkingMode: workingMode,
			Cache:       c.Cache,
		},
	}

	return c.createOperation, nil
}

func (c *CopyFolderBuilder) GetName() string {
	return "CopyFolder"
}

func (c *CopyFolderBuilder) IsValid() bool {
	return c.createOperation != nil
}

func (c *CopyFolderBuilder) setWorkingMode() (models.Options, error) {
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