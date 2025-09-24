package ast

// Type represents a type in the lambda calculus with simple types.
type Type interface {
	typeNode()
}

// BooleanType represents the boolean type.
type BooleanType struct{}

func (BooleanType) typeNode() {}

// IntType represents the integer type.
type IntType struct{}

func (IntType) typeNode() {}

// FuncType represents a function type from one type to another.
type FuncType struct {
	From Type
	To   Type
}

func (FuncType) typeNode() {}
