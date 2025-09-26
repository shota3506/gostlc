package types

import (
	"reflect"
	"testing"

	"github.com/shota3506/gostlc/internal/ast"
	"github.com/shota3506/gostlc/internal/token"
)

func pos(line, col int) token.Position {
	return token.Position{Line: line, Column: col}
}

func TestTypeChecker(t *testing.T) {
	tests := []struct {
		name     string
		input    ast.Expr
		expected ast.TypedExpr
	}{
		{
			name:     "boolean literal",
			input:    &ast.BoolExpr{Pos: pos(1, 1), Value: true},
			expected: ast.NewTypedBoolExpr(&ast.BoolExpr{Pos: pos(1, 1), Value: true}),
		},
		{
			name:     "integer literal",
			input:    &ast.IntExpr{Pos: pos(1, 1), Value: 42},
			expected: ast.NewTypedIntExpr(&ast.IntExpr{Pos: pos(1, 1), Value: 42}),
		},
		{
			name:  "builtin add function",
			input: &ast.VarExpr{Pos: pos(1, 1), Name: "add"},
			expected: ast.NewTypedVarExpr(
				&ast.FuncType{
					From: &ast.IntType{},
					To: &ast.FuncType{
						From: &ast.IntType{},
						To:   &ast.IntType{},
					},
				},
				&ast.VarExpr{Pos: pos(1, 1), Name: "add"},
			),
		},
		{
			name:  "builtin sub function",
			input: &ast.VarExpr{Pos: pos(1, 1), Name: "sub"},
			expected: ast.NewTypedVarExpr(
				&ast.FuncType{
					From: &ast.IntType{},
					To: &ast.FuncType{
						From: &ast.IntType{},
						To:   &ast.IntType{},
					},
				},
				&ast.VarExpr{Pos: pos(1, 1), Name: "sub"},
			),
		},
		{
			name: "builtin add with one argument",
			input: &ast.AppExpr{
				Pos:  pos(1, 1),
				Func: &ast.VarExpr{Pos: pos(1, 1), Name: "add"},
				Arg:  &ast.IntExpr{Pos: pos(1, 5), Value: 1},
			},
			expected: ast.NewTypedAppExpr(
				&ast.FuncType{
					From: &ast.IntType{},
					To:   &ast.IntType{},
				},
				pos(1, 1),
				ast.NewTypedVarExpr(
					&ast.FuncType{
						From: &ast.IntType{},
						To: &ast.FuncType{
							From: &ast.IntType{},
							To:   &ast.IntType{},
						},
					},
					&ast.VarExpr{Pos: pos(1, 1), Name: "add"},
				),
				ast.NewTypedIntExpr(&ast.IntExpr{Pos: pos(1, 5), Value: 1}),
			),
		},
		{
			name: "builtin add with two arguments",
			input: &ast.AppExpr{
				Pos: pos(1, 1),
				Func: &ast.AppExpr{
					Pos:  pos(1, 1),
					Func: &ast.VarExpr{Pos: pos(1, 1), Name: "add"},
					Arg:  &ast.IntExpr{Pos: pos(1, 5), Value: 1},
				},
				Arg: &ast.IntExpr{Pos: pos(1, 7), Value: 2},
			},
			expected: ast.NewTypedAppExpr(
				&ast.IntType{},
				pos(1, 1),
				ast.NewTypedAppExpr(
					&ast.FuncType{
						From: &ast.IntType{},
						To:   &ast.IntType{},
					},
					pos(1, 1),
					ast.NewTypedVarExpr(
						&ast.FuncType{
							From: &ast.IntType{},
							To: &ast.FuncType{
								From: &ast.IntType{},
								To:   &ast.IntType{},
							},
						},
						&ast.VarExpr{Pos: pos(1, 1), Name: "add"},
					),
					ast.NewTypedIntExpr(&ast.IntExpr{Pos: pos(1, 5), Value: 1}),
				),
				ast.NewTypedIntExpr(&ast.IntExpr{Pos: pos(1, 7), Value: 2}),
			),
		},
		{
			name: "builtin sub with two arguments",
			input: &ast.AppExpr{
				Pos: pos(1, 1),
				Func: &ast.AppExpr{
					Pos:  pos(1, 1),
					Func: &ast.VarExpr{Pos: pos(1, 1), Name: "sub"},
					Arg:  &ast.IntExpr{Pos: pos(1, 5), Value: 10},
				},
				Arg: &ast.IntExpr{Pos: pos(1, 8), Value: 3},
			},
			expected: ast.NewTypedAppExpr(
				&ast.IntType{},
				pos(1, 1),
				ast.NewTypedAppExpr(
					&ast.FuncType{
						From: &ast.IntType{},
						To:   &ast.IntType{},
					},
					pos(1, 1),
					ast.NewTypedVarExpr(
						&ast.FuncType{
							From: &ast.IntType{},
							To: &ast.FuncType{
								From: &ast.IntType{},
								To:   &ast.IntType{},
							},
						},
						&ast.VarExpr{Pos: pos(1, 1), Name: "sub"},
					),
					ast.NewTypedIntExpr(&ast.IntExpr{Pos: pos(1, 5), Value: 10}),
				),
				ast.NewTypedIntExpr(&ast.IntExpr{Pos: pos(1, 8), Value: 3}),
			),
		},
		{
			name: "identity function",
			input: &ast.AbsExpr{
				Pos:       pos(1, 1),
				Param:     "x",
				ParamType: &ast.BoolType{},
				Body:      &ast.VarExpr{Pos: pos(1, 10), Name: "x"},
			},
			expected: ast.NewTypedAbsExpr(
				&ast.FuncType{
					From: &ast.BoolType{},
					To:   &ast.BoolType{},
				},
				pos(1, 1),
				"x",
				&ast.BoolType{},
				ast.NewTypedVarExpr(&ast.BoolType{}, &ast.VarExpr{Pos: pos(1, 10), Name: "x"}),
			),
		},
		{
			name: "const function",
			input: &ast.AbsExpr{
				Pos:       pos(1, 1),
				Param:     "x",
				ParamType: &ast.IntType{},
				Body: &ast.AbsExpr{
					Pos:       pos(2, 1),
					Param:     "y",
					ParamType: &ast.BoolType{},
					Body:      &ast.VarExpr{Pos: pos(2, 10), Name: "x"},
				},
			},
			expected: ast.NewTypedAbsExpr(
				&ast.FuncType{
					From: &ast.IntType{},
					To: &ast.FuncType{
						From: &ast.BoolType{},
						To:   &ast.IntType{},
					},
				},
				pos(1, 1),
				"x",
				&ast.IntType{},
				ast.NewTypedAbsExpr(
					&ast.FuncType{
						From: &ast.BoolType{},
						To:   &ast.IntType{},
					},
					pos(2, 1),
					"y",
					&ast.BoolType{},
					ast.NewTypedVarExpr(&ast.IntType{}, &ast.VarExpr{Pos: pos(2, 10), Name: "x"}),
				),
			),
		},
		{
			name: "function application",
			input: &ast.AppExpr{
				Pos: pos(1, 1),
				Func: &ast.AbsExpr{
					Pos:       pos(1, 1),
					Param:     "x",
					ParamType: &ast.BoolType{},
					Body:      &ast.VarExpr{Pos: pos(1, 10), Name: "x"},
				},
				Arg: &ast.BoolExpr{Pos: pos(1, 15), Value: true},
			},
			expected: ast.NewTypedAppExpr(
				&ast.BoolType{},
				pos(1, 1),
				ast.NewTypedAbsExpr(
					&ast.FuncType{
						From: &ast.BoolType{},
						To:   &ast.BoolType{},
					},
					pos(1, 1),
					"x",
					&ast.BoolType{},
					ast.NewTypedVarExpr(&ast.BoolType{}, &ast.VarExpr{Pos: pos(1, 10), Name: "x"}),
				),
				ast.NewTypedBoolExpr(&ast.BoolExpr{Pos: pos(1, 15), Value: true}),
			),
		},
		{
			name: "if expression with booleans",
			input: &ast.IfExpr{
				Pos:  pos(1, 1),
				Cond: &ast.BoolExpr{Pos: pos(1, 4), Value: true},
				Then: &ast.BoolExpr{Pos: pos(1, 10), Value: false},
				Else: &ast.BoolExpr{Pos: pos(1, 20), Value: true},
			},
			expected: ast.NewTypedIfExpr(
				pos(1, 1),
				ast.NewTypedBoolExpr(&ast.BoolExpr{Pos: pos(1, 4), Value: true}),
				ast.NewTypedBoolExpr(&ast.BoolExpr{Pos: pos(1, 10), Value: false}),
				ast.NewTypedBoolExpr(&ast.BoolExpr{Pos: pos(1, 20), Value: true}),
			),
		},
		{
			name: "if expression with integers",
			input: &ast.IfExpr{
				Pos:  pos(1, 1),
				Cond: &ast.BoolExpr{Pos: pos(1, 4), Value: true},
				Then: &ast.IntExpr{Pos: pos(1, 10), Value: 1},
				Else: &ast.IntExpr{Pos: pos(1, 15), Value: 2},
			},
			expected: ast.NewTypedIfExpr(
				pos(1, 1),
				ast.NewTypedBoolExpr(&ast.BoolExpr{Pos: pos(1, 4), Value: true}),
				ast.NewTypedIntExpr(&ast.IntExpr{Pos: pos(1, 10), Value: 1}),
				ast.NewTypedIntExpr(&ast.IntExpr{Pos: pos(1, 15), Value: 2}),
			),
		},
		{
			name: "higher-order function",
			input: &ast.AbsExpr{
				Pos:   pos(1, 1),
				Param: "f",
				ParamType: &ast.FuncType{
					From: &ast.BoolType{},
					To:   &ast.IntType{},
				},
				Body: &ast.AbsExpr{
					Pos:       pos(2, 1),
					Param:     "x",
					ParamType: &ast.BoolType{},
					Body: &ast.AppExpr{
						Pos:  pos(3, 1),
						Func: &ast.VarExpr{Pos: pos(3, 1), Name: "f"},
						Arg:  &ast.VarExpr{Pos: pos(3, 10), Name: "x"},
					},
				},
			},
			expected: ast.NewTypedAbsExpr(
				&ast.FuncType{
					From: &ast.FuncType{
						From: &ast.BoolType{},
						To:   &ast.IntType{},
					},
					To: &ast.FuncType{
						From: &ast.BoolType{},
						To:   &ast.IntType{},
					},
				},
				pos(1, 1),
				"f",
				&ast.FuncType{
					From: &ast.BoolType{},
					To:   &ast.IntType{},
				},
				ast.NewTypedAbsExpr(
					&ast.FuncType{
						From: &ast.BoolType{},
						To:   &ast.IntType{},
					},
					pos(2, 1),
					"x",
					&ast.BoolType{},
					ast.NewTypedAppExpr(
						&ast.IntType{},
						pos(3, 1),
						ast.NewTypedVarExpr(
							&ast.FuncType{
								From: &ast.BoolType{},
								To:   &ast.IntType{},
							},
							&ast.VarExpr{Pos: pos(3, 1), Name: "f"},
						),
						ast.NewTypedVarExpr(
							&ast.BoolType{},
							&ast.VarExpr{Pos: pos(3, 10), Name: "x"},
						),
					),
				),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := Check(tt.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !compareTypedExprs(actual, tt.expected) {
				t.Errorf("typed AST mismatch:\ngot:  %#v\nwant: %#v", actual, tt.expected)
			}
		})
	}
}

func TestTypeCheckerErrors(t *testing.T) {
	tests := []struct {
		name          string
		input         ast.Expr
		expectedError string
	}{
		{
			name:          "undefined variable",
			input:         &ast.VarExpr{Pos: token.Position{Line: 1, Column: 1}, Name: "x"},
			expectedError: "1:1: undefined variable: x",
		},
		{
			name: "type mismatch in application",
			input: &ast.AppExpr{
				Pos: pos(1, 1),
				Func: &ast.AbsExpr{
					Pos:       pos(1, 1),
					Param:     "x",
					ParamType: &ast.BoolType{},
					Body:      &ast.VarExpr{Pos: pos(1, 10), Name: "x"},
				},
				Arg: &ast.IntExpr{Pos: pos(1, 15), Value: 42},
			},
			expectedError: "1:1: type mismatch in application: expected Bool, got Int",
		},
		{
			name: "applying non-function",
			input: &ast.AppExpr{
				Pos:  pos(1, 1),
				Func: &ast.IntExpr{Pos: pos(1, 1), Value: 42},
				Arg:  &ast.BoolExpr{Pos: pos(1, 4), Value: true},
			},
			expectedError: "1:1: cannot apply non-function type: Int",
		},
		{
			name: "non-boolean condition",
			input: &ast.IfExpr{
				Pos:  pos(1, 1),
				Cond: &ast.IntExpr{Pos: pos(1, 4), Value: 42},
				Then: &ast.BoolExpr{Pos: pos(1, 10), Value: true},
				Else: &ast.BoolExpr{Pos: pos(1, 20), Value: false},
			},
			expectedError: "1:1: condition must be boolean, got Int",
		},
		{
			name: "mismatched if branches",
			input: &ast.IfExpr{
				Pos:  pos(1, 1),
				Cond: &ast.BoolExpr{Pos: pos(1, 4), Value: true},
				Then: &ast.IntExpr{Pos: pos(1, 10), Value: 42},
				Else: &ast.BoolExpr{Pos: pos(1, 20), Value: false},
			},
			expectedError: "1:1: type mismatch in if-else branches: expected Int, got Bool",
		},
		{
			name: "undefined variable in abstraction body",
			input: &ast.AbsExpr{
				Pos:       pos(1, 1),
				Param:     "x",
				ParamType: &ast.BoolType{},
				Body: &ast.VarExpr{
					Pos:  pos(1, 10),
					Name: "y",
				},
			},
			expectedError: "1:10: undefined variable: y",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Check(tt.input)

			if err == nil {
				t.Errorf("expected error, but got nil")
			} else if err.Error() != tt.expectedError {
				t.Errorf("error mismatch: got %v, want %v", err.Error(), tt.expectedError)
			}
		})
	}
}

func TestTypesEqual(t *testing.T) {
	tests := []struct {
		name  string
		t1    ast.Type
		t2    ast.Type
		equal bool
	}{
		{
			name:  "same boolean types",
			t1:    &ast.BoolType{},
			t2:    &ast.BoolType{},
			equal: true,
		},
		{
			name:  "same integer types",
			t1:    &ast.IntType{},
			t2:    &ast.IntType{},
			equal: true,
		},
		{
			name:  "different base types",
			t1:    &ast.BoolType{},
			t2:    &ast.IntType{},
			equal: false,
		},
		{
			name: "same function types",
			t1: &ast.FuncType{
				From: &ast.BoolType{},
				To:   &ast.IntType{},
			},
			t2: &ast.FuncType{
				From: &ast.BoolType{},
				To:   &ast.IntType{},
			},
			equal: true,
		},
		{
			name: "different function parameter types",
			t1: &ast.FuncType{
				From: &ast.BoolType{},
				To:   &ast.IntType{},
			},
			t2: &ast.FuncType{
				From: &ast.IntType{},
				To:   &ast.IntType{},
			},
			equal: false,
		},
		{
			name: "different function return types",
			t1: &ast.FuncType{
				From: &ast.BoolType{},
				To:   &ast.IntType{},
			},
			t2: &ast.FuncType{
				From: &ast.BoolType{},
				To:   &ast.BoolType{},
			},
			equal: false,
		},
		{
			name: "nested function types",
			t1: &ast.FuncType{
				From: &ast.FuncType{
					From: &ast.BoolType{},
					To:   &ast.IntType{},
				},
				To: &ast.BoolType{},
			},
			t2: &ast.FuncType{
				From: &ast.FuncType{
					From: &ast.BoolType{},
					To:   &ast.IntType{},
				},
				To: &ast.BoolType{},
			},
			equal: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t1.Equal(tt.t2); got != tt.equal {
				t.Errorf("types.Equal() = %v, want %v", got, tt.equal)
			}
		})
	}
}

// compareTypedExprs compares two TypedExpr instances for deep equality
func compareTypedExprs(actual, expected ast.TypedExpr) bool {
	if actual == nil && expected == nil {
		return true
	}
	if actual == nil || expected == nil {
		return false
	}

	switch a := actual.(type) {
	case *ast.TypedBoolExpr:
		e, ok := expected.(*ast.TypedBoolExpr)
		if !ok {
			return false
		}
		return a.Value == e.Value && a.Pos == e.Pos

	case *ast.TypedIntExpr:
		e, ok := expected.(*ast.TypedIntExpr)
		if !ok {
			return false
		}
		return a.Value == e.Value && a.Pos == e.Pos

	case *ast.TypedVarExpr:
		e, ok := expected.(*ast.TypedVarExpr)
		if !ok {
			return false
		}
		return a.Name == e.Name && a.Pos == e.Pos && reflect.DeepEqual(a.Type(), e.Type())

	case *ast.TypedAbsExpr:
		e, ok := expected.(*ast.TypedAbsExpr)
		if !ok {
			return false
		}
		if a.Param != e.Param || a.Pos != e.Pos {
			return false
		}
		if !reflect.DeepEqual(a.ParamType, e.ParamType) {
			return false
		}
		return compareTypedExprs(a.Body, e.Body) && reflect.DeepEqual(a.Type(), e.Type())

	case *ast.TypedAppExpr:
		e, ok := expected.(*ast.TypedAppExpr)
		if !ok {
			return false
		}
		if a.Pos != e.Pos {
			return false
		}
		return compareTypedExprs(a.Func, e.Func) && compareTypedExprs(a.Arg, e.Arg) && reflect.DeepEqual(a.Type(), e.Type())

	case *ast.TypedIfExpr:
		e, ok := expected.(*ast.TypedIfExpr)
		if !ok {
			return false
		}
		if a.Pos != e.Pos {
			return false
		}
		return compareTypedExprs(a.Cond, e.Cond) &&
			compareTypedExprs(a.Then, e.Then) &&
			compareTypedExprs(a.Else, e.Else) &&
			reflect.DeepEqual(a.Type(), e.Type())

	default:
		return false
	}
}
