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

	"github.com/shota3506/gostlc/internal/ast"
	"github.com/shota3506/gostlc/internal/lexer"
)

type parser struct {
	lexer     *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
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
			Func: expr,
			Arg:  arg,
		}
	}
}

func (p *parser) canStartExpr() bool {
	switch p.curToken.Kind {
	case lexer.TokenKindLambda, lexer.TokenKindLParen,
		lexer.TokenKindTrue, lexer.TokenKindFalse,
		lexer.TokenKindIf, lexer.TokenKindInt,
		lexer.TokenKindIdent:
		return true
	default:
		return false
	}
}

// parsePrimary parses a primary expression (non-application)
func (p *parser) parsePrimary() (ast.Expr, error) {
	switch p.curToken.Kind {
	case lexer.TokenKindLambda:
		return p.parseAbstraction()
	case lexer.TokenKindLParen:
		return p.parseGrouping()
	case lexer.TokenKindTrue:
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		return &ast.BoolExpr{Value: true}, nil
	case lexer.TokenKindFalse:
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		return &ast.BoolExpr{Value: false}, nil
	case lexer.TokenKindIf:
		return p.parseIfExpr()
	case lexer.TokenKindInt:
		value := p.curToken.Value
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		var intVal int
		fmt.Sscanf(value, "%d", &intVal)
		return &ast.IntExpr{Value: intVal}, nil
	case lexer.TokenKindIdent:
		name := p.curToken.Value
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		return &ast.VarExpr{Name: name}, nil
	default:
		return nil, fmt.Errorf("unexpected token: %v", p.curToken)
	}
}

// parseAbstraction parses a lambda abstraction: \var:type. expr
func (p *parser) parseAbstraction() (ast.Expr, error) {
	// Consume '\'
	if err := p.nextToken(); err != nil {
		return nil, err
	}

	// Parse parameter name
	if p.curToken.Kind != lexer.TokenKindIdent {
		return nil, fmt.Errorf("expected identifier after '\\', got %v", p.curToken)
	}
	param := p.curToken.Value
	if err := p.nextToken(); err != nil {
		return nil, err
	}

	// Expect ':'
	if p.curToken.Kind != lexer.TokenKindColon {
		return nil, fmt.Errorf("expected ':' after parameter name, got %v", p.curToken)
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
	if p.curToken.Kind != lexer.TokenKindDot {
		return nil, fmt.Errorf("expected '.' after parameter type, got %v", p.curToken)
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
	if p.curToken.Kind != lexer.TokenKindRParen {
		return nil, fmt.Errorf("expected ')', got %v", p.curToken)
	}
	if err := p.nextToken(); err != nil {
		return nil, err
	}

	return expr, nil
}

// parseIfExpr parses a conditional expression
func (p *parser) parseIfExpr() (ast.Expr, error) {
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
	if p.curToken.Kind != lexer.TokenKindThen {
		return nil, fmt.Errorf("expected 'then', got %v", p.curToken)
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
	if p.curToken.Kind != lexer.TokenKindElse {
		return nil, fmt.Errorf("expected 'else', got %v", p.curToken)
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
	if p.curToken.Kind == lexer.TokenKindArrow {
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
	case lexer.TokenKindBoolType:
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		return &ast.BooleanType{}, nil
	case lexer.TokenKindIntType:
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		return &ast.IntType{}, nil
	case lexer.TokenKindLParen:
		// Grouped type
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		typ, err := p.parseType()
		if err != nil {
			return nil, err
		}
		if p.curToken.Kind != lexer.TokenKindRParen {
			return nil, fmt.Errorf("expected ')' after type, got %v", p.curToken)
		}
		if err := p.nextToken(); err != nil {
			return nil, err
		}
		return typ, nil
	default:
		return nil, fmt.Errorf("unexpected token in type: %v", p.curToken)
	}
}
