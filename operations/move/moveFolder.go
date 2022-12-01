package move

import (
	"ena/dirgod/models"
	"fmt"
)

type MoveFolder struct {
	options models.OperationOptions
}

func (m *MoveFolder) Exec() {
	fmt.Println("Moving a folder")
}
