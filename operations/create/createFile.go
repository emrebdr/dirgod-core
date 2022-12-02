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
		c.decideErrorOutput(err)
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

func (c *CreateFile) decideErrorOutput(err error) {
	if c.Options.WorkingMode == models.Force {
		c.Result.Completed = true
		c.Result.Err = err
	} else if c.Options.WorkingMode == models.Default {
		// TODO: check for dependency tree
		c.Result.Completed = false
		c.Result.Err = err
	} else {
		c.Result.Completed = false
		c.Result.Err = err
	}
}
