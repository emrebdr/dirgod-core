package create

import (
	"ena/dirgod/models"
	"ena/dirgod/operations"
	"os"
)

type CreateFile struct {
	Path           string
	Options        models.OperationOptions
	Result         operations.OperationResult
	RollbackResult operations.OperationResult
}

func (c *CreateFile) Exec() {
	file, err := os.Create(c.Path)
	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return
	}

	defer file.Close()

	c.Result.Completed = true
}

func (c *CreateFile) Rollback() {
	err := os.Remove(c.Path)
	if err != nil {
		c.RollbackResult.Completed = false
		c.RollbackResult.Err = err
		return
	}

	c.RollbackResult.Completed = true
}
