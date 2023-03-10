package delete

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/emrebdr/dirgod-core/constants"
	"github.com/emrebdr/dirgod-core/models"
	"github.com/emrebdr/dirgod-core/operations"
)

type DeleteFile struct {
	Source         string
	Options        models.OperationOptions
	Result         operations.OperationResult
	RollbackResult operations.OperationResult
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
	trashPath := filepath.Join(constants.EnaTmp, filename[len(filename)-1])
	err := os.Rename(trashPath, d.Source)
	if err != nil {
		d.RollbackResult.Completed = false
		d.RollbackResult.Err = err
		return
	}

	d.RollbackResult.Completed = true
}
