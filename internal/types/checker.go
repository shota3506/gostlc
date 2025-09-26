package types

import (
	"github.com/shota3506/gostlc/internal/ast"
	"github.com/shota3506/gostlc/internal/builtin"
)

// gamma represents the type environment for variable bindings.
type gamma struct {
	name   string
	typ    ast.Type
	parent *gamma
}

func newGamma() gamma {
	return gamma{}
}

func (g gamma) bind(name string, typ ast.Type) gamma {
	return gamma{
		name:   name,
		typ:    typ,
		parent: &g,
	}
}

func (g gamma) lookup(name string) (ast.Type, bool) {
	if name == "" {
		return nil, false
	}
	if g.name == name {
		return g.typ, true
	}
	if g.parent != nil {
		return g.parent.lookup(name)
	}
	return nil, false
}

// Check performs type checking and returns a typed AST.
func Check(expr ast.Expr) (ast.TypedExpr, error) {
	root := newGamma()
	for ident, typ := range builtin.FunctionTypes {
		root = root.bind(ident, typ)
	}
	return checkTyped(expr, root)
}

func checkTyped(expr ast.Expr, g gamma) (ast.TypedExpr, error) {
	switch e := expr.(type) {
	case *ast.VarExpr:
		return checkVar(e, g)
	case *ast.AbsExpr:
		return checkAbs(e, g)
	case *ast.AppExpr:
		return checkApp(e, g)
	case *ast.BoolExpr:
		return ast.NewTypedBoolExpr(e), nil
	case *ast.IntExpr:
		return ast.NewTypedIntExpr(e), nil
	case *ast.IfExpr:
		return checkIf(e, g)
	default:
		return nil, &UnknownExprTypeError{
			Pos:  expr.Position(),
			Expr: expr,
		}
	}
}

func checkVar(expr *ast.VarExpr, g gamma) (ast.TypedExpr, error) {
	typ, ok := g.lookup(expr.Name)
	if !ok {
		return nil, &UndefinedVariableError{
			Pos:  expr.Pos,
			Name: expr.Name,
		}
	}
	return ast.NewTypedVarExpr(typ, expr), nil
}

func checkAbs(expr *ast.AbsExpr, g gamma) (ast.TypedExpr, error) {
	typedBody, err := checkTyped(expr.Body, g.bind(expr.Param, expr.ParamType))
	if err != nil {
		return nil, err
	}

	funcType := &ast.FuncType{
		From: expr.ParamType,
		To:   typedBody.Type(),
	}
	return ast.NewTypedAbsExpr(funcType, expr.Pos, expr.Param, expr.ParamType, typedBody), nil
}

func checkApp(expr *ast.AppExpr, g gamma) (ast.TypedExpr, error) {
	typedFunc, err := checkTyped(expr.Func, g)
	if err != nil {
		return nil, err
	}

	ft, ok := typedFunc.Type().(*ast.FuncType)
	if !ok {
		return nil, &NotAFunctionError{
			Pos:  expr.Pos,
			Type: typedFunc.Type(),
		}
	}

	typedArg, err := checkTyped(expr.Arg, g)
	if err != nil {
		return nil, err
	}

	if !ft.From.Equal(typedArg.Type()) {
		return nil, &TypeMismatchError{
			Pos:      expr.Pos,
			Expected: ft.From,
			Actual:   typedArg.Type(),
			Context:  "application",
		}
	}

	return ast.NewTypedAppExpr(ft.To, expr.Pos, typedFunc, typedArg), nil
}

func checkIf(expr *ast.IfExpr, g gamma) (ast.TypedExpr, error) {
	typedCond, err := checkTyped(expr.Cond, g)
	if err != nil {
		return nil, err
	}

	if _, ok := typedCond.Type().(*ast.BoolType); !ok {
		return nil, &InvalidConditionTypeError{
			Pos:  expr.Pos,
			Type: typedCond.Type(),
		}
	}

	typedThen, err := checkTyped(expr.Then, g)
	if err != nil {
		return nil, err
	}

	typedElse, err := checkTyped(expr.Else, g)
	if err != nil {
		return nil, err
	}

	if !typedThen.Type().Equal(typedElse.Type()) {
		return nil, &TypeMismatchError{
			Pos:      expr.Pos,
			Expected: typedThen.Type(),
			Actual:   typedElse.Type(),
			Context:  "if-else branches",
		}
	}
	return ast.NewTypedIfExpr(expr.Pos, typedCond, typedThen, typedElse), nil
}
