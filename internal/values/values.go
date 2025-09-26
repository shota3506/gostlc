package values

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
	Env       *Rho
}

func (c *Closure) value() {}
func (c *Closure) String() string {
	return fmt.Sprintf("<closure:%s->%s>", c.ParamType, c.Body.Type())
}

type BuiltinFunc struct {
	Name       string
	ParamType  ast.Type
	ReturnType ast.Type
	Fn         func(args Value) (Value, error)
}

func (b *BuiltinFunc) value() {}
func (b *BuiltinFunc) String() string {
	return fmt.Sprintf("<builtin:%s:%s->%s>", b.Name, b.ParamType, b.ReturnType)
}

type PartialBuiltinFunc struct {
	Name       string // 元の関数名を保持
	ParamType  ast.Type
	ReturnType ast.Type
	Fn         func(args Value) (Value, error)
}

func (p *PartialBuiltinFunc) value() {}
func (p *PartialBuiltinFunc) String() string {
	return fmt.Sprintf("<builtin:%s[partial]:%s->%s>", p.Name, p.ParamType, p.ReturnType)
}
