package copy

import (
	"io"
	"os"

	"github.com/emrebdr/dirgod-core/models"
	"github.com/emrebdr/dirgod-core/operations"
)

type CopyFile struct {
	Source         string
	Destination    string
	Options        models.OperationOptions
	Result         operations.OperationResult
	RollbackResult operations.OperationResult
}

func (c *CopyFile) Exec() operations.OperationResult {
	srcFd, err := os.Open(c.Source)
	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return c.Result
	}
	defer srcFd.Close()
	dstFd, err := os.Create(c.Destination)
	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return c.Result
	}
	defer dstFd.Close()

	if _, err = io.Copy(dstFd, srcFd); err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return c.Result
	}
	srcInfo, err := os.Stat(c.Source)
	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return c.Result
	}

	err = os.Chmod(c.Destination, srcInfo.Mode())
	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return c.Result
	}

	c.Result.Completed = true
	return c.Result
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
