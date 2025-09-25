package types

import (
	"fmt"

	"github.com/shota3506/gostlc/internal/ast"
	"github.com/shota3506/gostlc/internal/token"
)

// UndefinedVariableError occurs when a variable is not found in the environment.
type UndefinedVariableError struct {
	Pos  token.Position
	Name string
}

func (e *UndefinedVariableError) Error() string {
	return fmt.Sprintf("%d:%d: undefined variable: %s", e.Pos.Line, e.Pos.Column, e.Name)
}

// TypeMismatchError occurs when expected and actual types don't match.
type TypeMismatchError struct {
	Pos      token.Position
	Expected ast.Type
	Actual   ast.Type
	Context  string
}

func (e *TypeMismatchError) Error() string {
	if e.Context != "" {
		return fmt.Sprintf("%d:%d: type mismatch in %s: expected %s, got %s", e.Pos.Line, e.Pos.Column, e.Context, e.Expected, e.Actual)
	}
	return fmt.Sprintf("%d:%d: type mismatch: expected %s, got %s", e.Pos.Line, e.Pos.Column, e.Expected, e.Actual)
}

// NotAFunctionError occurs when trying to apply a non-function value.
type NotAFunctionError struct {
	Pos  token.Position
	Type ast.Type
}

func (e *NotAFunctionError) Error() string {
	return fmt.Sprintf("%d:%d: cannot apply non-function type: %s", e.Pos.Line, e.Pos.Column, e.Type)
}

// InvalidConditionTypeError occurs when if-expression condition is not boolean.
type InvalidConditionTypeError struct {
	Pos  token.Position
	Type ast.Type
}

func (e *InvalidConditionTypeError) Error() string {
	return fmt.Sprintf("%d:%d: condition must be boolean, got %s", e.Pos.Line, e.Pos.Column, e.Type)
}

type UnknownExprTypeError struct {
	Pos  token.Position
	Expr ast.Expr
}

func (e *UnknownExprTypeError) Error() string {
	return fmt.Sprintf("%d:%d: unknown expression type: %T", e.Pos.Line, e.Pos.Column, e.Expr)
}
