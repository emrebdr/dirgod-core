package create

import (
	"ena/dirgod/models"
	"ena/dirgod/operations"
	"os"
)

var DEFAUL_CREATE_FOLDER_PERMISSION = os.FileMode(0755)

type CreateFolder struct {
	Path           string
	Options        models.OperationOptions
	Result         operations.OperationResult
	RollbackResult operations.OperationResult
}

func (c *CreateFolder) Exec() {
	var err error
	if c.Options.WorkingMode == models.Force {
		err = os.MkdirAll(c.Path, DEFAUL_CREATE_FOLDER_PERMISSION)
	} else {
		err = os.Mkdir(c.Path, DEFAUL_CREATE_FOLDER_PERMISSION)
	}

	if err != nil {
		c.decideErrorOutput(err)
		return
	}

	c.Result.Completed = true
}

func (c *CreateFolder) Rollback() {
	// Remove single empty folder or force remove all recursive child folders and files with checking WorkingMode
	err := os.RemoveAll(c.Path)
	if err != nil {
		c.RollbackResult.Completed = false
		c.RollbackResult.Err = err
		return
	}

	c.RollbackResult.Completed = true
}

func (c *CreateFolder) decideErrorOutput(err error) {
	if c.Options.WorkingMode == models.Force {
		c.Result.Completed = true
		c.Result.Err = err
	} else if c.Options.WorkingMode == models.Default {
		c.Result.Completed = false
		c.Result.Err = err
	} else {
		c.Result.Completed = false
		c.Result.Err = err
	}
}
