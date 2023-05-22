package copy

import (
	"io"
	"os"

	"github.com/emrebdr/dirgod-core/models"
	"github.com/emrebdr/dirgod-core/operations"
)

type CopyFolder struct {
	Source         string
	Destination    string
	Options        models.OperationOptions
	Result         operations.OperationResult
	RollbackResult operations.OperationResult
}

func (c *CopyFolder) Exec() operations.OperationResult {
	err := c.CopyDir(c.Source, c.Destination)
	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return c.Result
	}

	c.Result.Completed = true

	return c.Result
}

func (c *CopyFolder) Rollback() {
	err := os.RemoveAll(c.Destination)
	if err != nil {
		c.RollbackResult.Completed = false
		c.RollbackResult.Err = err
		return
	}

	c.RollbackResult.Completed = true
}

func (c *CopyFolder) CopyFile(source string, dest string) error {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destFile.Close()

	_, err = io.Copy(destFile, sourcefile)
	if err == nil {
		sourceInfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceInfo.Mode())
			return err
		}
		return nil
	}

	return err
}

func (c *CopyFolder) CopyDir(source string, dest string) error {
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dest, sourceInfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)
	if err != nil {
		return err
	}

	for _, obj := range objects {
		sourceFilePointer := source + "/" + obj.Name()
		destinationFilePointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			err = c.CopyDir(sourceFilePointer, destinationFilePointer)
			if err != nil {
				return err
			}
		} else {
			err = c.CopyFile(sourceFilePointer, destinationFilePointer)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
