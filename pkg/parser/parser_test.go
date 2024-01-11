package parser

import (
	"testing"

	"github.com/hculpan/kabkey/pkg/ast"
	"github.com/hculpan/kabkey/pkg/lexer"
)

func TestLetStatement(t *testing.T) {
	input := `
let x = 5;
let y = 10 ;
let foobar = 838383;
`
	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	checkParseErrors(t, p, []string{})

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements: expected 3 statements, got %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`
	l := lexer.NewLexer(input)
	p := NewParser(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	checkParseErrors(t, p, []string{})

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements: expected 3 statements, got %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement, got %T", stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("expected type *ast.LetStatement, got %T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("expected name %q, got %q", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		if letStmt.Name.Value != name {
			t.Errorf("expected token literal %q, got %q", name, letStmt.Name.TokenLiteral())
			return false
		}
	}

	return true
}

func checkParseErrors(t *testing.T, p *Parser, testErrors []string) {
	errors := p.Errors()
	if len(errors) != len(testErrors) {
		t.Errorf("expected %d errors, got %d errors", len(testErrors), len(errors))

		for _, msg := range errors {
			t.Errorf("parser error: %s", msg)
		}

		t.FailNow()
	}

	for i := range errors {
		if testErrors[i] != errors[i] {
			t.Errorf("expected error %q, got %q", testErrors[i], errors[i])
		}
		t.FailNow()
	}
}
