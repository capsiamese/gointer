package ast

import (
	"fmt"
	"gointer/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statment interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statments []Statment
}

func (p *Program) TokenLiteral() string {
	if len(p.Statments) > 0 {
		return p.Statments[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	buf := strings.Builder{}
	for _, v := range p.Statments {
		buf.WriteString(v.String())
	}
	return buf.String()
}

type LetStatment struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (*LetStatment) statementNode()         {}
func (l *LetStatment) TokenLiteral() string { return l.Token.Literal }
func (l *LetStatment) String() string {
	return fmt.Sprintf("%s %s = %s;",
		l.TokenLiteral(), l.Name.Value, l.Value.String())
}

type Identifier struct {
	Token token.Token
	Value string
}

func (*Identifier) expressionNode()        {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string {
	return i.Value
}

type ReturnStatment struct {
	Token       token.Token
	ReturnValue Expression
}

func (r *ReturnStatment) statementNode()       {}
func (r *ReturnStatment) TokenLiteral() string { return r.Token.Literal }
func (r *ReturnStatment) String() string {
	return fmt.Sprintf("%s %s;",
		r.TokenLiteral(), r.ReturnValue.String())
}

type ExpressionStatment struct {
	Token      token.Token
	Expression Expression
}

func (e *ExpressionStatment) statementNode()       {}
func (e *ExpressionStatment) TokenLiteral() string { return e.Token.Literal }
func (e *ExpressionStatment) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) expressionNode()      {}
func (i *IntegerLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	buf := strings.Builder{}
	buf.WriteString("(")
	buf.WriteString(pe.Operator)
	buf.WriteString(pe.Right.String())
	buf.WriteString(")")
	return buf.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	return fmt.Sprintf(
		"(%s%s%s)",
		ie.Left, ie.Operator, ie.Right)
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatment
	Alternative *BlockStatment
}
