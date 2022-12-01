package interfaces

import (
	"ena/dirgod/models"
)

type Operation interface {
	exec(o *models.OperationOptions)
}
