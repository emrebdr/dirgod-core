package chmod

import (
	"ena/dirgod/models"
	"fmt"
)

type ChmodFile struct {
	options models.OperationOptions
}

func (c *ChmodFile) Exec() {
	fmt.Println("Chmode file")
}
