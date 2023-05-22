package interfaces

type Object interface {
	Print()
	GetType() string
	GetName() string
	GetChecksum() string
}
