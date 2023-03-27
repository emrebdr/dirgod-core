package builder

import (
	"errors"

	"github.com/emrebdr/dirgod-core/interfaces"
	"github.com/emrebdr/dirgod-core/models"
	"github.com/emrebdr/dirgod-core/operations/write"
	"github.com/emrebdr/dirgod-core/utils"
)

type WriteBuilder struct {
	From        string `json:"from"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	WorkingMode string `json:"workingMode"`
	Cache       bool   `json:"cache"`
	operation   interfaces.Operation
}

func (w *WriteBuilder) Build() (interfaces.Operation, error) {
	if w.From == "" {
		return nil, errors.New("from type is empty")
	}

	if w.Source == "" {
		return nil, errors.New("source is empty")
	}

	if w.Destination == "" {
		return nil, errors.New("destination is empty")
	}

	workingMode, err := utils.SetWorkingMode(w.WorkingMode)
	if err != nil {
		return nil, err
	}

	w.operation = &write.Write{
		From:        write.FromType(w.From),
		Source:      w.Source,
		Destination: w.Destination,
		Options: models.OperationOptions{
			WorkingMode: workingMode,
			Cache:       w.Cache,
		},
	}

	return w.operation, nil
}

func (w *WriteBuilder) GetName() string {
	return "Write"
}

func (w *WriteBuilder) IsValid() bool {
	return w.operation != nil
}
