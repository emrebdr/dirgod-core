package copy

import (
	"ena/dirgod/models"
	"fmt"
)

type CopyFile struct {
	options models.OperationOptions
}

func (c *CopyFile) Exec() {
	fmt.Println("Copying a file")
}
