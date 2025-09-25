// package parser implements a parser for a simple typed lambda calculus.
//
// The grammar is defined as follows:
// ```
// expr ::= var
//        | "\" var ":" type "." expr         (* abstraction *)
//        | expr expr                         (* application *)
//        | "(" expr ")"                      (* grouping *)
//        | "true" | "false"                  (* boolean literals *)
//        | "if" expr "then" expr "else" expr (* conditional *)
//        | digit+                            (* integer literals *)
// type ::= "Bool"                            (* boolean type *)
//        | "Int"                             (* integer type *)
//        | type "->" type                    (* function type *)
//        | "(" type ")"                      (* grouping *)
// var  ::= letter (letter | digit)*          (* variable names *)
// ```

package parser

import (
	"fmt"
	"strconv"

	"github.com/shota3506/gostlc/internal/ast"
	"github.com/shota3506/gostlc/internal/lexer"
	"github.com/shota3506/gostlc/internal/token"
)

type parser struct {
	lexer     *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
}

// Parse parses the input string and returns the corresponding AST expression.
func Parse(s string) (ast.Expr, error) {
	l := lexer.New(s)
	p := &parser{lexer: l}

	if err := p.nextToken(); err != nil {
		return nil, err
	}
	if err := p.nextToken(); err != nil {
		return nil, err
	}
	return p.parseExpr()
}

func (p *parser) nextToken() error {
	p.curToken = p.peekToken
	tok, err := p.lexer.Next()
	if err != nil {
		return err
	}
	p.peekToken = tok
	return nil
}

// parseExpr parses an expression with left-associative application.
func (p *parser) parseExpr() (ast.Expr, error) {
	expr, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	// Handle application (left-associative)
	for {
		if !p.canStartExpr() {
			return expr, nil
		}
		arg, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		expr = &ast.AppExpr{
			Pos:  expr.Position(),
			Func: expr,
			Arg:  arg,
		}
	}
}

func (p *parser) canStartExpr() bool {
	switch p.curToken.Kind {
	case token.TokenKindLambda, token.TokenKindLParen,
		token.TokenKindTrue, token.TokenKindFalse,
		token.TokenKindIf, token.TokenKindInt,
		token.TokenKindIdent:
		return true
	default:
		return false
	}
}

// parsePrimary parses a primary expression (non-application)
func (p *parser) parsePrimary() (ast.Expr, error) {
	switch p.curToken.Kind {
	case token.TokenKindLambda:
		return p.parseAbstraction()
	case token.TokenKindLParen:
		return p.parseGrouping()
	case token.TokenKindTrue:
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		return &ast.BoolExpr{
			Pos:   p.curToken.Pos,
			Value: true,
		}, nil
	case token.TokenKindFalse:
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		return &ast.BoolExpr{
			Pos:   p.curToken.Pos,
			Value: false,
		}, nil
	case token.TokenKindIf:
		return p.parseIfExpr()
	case token.TokenKindInt:
		value := p.curToken.Value
		pos := p.curToken.Pos
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, newParseError(p.curToken, fmt.Sprintf("invalid integer literal: %v", value))
		}
		return &ast.IntExpr{
			Pos:   pos,
			Value: int(intVal),
		}, nil
	case token.TokenKindIdent:
		pos := p.curToken.Pos
		name := p.curToken.Value
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		return &ast.VarExpr{
			Pos:  pos,
			Name: name,
		}, nil
	default:
		return nil, newParseError(p.curToken, fmt.Sprintf("unexpected token: %v", p.curToken.Kind))
	}
}

// parseAbstraction parses a lambda abstraction: \var:type. expr
func (p *parser) parseAbstraction() (ast.Expr, error) {
	// Save position of lambda
	pos := p.curToken.Pos

	// Consume '\'
	if err := p.nextToken(); err != nil {
		return nil, err
	}

	// Parse parameter name
	if p.curToken.Kind != token.TokenKindIdent {
		return nil, newParseError(p.curToken, fmt.Sprintf("expected identifier after '\\': %v", p.curToken.Kind))
	}
	param := p.curToken.Value
	if err := p.nextToken(); err != nil {
		return nil, err
	}

	// Expect ':'
	if p.curToken.Kind != token.TokenKindColon {
		return nil, newParseError(p.curToken, fmt.Sprintf("expected ':' after parameter name: %v", p.curToken.Kind))
	}
	if err := p.nextToken(); err != nil {
		return nil, err
	}

	// Parse parameter type
	paramType, err := p.parseType()
	if err != nil {
		return nil, err
	}

	// Expect '.'
	if p.curToken.Kind != token.TokenKindDot {
		return nil, newParseError(p.curToken, fmt.Sprintf("expected '.' after parameter type: %v", p.curToken.Kind))
	}
	if err := p.nextToken(); err != nil {
		return nil, err
	}

	// Parse body
	body, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	return &ast.AbsExpr{
		Pos:       pos,
		Param:     param,
		ParamType: paramType,
		Body:      body,
	}, nil
}

// parseGrouping parses a parenthesized expression
func (p *parser) parseGrouping() (ast.Expr, error) {
	// Consume '('
	if err := p.nextToken(); err != nil {
		return nil, err
	}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	// Expect ')'
	if p.curToken.Kind != token.TokenKindRParen {
		return nil, newParseError(p.curToken, fmt.Sprintf("expected ')': %v", p.curToken.Kind))
	}
	if err := p.nextToken(); err != nil {
		return nil, err
	}

	return expr, nil
}

// parseIfExpr parses a conditional expression
func (p *parser) parseIfExpr() (ast.Expr, error) {
	// Save position of 'if'
	pos := p.curToken.Pos

	// Consume 'if'
	if err := p.nextToken(); err != nil {
		return nil, err
	}

	// Parse condition
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	// Expect 'then'
	if p.curToken.Kind != token.TokenKindThen {
		return nil, newParseError(p.curToken, fmt.Sprintf("expected 'then': %v", p.curToken.Kind))
	}
	if err := p.nextToken(); err != nil {
		return nil, err
	}

	// Parse then branch
	thenExpr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	// Expect 'else'
	if p.curToken.Kind != token.TokenKindElse {
		return nil, newParseError(p.curToken, fmt.Sprintf("expected 'else': %v", p.curToken.Kind))
	}
	if err := p.nextToken(); err != nil {
		return nil, err
	}

	// Parse else branch
	elseExpr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	return &ast.IfExpr{
		Pos:  pos,
		Cond: cond,
		Then: thenExpr,
		Else: elseExpr,
	}, nil
}

// parseType parses a type with right-associative arrow
func (p *parser) parseType() (ast.Type, error) {
	baseType, err := p.parseBaseType()
	if err != nil {
		return nil, err
	}

	// Check for function type (right-associative)
	if p.curToken.Kind == token.TokenKindArrow {
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		toType, err := p.parseType()
		if err != nil {
			return nil, err
		}
		return &ast.FuncType{
			From: baseType,
			To:   toType,
		}, nil
	}

	return baseType, nil
}

// parseBaseType parses a base type or grouped type
func (p *parser) parseBaseType() (ast.Type, error) {
	switch p.curToken.Kind {
	case token.TokenKindBoolType:
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		return &ast.BooleanType{}, nil
	case token.TokenKindIntType:
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		return &ast.IntType{}, nil
	case token.TokenKindLParen:
		// Grouped type
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		typ, err := p.parseType()
		if err != nil {
			return nil, err
		}
		if p.curToken.Kind != token.TokenKindRParen {
			return nil, newParseError(p.curToken, fmt.Sprintf("expected ')': %v", p.curToken.Kind))
		}
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		return typ, nil
	default:
		return nil, newParseError(p.curToken, fmt.Sprintf("unexpected token in type: %v", p.curToken.Kind))
	}
}
