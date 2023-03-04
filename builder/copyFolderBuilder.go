package builder

import (
	"errors"

	"github.com/emrebdr/dirgod-code/interfaces"
	"github.com/emrebdr/dirgod-code/models"
	"github.com/emrebdr/dirgod-code/operations/copy"
	"github.com/emrebdr/dirgod-code/utils"
)

type CopyFolderBuilder struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	WorkingMode string `json:"workingMode"`
	Cache       bool   `json:"cache"`
	operation   interfaces.Operation
}

func (c *CopyFolderBuilder) Build() (interfaces.Operation, error) {
	if c.Source == "" {
		return nil, errors.New("source is empty")
	}

	if c.Destination == "" {
		return nil, errors.New("destination is empty")
	}

	workingMode, err := utils.SetWorkingMode(c.WorkingMode)
	if err != nil {
		return nil, err
	}

	c.operation = &copy.CopyFolder{
		Source:      c.Source,
		Destination: c.Destination,
		Options: models.OperationOptions{
			WorkingMode: workingMode,
			Cache:       c.Cache,
		},
	}

	return c.operation, nil
}

func (c *CopyFolderBuilder) GetName() string {
	return "CopyFolder"
}

func (c *CopyFolderBuilder) IsValid() bool {
	return c.operation != nil
}
