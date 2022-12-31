package delete

import (
    "ena/dirgod/constants"
    "ena/dirgod/models"
	"ena/dirgod/operations"
	"os"
    "path/filepath"
    "strings"
)

type DeleteFile struct {
	Source          string
	Options         models.OperationOptions
	Result          operations.OperationResult
	RollbackResult  operations.OperationResult
}

func (d *DeleteFile) Exec() {
    err := os.Rename(d.Source, constants.EnaTmp)
    if err != nil {
        operations.DecideErrorOutput(&d.Options, &d.Result, err)
        return
    }

    d.Result.Completed = true
}

func (d *DeleteFile) Rollback() {
    filename := strings.Split(d.Source, "/")
    trashPath := filepath.Join(constants.EnaTmp, filename[len(filename) - 1])
    err := os.Rename(trashPath, d.Source)
    if err != nil {
        d.RollbackResult.Completed = false
        d.RollbackResult.Err = err
        return
    }

    d.RollbackResult.Completed = true
}
