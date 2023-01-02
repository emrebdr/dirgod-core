package builder

import (
	"ena/dirgod/interfaces"
	"ena/dirgod/models"
	"ena/dirgod/operations/ch"
	"ena/dirgod/utils"
	"errors"
)

type ChownBuilder struct {
	Source      string `json:"source"`
	UID         int    `json:"uid"`
	GID         int    `json:"gid"`
	Recursive   bool   `json:"recursive"`
	WorkingMode string `json:"workingMode"`
	Cache       bool   `json:"cache"`
	operation   interfaces.Operation
}

func (c *ChownBuilder) Build() (interfaces.Operation, error) {
	if c.Source == "" {
		return nil, errors.New("source is empty")
	}

	if c.GID <= 0 || c.UID <= 0 {
		return nil, errors.New("UID and GID must be greater than 0")
	}

	workingMode, err := utils.SetWorkingMode(c.WorkingMode)
	if err != nil {
		return nil, err
	}

	c.operation = &ch.Chown{
		Source:    c.Source,
		UID:       c.UID,
		GID:       c.GID,
		Recursive: c.Recursive,
		Options: models.OperationOptions{
			WorkingMode: workingMode,
			Cache:       c.Cache,
		},
	}

	return c.operation, nil
}

func (c *ChownBuilder) GetName() string {
	return "Chown"
}

func (c *ChownBuilder) IsValid() bool {
	return c.operation != nil
}
