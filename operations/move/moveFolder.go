package move

import (
	"ena/dirgod/models"
    "ena/dirgod/operations"
    "os"
)

type MoveFolder struct {
    Source         string
    Destination    string
    Options        models.OperationOptions
    Result         operations.OperationResult
    RollbackResult operations.OperationResult
}

func (m *MoveFolder) Exec() {
	err := os.Rename(m.Source, m.Destination)
    if err != nil {
        operations.DecideErrorOutput(&m.Options, &m.Result, err)
        return
    }

    m.Result.Completed = true
}

func (m *MoveFolder) Rollback() {
    err := os.Rename(m.Destination, m.Source)
    if err != nil {
        m.RollbackResult.Completed = false
        m.RollbackResult.Err = err
        return
    }

    m.RollbackResult.Completed = true
}