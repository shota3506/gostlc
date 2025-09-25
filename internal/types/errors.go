package types

import (
	"fmt"

	"github.com/shota3506/gostlc/internal/ast"
)

// UndefinedVariableError occurs when a variable is not found in the environment.
type UndefinedVariableError struct {
	Name string
}

func (e *UndefinedVariableError) Error() string {
	return fmt.Sprintf("undefined variable: %s", e.Name)
}

// TypeMismatchError occurs when expected and actual types don't match.
type TypeMismatchError struct {
	Expected ast.Type
	Actual   ast.Type
	Context  string
}

func (e *TypeMismatchError) Error() string {
	if e.Context != "" {
		return fmt.Sprintf("type mismatch in %s: expected %s, got %s", e.Context, e.Expected, e.Actual)
	}
	return fmt.Sprintf("type mismatch: expected %s, got %s", e.Expected, e.Actual)
}

// NotAFunctionError occurs when trying to apply a non-function value.
type NotAFunctionError struct {
	Type ast.Type
}

func (e *NotAFunctionError) Error() string {
	return fmt.Sprintf("cannot apply non-function type: %s", e.Type)
}

// InvalidConditionTypeError occurs when if-expression condition is not boolean.
type InvalidConditionTypeError struct {
	Type ast.Type
}

func (e *InvalidConditionTypeError) Error() string {
	return fmt.Sprintf("condition must be boolean: %s", e.Type)
}
