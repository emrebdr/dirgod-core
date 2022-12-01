package create

import (
	"ena/dirgod/models"
	"fmt"
)

type CreateFolder struct {
	options models.OperationOptions
}

func (c *CreateFolder) Exec() {
	fmt.Println("Creating a folder")
}
