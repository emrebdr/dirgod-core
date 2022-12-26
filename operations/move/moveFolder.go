package move

import (
	"ena/dirgod/models"
	"fmt"
)

type MoveFolder struct {
	Options models.OperationOptions
}

func (m *MoveFolder) Exec() {
	fmt.Println("Moving a folder")
}
