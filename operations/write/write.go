package write

import (
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/emrebdr/dirgod-core/models"
	"github.com/emrebdr/dirgod-core/operations"
)

type FromType string

const (
	Url  FromType = "url"
	File FromType = "file"
)

type Write struct {
	From           FromType
	Source         string
	Destination    string
	Options        models.OperationOptions
	Result         operations.OperationResult
	RollbackResult operations.OperationResult
}

func (w *Write) Exec() operations.OperationResult {
	content, err := w.getContent()
	if err != nil {
		operations.DecideErrorOutput(&w.Options, &w.Result, err)
		return w.Result
	}

	writeResult := os.WriteFile(w.Destination, content, 0644)
	if writeResult != nil {
		operations.DecideErrorOutput(&w.Options, &w.Result, writeResult)
		return w.Result
	}

	w.Result.Completed = true

	return w.Result
}

func (w *Write) Rollback() {
	// TBC
}

func (w *Write) getContent() ([]byte, error) {
	if w.From == "url" {
		return w.getContentFromUrl()
	} else if w.From == "file" {
		return w.getContentFromFile()
	}

	return nil, errors.New("unknown from type")
}

func (w *Write) getContentFromUrl() ([]byte, error) {
	resp, err := http.Get(w.Source)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (w *Write) getContentFromFile() ([]byte, error) {
	byt, err := os.ReadFile(w.Source)
	if err != nil {
		return nil, err
	}

	return byt, nil
}
