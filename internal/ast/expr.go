package ast

import (
	"github.com/shota3506/gostlc/internal/token"
)

// Expr represents an expression in the lambda calculus with simple types.
type Expr interface {
	exprNode()

	Position() token.Position
}

// VarExpr represents a variable expression.
type VarExpr struct {
	Pos  token.Position
	Name string
}

func (VarExpr) exprNode() {}
func (v VarExpr) Position() token.Position {
	return v.Pos
}

// AbsExpr represents a lambda abstraction expression.
type AbsExpr struct {
	Pos       token.Position
	Param     string
	ParamType Type
	Body      Expr
}

func (AbsExpr) exprNode() {}
func (v AbsExpr) Position() token.Position {
	return v.Pos
}

// AppExpr represents a function application expression.
type AppExpr struct {
	Pos  token.Position
	Func Expr
	Arg  Expr
}

func (AppExpr) exprNode() {}
func (v AppExpr) Position() token.Position {
	return v.Pos
}

// BoolExpr represents a boolean literal expression.
type BoolExpr struct {
	Pos   token.Position
	Value bool
}

func (BoolExpr) exprNode() {}
func (v BoolExpr) Position() token.Position {
	return v.Pos
}

// IntExpr represents an integer literal expression.
type IntExpr struct {
	Pos   token.Position
	Value int
}

func (IntExpr) exprNode() {}
func (v IntExpr) Position() token.Position {
	return v.Pos
}

// IfExpr represents an if-then-else expression.
type IfExpr struct {
	Pos  token.Position
	Cond Expr
	Then Expr
	Else Expr
}

func (IfExpr) exprNode() {}
func (v IfExpr) Position() token.Position {
	return v.Pos
}
