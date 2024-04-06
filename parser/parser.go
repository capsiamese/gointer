package parser

import (
	"fmt"
	"gointer/ast"
	"gointer/lexer"
	"gointer/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l, errors: make([]string, 0),
		prefixParseFns: make(map[token.TokenType]prefixParseFn),
		infixParseFns:  make(map[token.TokenType]infixParseFn),
	}
	p.nextToken()
	p.nextToken()
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	return p
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statments = make([]ast.Statment, 0)

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatment()
		if stmt != nil {
			program.Statments = append(program.Statments, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatment() ast.Statment {
	var result ast.Statment
	switch p.curToken.Type {
	case token.LET:
		result = p.parseLetStatment()
	case token.RETURN:
		result = p.parseReturnStatment()
	default:
		result = p.parseExpressionStatment()
	}
	if result != nil {
		return result
	}
	return nil
}

func (p *Parser) parseExpressionStatment() *ast.ExpressionStatment {
	stmt := &ast.ExpressionStatment{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}
	return prefix()
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseLetStatment() *ast.LetStatment {
	stmt := &ast.LetStatment{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	/*
		todo: parse expression
	*/

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) expectPeek(typ token.TokenType) bool {
	if p.peekTokenIs(typ) {
		p.nextToken()
		return true
	}
	p.peekError(typ)
	return false
}

func (p *Parser) peekTokenIs(typ token.TokenType) bool {
	return p.peekToken.Type == typ
}

func (p *Parser) curTokenIs(typ token.TokenType) bool {
	return p.curToken.Type == typ
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseReturnStatment() ast.Statment {
	rs := &ast.ReturnStatment{Token: p.curToken}

	p.nextToken()

	/*
		todo: parse expression
	*/
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return rs
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
	// sufixParseFn func(ast.Expression) ast.Exporession
)

const (
	_ int = iota
	LOWEST
	EQUALS       // ==
	LESSGRETATER // > <
	SUM          // +
	PRODUCT      // *
	PREFIX       // -X !X
	CALL         // Fn(x)
)
