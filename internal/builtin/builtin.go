package builtin

import (
	"errors"
	"fmt"

	"github.com/shota3506/gostlc/internal/ast"
	"github.com/shota3506/gostlc/internal/values"
)

func binaryOp[T, U values.Value](f func(a, b T) U) func(arg1 values.Value) (values.Value, error) {
	return func(arg1 values.Value) (values.Value, error) {
		tArg1, ok := arg1.(T)
		if !ok {
			return nil, fmt.Errorf("type mismatch")
		}
		return &values.PartialBuiltinFunc{
			Name:       "and",
			ParamType:  &ast.BoolType{},
			ReturnType: &ast.BoolType{},
			Fn: func(arg2 values.Value) (values.Value, error) {
				tArg2, ok := arg2.(T)
				if !ok {
					return nil, errors.New("type mismatch")
				}
				return f(tArg1, tArg2), nil
			},
		}, nil
	}
}

var FunctionTypes = map[string]ast.Type{
	// Arithmetic operations
	"add": &ast.FuncType{
		From: &ast.IntType{},
		To: &ast.FuncType{
			From: &ast.IntType{},
			To:   &ast.IntType{},
		},
	},
	"sub": &ast.FuncType{
		From: &ast.IntType{},
		To: &ast.FuncType{
			From: &ast.IntType{},
			To:   &ast.IntType{},
		},
	},
	// Boolean operations
	"and": &ast.FuncType{
		From: &ast.BoolType{},
		To: &ast.FuncType{
			From: &ast.BoolType{},
			To:   &ast.BoolType{},
		},
	},
	"or": &ast.FuncType{
		From: &ast.BoolType{},
		To: &ast.FuncType{
			From: &ast.BoolType{},
			To:   &ast.BoolType{},
		},
	},
	"not": &ast.FuncType{
		From: &ast.BoolType{},
		To:   &ast.BoolType{},
	},
	// Comparison operations
	"eq": &ast.FuncType{
		From: &ast.IntType{},
		To: &ast.FuncType{
			From: &ast.IntType{},
			To:   &ast.BoolType{},
		},
	},
	"ne": &ast.FuncType{
		From: &ast.IntType{},
		To: &ast.FuncType{
			From: &ast.IntType{},
			To:   &ast.BoolType{},
		},
	},
	"lt": &ast.FuncType{
		From: &ast.IntType{},
		To: &ast.FuncType{
			From: &ast.IntType{},
			To:   &ast.BoolType{},
		},
	},
	"le": &ast.FuncType{
		From: &ast.IntType{},
		To: &ast.FuncType{
			From: &ast.IntType{},
			To:   &ast.BoolType{},
		},
	},
	"gt": &ast.FuncType{
		From: &ast.IntType{},
		To: &ast.FuncType{
			From: &ast.IntType{},
			To:   &ast.BoolType{},
		},
	},
	"ge": &ast.FuncType{
		From: &ast.IntType{},
		To: &ast.FuncType{
			From: &ast.IntType{},
			To:   &ast.BoolType{},
		},
	},
}

var Functions = map[string]values.Value{
	"add": &values.BuiltinFunc{
		Name:      "add",
		ParamType: &ast.IntType{},
		ReturnType: &ast.FuncType{
			From: &ast.IntType{},
			To:   &ast.IntType{},
		},
		Fn: binaryOp(func(a, b *values.IntValue) *values.IntValue {
			return &values.IntValue{Value: a.Value + b.Value}
		}),
	},
	"sub": &values.BuiltinFunc{
		Name:      "sub",
		ParamType: &ast.IntType{},
		ReturnType: &ast.FuncType{
			From: &ast.IntType{},
			To:   &ast.IntType{},
		},
		Fn: binaryOp(func(a, b *values.IntValue) *values.IntValue {
			return &values.IntValue{Value: a.Value - b.Value}
		}),
	},
	"and": &values.BuiltinFunc{
		Name:      "and",
		ParamType: &ast.BoolType{},
		ReturnType: &ast.FuncType{
			From: &ast.BoolType{},
			To:   &ast.BoolType{},
		},
		Fn: binaryOp(func(a, b *values.BoolValue) *values.BoolValue {
			return &values.BoolValue{Value: a.Value && b.Value}
		}),
	},
	"or": &values.BuiltinFunc{
		Name:      "or",
		ParamType: &ast.BoolType{},
		ReturnType: &ast.FuncType{
			From: &ast.BoolType{},
			To:   &ast.BoolType{},
		},
		Fn: binaryOp(func(a, b *values.BoolValue) *values.BoolValue {
			return &values.BoolValue{Value: a.Value || b.Value}
		}),
	},
	"not": &values.BuiltinFunc{
		Name:       "not",
		ParamType:  &ast.BoolType{},
		ReturnType: &ast.BoolType{},
		Fn: func(arg values.Value) (values.Value, error) {
			boolArg, ok := arg.(*values.BoolValue)
			if !ok {
				return nil, errors.New("type mismatch: expected Bool")
			}
			return &values.BoolValue{Value: !boolArg.Value}, nil
		},
	},
	"eq": &values.BuiltinFunc{
		Name:      "eq",
		ParamType: &ast.IntType{},
		ReturnType: &ast.FuncType{
			From: &ast.IntType{},
			To:   &ast.BoolType{},
		},
		Fn: binaryOp(func(a, b *values.IntValue) *values.BoolValue {
			return &values.BoolValue{Value: a.Value == b.Value}
		}),
	},
	"ne": &values.BuiltinFunc{
		Name:      "ne",
		ParamType: &ast.IntType{},
		ReturnType: &ast.FuncType{
			From: &ast.IntType{},
			To:   &ast.BoolType{},
		},
		Fn: binaryOp(func(a, b *values.IntValue) *values.BoolValue {
			return &values.BoolValue{Value: a.Value != b.Value}
		}),
	},
	"lt": &values.BuiltinFunc{
		Name:      "lt",
		ParamType: &ast.IntType{},
		ReturnType: &ast.FuncType{
			From: &ast.IntType{},
			To:   &ast.BoolType{},
		},
		Fn: binaryOp(func(a, b *values.IntValue) *values.BoolValue {
			return &values.BoolValue{Value: a.Value < b.Value}
		}),
	},
	"le": &values.BuiltinFunc{
		Name:      "le",
		ParamType: &ast.IntType{},
		ReturnType: &ast.FuncType{
			From: &ast.IntType{},
			To:   &ast.BoolType{},
		},
		Fn: binaryOp(func(a, b *values.IntValue) *values.BoolValue {
			return &values.BoolValue{Value: a.Value <= b.Value}
		}),
	},
	"gt": &values.BuiltinFunc{
		Name:      "gt",
		ParamType: &ast.IntType{},
		ReturnType: &ast.FuncType{
			From: &ast.IntType{},
			To:   &ast.BoolType{},
		},
		Fn: binaryOp(func(a, b *values.IntValue) *values.BoolValue {
			return &values.BoolValue{Value: a.Value > b.Value}
		}),
	},
	"ge": &values.BuiltinFunc{
		Name:      "ge",
		ParamType: &ast.IntType{},
		ReturnType: &ast.FuncType{
			From: &ast.IntType{},
			To:   &ast.BoolType{},
		},
		Fn: binaryOp(func(a, b *values.IntValue) *values.BoolValue {
			return &values.BoolValue{Value: a.Value >= b.Value}
		}),
	},
}
