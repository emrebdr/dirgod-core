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

func (c *CreateFile) Exec() {
	file, err := os.Create(c.Source)
	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return
	}

	defer file.Close()

	c.Result.Completed = true
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
