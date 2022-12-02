package interfaces

type Operation interface {
	Exec()
	Rollback()
}
