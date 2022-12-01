package copy

import (
	"ena/dirgod/models"
	"fmt"
)

type CopyFolder struct {
	options models.OperationOptions
}

func (c *CopyFolder) Exec() {
	fmt.Println("Copying a folder")
}
