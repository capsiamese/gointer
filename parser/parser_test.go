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
	checkParseError(t, p)
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

func checkParseError(t *testing.T, p *Parser) {
	es := p.Errors()
	if len(es) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(es))
	for _, v := range es {
		t.Errorf("parse error: %q", v)
	}
	t.FailNow()
}

func TestReturnStatment(t *testing.T) {
	input := `
	return 5;
	return 10;
	return add(1, 2);`

	p := New(lexer.New(input))
	program := p.ParseProgram()

	if len(program.Statments) != 3 {
		t.Fatalf("program staments dose not contain 3 statements got=%d", len(program.Statments))
	}
	for _, stmt := range program.Statments {
		rs, ok := stmt.(*ast.ReturnStatment)
		if !ok {
			t.Errorf("stmt not ast.ReturnStatment got=%T", stmt)
		}
		if rs.TokenLiteral() != "return" {
			t.Errorf("return statment token literial not 'return' got=%q", rs.TokenLiteral())
		}
	}
}
