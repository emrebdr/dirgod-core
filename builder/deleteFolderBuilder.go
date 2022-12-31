package builder

import (
    "ena/dirgod/interfaces"
    "ena/dirgod/models"
    "ena/dirgod/operations/delete"
    "errors"
)

type DeleteFolderBuilder struct {
    Source          string               `json:"source"`
    WorkingMode     string               `json:"workingMode"`
    Cache           bool                 `json:"cache"`
    createOperation interfaces.Operation
}

func (d *DeleteFolderBuilder) Build() (interfaces.Operation, error) {
    if d.Source == "" {
        return nil, errors.New("source is empty")
    }

    workingMode, err := d.setWorkingMode()
    if err != nil {
        return nil, err
    }

    d.createOperation = &delete.DeleteFolder{
        Source:          d.Source,
        Options: models.OperationOptions{
            WorkingMode: workingMode,
            Cache:       d.Cache,
        },
    }

    return d.createOperation, nil
}

func (d *DeleteFolderBuilder) GetName() string {
    return "DeleteFolder"
}

func (d *DeleteFolderBuilder) IsValid() bool {
    return d.createOperation != nil
}

func (d *DeleteFolderBuilder) setWorkingMode() (models.Options, error) {
    if d.WorkingMode != "" {
        switch d.WorkingMode {
        case "force":
            return models.Force, nil
        case "strict":
            return models.Strict, nil
        case "default":
            return models.Default, nil
        default:
            return -1, errors.New("unknown working mode")
        }
    }

    return models.Default, nil
}