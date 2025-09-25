package eval

import (
	"fmt"

	"github.com/shota3506/gostlc/internal/ast"
)

type rho struct {
	name   string
	value  Value
	parent *rho
}

func newRho() rho {
	return rho{}
}

func (r rho) lookup(name string) (Value, bool) {
	if name == "" {
		return nil, false
	}
	if r.name == name {
		return r.value, true
	}
	if r.parent != nil {
		return r.parent.lookup(name)
	}
	return nil, false
}

func (r rho) bind(name string, value Value) rho {
	return rho{
		name:   name,
		value:  value,
		parent: &r,
	}
}

type Evaluator struct {
}

func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

func (ev *Evaluator) Eval(expr ast.TypedExpr) (Value, error) {
	root := newRho()
	return ev.evalExpr(expr, root)
}

func (ev *Evaluator) evalExpr(expr ast.TypedExpr, rho rho) (Value, error) {
	switch e := expr.(type) {
	case *ast.TypedIntExpr:
		return &IntValue{Value: e.Value}, nil

	case *ast.TypedBoolExpr:
		return &BoolValue{Value: e.Value}, nil

	case *ast.TypedVarExpr:
		val, ok := rho.lookup(e.Name)
		if !ok {
			return nil, fmt.Errorf("undefined variable: %s at line %d, col %d", e.Name, e.Pos.Line, e.Pos.Column)
		}
		return val, nil

	case *ast.TypedAbsExpr:
		return &Closure{
			Param:     e.Param,
			ParamType: e.ParamType,
			Body:      e.Body,
			Rho:       rho,
		}, nil

	case *ast.TypedAppExpr:
		fnVal, err := ev.evalExpr(e.Func, rho)
		if err != nil {
			return nil, err
		}

		argVal, err := ev.evalExpr(e.Arg, rho)
		if err != nil {
			return nil, err
		}

		switch fn := fnVal.(type) {
		case *Closure:
			return ev.evalExpr(fn.Body, fn.Rho.bind(fn.Param, argVal))
		default:
			return nil, fmt.Errorf("expected function value at line %d, col %d", e.Pos.Line, e.Pos.Column)
		}

	case *ast.TypedIfExpr:
		condVal, err := ev.evalExpr(e.Cond, rho)
		if err != nil {
			return nil, err
		}

		boolVal, ok := condVal.(*BoolValue)
		if !ok {
			return nil, fmt.Errorf("expected boolean value in if condition at line %d, col %d", e.Pos.Line, e.Pos.Column)
		}

		if boolVal.Value {
			return ev.evalExpr(e.Then, rho)
		}
		return ev.evalExpr(e.Else, rho)

	default:
		return nil, fmt.Errorf("unsupported expression type: %T", expr)
	}
}
