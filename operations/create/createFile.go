package create

import (
	"ena/dirgod/models"
	"fmt"
)

type CreateFile struct {
	options models.OperationOptions
}

func (c *CreateFile) Exec() {
	fmt.Println("Creating a file")
}
