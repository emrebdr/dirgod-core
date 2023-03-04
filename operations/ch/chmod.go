package ch

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/emrebdr/dirgod-code/models"
	"github.com/emrebdr/dirgod-code/operations"
)

type Chmod struct {
	Source           string
	PermCode         fs.FileMode
	PreviousPermCode fs.FileMode
	Recursive        bool
	Options          models.OperationOptions
	Result           operations.OperationResult
	RollbackResult   operations.OperationResult
}

func (c *Chmod) Exec() {
	stats, err := os.Stat(c.Source)
	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return
	}

	c.PreviousPermCode = stats.Mode()

	if c.Recursive {
		err = filepath.WalkDir(c.Source, c.chmodRecursive)
	} else {
		err = os.Chmod(c.Source, c.PermCode)
	}

	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return
	}

	c.Result.Completed = true
}

func (c *Chmod) Rollback() {
	var err error

	if c.Recursive {
		err = filepath.WalkDir(c.Source, c.chmodRecursiveRollback)
	} else {
		err = os.Chmod(c.Source, c.PreviousPermCode)
	}

	if err != nil {
		c.RollbackResult.Completed = false
		c.RollbackResult.Err = err
		return
	}

	c.RollbackResult.Completed = true
}

func (c *Chmod) chmodRecursive(path string, de fs.DirEntry, err error) error {
	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return err
	}

	err = os.Chmod(path, c.PermCode)
	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return err
	}

	return nil
}

func (c *Chmod) chmodRecursiveRollback(path string, de fs.DirEntry, err error) error {
	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return err
	}

	err = os.Chmod(path, c.PreviousPermCode)
	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return err
	}

	return nil
}
