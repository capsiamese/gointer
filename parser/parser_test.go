package parser

import (
	"fmt"
	"gointer/ast"
	"gointer/lexer"
	"reflect"
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

func TestIdentifierExpression(t *testing.T) {
	p := New(lexer.New("foobar;"))
	program := p.ParseProgram()
	checkParseError(t, p)

	if len(program.Statments) != 1 {
		t.Fatalf("program staments dose not contain 1 statements got=%d", len(program.Statments))
	}

	stmt, ok := program.Statments[0].(*ast.ExpressionStatment)
	if !ok {
		t.Fatalf("statment node not expression statment got=%T", program.Statments[0])
	}
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not identifier got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Fatalf("ident value not foobar got=%s", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Fatalf("ident token literal not foobar got=%s", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5"

	p := New(lexer.New(input))
	program := p.ParseProgram()
	checkParseError(t, p)

	if len(program.Statments) != 1 {
		t.Fatalf("")
	}
	stmt, ok := program.Statments[0].(*ast.ExpressionStatment)
	if !ok {
		t.Fatalf("")
	}
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("")
	}
	if literal.Value != 5 {
		t.Fatalf("")
	}
	if literal.TokenLiteral() != "5" {
		t.Fatalf("")
	}
}

func TestParseingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTests {
		p := New(lexer.New(tt.input))
		prog := p.ParseProgram()
		checkParseError(t, p)
		AssertStmentCount(t, prog, 1)
		stmt := AssertStmentType[*ast.ExpressionStatment](t, prog, 0)
		exp := AssertExprType[*ast.PrefixExpression](t, stmt.Expression)
		if exp.Operator != tt.operator {
			t.Fatalf("")
		}
		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			t.Fatalf("")
		}
	}
}

func testIntegerLiteral(t *testing.T, l ast.Expression, val int64) bool {
	intVal, ok := l.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("")
		return false
	}
	if intVal.Value != val {
		t.Errorf("")
		return false
	}
	if intVal.TokenLiteral() != fmt.Sprintf("%d", val) {
		t.Errorf("")
		return false
	}
	return true
}

func AssertStmentType[T ast.Statment](t *testing.T, p *ast.Program, n int) T {
	stmt, ok := p.Statments[n].(T)
	if !ok {
		t.Fatalf("program statments[%d] not %s got %T", n, reflect.TypeFor[T](), p.Statments[n])
	}
	return stmt
}

func AssertStmentCount(t *testing.T, p *ast.Program, n int) {
	if len(p.Statments) != n {
		t.Fatalf("program statments count expect %d actual %d", n, len(p.Statments))
	}
}

func AssertExprType[T ast.Expression](t *testing.T, e ast.Expression) T {
	exp, ok := e.(T)
	if !ok {
		t.Fatalf("stmt is not %T got %T", reflect.TypeFor[T](), e)
	}
	return exp
}
