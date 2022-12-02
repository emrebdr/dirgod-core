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
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
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
