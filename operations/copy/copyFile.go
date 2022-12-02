package copy

import (
	"ena/dirgod/models"
	"ena/dirgod/operations"
	"os"
)

var DEFAULT_FILE_COPY_PERMS = os.FileMode(0644)

type CopyFile struct {
	Source         string
	Destination    string
	Options        models.OperationOptions
	Result         operations.OperationResult
	RollbackResult operations.OperationResult
}

func (c *CopyFile) Exec() {
	b, err := os.ReadFile(c.Source)
	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return
	}
	werr := os.WriteFile(c.Destination, b, DEFAULT_FILE_COPY_PERMS)
	if werr != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, werr)
		return
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
