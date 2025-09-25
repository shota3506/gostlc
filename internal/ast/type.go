package ast

import "fmt"

// Type represents a type in the lambda calculus with simple types.
type Type interface {
	fmt.Stringer

	typeNode()
}

// BooleanType represents the boolean type.
type BooleanType struct{}

func (BooleanType) typeNode() {}

func (BooleanType) String() string {
	return "Bool"
}

// IntType represents the integer type.
type IntType struct{}

func (IntType) typeNode() {}

func (IntType) String() string {
	return "Int"
}

// FuncType represents a function type from one type to another.
type FuncType struct {
	From Type
	To   Type
}

func (FuncType) typeNode() {}

func (f FuncType) String() string {
	return fmt.Sprintf("(%s -> %s)", f.From, f.To)
}
