package lexer

type TokenKind int

const (
	TokenKindEOF      TokenKind = iota
	TokenKindIdent              // x, y, foo
	TokenKindInt                // 42, 0
	TokenKindTrue               // true
	TokenKindFalse              // false
	TokenKindIf                 // if
	TokenKindThen               // then
	TokenKindElse               // else
	TokenKindBoolType           // Bool (type)
	TokenKindIntType            // Int (type)
	TokenKindLambda             // \
	TokenKindDot                // .
	TokenKindColon              // :
	TokenKindArrow              // ->
	TokenKindLParen             // (
	TokenKindRParen             // )
)

func (k TokenKind) String() string {
	switch k {
	case TokenKindEOF:
		return "EOF"
	case TokenKindIdent:
		return "Ident"
	case TokenKindInt:
		return "Int"
	case TokenKindTrue:
		return "True"
	case TokenKindFalse:
		return "False"
	case TokenKindIf:
		return "If"
	case TokenKindThen:
		return "Then"
	case TokenKindElse:
		return "Else"
	case TokenKindBoolType:
		return "BoolType"
	case TokenKindIntType:
		return "IntType"
	case TokenKindLambda:
		return "Lambda"
	case TokenKindDot:
		return "Dot"
	case TokenKindColon:
		return "Colon"
	case TokenKindArrow:
		return "Arrow"
	case TokenKindLParen:
		return "LParen"
	case TokenKindRParen:
		return "RParen"
	default:
		return "Unknown"
	}
}

type Token struct {
	Kind  TokenKind
	Value string
}
