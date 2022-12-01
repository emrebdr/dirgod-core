package delete

import (
	"ena/dirgod/models"
	"fmt"
)

type DeleteFile struct {
	options models.OperationOptions
}

func (d *DeleteFile) Exec() {
	fmt.Println("Deleting a file")
}
