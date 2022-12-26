package delete

import (
	"ena/dirgod/models"
	"ena/dirgod/operations"
	"os"
)

type DeleteFile struct {
	Source          string
	PrevFileContent []byte
	Options         models.OperationOptions
	Result          operations.OperationResult
	RollbackResult  operations.OperationResult
}

func (d *DeleteFile) Exec() {
	file, err := os.ReadFile(d.Source)
	if err != nil {
		operations.DecideErrorOutput(&d.Options, &d.Result, err)
		return
	}

	d.PrevFileContent = file

	err = os.Remove(d.Source)
	if err != nil {
		operations.DecideErrorOutput(&d.Options, &d.Result, err)
		return
	}

	d.Result.Completed = true
}

func (d *DeleteFile) Rollback() {
	file, err := os.Create(d.Source)
	if err != nil {
		d.RollbackResult.Completed = false
		d.RollbackResult.Err = err
		return
	}

	_, err = file.Write(d.PrevFileContent)
	if err != nil {
		d.RollbackResult.Completed = false
		d.RollbackResult.Err = err
		return
	}

	d.RollbackResult.Completed = true
}
