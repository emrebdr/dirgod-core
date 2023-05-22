package create

import (
	"os"

	"github.com/emrebdr/dirgod-core/models"
	"github.com/emrebdr/dirgod-core/operations"
)

var DEFAULT_CREATE_FOLDER_PERMISSION = os.FileMode(0755)

type CreateFolder struct {
	Source         string
	Options        models.OperationOptions
	Result         operations.OperationResult
	RollbackResult operations.OperationResult
}

func (c *CreateFolder) Exec() operations.OperationResult {
	var err error
	if c.Options.WorkingMode == models.Force {
		err = os.MkdirAll(c.Source, DEFAULT_CREATE_FOLDER_PERMISSION)
	} else {
		err = os.Mkdir(c.Source, DEFAULT_CREATE_FOLDER_PERMISSION)
	}

	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return c.Result
	}

	c.Result.Completed = true

	return c.Result
}

func (c *CreateFolder) Rollback() {
	// Remove single empty folder or force remove all recursive child folders and files with checking WorkingMode
	err := os.RemoveAll(c.Source)
	if err != nil {
		c.RollbackResult.Completed = false
		c.RollbackResult.Err = err
		return
	}

	c.RollbackResult.Completed = true
}
