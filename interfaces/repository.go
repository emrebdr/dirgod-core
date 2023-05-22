package interfaces

type Repository interface {
	Commit() error
	Rollback(string) error
}
