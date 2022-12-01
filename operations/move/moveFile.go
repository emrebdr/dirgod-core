package move

import (
	"ena/dirgod/models"
	"fmt"
)

type MoveFile struct {
	options models.OperationOptions
}

func (m *MoveFile) Exec() {
	fmt.Println("Moving a file")
}
