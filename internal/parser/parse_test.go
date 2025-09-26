package parser_test

import (
	"testing"

	"github.com/shota3506/gostlc/internal/ast"
	"github.com/shota3506/gostlc/internal/parser"
)

func TestParser(t *testing.T) {
	for _, tt := range []struct {
		name     string
		input    string
		expected ast.Expr
	}{
		{
			name:  "Identity function",
			input: `(\x:Bool. x) true`,
			expected: &ast.AppExpr{
				Func: &ast.AbsExpr{
					Param:     "x",
					ParamType: &ast.BoolType{},
					Body:      &ast.VarExpr{Name: "x"},
				},
				Arg: &ast.BoolExpr{Value: true},
			},
		},
		{
			name:  "Identity function with integer",
			input: `(\x:Int. x) 42`,
			expected: &ast.AppExpr{
				Func: &ast.AbsExpr{
					Param:     "x",
					ParamType: &ast.IntType{},
					Body:      &ast.VarExpr{Name: "x"},
				},
				Arg: &ast.IntExpr{Value: 42},
			},
		},
		{
			name:  "Constant function",
			input: `(\x:Int. \y:Int. x) 10 20`,
			expected: &ast.AppExpr{
				Func: &ast.AppExpr{
					Func: &ast.AbsExpr{
						Param:     "x",
						ParamType: &ast.IntType{},
						Body: &ast.AbsExpr{
							Param:     "y",
							ParamType: &ast.IntType{},
							Body:      &ast.VarExpr{Name: "x"},
						},
					},
					Arg: &ast.IntExpr{Value: 10},
				},
				Arg: &ast.IntExpr{Value: 20},
			},
		},
		{
			name:     "Boolean literals true",
			input:    `true`,
			expected: &ast.BoolExpr{Value: true},
		},
		{
			name:     "Boolean literals false",
			input:    `false`,
			expected: &ast.BoolExpr{Value: false},
		},
		{
			name:  "Simple conditional",
			input: `if true then 1 else 0`,
			expected: &ast.IfExpr{
				Cond: &ast.BoolExpr{Value: true},
				Then: &ast.IntExpr{Value: 1},
				Else: &ast.IntExpr{Value: 0},
			},
		},
		{
			name:  "Conditional with false condition",
			input: `if false then 100 else 200`,
			expected: &ast.IfExpr{
				Cond: &ast.BoolExpr{Value: false},
				Then: &ast.IntExpr{Value: 100},
				Else: &ast.IntExpr{Value: 200},
			},
		},
		{
			name:  "Nested function application",
			input: `(\x:Bool. \y:Bool. x) true false`,
			expected: &ast.AppExpr{
				Func: &ast.AppExpr{
					Func: &ast.AbsExpr{
						Param:     "x",
						ParamType: &ast.BoolType{},
						Body: &ast.AbsExpr{
							Param:     "y",
							ParamType: &ast.BoolType{},
							Body:      &ast.VarExpr{Name: "x"},
						},
					},
					Arg: &ast.BoolExpr{Value: true},
				},
				Arg: &ast.BoolExpr{Value: false},
			},
		},
		{
			name:  "Select second argument",
			input: `(\x:Int. \y:Int. y) 5 7`,
			expected: &ast.AppExpr{
				Func: &ast.AppExpr{
					Func: &ast.AbsExpr{
						Param:     "x",
						ParamType: &ast.IntType{},
						Body: &ast.AbsExpr{
							Param:     "y",
							ParamType: &ast.IntType{},
							Body:      &ast.VarExpr{Name: "y"},
						},
					},
					Arg: &ast.IntExpr{Value: 5},
				},
				Arg: &ast.IntExpr{Value: 7},
			},
		},
		{
			name:  "Apply identity to itself",
			input: `(\f:Bool->Bool. f true) (\x:Bool. x)`,
			expected: &ast.AppExpr{
				Func: &ast.AbsExpr{
					Param: "f",
					ParamType: &ast.FuncType{
						From: &ast.BoolType{},
						To:   &ast.BoolType{},
					},
					Body: &ast.AppExpr{
						Func: &ast.VarExpr{Name: "f"},
						Arg:  &ast.BoolExpr{Value: true},
					},
				},
				Arg: &ast.AbsExpr{
					Param:     "x",
					ParamType: &ast.BoolType{},
					Body:      &ast.VarExpr{Name: "x"},
				},
			},
		},
		{
			name:     "Simple integer literal",
			input:    `42`,
			expected: &ast.IntExpr{Value: 42},
		},
		{
			name:     "Negative integer literal",
			input:    `-456`,
			expected: &ast.IntExpr{Value: -456},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.Parse(tt.input)
			if err != nil {
				t.Errorf("Parse() error = %v", err)
				return
			}
			if !equalAST(result, tt.expected) {
				t.Errorf("Parse() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func equalAST(a, b ast.Expr) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	switch x := a.(type) {
	case *ast.VarExpr:
		y, ok := b.(*ast.VarExpr)
		return ok && x.Name == y.Name

	case *ast.AbsExpr:
		y, ok := b.(*ast.AbsExpr)
		return ok && x.Param == y.Param && equalType(x.ParamType, y.ParamType) && equalAST(x.Body, y.Body)

	case *ast.AppExpr:
		y, ok := b.(*ast.AppExpr)
		return ok && equalAST(x.Func, y.Func) && equalAST(x.Arg, y.Arg)

	case *ast.BoolExpr:
		y, ok := b.(*ast.BoolExpr)
		return ok && x.Value == y.Value

	case *ast.IntExpr:
		y, ok := b.(*ast.IntExpr)
		return ok && x.Value == y.Value

	case *ast.IfExpr:
		y, ok := b.(*ast.IfExpr)
		return ok && equalAST(x.Cond, y.Cond) && equalAST(x.Then, y.Then) && equalAST(x.Else, y.Else)

	default:
		return false
	}
}

func equalType(a, b ast.Type) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	switch x := a.(type) {
	case *ast.BoolType:
		_, ok := b.(*ast.BoolType)
		return ok

	case *ast.IntType:
		_, ok := b.(*ast.IntType)
		return ok

	case *ast.FuncType:
		y, ok := b.(*ast.FuncType)
		return ok && equalType(x.From, y.From) && equalType(x.To, y.To)

	default:
		return false
	}
}
