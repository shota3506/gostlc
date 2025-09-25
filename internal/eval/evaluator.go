package eval

import (
	"fmt"

	"github.com/shota3506/gostlc/internal/ast"
	"github.com/shota3506/gostlc/internal/builtin"
	"github.com/shota3506/gostlc/internal/values"
)

func Eval(expr ast.TypedExpr) (values.Value, error) {
	root := values.NewRho()
	for ident, val := range builtin.Functions {
		root = root.Bind(ident, val)
	}
	return evalExpr(expr, root)
}

func evalExpr(expr ast.TypedExpr, env values.Rho) (values.Value, error) {
	switch e := expr.(type) {
	case *ast.TypedIntExpr:
		return &values.IntValue{Value: e.Value}, nil

	case *ast.TypedBoolExpr:
		return &values.BoolValue{Value: e.Value}, nil

	case *ast.TypedVarExpr:
		val, ok := env.Lookup(e.Name)
		if !ok {
			return nil, fmt.Errorf("undefined variable: %s at line %d, col %d", e.Name, e.Pos.Line, e.Pos.Column)
		}
		return val, nil

	case *ast.TypedAbsExpr:
		return &values.Closure{
			Param:     e.Param,
			ParamType: e.ParamType,
			Body:      e.Body,
			Env:       env,
		}, nil

	case *ast.TypedAppExpr:
		fnVal, err := evalExpr(e.Func, env)
		if err != nil {
			return nil, err
		}

		argVal, err := evalExpr(e.Arg, env)
		if err != nil {
			return nil, err
		}

		switch fn := fnVal.(type) {
		case *values.Closure:
			return evalExpr(fn.Body, fn.Env.Bind(fn.Param, argVal))
		case *values.BuiltinFunc:
			return fn.Fn(argVal)
		case *values.PartialBuiltinFunc:
			return fn.Fn(argVal)
		default:
			return nil, fmt.Errorf("expected function value at line %d, col %d", e.Pos.Line, e.Pos.Column)
		}

	case *ast.TypedIfExpr:
		condVal, err := evalExpr(e.Cond, env)
		if err != nil {
			return nil, err
		}

		boolVal, ok := condVal.(*values.BoolValue)
		if !ok {
			return nil, fmt.Errorf("expected boolean value in if condition at line %d, col %d", e.Pos.Line, e.Pos.Column)
		}

		if boolVal.Value {
			return evalExpr(e.Then, env)
		}
		return evalExpr(e.Else, env)

	default:
		return nil, fmt.Errorf("unsupported expression type: %T", expr)
	}
}
