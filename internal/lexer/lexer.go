package lexer

import (
	"fmt"
	"io"
	"log"
	"strings"
	"unicode"
	"unicode/utf8"
)

type bufferedRuneReader struct {
	r   io.RuneReader
	buf *rune
}

func newBufferedRuneReader(r io.RuneReader) *bufferedRuneReader {
	return &bufferedRuneReader{
		r: r,
	}
}

func (r *bufferedRuneReader) Peek() (rune, int, error) {
	if r.buf != nil {
		ru := *r.buf
		return ru, utf8.RuneLen(ru), nil
	}
	ru, size, err := r.r.ReadRune()
	if err != nil {
		return 0, 0, err
	}
	r.buf = &ru
	return ru, size, nil
}

func (r *bufferedRuneReader) ReadRune() (rune, int, error) {
	if r.buf != nil {
		ru := *r.buf
		r.buf = nil
		return ru, utf8.RuneLen(ru), nil
	}
	return r.r.ReadRune()
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

	ch, _, err := l.reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			return Token{
				Kind: TokenKindEOF,
			}, nil
		}
		return Token{}, err
	}

	switch ch {
	case '\\':
		return Token{Kind: TokenKindLambda, Value: string(ch)}, nil
	case '.':
		return Token{Kind: TokenKindDot, Value: string(ch)}, nil
	case ':':
		return Token{Kind: TokenKindColon, Value: string(ch)}, nil
	case '(':
		return Token{Kind: TokenKindLParen, Value: string(ch)}, nil
	case ')':
		return Token{Kind: TokenKindRParen, Value: string(ch)}, nil
	case '-':
		nextCh, _, err := l.reader.Peek()
		if err != nil {
			return Token{}, err
		}
		if nextCh == '>' {
			_, _, _ = l.reader.ReadRune()
			return Token{Kind: TokenKindArrow, Value: "->"}, nil
		}
		return Token{}, fmt.Errorf("unexpected character after '-': %q", nextCh)
	}

	if unicode.IsDigit(ch) {
		return Token{
			Kind:  TokenKindInt,
			Value: l.readInteger(ch),
		}, nil
	}

	if isAlphabetOrUnderscore(ch) {
		ident := l.readIdentifier(ch)
		switch ident {
		case "true":
			return Token{Kind: TokenKindTrue, Value: ident}, nil
		case "false":
			return Token{Kind: TokenKindFalse, Value: ident}, nil
		case "if":
			return Token{Kind: TokenKindIf, Value: ident}, nil
		case "then":
			return Token{Kind: TokenKindThen, Value: ident}, nil
		case "else":
			return Token{Kind: TokenKindElse, Value: ident}, nil
		case "Bool":
			return Token{Kind: TokenKindBoolType, Value: ident}, nil
		case "Int":
			return Token{Kind: TokenKindIntType, Value: ident}, nil
		default:
			return Token{
				Kind:  TokenKindIdent,
				Value: ident,
			}, nil
		}
	}

	return Token{}, fmt.Errorf("unexpected character: %q", ch)
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
		_, _, _ = l.reader.ReadRune() // ignore error because we already peeked
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
		if !unicode.IsDigit(next) {
			break
		}
		b.WriteRune(next)
		_, _, _ = l.reader.ReadRune() // ignore error because we already peeked
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
		_, _, _ = l.reader.ReadRune() // ignore error because we already peeked
	}

	return b.String()
}

func isAlphabetOrUnderscore(ch rune) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func isAlphanumericOrUnderscore(ch rune) bool {
	return isAlphabetOrUnderscore(ch) || ('0' <= ch && ch <= '9')
}
