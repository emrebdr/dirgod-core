package builder

import(
    "ena/dirgod/interfaces"
    "ena/dirgod/operations/ch"
    "ena/dirgod/models"
    "errors"
)

type ChownBuilder struct {
    Source          string               `json:"source"`
    UID             int                  `json:"uid"`
    GID             int                  `json:"gid"`
    Recursive       bool                 `json:"recursive"`
    WorkingMode     string               `json:"workingMode"`
    Cache           bool                 `json:"cache"`
    createOperation interfaces.Operation
}

func (c *ChownBuilder) Build() (interfaces.Operation, error) {
    if c.Source == "" {
        return nil, errors.New("source is empty")
    }

    if c.GID <= 0 || c.UID <= 0 {
        return nil, errors.New("UID and GID must be greater than 0")
    }

    workingMode, err := c.setWorkingMode()
    if err != nil {
        return nil, err
    }

    c.createOperation = &ch.Chown{
        Source: c.Source,
        UID: c.UID,
        GID: c.GID,
        Recursive: c.Recursive,
        Options: models.OperationOptions{
            WorkingMode: workingMode,
            Cache:       c.Cache,
        },
    }

    return c.createOperation, nil
}

func (c *ChownBuilder) GetName() string {
    return "Chown"
}

func (c *ChownBuilder) IsValid() bool {
    return c.createOperation != nil
}

func (c *ChownBuilder) setWorkingMode() (models.Options, error) {
    if c.WorkingMode != "" {
        switch c.WorkingMode {
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