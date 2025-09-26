package builtin

import (
	"errors"

	"github.com/shota3506/gostlc/internal/ast"
	"github.com/shota3506/gostlc/internal/values"
)

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
}

var Functions = map[string]values.Value{
	"add": &values.BuiltinFunc{
		Name:      "add",
		ParamType: &ast.IntType{},
		ReturnType: &ast.FuncType{
			From: &ast.IntType{},
			To:   &ast.IntType{},
		},
		Fn: func(arg1 values.Value) (values.Value, error) {
			intArg1, ok := arg1.(*values.IntValue)
			if !ok {
				return nil, errors.New("type mismatch: expected Int")
			}
			return &values.PartialBuiltinFunc{
				Name:       "add",
				ParamType:  &ast.IntType{},
				ReturnType: &ast.IntType{},
				Fn: func(arg2 values.Value) (values.Value, error) {
					intArg2, ok := arg2.(*values.IntValue)
					if !ok {
						return nil, errors.New("type mismatch: expected Int")
					}
					return &values.IntValue{Value: intArg1.Value + intArg2.Value}, nil
				},
			}, nil
		},
	},
	"sub": &values.BuiltinFunc{
		Name:      "sub",
		ParamType: &ast.IntType{},
		ReturnType: &ast.FuncType{
			From: &ast.IntType{},
			To:   &ast.IntType{},
		},
		Fn: func(arg1 values.Value) (values.Value, error) {
			intArg1, ok := arg1.(*values.IntValue)
			if !ok {
				return nil, errors.New("type mismatch: expected Int")
			}
			return &values.PartialBuiltinFunc{
				Name:       "sub",
				ParamType:  &ast.IntType{},
				ReturnType: &ast.IntType{},
				Fn: func(arg2 values.Value) (values.Value, error) {
					intArg2, ok := arg2.(*values.IntValue)
					if !ok {
						return nil, errors.New("type mismatch: expected Int")
					}
					return &values.IntValue{Value: intArg1.Value - intArg2.Value}, nil
				},
			}, nil
		},
	},
	"and": &values.BuiltinFunc{
		Name:      "and",
		ParamType: &ast.BoolType{},
		ReturnType: &ast.FuncType{
			From: &ast.BoolType{},
			To:   &ast.BoolType{},
		},
		Fn: func(arg1 values.Value) (values.Value, error) {
			boolArg1, ok := arg1.(*values.BoolValue)
			if !ok {
				return nil, errors.New("type mismatch: expected Bool")
			}
			return &values.PartialBuiltinFunc{
				Name:       "and",
				ParamType:  &ast.BoolType{},
				ReturnType: &ast.BoolType{},
				Fn: func(arg2 values.Value) (values.Value, error) {
					boolArg2, ok := arg2.(*values.BoolValue)
					if !ok {
						return nil, errors.New("type mismatch: expected Bool")
					}
					return &values.BoolValue{Value: boolArg1.Value && boolArg2.Value}, nil
				},
			}, nil
		},
	},
	"or": &values.BuiltinFunc{
		Name:      "or",
		ParamType: &ast.BoolType{},
		ReturnType: &ast.FuncType{
			From: &ast.BoolType{},
			To:   &ast.BoolType{},
		},
		Fn: func(arg1 values.Value) (values.Value, error) {
			boolArg1, ok := arg1.(*values.BoolValue)
			if !ok {
				return nil, errors.New("type mismatch: expected Bool")
			}
			return &values.PartialBuiltinFunc{
				Name:       "or",
				ParamType:  &ast.BoolType{},
				ReturnType: &ast.BoolType{},
				Fn: func(arg2 values.Value) (values.Value, error) {
					boolArg2, ok := arg2.(*values.BoolValue)
					if !ok {
						return nil, errors.New("type mismatch: expected Bool")
					}
					return &values.BoolValue{Value: boolArg1.Value || boolArg2.Value}, nil
				},
			}, nil
		},
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
}
