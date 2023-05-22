package create

import (
	"os"

	"github.com/emrebdr/dirgod-core/models"
	"github.com/emrebdr/dirgod-core/operations"
)

type CreateFile struct {
	Source         string
	Options        models.OperationOptions
	Result         operations.OperationResult
	RollbackResult operations.OperationResult
}

func (c *CreateFile) Exec() operations.OperationResult {
	file, err := os.Create(c.Source)
	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return c.Result
	}

	defer file.Close()

	c.Result.Completed = true

	return c.Result
}

func (c *CreateFile) Rollback() {
	err := os.Remove(c.Source)
	if err != nil {
		c.RollbackResult.Completed = false
		c.RollbackResult.Err = err
		return
	}

	c.RollbackResult.Completed = true
}
