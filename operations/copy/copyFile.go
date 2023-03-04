package copy

import (
	"io"
	"os"

	"github.com/emrebdr/dirgod-code/models"
	"github.com/emrebdr/dirgod-code/operations"
)

type CopyFile struct {
	Source         string
	Destination    string
	Options        models.OperationOptions
	Result         operations.OperationResult
	RollbackResult operations.OperationResult
}

func (c *CopyFile) Exec() {
	srcFd, err := os.Open(c.Source)
	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return
	}
	defer srcFd.Close()
	dstFd, err := os.Create(c.Destination)
	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return
	}
	defer dstFd.Close()

	if _, err = io.Copy(dstFd, srcFd); err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return
	}
	srcInfo, err := os.Stat(c.Source)
	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return
	}

	err = os.Chmod(c.Destination, srcInfo.Mode())
	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
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
