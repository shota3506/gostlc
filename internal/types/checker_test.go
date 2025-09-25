package types

import (
	"testing"

	"github.com/shota3506/gostlc/internal/ast"
)

func TestTypeChecker(t *testing.T) {
	tests := []struct {
		name         string
		input        ast.Expr
		expectedType ast.Type
	}{
		{
			name:         "boolean literal",
			input:        &ast.BoolExpr{Value: true},
			expectedType: &ast.BooleanType{},
		},
		{
			name:         "integer literal",
			input:        &ast.IntExpr{Value: 42},
			expectedType: &ast.IntType{},
		},
		{
			name: "identity function",
			input: &ast.AbsExpr{
				Param:     "x",
				ParamType: &ast.BooleanType{},
				Body:      &ast.VarExpr{Name: "x"},
			},
			expectedType: &ast.FuncType{
				From: &ast.BooleanType{},
				To:   &ast.BooleanType{},
			},
		},
		{
			name: "const function",
			input: &ast.AbsExpr{
				Param:     "x",
				ParamType: &ast.IntType{},
				Body: &ast.AbsExpr{
					Param:     "y",
					ParamType: &ast.BooleanType{},
					Body:      &ast.VarExpr{Name: "x"},
				},
			},
			expectedType: &ast.FuncType{
				From: &ast.IntType{},
				To: &ast.FuncType{
					From: &ast.BooleanType{},
					To:   &ast.IntType{},
				},
			},
		},
		{
			name: "function application",
			input: &ast.AppExpr{
				Func: &ast.AbsExpr{
					Param:     "x",
					ParamType: &ast.BooleanType{},
					Body:      &ast.VarExpr{Name: "x"},
				},
				Arg: &ast.BoolExpr{Value: true},
			},
			expectedType: &ast.BooleanType{},
		},
		{
			name: "if expression with booleans",
			input: &ast.IfExpr{
				Cond: &ast.BoolExpr{Value: true},
				Then: &ast.BoolExpr{Value: false},
				Else: &ast.BoolExpr{Value: true},
			},
			expectedType: &ast.BooleanType{},
		},
		{
			name: "if expression with integers",
			input: &ast.IfExpr{
				Cond: &ast.BoolExpr{Value: true},
				Then: &ast.IntExpr{Value: 1},
				Else: &ast.IntExpr{Value: 2},
			},
			expectedType: &ast.IntType{},
		},
		{
			name: "higher-order function",
			input: &ast.AbsExpr{
				Param: "f",
				ParamType: &ast.FuncType{
					From: &ast.BooleanType{},
					To:   &ast.IntType{},
				},
				Body: &ast.AbsExpr{
					Param:     "x",
					ParamType: &ast.BooleanType{},
					Body: &ast.AppExpr{
						Func: &ast.VarExpr{Name: "f"},
						Arg:  &ast.VarExpr{Name: "x"},
					},
				},
			},
			expectedType: &ast.FuncType{
				From: &ast.FuncType{
					From: &ast.BooleanType{},
					To:   &ast.IntType{},
				},
				To: &ast.FuncType{
					From: &ast.BooleanType{},
					To:   &ast.IntType{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := NewTypeChecker()
			typ, err := tc.Check(tt.input)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if typ.String() != tt.expectedType.String() {
				t.Errorf("type mismatch: got %s, want %s", typ, tt.expectedType)
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
			input:         &ast.VarExpr{Name: "x"},
			expectedError: "undefined variable: x",
		},
		{
			name: "type mismatch in application",
			input: &ast.AppExpr{
				Func: &ast.AbsExpr{
					Param:     "x",
					ParamType: &ast.BooleanType{},
					Body:      &ast.VarExpr{Name: "x"},
				},
				Arg: &ast.IntExpr{Value: 42},
			},
			expectedError: "type mismatch in application: expected Bool, got Int",
		},
		{
			name: "applying non-function",
			input: &ast.AppExpr{
				Func: &ast.IntExpr{Value: 42},
				Arg:  &ast.BoolExpr{Value: true},
			},
			expectedError: "cannot apply non-function type: Int",
		},
		{
			name: "non-boolean condition",
			input: &ast.IfExpr{
				Cond: &ast.IntExpr{Value: 42},
				Then: &ast.BoolExpr{Value: true},
				Else: &ast.BoolExpr{Value: false},
			},
			expectedError: "condition must be boolean: Int",
		},
		{
			name: "mismatched if branches",
			input: &ast.IfExpr{
				Cond: &ast.BoolExpr{Value: true},
				Then: &ast.IntExpr{Value: 42},
				Else: &ast.BoolExpr{Value: false},
			},
			expectedError: "type mismatch in if-else branches: expected Int, got Bool",
		},
		{
			name: "undefined variable in abstraction body",
			input: &ast.AbsExpr{
				Param:     "x",
				ParamType: &ast.BooleanType{},
				Body:      &ast.VarExpr{Name: "y"},
			},
			expectedError: "undefined variable: y",
		},
		{
			name: "type mismatch in nested application",
			input: &ast.AppExpr{
				Func: &ast.AbsExpr{
					Param:     "f",
					ParamType: &ast.FuncType{From: &ast.BooleanType{}, To: &ast.IntType{}},
					Body: &ast.AppExpr{
						Func: &ast.VarExpr{Name: "f"},
						Arg:  &ast.IntExpr{Value: 42},
					},
				},
				Arg: &ast.AbsExpr{
					Param:     "x",
					ParamType: &ast.BooleanType{},
					Body:      &ast.IntExpr{Value: 0},
				},
			},
			expectedError: "type mismatch in application: expected Bool, got Int",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := NewTypeChecker()
			_, err := tc.Check(tt.input)

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
			t1:    &ast.BooleanType{},
			t2:    &ast.BooleanType{},
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
			t1:    &ast.BooleanType{},
			t2:    &ast.IntType{},
			equal: false,
		},
		{
			name: "same function types",
			t1: &ast.FuncType{
				From: &ast.BooleanType{},
				To:   &ast.IntType{},
			},
			t2: &ast.FuncType{
				From: &ast.BooleanType{},
				To:   &ast.IntType{},
			},
			equal: true,
		},
		{
			name: "different function parameter types",
			t1: &ast.FuncType{
				From: &ast.BooleanType{},
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
				From: &ast.BooleanType{},
				To:   &ast.IntType{},
			},
			t2: &ast.FuncType{
				From: &ast.BooleanType{},
				To:   &ast.BooleanType{},
			},
			equal: false,
		},
		{
			name: "nested function types",
			t1: &ast.FuncType{
				From: &ast.FuncType{
					From: &ast.BooleanType{},
					To:   &ast.IntType{},
				},
				To: &ast.BooleanType{},
			},
			t2: &ast.FuncType{
				From: &ast.FuncType{
					From: &ast.BooleanType{},
					To:   &ast.IntType{},
				},
				To: &ast.BooleanType{},
			},
			equal: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := typesEqual(tt.t1, tt.t2)
			if got != tt.equal {
				t.Errorf("typesEqual() = %v, want %v", got, tt.equal)
			}
		})
	}
}
