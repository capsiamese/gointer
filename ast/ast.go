package ast

import "gointer/token"

type Node interface {
	TokenLiteral() string
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

type LetStatment struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (*LetStatment) statementNode()         {}
func (l *LetStatment) TokenLiteral() string { return l.Token.Literal }

type Identifier struct {
	Token token.Token
	Value string
}

func (*Identifier) expressionNode()        {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

type ReturnStatment struct {
	Token       token.Token
	ReturnValue Expression
}

func (r *ReturnStatment) statementNode()       {}
func (r *ReturnStatment) TokenLiteral() string { return r.Token.Literal }
