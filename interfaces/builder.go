package interfaces

type Builder interface {
	Build() (Operation, error);
}