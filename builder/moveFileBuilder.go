package builder

import(
    "ena/dirgod/interfaces"
    "ena/dirgod/models"
    "ena/dirgod/operations/move"
    "errors"
)

type MoveFileBuilder struct {
    Source string                        `json:"source"`
    WorkingMode     string               `json:"workingMode"`
    Cache           bool                 `json:"cache"`
    createOperation interfaces.Operation
}

func (m *MoveFileBuilder) Build() (interfaces.Operation, error) {
    if m.Source == "" {
        return nil, errors.New("source is empty")
    }

    workingMode, err := m.setWorkingMode()
    if err != nil {
        return nil, err
    }

    m.createOperation = &move.MoveFile{
        Source:          m.Source,
        Options: models.OperationOptions{
            WorkingMode: workingMode,
            Cache:       m.Cache,
        },
    }

    return m.createOperation, nil
}

func (m *MoveFileBuilder) GetName() string {
    return "MoveFile"
}

func (m *MoveFileBuilder) IsValid() bool {
    return m.createOperation != nil
}

func (m *MoveFileBuilder) setWorkingMode() (models.Options, error) {
    if m.WorkingMode != "" {
        switch m.WorkingMode {
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