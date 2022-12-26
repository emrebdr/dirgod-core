package copy

import (
	"ena/dirgod/models"
	"ena/dirgod/operations"
	"fmt"
)

type CopyFolder struct {
	Source         string
	Destination    string
	Options        models.OperationOptions
	Result         operations.OperationResult
	RollbackResult operations.OperationResult
}

func (c *CopyFolder) Exec() {
	fmt.Println("Copying a folder")
}

func (c *CopyFolder) Rollback() {
	fmt.Println("Copying a folder")
}
