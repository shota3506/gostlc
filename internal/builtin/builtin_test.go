package builtin

import (
	"testing"

	"github.com/shota3506/gostlc/internal/ast"
	"github.com/shota3506/gostlc/internal/values"
)

func TestFunctionTypes(t *testing.T) {
	composeType := func(paramType, returnType ast.Type) ast.Type {
		return &ast.FuncType{
			From: paramType,
			To:   returnType,
		}
	}

	for name, funcType := range FunctionTypes {
		t.Run(name, func(t *testing.T) {
			fn, ok := Functions[name]
			if !ok {
				t.Fatalf("Function %s exists in FunctionTypes but not in Functions", name)
			}

			builtinFunc, ok := fn.(*values.BuiltinFunc)
			if !ok {
				t.Fatalf("Function %s is not a BuiltinFunc", name)
			}

			if !funcType.Equal(composeType(builtinFunc.ParamType, builtinFunc.ReturnType)) {
				t.Errorf("Type mismatch for function %s:\nFunctionTypes: %v\nFunctions: %v->%v",
					name, funcType, builtinFunc.ParamType, builtinFunc.ReturnType)
			}
		})
	}
}

func TestAddFunction(t *testing.T) {
	addFunc := Functions["add"].(*values.BuiltinFunc)

	tests := []struct {
		name     string
		arg1     values.Value
		arg2     values.Value
		expected int
	}{
		{
			name:     "positive numbers",
			arg1:     &values.IntValue{Value: 5},
			arg2:     &values.IntValue{Value: 3},
			expected: 8,
		},
		{
			name:     "negative numbers",
			arg1:     &values.IntValue{Value: -5},
			arg2:     &values.IntValue{Value: -3},
			expected: -8,
		},
		{
			name:     "mixed sign",
			arg1:     &values.IntValue{Value: 10},
			arg2:     &values.IntValue{Value: -7},
			expected: 3,
		},
		{
			name:     "zero",
			arg1:     &values.IntValue{Value: 0},
			arg2:     &values.IntValue{Value: 42},
			expected: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result1, err := addFunc.Fn(tt.arg1)
			if err != nil {
				t.Fatalf("Unexpected error on first application: %v", err)
			}

			partialFunc, ok := result1.(*values.PartialBuiltinFunc)
			if !ok {
				t.Fatalf("First application did not return PartialBuiltinFunc")
			}

			result2, err := partialFunc.Fn(tt.arg2)
			if err != nil {
				t.Fatalf("Unexpected error on second application: %v", err)
			}

			intResult, ok := result2.(*values.IntValue)
			if !ok {
				t.Fatalf("Result is not IntValue")
			}

			if intResult.Value != tt.expected {
				t.Errorf("add(%d, %d) = %d, expected %d",
					tt.arg1.(*values.IntValue).Value,
					tt.arg2.(*values.IntValue).Value,
					intResult.Value,
					tt.expected,
				)
			}
		})
	}
}

func TestSubFunction(t *testing.T) {
	subFunc := Functions["sub"].(*values.BuiltinFunc)

	tests := []struct {
		name     string
		arg1     values.Value
		arg2     values.Value
		expected int
	}{
		{
			name:     "positive numbers",
			arg1:     &values.IntValue{Value: 10},
			arg2:     &values.IntValue{Value: 3},
			expected: 7,
		},
		{
			name:     "negative numbers",
			arg1:     &values.IntValue{Value: -5},
			arg2:     &values.IntValue{Value: -3},
			expected: -2,
		},
		{
			name:     "result is negative",
			arg1:     &values.IntValue{Value: 3},
			arg2:     &values.IntValue{Value: 7},
			expected: -4,
		},
		{
			name:     "zero result",
			arg1:     &values.IntValue{Value: 42},
			arg2:     &values.IntValue{Value: 42},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result1, err := subFunc.Fn(tt.arg1)
			if err != nil {
				t.Fatalf("Unexpected error on first application: %v", err)
			}

			partialFunc, ok := result1.(*values.PartialBuiltinFunc)
			if !ok {
				t.Fatalf("First application did not return PartialBuiltinFunc")
			}

			result2, err := partialFunc.Fn(tt.arg2)
			if err != nil {
				t.Fatalf("Unexpected error on second application: %v", err)
			}

			intResult, ok := result2.(*values.IntValue)
			if !ok {
				t.Fatalf("Result is not IntValue")
			}

			if intResult.Value != tt.expected {
				t.Errorf("sub(%d, %d) = %d, expected %d",
					tt.arg1.(*values.IntValue).Value,
					tt.arg2.(*values.IntValue).Value,
					intResult.Value,
					tt.expected,
				)
			}
		})
	}
}

func TestAndFunction(t *testing.T) {
	andFunc := Functions["and"].(*values.BuiltinFunc)

	tests := []struct {
		name     string
		arg1     values.Value
		arg2     values.Value
		expected bool
	}{
		{
			name:     "true and true",
			arg1:     &values.BoolValue{Value: true},
			arg2:     &values.BoolValue{Value: true},
			expected: true,
		},
		{
			name:     "true and false",
			arg1:     &values.BoolValue{Value: true},
			arg2:     &values.BoolValue{Value: false},
			expected: false,
		},
		{
			name:     "false and true",
			arg1:     &values.BoolValue{Value: false},
			arg2:     &values.BoolValue{Value: true},
			expected: false,
		},
		{
			name:     "false and false",
			arg1:     &values.BoolValue{Value: false},
			arg2:     &values.BoolValue{Value: false},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result1, err := andFunc.Fn(tt.arg1)
			if err != nil {
				t.Fatalf("Unexpected error on first application: %v", err)
			}

			partialFunc, ok := result1.(*values.PartialBuiltinFunc)
			if !ok {
				t.Fatalf("First application did not return PartialBuiltinFunc")
			}

			result2, err := partialFunc.Fn(tt.arg2)
			if err != nil {
				t.Fatalf("Unexpected error on second application: %v", err)
			}

			boolResult, ok := result2.(*values.BoolValue)
			if !ok {
				t.Fatalf("Result is not BoolValue")
			}

			if boolResult.Value != tt.expected {
				t.Errorf("and(%v, %v) = %v, expected %v",
					tt.arg1.(*values.BoolValue).Value,
					tt.arg2.(*values.BoolValue).Value,
					boolResult.Value,
					tt.expected,
				)
			}
		})
	}
}

func TestOrFunction(t *testing.T) {
	orFunc := Functions["or"].(*values.BuiltinFunc)

	tests := []struct {
		name     string
		arg1     values.Value
		arg2     values.Value
		expected bool
	}{
		{
			name:     "true or true",
			arg1:     &values.BoolValue{Value: true},
			arg2:     &values.BoolValue{Value: true},
			expected: true,
		},
		{
			name:     "true or false",
			arg1:     &values.BoolValue{Value: true},
			arg2:     &values.BoolValue{Value: false},
			expected: true,
		},
		{
			name:     "false or true",
			arg1:     &values.BoolValue{Value: false},
			arg2:     &values.BoolValue{Value: true},
			expected: true,
		},
		{
			name:     "false or false",
			arg1:     &values.BoolValue{Value: false},
			arg2:     &values.BoolValue{Value: false},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result1, err := orFunc.Fn(tt.arg1)
			if err != nil {
				t.Fatalf("Unexpected error on first application: %v", err)
			}

			partialFunc, ok := result1.(*values.PartialBuiltinFunc)
			if !ok {
				t.Fatalf("First application did not return PartialBuiltinFunc")
			}

			result2, err := partialFunc.Fn(tt.arg2)
			if err != nil {
				t.Fatalf("Unexpected error on second application: %v", err)
			}

			boolResult, ok := result2.(*values.BoolValue)
			if !ok {
				t.Fatalf("Result is not BoolValue")
			}

			if boolResult.Value != tt.expected {
				t.Errorf("or(%v, %v) = %v, expected %v",
					tt.arg1.(*values.BoolValue).Value,
					tt.arg2.(*values.BoolValue).Value,
					boolResult.Value,
					tt.expected,
				)
			}
		})
	}
}

func TestNotFunction(t *testing.T) {
	notFunc := Functions["not"].(*values.BuiltinFunc)

	tests := []struct {
		name     string
		arg      values.Value
		expected bool
	}{
		{
			name:     "not true",
			arg:      &values.BoolValue{Value: true},
			expected: false,
		},
		{
			name:     "not false",
			arg:      &values.BoolValue{Value: false},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := notFunc.Fn(tt.arg)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			boolResult, ok := result.(*values.BoolValue)
			if !ok {
				t.Fatalf("Result is not BoolValue")
			}

			if boolResult.Value != tt.expected {
				t.Errorf("not(%v) = %v, expected %v",
					tt.arg.(*values.BoolValue).Value,
					boolResult.Value,
					tt.expected,
				)
			}
		})
	}
}
