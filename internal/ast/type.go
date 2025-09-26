package ast

import "fmt"

// Type represents a type in the lambda calculus with simple types.
type Type interface {
	fmt.Stringer
	Equal(u Type) bool

	typeNode()
}

// BoolType represents the boolean type.
type BoolType struct{}

func (*BoolType) typeNode() {}

func (*BoolType) String() string {
	return "Bool"
}

func (b *BoolType) Equal(u Type) bool {
	_, ok := u.(*BoolType)
	return ok
}

// IntType represents the integer type.
type IntType struct{}

func (*IntType) typeNode() {}

func (*IntType) String() string {
	return "Int"
}

func (i *IntType) Equal(u Type) bool {
	_, ok := u.(*IntType)
	return ok
}

// FuncType represents a function type from one type to another.
type FuncType struct {
	From Type
	To   Type
}

func (*FuncType) typeNode() {}

func (f *FuncType) String() string {
	return fmt.Sprintf("(%s->%s)", f.From, f.To)
}

func (f *FuncType) Equal(u Type) bool {
	v, ok := u.(*FuncType)
	if !ok {
		return false
	}
	return f.From.Equal(v.From) && f.To.Equal(v.To)
}
