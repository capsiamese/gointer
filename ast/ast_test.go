package ast

import (
	"gointer/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statments: []Statment{
			&LetStatment{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}
	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program wrong got=%q", program.String())
	}
}
