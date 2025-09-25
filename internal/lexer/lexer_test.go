package lexer_test

import (
	"errors"
	"testing"

	"github.com/shota3506/gostlc/internal/lexer"
	"github.com/shota3506/gostlc/internal/token"
)

func TestLexer(t *testing.T) {
	for _, tt := range []struct {
		name     string
		input    string
		expected []token.Token
	}{
		{
			name:  "Identity function",
			input: `(\x:Bool. x) true`,
			expected: []token.Token{
				{Kind: token.TokenKindLParen, Value: "(", Pos: token.Position{Line: 1, Column: 1}},
				{Kind: token.TokenKindLambda, Value: "\\", Pos: token.Position{Line: 1, Column: 2}},
				{Kind: token.TokenKindIdent, Value: "x", Pos: token.Position{Line: 1, Column: 3}},
				{Kind: token.TokenKindColon, Value: ":", Pos: token.Position{Line: 1, Column: 4}},
				{Kind: token.TokenKindBoolType, Value: "Bool", Pos: token.Position{Line: 1, Column: 5}},
				{Kind: token.TokenKindDot, Value: ".", Pos: token.Position{Line: 1, Column: 9}},
				{Kind: token.TokenKindIdent, Value: "x", Pos: token.Position{Line: 1, Column: 11}},
				{Kind: token.TokenKindRParen, Value: ")", Pos: token.Position{Line: 1, Column: 12}},
				{Kind: token.TokenKindTrue, Value: "true", Pos: token.Position{Line: 1, Column: 14}},
				{Kind: token.TokenKindEOF, Value: "", Pos: token.Position{Line: 1, Column: 18}},
			},
		},
		{
			name:  "Identity function with integer",
			input: `(\x:Int. x) 42`,
			expected: []token.Token{
				{Kind: token.TokenKindLParen, Value: "(", Pos: token.Position{Line: 1, Column: 1}},
				{Kind: token.TokenKindLambda, Value: "\\", Pos: token.Position{Line: 1, Column: 2}},
				{Kind: token.TokenKindIdent, Value: "x", Pos: token.Position{Line: 1, Column: 3}},
				{Kind: token.TokenKindColon, Value: ":", Pos: token.Position{Line: 1, Column: 4}},
				{Kind: token.TokenKindIntType, Value: "Int", Pos: token.Position{Line: 1, Column: 5}},
				{Kind: token.TokenKindDot, Value: ".", Pos: token.Position{Line: 1, Column: 8}},
				{Kind: token.TokenKindIdent, Value: "x", Pos: token.Position{Line: 1, Column: 10}},
				{Kind: token.TokenKindRParen, Value: ")", Pos: token.Position{Line: 1, Column: 11}},
				{Kind: token.TokenKindInt, Value: "42", Pos: token.Position{Line: 1, Column: 13}},
				{Kind: token.TokenKindEOF, Value: "", Pos: token.Position{Line: 1, Column: 15}},
			},
		},
		{
			name:  "Constant function",
			input: `(\x:Int. \y:Int. x) 10 20`,
			expected: []token.Token{
				{Kind: token.TokenKindLParen, Value: "(", Pos: token.Position{Line: 1, Column: 1}},
				{Kind: token.TokenKindLambda, Value: "\\", Pos: token.Position{Line: 1, Column: 2}},
				{Kind: token.TokenKindIdent, Value: "x", Pos: token.Position{Line: 1, Column: 3}},
				{Kind: token.TokenKindColon, Value: ":", Pos: token.Position{Line: 1, Column: 4}},
				{Kind: token.TokenKindIntType, Value: "Int", Pos: token.Position{Line: 1, Column: 5}},
				{Kind: token.TokenKindDot, Value: ".", Pos: token.Position{Line: 1, Column: 8}},
				{Kind: token.TokenKindLambda, Value: "\\", Pos: token.Position{Line: 1, Column: 10}},
				{Kind: token.TokenKindIdent, Value: "y", Pos: token.Position{Line: 1, Column: 11}},
				{Kind: token.TokenKindColon, Value: ":", Pos: token.Position{Line: 1, Column: 12}},
				{Kind: token.TokenKindIntType, Value: "Int", Pos: token.Position{Line: 1, Column: 13}},
				{Kind: token.TokenKindDot, Value: ".", Pos: token.Position{Line: 1, Column: 16}},
				{Kind: token.TokenKindIdent, Value: "x", Pos: token.Position{Line: 1, Column: 18}},
				{Kind: token.TokenKindRParen, Value: ")", Pos: token.Position{Line: 1, Column: 19}},
				{Kind: token.TokenKindInt, Value: "10", Pos: token.Position{Line: 1, Column: 21}},
				{Kind: token.TokenKindInt, Value: "20", Pos: token.Position{Line: 1, Column: 24}},
				{Kind: token.TokenKindEOF, Value: "", Pos: token.Position{Line: 1, Column: 26}},
			},
		},
		{
			name:  "Boolean literals true",
			input: `true`,
			expected: []token.Token{
				{Kind: token.TokenKindTrue, Value: "true", Pos: token.Position{Line: 1, Column: 1}},
				{Kind: token.TokenKindEOF, Value: "", Pos: token.Position{Line: 1, Column: 5}},
			},
		},
		{
			name:  "Boolean literals false",
			input: `false`,
			expected: []token.Token{
				{Kind: token.TokenKindFalse, Value: "false", Pos: token.Position{Line: 1, Column: 1}},
				{Kind: token.TokenKindEOF, Value: "", Pos: token.Position{Line: 1, Column: 6}},
			},
		},
		{
			name:  "Simple conditional",
			input: `if true then 1 else 0`,
			expected: []token.Token{
				{Kind: token.TokenKindIf, Value: "if", Pos: token.Position{Line: 1, Column: 1}},
				{Kind: token.TokenKindTrue, Value: "true", Pos: token.Position{Line: 1, Column: 4}},
				{Kind: token.TokenKindThen, Value: "then", Pos: token.Position{Line: 1, Column: 9}},
				{Kind: token.TokenKindInt, Value: "1", Pos: token.Position{Line: 1, Column: 14}},
				{Kind: token.TokenKindElse, Value: "else", Pos: token.Position{Line: 1, Column: 16}},
				{Kind: token.TokenKindInt, Value: "0", Pos: token.Position{Line: 1, Column: 21}},
				{Kind: token.TokenKindEOF, Value: "", Pos: token.Position{Line: 1, Column: 22}},
			},
		},
		{
			name:  "Conditional with false condition",
			input: `if false then 100 else 200`,
			expected: []token.Token{
				{Kind: token.TokenKindIf, Value: "if", Pos: token.Position{Line: 1, Column: 1}},
				{Kind: token.TokenKindFalse, Value: "false", Pos: token.Position{Line: 1, Column: 4}},
				{Kind: token.TokenKindThen, Value: "then", Pos: token.Position{Line: 1, Column: 10}},
				{Kind: token.TokenKindInt, Value: "100", Pos: token.Position{Line: 1, Column: 15}},
				{Kind: token.TokenKindElse, Value: "else", Pos: token.Position{Line: 1, Column: 19}},
				{Kind: token.TokenKindInt, Value: "200", Pos: token.Position{Line: 1, Column: 24}},
				{Kind: token.TokenKindEOF, Value: "", Pos: token.Position{Line: 1, Column: 27}},
			},
		},
		{
			name:  "Nested function application",
			input: `(\x:Bool. \y:Bool. x) true false`,
			expected: []token.Token{
				{Kind: token.TokenKindLParen, Value: "(", Pos: token.Position{Line: 1, Column: 1}},
				{Kind: token.TokenKindLambda, Value: "\\", Pos: token.Position{Line: 1, Column: 2}},
				{Kind: token.TokenKindIdent, Value: "x", Pos: token.Position{Line: 1, Column: 3}},
				{Kind: token.TokenKindColon, Value: ":", Pos: token.Position{Line: 1, Column: 4}},
				{Kind: token.TokenKindBoolType, Value: "Bool", Pos: token.Position{Line: 1, Column: 5}},
				{Kind: token.TokenKindDot, Value: ".", Pos: token.Position{Line: 1, Column: 9}},
				{Kind: token.TokenKindLambda, Value: "\\", Pos: token.Position{Line: 1, Column: 11}},
				{Kind: token.TokenKindIdent, Value: "y", Pos: token.Position{Line: 1, Column: 12}},
				{Kind: token.TokenKindColon, Value: ":", Pos: token.Position{Line: 1, Column: 13}},
				{Kind: token.TokenKindBoolType, Value: "Bool", Pos: token.Position{Line: 1, Column: 14}},
				{Kind: token.TokenKindDot, Value: ".", Pos: token.Position{Line: 1, Column: 18}},
				{Kind: token.TokenKindIdent, Value: "x", Pos: token.Position{Line: 1, Column: 20}},
				{Kind: token.TokenKindRParen, Value: ")", Pos: token.Position{Line: 1, Column: 21}},
				{Kind: token.TokenKindTrue, Value: "true", Pos: token.Position{Line: 1, Column: 23}},
				{Kind: token.TokenKindFalse, Value: "false", Pos: token.Position{Line: 1, Column: 28}},
				{Kind: token.TokenKindEOF, Value: "", Pos: token.Position{Line: 1, Column: 33}},
			},
		},
		{
			name:  "Select second argument",
			input: `(\x:Int. \y:Int. y) 5 7`,
			expected: []token.Token{
				{Kind: token.TokenKindLParen, Value: "(", Pos: token.Position{Line: 1, Column: 1}},
				{Kind: token.TokenKindLambda, Value: "\\", Pos: token.Position{Line: 1, Column: 2}},
				{Kind: token.TokenKindIdent, Value: "x", Pos: token.Position{Line: 1, Column: 3}},
				{Kind: token.TokenKindColon, Value: ":", Pos: token.Position{Line: 1, Column: 4}},
				{Kind: token.TokenKindIntType, Value: "Int", Pos: token.Position{Line: 1, Column: 5}},
				{Kind: token.TokenKindDot, Value: ".", Pos: token.Position{Line: 1, Column: 8}},
				{Kind: token.TokenKindLambda, Value: "\\", Pos: token.Position{Line: 1, Column: 10}},
				{Kind: token.TokenKindIdent, Value: "y", Pos: token.Position{Line: 1, Column: 11}},
				{Kind: token.TokenKindColon, Value: ":", Pos: token.Position{Line: 1, Column: 12}},
				{Kind: token.TokenKindIntType, Value: "Int", Pos: token.Position{Line: 1, Column: 13}},
				{Kind: token.TokenKindDot, Value: ".", Pos: token.Position{Line: 1, Column: 16}},
				{Kind: token.TokenKindIdent, Value: "y", Pos: token.Position{Line: 1, Column: 18}},
				{Kind: token.TokenKindRParen, Value: ")", Pos: token.Position{Line: 1, Column: 19}},
				{Kind: token.TokenKindInt, Value: "5", Pos: token.Position{Line: 1, Column: 21}},
				{Kind: token.TokenKindInt, Value: "7", Pos: token.Position{Line: 1, Column: 23}},
				{Kind: token.TokenKindEOF, Value: "", Pos: token.Position{Line: 1, Column: 24}},
			},
		},
		{
			name:  "Apply identity to itself",
			input: `(\f:Bool->Bool. f true) (\x:Bool. x)`,
			expected: []token.Token{
				{Kind: token.TokenKindLParen, Value: "(", Pos: token.Position{Line: 1, Column: 1}},
				{Kind: token.TokenKindLambda, Value: "\\", Pos: token.Position{Line: 1, Column: 2}},
				{Kind: token.TokenKindIdent, Value: "f", Pos: token.Position{Line: 1, Column: 3}},
				{Kind: token.TokenKindColon, Value: ":", Pos: token.Position{Line: 1, Column: 4}},
				{Kind: token.TokenKindBoolType, Value: "Bool", Pos: token.Position{Line: 1, Column: 5}},
				{Kind: token.TokenKindArrow, Value: "->", Pos: token.Position{Line: 1, Column: 9}},
				{Kind: token.TokenKindBoolType, Value: "Bool", Pos: token.Position{Line: 1, Column: 11}},
				{Kind: token.TokenKindDot, Value: ".", Pos: token.Position{Line: 1, Column: 15}},
				{Kind: token.TokenKindIdent, Value: "f", Pos: token.Position{Line: 1, Column: 17}},
				{Kind: token.TokenKindTrue, Value: "true", Pos: token.Position{Line: 1, Column: 19}},
				{Kind: token.TokenKindRParen, Value: ")", Pos: token.Position{Line: 1, Column: 23}},
				{Kind: token.TokenKindLParen, Value: "(", Pos: token.Position{Line: 1, Column: 25}},
				{Kind: token.TokenKindLambda, Value: "\\", Pos: token.Position{Line: 1, Column: 26}},
				{Kind: token.TokenKindIdent, Value: "x", Pos: token.Position{Line: 1, Column: 27}},
				{Kind: token.TokenKindColon, Value: ":", Pos: token.Position{Line: 1, Column: 28}},
				{Kind: token.TokenKindBoolType, Value: "Bool", Pos: token.Position{Line: 1, Column: 29}},
				{Kind: token.TokenKindDot, Value: ".", Pos: token.Position{Line: 1, Column: 33}},
				{Kind: token.TokenKindIdent, Value: "x", Pos: token.Position{Line: 1, Column: 35}},
				{Kind: token.TokenKindRParen, Value: ")", Pos: token.Position{Line: 1, Column: 36}},
				{Kind: token.TokenKindEOF, Value: "", Pos: token.Position{Line: 1, Column: 37}},
			},
		},
		{
			name:  "Simple integer literal",
			input: `42`,
			expected: []token.Token{
				{Kind: token.TokenKindInt, Value: "42", Pos: token.Position{Line: 1, Column: 1}},
				{Kind: token.TokenKindEOF, Value: "", Pos: token.Position{Line: 1, Column: 3}},
			},
		},
		{
			name:  "Negative integer literal",
			input: `-123`,
			expected: []token.Token{
				{Kind: token.TokenKindInt, Value: "-123", Pos: token.Position{Line: 1, Column: 1}},
				{Kind: token.TokenKindEOF, Value: "", Pos: token.Position{Line: 1, Column: 5}},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			tokens := []token.Token{}
			for {
				tok, err := l.Next()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				tokens = append(tokens, tok)
				if tok.Kind == token.TokenKindEOF {
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
		expectedPos   token.Position
	}{
		{
			name:          "Unexpected character @",
			input:         `@`,
			expectedError: `1:1: unexpected character: '@'`,
			expectedPos:   token.Position{Line: 1, Column: 1},
		},
		{
			name:          "Unexpected character # in middle",
			input:         `true # false`,
			expectedError: `1:6: unexpected character: '#'`,
			expectedPos:   token.Position{Line: 1, Column: 6},
		},
		{
			name:          "Single dash without arrow",
			input:         `- 5`,
			expectedError: `1:2: unexpected character after '-': ' '`,
			expectedPos:   token.Position{Line: 1, Column: 2},
		},
		{
			name:          "Dash at end of input",
			input:         `(\x:Int. x) -`,
			expectedError: `1:14: unexpected eof after '-'`,
			expectedPos:   token.Position{Line: 1, Column: 14},
		},
		{
			name:          "Unexpected character after dash",
			input:         `-a`,
			expectedError: `1:2: unexpected character after '-': 'a'`,
			expectedPos:   token.Position{Line: 1, Column: 2},
		},
		{
			name:          "Multiple unexpected characters",
			input:         `true & false`,
			expectedError: `1:6: unexpected character: '&'`,
			expectedPos:   token.Position{Line: 1, Column: 6},
		},
		{
			name:          "Unexpected character on new line",
			input:         "true\n$",
			expectedError: `2:1: unexpected character: '$'`,
			expectedPos:   token.Position{Line: 2, Column: 1},
		},
		{
			name:          "Unexpected character after whitespace",
			input:         `   !`,
			expectedError: `1:4: unexpected character: '!'`,
			expectedPos:   token.Position{Line: 1, Column: 4},
		},
		{
			name:          "Unicode character",
			input:         `λ`,
			expectedError: `1:1: unexpected character: 'λ'`,
			expectedPos:   token.Position{Line: 1, Column: 1},
		},
		{
			name:          "Unexpected character in expression",
			input:         `(\x:Bool. x) % true`,
			expectedError: `1:14: unexpected character: '%'`,
			expectedPos:   token.Position{Line: 1, Column: 14},
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
				if tok.Kind == token.TokenKindEOF {
					t.Fatalf("expected error but reached EOF successfully")
				}
			}
		})
	}
}
