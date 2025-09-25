package lexer_test

import (
	"errors"
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
				{Kind: lexer.TokenKindLParen, Value: "(", Pos: lexer.Position{Line: 1, Column: 1}},
				{Kind: lexer.TokenKindLambda, Value: "\\", Pos: lexer.Position{Line: 1, Column: 2}},
				{Kind: lexer.TokenKindIdent, Value: "x", Pos: lexer.Position{Line: 1, Column: 3}},
				{Kind: lexer.TokenKindColon, Value: ":", Pos: lexer.Position{Line: 1, Column: 4}},
				{Kind: lexer.TokenKindBoolType, Value: "Bool", Pos: lexer.Position{Line: 1, Column: 5}},
				{Kind: lexer.TokenKindDot, Value: ".", Pos: lexer.Position{Line: 1, Column: 9}},
				{Kind: lexer.TokenKindIdent, Value: "x", Pos: lexer.Position{Line: 1, Column: 11}},
				{Kind: lexer.TokenKindRParen, Value: ")", Pos: lexer.Position{Line: 1, Column: 12}},
				{Kind: lexer.TokenKindTrue, Value: "true", Pos: lexer.Position{Line: 1, Column: 14}},
				{Kind: lexer.TokenKindEOF, Value: "", Pos: lexer.Position{Line: 1, Column: 18}},
			},
		},
		{
			name:  "Identity function with integer",
			input: `(\x:Int. x) 42`,
			expected: []lexer.Token{
				{Kind: lexer.TokenKindLParen, Value: "(", Pos: lexer.Position{Line: 1, Column: 1}},
				{Kind: lexer.TokenKindLambda, Value: "\\", Pos: lexer.Position{Line: 1, Column: 2}},
				{Kind: lexer.TokenKindIdent, Value: "x", Pos: lexer.Position{Line: 1, Column: 3}},
				{Kind: lexer.TokenKindColon, Value: ":", Pos: lexer.Position{Line: 1, Column: 4}},
				{Kind: lexer.TokenKindIntType, Value: "Int", Pos: lexer.Position{Line: 1, Column: 5}},
				{Kind: lexer.TokenKindDot, Value: ".", Pos: lexer.Position{Line: 1, Column: 8}},
				{Kind: lexer.TokenKindIdent, Value: "x", Pos: lexer.Position{Line: 1, Column: 10}},
				{Kind: lexer.TokenKindRParen, Value: ")", Pos: lexer.Position{Line: 1, Column: 11}},
				{Kind: lexer.TokenKindInt, Value: "42", Pos: lexer.Position{Line: 1, Column: 13}},
				{Kind: lexer.TokenKindEOF, Value: "", Pos: lexer.Position{Line: 1, Column: 15}},
			},
		},
		{
			name:  "Constant function",
			input: `(\x:Int. \y:Int. x) 10 20`,
			expected: []lexer.Token{
				{Kind: lexer.TokenKindLParen, Value: "(", Pos: lexer.Position{Line: 1, Column: 1}},
				{Kind: lexer.TokenKindLambda, Value: "\\", Pos: lexer.Position{Line: 1, Column: 2}},
				{Kind: lexer.TokenKindIdent, Value: "x", Pos: lexer.Position{Line: 1, Column: 3}},
				{Kind: lexer.TokenKindColon, Value: ":", Pos: lexer.Position{Line: 1, Column: 4}},
				{Kind: lexer.TokenKindIntType, Value: "Int", Pos: lexer.Position{Line: 1, Column: 5}},
				{Kind: lexer.TokenKindDot, Value: ".", Pos: lexer.Position{Line: 1, Column: 8}},
				{Kind: lexer.TokenKindLambda, Value: "\\", Pos: lexer.Position{Line: 1, Column: 10}},
				{Kind: lexer.TokenKindIdent, Value: "y", Pos: lexer.Position{Line: 1, Column: 11}},
				{Kind: lexer.TokenKindColon, Value: ":", Pos: lexer.Position{Line: 1, Column: 12}},
				{Kind: lexer.TokenKindIntType, Value: "Int", Pos: lexer.Position{Line: 1, Column: 13}},
				{Kind: lexer.TokenKindDot, Value: ".", Pos: lexer.Position{Line: 1, Column: 16}},
				{Kind: lexer.TokenKindIdent, Value: "x", Pos: lexer.Position{Line: 1, Column: 18}},
				{Kind: lexer.TokenKindRParen, Value: ")", Pos: lexer.Position{Line: 1, Column: 19}},
				{Kind: lexer.TokenKindInt, Value: "10", Pos: lexer.Position{Line: 1, Column: 21}},
				{Kind: lexer.TokenKindInt, Value: "20", Pos: lexer.Position{Line: 1, Column: 24}},
				{Kind: lexer.TokenKindEOF, Value: "", Pos: lexer.Position{Line: 1, Column: 26}},
			},
		},
		{
			name:  "Boolean literals true",
			input: `true`,
			expected: []lexer.Token{
				{Kind: lexer.TokenKindTrue, Value: "true", Pos: lexer.Position{Line: 1, Column: 1}},
				{Kind: lexer.TokenKindEOF, Value: "", Pos: lexer.Position{Line: 1, Column: 5}},
			},
		},
		{
			name:  "Boolean literals false",
			input: `false`,
			expected: []lexer.Token{
				{Kind: lexer.TokenKindFalse, Value: "false", Pos: lexer.Position{Line: 1, Column: 1}},
				{Kind: lexer.TokenKindEOF, Value: "", Pos: lexer.Position{Line: 1, Column: 6}},
			},
		},
		{
			name:  "Simple conditional",
			input: `if true then 1 else 0`,
			expected: []lexer.Token{
				{Kind: lexer.TokenKindIf, Value: "if", Pos: lexer.Position{Line: 1, Column: 1}},
				{Kind: lexer.TokenKindTrue, Value: "true", Pos: lexer.Position{Line: 1, Column: 4}},
				{Kind: lexer.TokenKindThen, Value: "then", Pos: lexer.Position{Line: 1, Column: 9}},
				{Kind: lexer.TokenKindInt, Value: "1", Pos: lexer.Position{Line: 1, Column: 14}},
				{Kind: lexer.TokenKindElse, Value: "else", Pos: lexer.Position{Line: 1, Column: 16}},
				{Kind: lexer.TokenKindInt, Value: "0", Pos: lexer.Position{Line: 1, Column: 21}},
				{Kind: lexer.TokenKindEOF, Value: "", Pos: lexer.Position{Line: 1, Column: 22}},
			},
		},
		{
			name:  "Conditional with false condition",
			input: `if false then 100 else 200`,
			expected: []lexer.Token{
				{Kind: lexer.TokenKindIf, Value: "if", Pos: lexer.Position{Line: 1, Column: 1}},
				{Kind: lexer.TokenKindFalse, Value: "false", Pos: lexer.Position{Line: 1, Column: 4}},
				{Kind: lexer.TokenKindThen, Value: "then", Pos: lexer.Position{Line: 1, Column: 10}},
				{Kind: lexer.TokenKindInt, Value: "100", Pos: lexer.Position{Line: 1, Column: 15}},
				{Kind: lexer.TokenKindElse, Value: "else", Pos: lexer.Position{Line: 1, Column: 19}},
				{Kind: lexer.TokenKindInt, Value: "200", Pos: lexer.Position{Line: 1, Column: 24}},
				{Kind: lexer.TokenKindEOF, Value: "", Pos: lexer.Position{Line: 1, Column: 27}},
			},
		},
		{
			name:  "Nested function application",
			input: `(\x:Bool. \y:Bool. x) true false`,
			expected: []lexer.Token{
				{Kind: lexer.TokenKindLParen, Value: "(", Pos: lexer.Position{Line: 1, Column: 1}},
				{Kind: lexer.TokenKindLambda, Value: "\\", Pos: lexer.Position{Line: 1, Column: 2}},
				{Kind: lexer.TokenKindIdent, Value: "x", Pos: lexer.Position{Line: 1, Column: 3}},
				{Kind: lexer.TokenKindColon, Value: ":", Pos: lexer.Position{Line: 1, Column: 4}},
				{Kind: lexer.TokenKindBoolType, Value: "Bool", Pos: lexer.Position{Line: 1, Column: 5}},
				{Kind: lexer.TokenKindDot, Value: ".", Pos: lexer.Position{Line: 1, Column: 9}},
				{Kind: lexer.TokenKindLambda, Value: "\\", Pos: lexer.Position{Line: 1, Column: 11}},
				{Kind: lexer.TokenKindIdent, Value: "y", Pos: lexer.Position{Line: 1, Column: 12}},
				{Kind: lexer.TokenKindColon, Value: ":", Pos: lexer.Position{Line: 1, Column: 13}},
				{Kind: lexer.TokenKindBoolType, Value: "Bool", Pos: lexer.Position{Line: 1, Column: 14}},
				{Kind: lexer.TokenKindDot, Value: ".", Pos: lexer.Position{Line: 1, Column: 18}},
				{Kind: lexer.TokenKindIdent, Value: "x", Pos: lexer.Position{Line: 1, Column: 20}},
				{Kind: lexer.TokenKindRParen, Value: ")", Pos: lexer.Position{Line: 1, Column: 21}},
				{Kind: lexer.TokenKindTrue, Value: "true", Pos: lexer.Position{Line: 1, Column: 23}},
				{Kind: lexer.TokenKindFalse, Value: "false", Pos: lexer.Position{Line: 1, Column: 28}},
				{Kind: lexer.TokenKindEOF, Value: "", Pos: lexer.Position{Line: 1, Column: 33}},
			},
		},
		{
			name:  "Select second argument",
			input: `(\x:Int. \y:Int. y) 5 7`,
			expected: []lexer.Token{
				{Kind: lexer.TokenKindLParen, Value: "(", Pos: lexer.Position{Line: 1, Column: 1}},
				{Kind: lexer.TokenKindLambda, Value: "\\", Pos: lexer.Position{Line: 1, Column: 2}},
				{Kind: lexer.TokenKindIdent, Value: "x", Pos: lexer.Position{Line: 1, Column: 3}},
				{Kind: lexer.TokenKindColon, Value: ":", Pos: lexer.Position{Line: 1, Column: 4}},
				{Kind: lexer.TokenKindIntType, Value: "Int", Pos: lexer.Position{Line: 1, Column: 5}},
				{Kind: lexer.TokenKindDot, Value: ".", Pos: lexer.Position{Line: 1, Column: 8}},
				{Kind: lexer.TokenKindLambda, Value: "\\", Pos: lexer.Position{Line: 1, Column: 10}},
				{Kind: lexer.TokenKindIdent, Value: "y", Pos: lexer.Position{Line: 1, Column: 11}},
				{Kind: lexer.TokenKindColon, Value: ":", Pos: lexer.Position{Line: 1, Column: 12}},
				{Kind: lexer.TokenKindIntType, Value: "Int", Pos: lexer.Position{Line: 1, Column: 13}},
				{Kind: lexer.TokenKindDot, Value: ".", Pos: lexer.Position{Line: 1, Column: 16}},
				{Kind: lexer.TokenKindIdent, Value: "y", Pos: lexer.Position{Line: 1, Column: 18}},
				{Kind: lexer.TokenKindRParen, Value: ")", Pos: lexer.Position{Line: 1, Column: 19}},
				{Kind: lexer.TokenKindInt, Value: "5", Pos: lexer.Position{Line: 1, Column: 21}},
				{Kind: lexer.TokenKindInt, Value: "7", Pos: lexer.Position{Line: 1, Column: 23}},
				{Kind: lexer.TokenKindEOF, Value: "", Pos: lexer.Position{Line: 1, Column: 24}},
			},
		},
		{
			name:  "Apply identity to itself",
			input: `(\f:Bool->Bool. f true) (\x:Bool. x)`,
			expected: []lexer.Token{
				{Kind: lexer.TokenKindLParen, Value: "(", Pos: lexer.Position{Line: 1, Column: 1}},
				{Kind: lexer.TokenKindLambda, Value: "\\", Pos: lexer.Position{Line: 1, Column: 2}},
				{Kind: lexer.TokenKindIdent, Value: "f", Pos: lexer.Position{Line: 1, Column: 3}},
				{Kind: lexer.TokenKindColon, Value: ":", Pos: lexer.Position{Line: 1, Column: 4}},
				{Kind: lexer.TokenKindBoolType, Value: "Bool", Pos: lexer.Position{Line: 1, Column: 5}},
				{Kind: lexer.TokenKindArrow, Value: "->", Pos: lexer.Position{Line: 1, Column: 9}},
				{Kind: lexer.TokenKindBoolType, Value: "Bool", Pos: lexer.Position{Line: 1, Column: 11}},
				{Kind: lexer.TokenKindDot, Value: ".", Pos: lexer.Position{Line: 1, Column: 15}},
				{Kind: lexer.TokenKindIdent, Value: "f", Pos: lexer.Position{Line: 1, Column: 17}},
				{Kind: lexer.TokenKindTrue, Value: "true", Pos: lexer.Position{Line: 1, Column: 19}},
				{Kind: lexer.TokenKindRParen, Value: ")", Pos: lexer.Position{Line: 1, Column: 23}},
				{Kind: lexer.TokenKindLParen, Value: "(", Pos: lexer.Position{Line: 1, Column: 25}},
				{Kind: lexer.TokenKindLambda, Value: "\\", Pos: lexer.Position{Line: 1, Column: 26}},
				{Kind: lexer.TokenKindIdent, Value: "x", Pos: lexer.Position{Line: 1, Column: 27}},
				{Kind: lexer.TokenKindColon, Value: ":", Pos: lexer.Position{Line: 1, Column: 28}},
				{Kind: lexer.TokenKindBoolType, Value: "Bool", Pos: lexer.Position{Line: 1, Column: 29}},
				{Kind: lexer.TokenKindDot, Value: ".", Pos: lexer.Position{Line: 1, Column: 33}},
				{Kind: lexer.TokenKindIdent, Value: "x", Pos: lexer.Position{Line: 1, Column: 35}},
				{Kind: lexer.TokenKindRParen, Value: ")", Pos: lexer.Position{Line: 1, Column: 36}},
				{Kind: lexer.TokenKindEOF, Value: "", Pos: lexer.Position{Line: 1, Column: 37}},
			},
		},
		{
			name:  "Simple integer literal",
			input: `42`,
			expected: []lexer.Token{
				{Kind: lexer.TokenKindInt, Value: "42", Pos: lexer.Position{Line: 1, Column: 1}},
				{Kind: lexer.TokenKindEOF, Value: "", Pos: lexer.Position{Line: 1, Column: 3}},
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

func TestLexerErrors(t *testing.T) {
	for _, tt := range []struct {
		name          string
		input         string
		expectedError string
		expectedPos   lexer.Position
	}{
		{
			name:          "Unexpected character @",
			input:         `@`,
			expectedError: `1:1: unexpected character: '@'`,
			expectedPos:   lexer.Position{Line: 1, Column: 1},
		},
		{
			name:          "Unexpected character # in middle",
			input:         `true # false`,
			expectedError: `1:6: unexpected character: '#'`,
			expectedPos:   lexer.Position{Line: 1, Column: 6},
		},
		{
			name:          "Single dash without arrow",
			input:         `- 5`,
			expectedError: `1:2: unexpected character after '-': ' '`,
			expectedPos:   lexer.Position{Line: 1, Column: 2},
		},
		{
			name:          "Dash at end of input",
			input:         `(\x:Int. x) -`,
			expectedError: `1:14: unexpected eof after '-'`,
			expectedPos:   lexer.Position{Line: 1, Column: 14},
		},
		{
			name:          "Unexpected character after dash",
			input:         `-a`,
			expectedError: `1:2: unexpected character after '-': 'a'`,
			expectedPos:   lexer.Position{Line: 1, Column: 2},
		},
		{
			name:          "Multiple unexpected characters",
			input:         `true & false`,
			expectedError: `1:6: unexpected character: '&'`,
			expectedPos:   lexer.Position{Line: 1, Column: 6},
		},
		{
			name:          "Unexpected character on new line",
			input:         "true\n$",
			expectedError: `2:1: unexpected character: '$'`,
			expectedPos:   lexer.Position{Line: 2, Column: 1},
		},
		{
			name:          "Unexpected character after whitespace",
			input:         `   !`,
			expectedError: `1:4: unexpected character: '!'`,
			expectedPos:   lexer.Position{Line: 1, Column: 4},
		},
		{
			name:          "Unicode character",
			input:         `λ`,
			expectedError: `1:1: unexpected character: 'λ'`,
			expectedPos:   lexer.Position{Line: 1, Column: 1},
		},
		{
			name:          "Unexpected character in expression",
			input:         `(\x:Bool. x) % true`,
			expectedError: `1:14: unexpected character: '%'`,
			expectedPos:   lexer.Position{Line: 1, Column: 14},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)

			for {
				tok, err := l.Next()
				if err != nil {
					var lexErr *lexer.LexerError
					if !errors.As(err, &lexErr) {
						t.Fatalf("expected *lexer.LexerError, got %T", err)
					}

					if err.Error() != tt.expectedError {
						t.Errorf("error message: expected %q, got %q", tt.expectedError, err.Error())
					}
					return
				}
				if tok.Kind == lexer.TokenKindEOF {
					t.Fatalf("expected error but reached EOF successfully")
				}
			}
		})
	}
}
