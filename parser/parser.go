package parser

import (
	"fmt"
	"gointer/ast"
	"gointer/lexer"
	"gointer/token"
	"strconv"
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
	{
		p.registerPrefix(token.IDENT, p.parseIdentifier)
		p.registerPrefix(token.INT, p.parseIntegerLiteral)
		p.registerPrefix(token.BANG, p.parsePrefixExpression)
		p.registerPrefix(token.MINUS, p.parsePrefixExpression)
		p.registerPrefix(token.TRUE, p.parseBoolean)
		p.registerPrefix(token.FALSE, p.parseBoolean)
		p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	}
	{
		p.registerInfix(token.PLUS, p.parseInfixExpression)
		p.registerInfix(token.MINUS, p.parseInfixExpression)
		p.registerInfix(token.SLASH, p.parseInfixExpression)
		p.registerInfix(token.ASTERISK, p.parseInfixExpression)
		p.registerInfix(token.EQ, p.parseInfixExpression)
		p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
		p.registerInfix(token.LT, p.parseInfixExpression)
		p.registerInfix(token.GT, p.parseInfixExpression)
	}
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

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
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
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expr.Right = p.parseExpression(PREFIX)
	return expr
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	exp.Right = p.parseExpression(precedence)
	return exp
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

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	val, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("could not parse %q as integer",
			p.curToken.Literal))
		return nil
	}
	lit.Value = val
	return lit
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.curToken,
		Value: p.curTokenIs(token.TRUE),
	}
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

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	p.errors = append(p.errors,
		fmt.Sprintf("no prefix parse function for %s found", t))
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

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGRETATER,
	token.GT:       LESSGRETATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}
