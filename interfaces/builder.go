package interfaces

type Builder interface {
	Build() (Operation, error)
	GetName() string
	IsValid() bool
}
