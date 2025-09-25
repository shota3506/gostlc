package lexer

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/shota3506/gostlc/internal/token"
)

type LexerError struct {
	message string
	pos     token.Position
	err     error
}

func (e *LexerError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%d:%d: %s: %v", e.pos.Line, e.pos.Column, e.message, e.err)
	}
	return fmt.Sprintf("%d:%d: %s", e.pos.Line, e.pos.Column, e.message)
}

func (e *LexerError) Pos() token.Position {
	return e.pos
}

func (e *LexerError) Unwrap() error {
	return e.err
}

type bufferedRuneReader struct {
	r   io.RuneReader
	buf *rune

	line   int
	colume int
}

func newBufferedRuneReader(r io.RuneReader) *bufferedRuneReader {
	return &bufferedRuneReader{
		r: r,

		line:   1,
		colume: 1,
	}
}

func (r *bufferedRuneReader) pos() token.Position {
	return token.Position{Line: r.line, Column: r.colume}
}

func (r *bufferedRuneReader) Peek() (rune, token.Position, error) {
	if r.buf != nil {
		ru := *r.buf
		return ru, r.pos(), nil
	}
	ru, _, err := r.r.ReadRune()
	if err != nil {
		return 0, r.pos(), err
	}
	r.buf = &ru
	return ru, r.pos(), nil
}

func (r *bufferedRuneReader) Read() (ru rune, pos token.Position, err error) {
	defer func() {
		if err == nil {
			if ru == '\n' {
				r.line++
				r.colume = 1
			} else {
				r.colume++
			}
		}
	}()

	if r.buf != nil {
		ru = *r.buf
		r.buf = nil
		return ru, r.pos(), nil
	}
	ru, _, err = r.r.ReadRune()
	return ru, r.pos(), err
}

// Lexer is a lexical analyzer for the lambda calculus with simple types.
type Lexer struct {
	reader *bufferedRuneReader
}

func New(s string) *Lexer {
	return &Lexer{
		reader: newBufferedRuneReader(strings.NewReader(s)),
	}
}

func (l *Lexer) Next() (token.Token, error) {
	l.skipWhitespace()

	ch, pos, err := l.reader.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return token.Token{
				Kind: token.TokenKindEOF,
				Pos:  l.reader.pos(),
			}, nil
		}
		return token.Token{}, &LexerError{
			message: "read character",
			pos:     pos,
			err:     err,
		}
	}

	switch ch {
	case '\\':
		return token.Token{Kind: token.TokenKindLambda, Value: string(ch), Pos: pos}, nil
	case '.':
		return token.Token{Kind: token.TokenKindDot, Value: string(ch), Pos: pos}, nil
	case ':':
		return token.Token{Kind: token.TokenKindColon, Value: string(ch), Pos: pos}, nil
	case '(':
		return token.Token{Kind: token.TokenKindLParen, Value: string(ch), Pos: pos}, nil
	case ')':
		return token.Token{Kind: token.TokenKindRParen, Value: string(ch), Pos: pos}, nil
	case '-':
		nextCh, nextPos, err := l.reader.Peek()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return token.Token{}, &LexerError{
					message: "unexpected eof after '-'",
					pos:     nextPos,
				}
			}
			return token.Token{}, &LexerError{
				message: "read character",
				pos:     nextPos,
				err:     err,
			}
		}
		if nextCh == '>' {
			_, _, _ = l.reader.Read()
			return token.Token{Kind: token.TokenKindArrow, Value: "->", Pos: pos}, nil
		}
		return token.Token{}, &LexerError{
			message: fmt.Sprintf("unexpected character after '-': %q", nextCh),
			pos:     nextPos,
		}
	}

	if isDigit(ch) {
		return token.Token{
			Kind:  token.TokenKindInt,
			Value: l.readInteger(ch),
			Pos:   pos,
		}, nil
	}

	if isAlphabetOrUnderscore(ch) {
		ident := l.readIdentifier(ch)
		switch ident {
		case "true":
			return token.Token{Kind: token.TokenKindTrue, Value: ident, Pos: pos}, nil
		case "false":
			return token.Token{Kind: token.TokenKindFalse, Value: ident, Pos: pos}, nil
		case "if":
			return token.Token{Kind: token.TokenKindIf, Value: ident, Pos: pos}, nil
		case "then":
			return token.Token{Kind: token.TokenKindThen, Value: ident, Pos: pos}, nil
		case "else":
			return token.Token{Kind: token.TokenKindElse, Value: ident, Pos: pos}, nil
		case "Bool":
			return token.Token{Kind: token.TokenKindBoolType, Value: ident, Pos: pos}, nil
		case "Int":
			return token.Token{Kind: token.TokenKindIntType, Value: ident, Pos: pos}, nil
		default:
			return token.Token{
				Kind:  token.TokenKindIdent,
				Value: ident,
				Pos:   pos,
			}, nil
		}
	}

	return token.Token{}, &LexerError{
		message: fmt.Sprintf("unexpected character: %q", ch),
		pos:     pos,
	}
}

func (l *Lexer) skipWhitespace() {
	for {
		ch, _, err := l.reader.Peek()
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Fatal(err)
		}
		if ch != ' ' && ch != '\t' && ch != '\n' && ch != '\r' {
			return
		}
		_, _, _ = l.reader.Read() // ignore error because we already peeked
	}
}

func (l *Lexer) readInteger(ch rune) string {
	var b strings.Builder
	b.WriteRune(ch)
	for {
		next, _, err := l.reader.Peek()
		if err != nil {
			break
		}
		if !isDigit(next) {
			break
		}
		b.WriteRune(next)
		_, _, _ = l.reader.Read() // ignore error because we already peeked
	}

	return b.String()
}

func (l *Lexer) readIdentifier(ch rune) string {
	var b strings.Builder
	b.WriteRune(ch)
	for {
		next, _, err := l.reader.Peek()
		if err != nil {
			break
		}
		if !isAlphanumericOrUnderscore(next) {
			break
		}
		b.WriteRune(next)
		_, _, _ = l.reader.Read() // ignore error because we already peeked
	}

	return b.String()
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func isAlphabetOrUnderscore(ch rune) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func isAlphanumericOrUnderscore(ch rune) bool {
	return isAlphabetOrUnderscore(ch) || ('0' <= ch && ch <= '9')
}
