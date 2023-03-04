package ch

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/emrebdr/dirgod-code/models"
	"github.com/emrebdr/dirgod-code/operations"
)

type Chown struct {
	Source         string
	UID            int
	GID            int
	PreviousUID    int
	PreviousGID    int
	Recursive      bool
	Options        models.OperationOptions
	Result         operations.OperationResult
	RollbackResult operations.OperationResult
}

func (c *Chown) Exec() {
	_, err := os.Stat(c.Source)
	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return
	}

	c.PreviousGID = os.Getgid()
	c.PreviousUID = os.Getuid()

	if c.Recursive {
		err = filepath.WalkDir(c.Source, c.chownRecursive)
	} else {
		err = os.Chown(c.Source, c.UID, c.GID)
	}

	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return
	}

	c.Result.Completed = true
}

func (c *Chown) Rollback() {
	var err error
	if c.Recursive {
		err = filepath.WalkDir(c.Source, c.chownRecursiveRollback)
	} else {
		err = os.Chown(c.Source, c.PreviousUID, c.PreviousGID)
	}

	if err != nil {
		c.RollbackResult.Completed = false
		c.RollbackResult.Err = err
		return
	}

	c.RollbackResult.Completed = true
}

func (c *Chown) chownRecursive(path string, de fs.DirEntry, err error) error {
	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return err
	}

	err = os.Chown(path, c.UID, c.GID)
	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return err
	}

	return nil
}

func (c *Chown) chownRecursiveRollback(path string, de fs.DirEntry, err error) error {
	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return err
	}

	err = os.Chown(path, c.PreviousUID, c.PreviousGID)
	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return err
	}

	return nil
}
