package models

type Options int

const (
	Default Options = iota
	Force
	Strict
)

type OperationOptions struct {
	WorkingMode Options
	Cache       bool
}
