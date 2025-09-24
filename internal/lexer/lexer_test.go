package lexer_test

import (
	"testing"

	"github.com/shota3506/gostlc/internal/lexer"
)

func TestLexer(t *testing.T) {
	for _, tt := range []struct {
		name     string
		input    string
		expected []lexer.Token
	}{
		{
			name:  "Identity function",
			input: `(\x:Bool. x) true`,
			expected: []lexer.Token{
				{Kind: lexer.TokenKindLParen, Value: "("},
				{Kind: lexer.TokenKindLambda, Value: "\\"},
				{Kind: lexer.TokenKindIdent, Value: "x"},
				{Kind: lexer.TokenKindColon, Value: ":"},
				{Kind: lexer.TokenKindBoolType, Value: "Bool"},
				{Kind: lexer.TokenKindDot, Value: "."},
				{Kind: lexer.TokenKindIdent, Value: "x"},
				{Kind: lexer.TokenKindRParen, Value: ")"},
				{Kind: lexer.TokenKindTrue, Value: "true"},
				{Kind: lexer.TokenKindEOF, Value: ""},
			},
		},
		{
			name:  "Identity function with integer",
			input: `(\x:Int. x) 42`,
			expected: []lexer.Token{
				{Kind: lexer.TokenKindLParen, Value: "("},
				{Kind: lexer.TokenKindLambda, Value: "\\"},
				{Kind: lexer.TokenKindIdent, Value: "x"},
				{Kind: lexer.TokenKindColon, Value: ":"},
				{Kind: lexer.TokenKindIntType, Value: "Int"},
				{Kind: lexer.TokenKindDot, Value: "."},
				{Kind: lexer.TokenKindIdent, Value: "x"},
				{Kind: lexer.TokenKindRParen, Value: ")"},
				{Kind: lexer.TokenKindInt, Value: "42"},
				{Kind: lexer.TokenKindEOF, Value: ""},
			},
		},
		{
			name:  "Constant function",
			input: `(\x:Int. \y:Int. x) 10 20`,
			expected: []lexer.Token{
				{Kind: lexer.TokenKindLParen, Value: "("},
				{Kind: lexer.TokenKindLambda, Value: "\\"},
				{Kind: lexer.TokenKindIdent, Value: "x"},
				{Kind: lexer.TokenKindColon, Value: ":"},
				{Kind: lexer.TokenKindIntType, Value: "Int"},
				{Kind: lexer.TokenKindDot, Value: "."},
				{Kind: lexer.TokenKindLambda, Value: "\\"},
				{Kind: lexer.TokenKindIdent, Value: "y"},
				{Kind: lexer.TokenKindColon, Value: ":"},
				{Kind: lexer.TokenKindIntType, Value: "Int"},
				{Kind: lexer.TokenKindDot, Value: "."},
				{Kind: lexer.TokenKindIdent, Value: "x"},
				{Kind: lexer.TokenKindRParen, Value: ")"},
				{Kind: lexer.TokenKindInt, Value: "10"},
				{Kind: lexer.TokenKindInt, Value: "20"},
				{Kind: lexer.TokenKindEOF, Value: ""},
			},
		},
		{
			name:  "Boolean literals true",
			input: `true`,
			expected: []lexer.Token{
				{Kind: lexer.TokenKindTrue, Value: "true"},
				{Kind: lexer.TokenKindEOF, Value: ""},
			},
		},
		{
			name:  "Boolean literals false",
			input: `false`,
			expected: []lexer.Token{
				{Kind: lexer.TokenKindFalse, Value: "false"},
				{Kind: lexer.TokenKindEOF, Value: ""},
			},
		},
		{
			name:  "Simple conditional",
			input: `if true then 1 else 0`,
			expected: []lexer.Token{
				{Kind: lexer.TokenKindIf, Value: "if"},
				{Kind: lexer.TokenKindTrue, Value: "true"},
				{Kind: lexer.TokenKindThen, Value: "then"},
				{Kind: lexer.TokenKindInt, Value: "1"},
				{Kind: lexer.TokenKindElse, Value: "else"},
				{Kind: lexer.TokenKindInt, Value: "0"},
				{Kind: lexer.TokenKindEOF, Value: ""},
			},
		},
		{
			name:  "Conditional with false condition",
			input: `if false then 100 else 200`,
			expected: []lexer.Token{
				{Kind: lexer.TokenKindIf, Value: "if"},
				{Kind: lexer.TokenKindFalse, Value: "false"},
				{Kind: lexer.TokenKindThen, Value: "then"},
				{Kind: lexer.TokenKindInt, Value: "100"},
				{Kind: lexer.TokenKindElse, Value: "else"},
				{Kind: lexer.TokenKindInt, Value: "200"},
				{Kind: lexer.TokenKindEOF, Value: ""},
			},
		},
		{
			name:  "Nested function application",
			input: `(\x:Bool. \y:Bool. x) true false`,
			expected: []lexer.Token{
				{Kind: lexer.TokenKindLParen, Value: "("},
				{Kind: lexer.TokenKindLambda, Value: "\\"},
				{Kind: lexer.TokenKindIdent, Value: "x"},
				{Kind: lexer.TokenKindColon, Value: ":"},
				{Kind: lexer.TokenKindBoolType, Value: "Bool"},
				{Kind: lexer.TokenKindDot, Value: "."},
				{Kind: lexer.TokenKindLambda, Value: "\\"},
				{Kind: lexer.TokenKindIdent, Value: "y"},
				{Kind: lexer.TokenKindColon, Value: ":"},
				{Kind: lexer.TokenKindBoolType, Value: "Bool"},
				{Kind: lexer.TokenKindDot, Value: "."},
				{Kind: lexer.TokenKindIdent, Value: "x"},
				{Kind: lexer.TokenKindRParen, Value: ")"},
				{Kind: lexer.TokenKindTrue, Value: "true"},
				{Kind: lexer.TokenKindFalse, Value: "false"},
				{Kind: lexer.TokenKindEOF, Value: ""},
			},
		},
		{
			name:  "Select second argument",
			input: `(\x:Int. \y:Int. y) 5 7`,
			expected: []lexer.Token{
				{Kind: lexer.TokenKindLParen, Value: "("},
				{Kind: lexer.TokenKindLambda, Value: "\\"},
				{Kind: lexer.TokenKindIdent, Value: "x"},
				{Kind: lexer.TokenKindColon, Value: ":"},
				{Kind: lexer.TokenKindIntType, Value: "Int"},
				{Kind: lexer.TokenKindDot, Value: "."},
				{Kind: lexer.TokenKindLambda, Value: "\\"},
				{Kind: lexer.TokenKindIdent, Value: "y"},
				{Kind: lexer.TokenKindColon, Value: ":"},
				{Kind: lexer.TokenKindIntType, Value: "Int"},
				{Kind: lexer.TokenKindDot, Value: "."},
				{Kind: lexer.TokenKindIdent, Value: "y"},
				{Kind: lexer.TokenKindRParen, Value: ")"},
				{Kind: lexer.TokenKindInt, Value: "5"},
				{Kind: lexer.TokenKindInt, Value: "7"},
				{Kind: lexer.TokenKindEOF, Value: ""},
			},
		},
		{
			name:  "Apply identity to itself",
			input: `(\f:Bool->Bool. f true) (\x:Bool. x)`,
			expected: []lexer.Token{
				{Kind: lexer.TokenKindLParen, Value: "("},
				{Kind: lexer.TokenKindLambda, Value: "\\"},
				{Kind: lexer.TokenKindIdent, Value: "f"},
				{Kind: lexer.TokenKindColon, Value: ":"},
				{Kind: lexer.TokenKindBoolType, Value: "Bool"},
				{Kind: lexer.TokenKindArrow, Value: "->"},
				{Kind: lexer.TokenKindBoolType, Value: "Bool"},
				{Kind: lexer.TokenKindDot, Value: "."},
				{Kind: lexer.TokenKindIdent, Value: "f"},
				{Kind: lexer.TokenKindTrue, Value: "true"},
				{Kind: lexer.TokenKindRParen, Value: ")"},
				{Kind: lexer.TokenKindLParen, Value: "("},
				{Kind: lexer.TokenKindLambda, Value: "\\"},
				{Kind: lexer.TokenKindIdent, Value: "x"},
				{Kind: lexer.TokenKindColon, Value: ":"},
				{Kind: lexer.TokenKindBoolType, Value: "Bool"},
				{Kind: lexer.TokenKindDot, Value: "."},
				{Kind: lexer.TokenKindIdent, Value: "x"},
				{Kind: lexer.TokenKindRParen, Value: ")"},
				{Kind: lexer.TokenKindEOF, Value: ""},
			},
		},
		{
			name:  "Simple integer literal",
			input: `42`,
			expected: []lexer.Token{
				{Kind: lexer.TokenKindInt, Value: "42"},
				{Kind: lexer.TokenKindEOF, Value: ""},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			tokens := []lexer.Token{}
			for {
				tok, err := l.Next()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				tokens = append(tokens, tok)
				if tok.Kind == lexer.TokenKindEOF {
					break
				}
			}
			if len(tokens) != len(tt.expected) {
				t.Fatalf("expected %d tokens, got %d", len(tt.expected), len(tokens))
			}
			for i, expectedTok := range tt.expected {
				if tokens[i] != expectedTok {
					t.Errorf("token %d: expected %+v, got %+v", i, expectedTok, tokens[i])
				}
			}
		})
	}
}
