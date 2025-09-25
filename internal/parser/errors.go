package parser

import (
	"fmt"

	"github.com/shota3506/gostlc/internal/token"
)

// ParseError represents an error that occurred during parsing.
type ParseError struct {
	Pos     token.Position
	Message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("%d:%d: %s", e.Pos.Line, e.Pos.Column, e.Message)
}

func newParseError(tok token.Token, message string) error {
	return &ParseError{
		Pos:     tok.Pos,
		Message: message,
	}
}
