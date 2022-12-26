package move

import (
	"ena/dirgod/models"
	"ena/dirgod/operations"
	"os"
)

var DEFAULT_FILE_WRITE_PERMS = os.FileMode(0644)

type MoveFile struct {
	Source         string
	Destination    string
	Options        models.OperationOptions
	Result         operations.OperationResult
	RollbackResult operations.OperationResult
}

func (m *MoveFile) Exec() {
	err := m.moveFile(m.Source, m.Destination)
	if err != nil {
		operations.DecideErrorOutput(&m.Options, &m.Result, err)
		return
	}

	m.Result.Completed = true
}

func (m *MoveFile) Rollback() {
	err := m.moveFile(m.Destination, m.Source)
	if err != nil {
		m.RollbackResult.Completed = false
		m.RollbackResult.Err = err
		return
	}

	m.RollbackResult.Completed = true
}

func (m *MoveFile) moveFile(src, destination string) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	err = os.WriteFile(destination, content, DEFAULT_FILE_WRITE_PERMS)
	if err != nil {
		return err
	}

	err = os.Remove(src)
	if err != nil {
		return err
	}

	return nil
}
