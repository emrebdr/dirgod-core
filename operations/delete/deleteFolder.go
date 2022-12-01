package delete

import (
	"ena/dirgod/models"
	"fmt"
)

type DeleteFolder struct {
	options models.OperationOptions
}

func (d *DeleteFolder) Exec() {
	fmt.Println("Deleting a folder")
}
