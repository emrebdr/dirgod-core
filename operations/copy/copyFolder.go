package copy

import (
	"ena/dirgod/models"
	"ena/dirgod/operations"
	"io"
	"os"
	"path"
)

type CopyFolder struct {
	Source         string
	Destination    string
	Options        models.OperationOptions
	Result         operations.OperationResult
	RollbackResult operations.OperationResult
}

func (c *CopyFolder) Exec() {
	srcInfo, err := os.Stat(c.Source)
	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return
	}

	err = os.MkdirAll(c.Destination, srcInfo.Mode())
	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return
	}

	fds, err := os.ReadDir(c.Source)
	if err != nil {
		operations.DecideErrorOutput(&c.Options, &c.Result, err)
		return
	}

	for _, fd := range fds {
		srcFp := path.Join(c.Source, fd.Name())
		dstFp := path.Join(c.Destination, fd.Name())

		if fd.IsDir() {
			c.Exec()
		} else {
			if err = c.File(srcFp, dstFp); err != nil {
				operations.DecideErrorOutput(&c.Options, &c.Result, err)
				return
			}
		}
	}

	c.Result.Completed = true
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

func (c *CopyFolder) File(src, dst string) error {
	var err error
	var srcFd *os.File
	var dstFd *os.File
	var srcInfo os.FileInfo

	if srcFd, err = os.Open(src); err != nil {
		return err
	}
	defer srcFd.Close()

	if dstFd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstFd.Close()

	if _, err = io.Copy(dstFd, srcFd); err != nil {
		return err
	}
	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}
