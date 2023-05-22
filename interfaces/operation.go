package interfaces

import "github.com/emrebdr/dirgod-core/operations"

type Operation interface {
	Exec() operations.OperationResult
	Rollback()
}
