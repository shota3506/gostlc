package ast

// Expr represents an expression in the lambda calculus with simple types.
type Expr interface {
	exprNode()
}

// VarExpr represents a variable expression.
type VarExpr struct {
	Name string
}

func (VarExpr) exprNode() {}

// AbsExpr represents a lambda abstraction expression.
type AbsExpr struct {
	Param     string
	ParamType Type
	Body      Expr
}

func (AbsExpr) exprNode() {}

// AppExpr represents a function application expression.
type AppExpr struct {
	Func Expr
	Arg  Expr
}

func (AppExpr) exprNode() {}

// BoolExpr represents a boolean literal expression.
type BoolExpr struct {
	Value bool
}

func (BoolExpr) exprNode() {}

// IntExpr represents an integer literal expression.
type IntExpr struct {
	Value int
}

func (IntExpr) exprNode() {}

// IfExpr represents an if-then-else expression.
type IfExpr struct {
	Cond Expr
	Then Expr
	Else Expr
}

func (IfExpr) exprNode() {}
