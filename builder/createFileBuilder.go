package builder

import (
	"ena/dirgod/interfaces"
	"ena/dirgod/models"
	"ena/dirgod/operations/create"
	"errors"
)

type CreateFileBuilder struct {
	Path            string               `json:"path"`
	WorkingMode     string               `json:"workingMode"`
	Cache           string               `json:"cache"`
	createOperation interfaces.Operation `json:"-"`
}

func (c *CreateFileBuilder) Build() (interfaces.Operation, error) {
	if c.Path == "" {
		return nil, errors.New("path is empty")
	}

	workingMode, err := c.setWorkingMode()
	if err != nil {
		return nil, err
	}

	cache, err := c.setCache()
	if err != nil {
		return nil, err
	}

	c.createOperation = &create.CreateFile{
		Path: c.Path,
		Options: models.OperationOptions{
			WorkingMode: workingMode,
			Cache:       cache,
		},
	}

	return c.createOperation, nil
}

func (c *CreateFileBuilder) GetName() string {
	return "CreateFile"
}

func (c *CreateFileBuilder) setWorkingMode() (models.Options, error) {
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

func (c *CreateFileBuilder) setCache() (bool, error) {
	if c.Cache != "" {
		switch c.Cache {
		case "enable":
			return true, nil
		case "disable":
			return false, nil
		default:
			return false, errors.New("unknown cache mode")
		}
	}

	return false, nil
}
