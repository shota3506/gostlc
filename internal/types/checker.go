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

// Check performs type checking and returns a typed AST.
func (tc *TypeChecker) Check(expr ast.Expr) (ast.TypedExpr, error) {
	root := &gamma{bindings: make(map[string]ast.Type)}
	return tc.checkTyped(expr, root)
}

func (tc *TypeChecker) checkTyped(expr ast.Expr, g *gamma) (ast.TypedExpr, error) {
	switch e := expr.(type) {
	case *ast.VarExpr:
		return tc.checkVar(e, g)
	case *ast.AbsExpr:
		return tc.checkAbs(e, g)
	case *ast.AppExpr:
		return tc.checkApp(e, g)
	case *ast.BoolExpr:
		return ast.NewTypedBoolExpr(e), nil
	case *ast.IntExpr:
		return ast.NewTypedIntExpr(e), nil
	case *ast.IfExpr:
		return tc.checkIf(e, g)
	default:
		return nil, &UnknownExprTypeError{
			Pos:  expr.Position(),
			Expr: expr,
		}
	}
}

func (tc *TypeChecker) checkVar(expr *ast.VarExpr, g *gamma) (ast.TypedExpr, error) {
	typ, ok := g.lookup(expr.Name)
	if !ok {
		return nil, &UndefinedVariableError{
			Pos:  expr.Pos,
			Name: expr.Name,
		}
	}
	return ast.NewTypedVarExpr(typ, expr), nil
}

func (tc *TypeChecker) checkAbs(expr *ast.AbsExpr, g *gamma) (ast.TypedExpr, error) {
	child := g.child()
	child.bind(expr.Param, expr.ParamType)

	typedBody, err := tc.checkTyped(expr.Body, child)
	if err != nil {
		return nil, err
	}

	funcType := &ast.FuncType{
		From: expr.ParamType,
		To:   typedBody.Type(),
	}
	return ast.NewTypedAbsExpr(funcType, expr.Pos, expr.Param, expr.ParamType, typedBody), nil
}

func (tc *TypeChecker) checkApp(expr *ast.AppExpr, g *gamma) (ast.TypedExpr, error) {
	typedFunc, err := tc.checkTyped(expr.Func, g)
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

	typedArg, err := tc.checkTyped(expr.Arg, g)
	if err != nil {
		return nil, err
	}

	if !typesEqual(ft.From, typedArg.Type()) {
		return nil, &TypeMismatchError{
			Pos:      expr.Pos,
			Expected: ft.From,
			Actual:   typedArg.Type(),
			Context:  "application",
		}
	}

	return ast.NewTypedAppExpr(ft.To, expr.Pos, typedFunc, typedArg), nil
}

func (tc *TypeChecker) checkIf(expr *ast.IfExpr, g *gamma) (ast.TypedExpr, error) {
	typedCond, err := tc.checkTyped(expr.Cond, g)
	if err != nil {
		return nil, err
	}

	if _, ok := typedCond.Type().(*ast.BooleanType); !ok {
		return nil, &InvalidConditionTypeError{
			Pos:  expr.Pos,
			Type: typedCond.Type(),
		}
	}

	typedThen, err := tc.checkTyped(expr.Then, g)
	if err != nil {
		return nil, err
	}

	typedElse, err := tc.checkTyped(expr.Else, g)
	if err != nil {
		return nil, err
	}

	if !typesEqual(typedThen.Type(), typedElse.Type()) {
		return nil, &TypeMismatchError{
			Pos:      expr.Pos,
			Expected: typedThen.Type(),
			Actual:   typedElse.Type(),
			Context:  "if-else branches",
		}
	}
	return ast.NewTypedIfExpr(expr.Pos, typedCond, typedThen, typedElse), nil
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
