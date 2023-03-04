package builder

import (
	"errors"
	"os"

	"github.com/emrebdr/dirgod-code/interfaces"
	"github.com/emrebdr/dirgod-code/models"
	"github.com/emrebdr/dirgod-code/operations/ch"
	"github.com/emrebdr/dirgod-code/utils"
)

type ChmodBuilder struct {
	Source      string `json:"source"`
	Perm        int    `json:"perm"`
	Recursive   bool   `json:"recursive"`
	WorkingMode string `json:"workingMode"`
	Cache       bool   `json:"cache"`
	operation   interfaces.Operation
}

func (c *ChmodBuilder) Build() (interfaces.Operation, error) {
	if c.Source == "" {
		return nil, errors.New("source is empty")
	}

	if c.Perm <= 0 {
		return nil, errors.New("perm must be greater than 0")
	}

	workingMode, err := utils.SetWorkingMode(c.WorkingMode)
	if err != nil {
		return nil, err
	}

	c.operation = &ch.Chmod{
		Source:    c.Source,
		PermCode:  os.FileMode(c.Perm),
		Recursive: c.Recursive,
		Options: models.OperationOptions{
			WorkingMode: workingMode,
			Cache:       c.Cache,
		},
	}

	return c.operation, nil
}

func (c *ChmodBuilder) GetName() string {
	return "Chmod"
}

func (c *ChmodBuilder) IsValid() bool {
	return c.operation != nil
}
