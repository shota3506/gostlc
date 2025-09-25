package builtin

import (
	"errors"

	"github.com/shota3506/gostlc/internal/ast"
	"github.com/shota3506/gostlc/internal/values"
)

var FunctionTypes = map[string]ast.Type{
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
}
