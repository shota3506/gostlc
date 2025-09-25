package types

import (
	"github.com/shota3506/gostlc/internal/ast"
)

// gamma represents the type environment for variable bindings.
type gamma struct {
	bindings map[string]ast.Type
	parent   *gamma
}

func (g *gamma) child() *gamma {
	return &gamma{
		bindings: make(map[string]ast.Type),
		parent:   g,
	}
}

func (g *gamma) bind(name string, typ ast.Type) {
	g.bindings[name] = typ
}

func (g *gamma) lookup(name string) (ast.Type, bool) {
	if typ, ok := g.bindings[name]; ok {
		return typ, true
	}
	if g.parent != nil {
		return g.parent.lookup(name)
	}
	return nil, false
}

// TypeChecker provides type checking functionality for expressions.
type TypeChecker struct{}

func NewTypeChecker() *TypeChecker {
	return &TypeChecker{}
}

// Check performs type checking on an expression and returns its type.
func (tc *TypeChecker) Check(expr ast.Expr) (ast.Type, error) {
	root := &gamma{bindings: make(map[string]ast.Type)}
	return tc.check(expr, root)
}

func (tc *TypeChecker) check(expr ast.Expr, g *gamma) (ast.Type, error) {
	switch e := expr.(type) {
	case *ast.VarExpr:
		return tc.checkVar(e, g)
	case *ast.AbsExpr:
		return tc.checkAbs(e, g)
	case *ast.AppExpr:
		return tc.checkApp(e, g)
	case *ast.BoolExpr:
		return &ast.BooleanType{}, nil
	case *ast.IntExpr:
		return &ast.IntType{}, nil
	case *ast.IfExpr:
		return tc.checkIf(e, g)
	default:
		return nil, &UnknownExprTypeError{
			Pos:  expr.Position(),
			Expr: expr,
		}
	}
}

func (tc *TypeChecker) checkVar(expr *ast.VarExpr, g *gamma) (ast.Type, error) {
	typ, ok := g.lookup(expr.Name)
	if !ok {
		return nil, &UndefinedVariableError{
			Pos:  expr.Pos,
			Name: expr.Name,
		}
	}
	return typ, nil
}

func (tc *TypeChecker) checkAbs(expr *ast.AbsExpr, g *gamma) (ast.Type, error) {
	child := g.child()
	child.bind(expr.Param, expr.ParamType)

	bodyType, err := tc.check(expr.Body, child)
	if err != nil {
		return nil, err
	}
	return &ast.FuncType{
		From: expr.ParamType,
		To:   bodyType,
	}, nil
}

func (tc *TypeChecker) checkApp(expr *ast.AppExpr, g *gamma) (ast.Type, error) {
	funcType, err := tc.check(expr.Func, g)
	if err != nil {
		return nil, err
	}

	ft, ok := funcType.(*ast.FuncType)
	if !ok {
		return nil, &NotAFunctionError{
			Pos:  expr.Pos,
			Type: funcType,
		}
	}

	argType, err := tc.check(expr.Arg, g)
	if err != nil {
		return nil, err
	}

	if !typesEqual(ft.From, argType) {
		return nil, &TypeMismatchError{
			Pos:      expr.Pos,
			Expected: ft.From,
			Actual:   argType,
			Context:  "application",
		}
	}

	return ft.To, nil
}

func (tc *TypeChecker) checkIf(expr *ast.IfExpr, g *gamma) (ast.Type, error) {
	condType, err := tc.check(expr.Cond, g)
	if err != nil {
		return nil, err
	}

	if _, ok := condType.(*ast.BooleanType); !ok {
		return nil, &InvalidConditionTypeError{
			Pos:  expr.Pos,
			Type: condType,
		}
	}

	thenType, err := tc.check(expr.Then, g)
	if err != nil {
		return nil, err
	}

	elseType, err := tc.check(expr.Else, g)
	if err != nil {
		return nil, err
	}

	if !typesEqual(thenType, elseType) {
		return nil, &TypeMismatchError{
			Pos:      expr.Pos,
			Expected: thenType,
			Actual:   elseType,
			Context:  "if-else branches",
		}
	}
	return thenType, nil
}

func typesEqual(t1, t2 ast.Type) bool {
	switch typ1 := t1.(type) {
	case *ast.BooleanType:
		_, ok := t2.(*ast.BooleanType)
		return ok
	case *ast.IntType:
		_, ok := t2.(*ast.IntType)
		return ok
	case *ast.FuncType:
		typ2, ok := t2.(*ast.FuncType)
		if !ok {
			return false
		}
		return typesEqual(typ1.From, typ2.From) && typesEqual(typ1.To, typ2.To)
	default:
		return false
	}
}
