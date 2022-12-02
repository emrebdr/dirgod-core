package copy

import (
	"ena/dirgod/models"
	. "ena/dirgod/operations"
	"os"
)

var DEFAULT_FILE_COPY_PERMS = os.FileMode(0644)

type CopyFile struct {
	Source         string
	Destination    string
	Options        models.OperationOptions
	Result         OperationResult
	RollbackResult OperationResult
}

func (c *CopyFile) Exec() {
	b, err := os.ReadFile(c.Source)
	if err != nil {
		c.decideErrorOutput(err)
		return
	}
	werr := os.WriteFile(c.Destination, b, DEFAULT_FILE_COPY_PERMS)
	if werr != nil {
		c.decideErrorOutput(werr)
	}

	c.Result.Completed = true
}

// Rolling an copy operation is just deleting the copied file
func (c *CopyFile) Rollback() {
	err := os.Remove(c.Destination)
	if err != nil {
		c.RollbackResult.Completed = false
		c.RollbackResult.Err = err
	} else {
		c.RollbackResult.Completed = true
	}
}

func (c *CopyFile) decideErrorOutput(err error) {
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
