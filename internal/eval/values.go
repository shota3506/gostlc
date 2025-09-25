package eval

import (
	"fmt"

	"github.com/shota3506/gostlc/internal/ast"
)

type Value interface {
	fmt.Stringer

	value()
}

type IntValue struct {
	Value int
}

func (v *IntValue) value() {}
func (v *IntValue) String() string {
	return fmt.Sprintf("%d", v.Value)
}

type BoolValue struct {
	Value bool
}

func (v *BoolValue) value() {}
func (v *BoolValue) String() string {
	if v.Value {
		return "true"
	}
	return "false"
}

type Closure struct {
	Param     string
	ParamType ast.Type
	Body      ast.TypedExpr
	Rho       rho
}

func (c *Closure) value() {}
func (c *Closure) String() string {
	return fmt.Sprintf("<λ%s:%s.%s>", c.Param, c.ParamType, c.Body.Type())
}
