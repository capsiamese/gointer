package parser

import (
	"gointer/ast"
	"gointer/lexer"
	"testing"
)

func TestLetStatments(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 123456;`

	p := New(lexer.New(input))
	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram returned nil")
	}
	if len(program.Statments) != 3 {
		t.Fatalf("program staments dose not contain 3 statements got=%d", len(program.Statments))
	}
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"}, {"y"}, {"foobar"},
	}
	for i, tt := range tests {
		stmt := program.Statments[i]
		if !testLetStatment(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatment(t *testing.T, s ast.Statment, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("TokenLiterial not let got=%q", s.TokenLiteral())
		return false
	}
	letStmt, ok := s.(*ast.LetStatment)
	if !ok {
		t.Errorf("s not *ast.LetStatment got=%T", s)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf("letStmt name not %s got=%s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt token literial not %s got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}
	return true
}
