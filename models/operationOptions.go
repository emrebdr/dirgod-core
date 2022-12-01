package models

type Options int

const (
	Default Options = iota
	Force
	Strict
)

type OperationOptions struct {
	Options Options
	Cache   bool
}
