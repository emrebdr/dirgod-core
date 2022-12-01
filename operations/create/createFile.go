package create

import (
	"ena/dirgod/models"
	"fmt"
)

type CreateFile struct {
}

func (c *CreateFile) Exec(o *models.OperationOptions) {
	fmt.Println("Creating a file...")
}
