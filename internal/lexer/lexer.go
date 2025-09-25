package lexer

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

type LexerError struct {
	message string
	pos     Position
	err     error
}

func (e *LexerError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%d:%d: %s: %v", e.pos.Line, e.pos.Column, e.message, e.err)
	}
	return fmt.Sprintf("%d:%d: %s", e.pos.Line, e.pos.Column, e.message)
}

func (e *LexerError) Pos() Position {
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

func (r *bufferedRuneReader) pos() Position {
	return Position{Line: r.line, Column: r.colume}
}

func (r *bufferedRuneReader) Peek() (rune, Position, error) {
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

func (r *bufferedRuneReader) Read() (ru rune, pos Position, err error) {
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

func (l *Lexer) Next() (Token, error) {
	l.skipWhitespace()

	ch, pos, err := l.reader.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return Token{
				Kind: TokenKindEOF,
				Pos:  l.reader.pos(),
			}, nil
		}
		return Token{}, &LexerError{
			message: "read character",
			pos:     pos,
			err:     err,
		}
	}

	switch ch {
	case '\\':
		return Token{Kind: TokenKindLambda, Value: string(ch), Pos: pos}, nil
	case '.':
		return Token{Kind: TokenKindDot, Value: string(ch), Pos: pos}, nil
	case ':':
		return Token{Kind: TokenKindColon, Value: string(ch), Pos: pos}, nil
	case '(':
		return Token{Kind: TokenKindLParen, Value: string(ch), Pos: pos}, nil
	case ')':
		return Token{Kind: TokenKindRParen, Value: string(ch), Pos: pos}, nil
	case '-':
		nextCh, nextPos, err := l.reader.Peek()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return Token{}, &LexerError{
					message: "unexpected eof after '-'",
					pos:     nextPos,
				}
			}
			return Token{}, &LexerError{
				message: "read character",
				pos:     nextPos,
				err:     err,
			}
		}
		if nextCh == '>' {
			_, _, _ = l.reader.Read()
			return Token{Kind: TokenKindArrow, Value: "->", Pos: pos}, nil
		}
		return Token{}, &LexerError{
			message: fmt.Sprintf("unexpected character after '-': %q", nextCh),
			pos:     nextPos,
		}
	}

	if isDigit(ch) {
		return Token{
			Kind:  TokenKindInt,
			Value: l.readInteger(ch),
			Pos:   pos,
		}, nil
	}

	if isAlphabetOrUnderscore(ch) {
		ident := l.readIdentifier(ch)
		switch ident {
		case "true":
			return Token{Kind: TokenKindTrue, Value: ident, Pos: pos}, nil
		case "false":
			return Token{Kind: TokenKindFalse, Value: ident, Pos: pos}, nil
		case "if":
			return Token{Kind: TokenKindIf, Value: ident, Pos: pos}, nil
		case "then":
			return Token{Kind: TokenKindThen, Value: ident, Pos: pos}, nil
		case "else":
			return Token{Kind: TokenKindElse, Value: ident, Pos: pos}, nil
		case "Bool":
			return Token{Kind: TokenKindBoolType, Value: ident, Pos: pos}, nil
		case "Int":
			return Token{Kind: TokenKindIntType, Value: ident, Pos: pos}, nil
		default:
			return Token{
				Kind:  TokenKindIdent,
				Value: ident,
				Pos:   pos,
			}, nil
		}
	}

	return Token{}, &LexerError{
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
