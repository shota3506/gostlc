package eval

import (
	"testing"

	"github.com/shota3506/gostlc/internal/parser"
	"github.com/shota3506/gostlc/internal/types"
	"github.com/shota3506/gostlc/internal/values"
)

func TestEvalIntegerExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"int literal 42", "42", 42},
		{"int literal 0", "0", 0},
		{"int literal 100", "100", 100},
		{"identity function", "(\\x:Int.x) 42", 42},
		{"const function first", "((\\x:Int.\\y:Int.x) 5) 10", 5},
		{"const function second", "((\\x:Int.\\y:Int.y) 5) 10", 10},
		{"if true branch", "if true then 10 else 20", 10},
		{"if false branch", "if false then 10 else 20", 20},
		{"lambda with if true", "(\\x:Bool.if x then 100 else 200) true", 100},
		{"lambda with if false", "(\\x:Bool.if x then 100 else 200) false", 200},
		{"higher order identity", "(\\f:Int->Int.\\x:Int.f x) (\\y:Int.y) 42", 42},
		{"complex nesting", "(\\f:Int->Int->Int.\\x:Int.\\y:Int.f x y) (\\a:Int.\\b:Int.a) 10 20", 10},
		{"complex nesting swap", "(\\f:Int->Int->Int.\\x:Int.\\y:Int.f x y) (\\a:Int.\\b:Int.b) 10 20", 20},
		{"application chain", "(\\x:Int.\\y:Int.\\z:Int.z) 1 2 3", 3},
		{"function composition", "(\\f:Int->Int.\\g:Int->Int.\\x:Int.f (g x)) (\\a:Int.a) (\\b:Int.b) 42", 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := eval(t, tt.input)

			intVal, ok := val.(*values.IntValue)
			if !ok {
				t.Fatalf("expected IntValue, got %T", val)
			}

			if intVal.Value != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, intVal.Value)
			}
		})
	}
}

func TestEvalBooleanExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"bool literal true", "true", true},
		{"bool literal false", "false", false},
		{"identity bool true", "(\\x:Bool.x) true", true},
		{"identity bool false", "(\\x:Bool.x) false", false},
		{"const bool first true", "((\\x:Bool.\\y:Bool.x) true) false", true},
		{"const bool first false", "((\\x:Bool.\\y:Bool.x) false) true", false},
		{"const bool second true", "((\\x:Bool.\\y:Bool.y) false) true", true},
		{"const bool second false", "((\\x:Bool.\\y:Bool.y) true) false", false},
		{"if in lambda true", "(\\x:Int.if true then x else x) 5", true}, // This should return int, not bool
	}

	// Remove the last incorrect test case
	tests = tests[:len(tests)-1]

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := eval(t, tt.input)

			boolVal, ok := val.(*values.BoolValue)
			if !ok {
				t.Fatalf("expected BoolValue, got %T", val)
			}

			if boolVal.Value != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, boolVal.Value)
			}
		})
	}
}

func TestEvalClosures(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedParam  string
		expectedString string
	}{
		{"simple lambda", "\\x:Int.x", "x", "<λx:Int.Int>"},
		{"lambda with bool param", "\\b:Bool.b", "b", "<λb:Bool.Bool>"},
		{"nested lambda outer", "\\x:Int.\\y:Int.x", "x", "<λx:Int.(Int -> Int)>"},
		{"higher order function", "\\f:Int->Int.f", "f", "<λf:(Int -> Int).(Int -> Int)>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := eval(t, tt.input)

			closure, ok := val.(*values.Closure)
			if !ok {
				t.Fatalf("expected Closure, got %T", val)
			}

			str := closure.String()
			if str != tt.expectedString {
				t.Errorf("expected string '%s', got '%s'", tt.expectedString, str)
			}
		})
	}
}

func eval(t *testing.T, input string) values.Value {
	t.Helper()

	expr, err := parser.Parse(input)
	if err != nil {
		t.Fatalf("parser error: %v", err)
	}

	typedExpr, err := types.Check(expr)
	if err != nil {
		t.Fatalf("type checker error: %v", err)
	}

	val, err := Eval(typedExpr)
	if err != nil {
		t.Fatalf("evaluator error: %v", err)
	}
	return val
}
