package builder

import (
	"errors"
	"ena/dirgod/interfaces"
	"ena/dirgod/operations/create"
	"ena/dirgod/models"
)

type CreateFolderBuilder struct {
	Path string `json:"path"`
	WorkingMode string `json:"workingMode"`
	Cache string `json:"cache"`
	createOperation interfaces.Operation `json:"-"`
}

func (c *CreateFolderBuilder) Build() (interfaces.Operation, error) {
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

	c.createOperation = &create.CreateFolder{
		Path: c.Path,
		Options: models.OperationOptions{
			WorkingMode: workingMode,
			Cache: cache,
		},
	}

	return c.createOperation, nil
}

func (c *CreateFolderBuilder) setWorkingMode() (models.Options, error) {
	if c.WorkingMode != "" {
		switch c.WorkingMode {
		case "Force":
			return models.Force, nil
		case "Strict":
			return models.Strict, nil
		case "Default":
			return models.Default, nil
		default:
			return -1, errors.New("Unknown working mode")
		}
	}

	return models.Default, nil
}

func (c *CreateFolderBuilder) setCache() (bool, error) {
	if c.Cache != "" {
		switch c.Cache {
		case "Enable":
			return true, nil
		case "Disable":
			return false, nil
		default:
			return false, errors.New("Unknown cache option")
		}
	}

	return false, nil
}