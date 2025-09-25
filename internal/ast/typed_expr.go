package ast

import (
	"github.com/shota3506/gostlc/internal/token"
)

type TypedExpr interface {
	typedExprNode()

	Position() token.Position
	Type() Type
}

type TypedVarExpr struct {
	VarExpr

	typ Type
}

func NewTypedVarExpr(typ Type, expr *VarExpr) *TypedVarExpr {
	return &TypedVarExpr{
		VarExpr: *expr,
		typ:     typ,
	}
}

func (TypedVarExpr) typedExprNode()              {}
func (e *TypedVarExpr) Position() token.Position { return e.Pos }
func (e *TypedVarExpr) Type() Type               { return e.typ }

type TypedAbsExpr struct {
	Pos       token.Position
	Param     string
	ParamType Type
	Body      TypedExpr

	typ Type
}

func NewTypedAbsExpr(typ Type, pos token.Position, param string, paramType Type, body TypedExpr) *TypedAbsExpr {
	return &TypedAbsExpr{
		Pos:       pos,
		Param:     param,
		ParamType: paramType,
		Body:      body,
		typ:       typ,
	}
}

func (TypedAbsExpr) typedExprNode()              {}
func (e *TypedAbsExpr) Position() token.Position { return e.Pos }
func (e *TypedAbsExpr) Type() Type               { return e.typ }

type TypedAppExpr struct {
	Pos  token.Position
	Func TypedExpr
	Arg  TypedExpr

	typ Type
}

func NewTypedAppExpr(typ Type, pos token.Position, fn, arg TypedExpr) *TypedAppExpr {
	return &TypedAppExpr{
		Pos:  pos,
		Func: fn,
		Arg:  arg,
		typ:  typ,
	}
}

func (TypedAppExpr) typedExprNode()              {}
func (e *TypedAppExpr) Position() token.Position { return e.Pos }
func (e *TypedAppExpr) Type() Type               { return e.typ }

type TypedBoolExpr struct {
	BoolExpr
}

func NewTypedBoolExpr(expr *BoolExpr) *TypedBoolExpr {
	return &TypedBoolExpr{
		BoolExpr: *expr,
	}
}

func (TypedBoolExpr) typedExprNode()              {}
func (e *TypedBoolExpr) Position() token.Position { return e.Pos }
func (e *TypedBoolExpr) Type() Type               { return &BooleanType{} }

type TypedIntExpr struct {
	IntExpr
}

func NewTypedIntExpr(expr *IntExpr) *TypedIntExpr {
	return &TypedIntExpr{
		IntExpr: *expr,
	}
}

func (TypedIntExpr) typedExprNode()              {}
func (e *TypedIntExpr) Position() token.Position { return e.Pos }
func (e *TypedIntExpr) Type() Type               { return &IntType{} }

type TypedIfExpr struct {
	Pos  token.Position
	Cond TypedExpr
	Then TypedExpr
	Else TypedExpr
}

func NewTypedIfExpr(pos token.Position, cond, then, elseExpr TypedExpr) *TypedIfExpr {
	return &TypedIfExpr{
		Pos:  pos,
		Cond: cond,
		Then: then,
		Else: elseExpr,
	}
}

func (TypedIfExpr) typedExprNode()              {}
func (e *TypedIfExpr) Position() token.Position { return e.Pos }
func (e *TypedIfExpr) Type() Type               { return e.Then.Type() }
