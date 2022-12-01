package chmod

import (
	"ena/dirgod/models"
	"fmt"
)

type ChmodFolder struct {
	options models.OperationOptions
}

func (c *ChmodFolder) Exec() {
	fmt.Println("Chmode folder")
}
