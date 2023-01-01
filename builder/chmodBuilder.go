package builder

import (
	"ena/dirgod/interfaces"
	"ena/dirgod/models"
	"ena/dirgod/operations/ch"
	"errors"
	"os"
)

type ChmodBuilder struct {
	Source          string `json:"source"`
	Perm            int    `json:"perm"`
	Recursive       bool   `json:"recursive"`
	WorkingMode     string `json:"workingMode"`
	Cache           bool   `json:"cache"`
	createOperation interfaces.Operation
}

func (c *ChmodBuilder) Build() (interfaces.Operation, error) {
	if c.Source == "" {
		return nil, errors.New("source is empty")
	}

	if c.Perm <= 0 {
		return nil, errors.New("perm must be greater than 0")
	}

	workingMode, err := c.setWorkingMode()
	if err != nil {
		return nil, err
	}

	c.createOperation = &ch.Chmod{
		Source:    c.Source,
		PermCode:  os.FileMode(c.Perm),
		Recursive: c.Recursive,
		Options: models.OperationOptions{
			WorkingMode: workingMode,
			Cache:       c.Cache,
		},
	}

	return c.createOperation, nil
}

func (c *ChmodBuilder) GetName() string {
	return "Chmod"
}

func (c *ChmodBuilder) IsValid() bool {
	return c.createOperation != nil
}

func (c *ChmodBuilder) setWorkingMode() (models.Options, error) {
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
