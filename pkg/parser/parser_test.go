package parser

import (
	"fmt"
	"testing"

	"github.com/hculpan/kabkey/pkg/ast"
	"github.com/hculpan/kabkey/pkg/lexer"
)

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"
	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParseErrors(t, p, []string{})

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, 0, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, 0, exp.Arguments[0], 1)
	testInfixExpression(t, 0, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, 0, exp.Arguments[2], 4, "+", 5)
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for testIndex, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParseErrors(t, p, []string{})

		stmt := program.Statements[0].(*ast.ExpressionStatement)

		function := stmt.Expression.(*ast.FunctionLiteral)
		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("[test %d] length parameters wrong. want %d, got=%d\n", testIndex, len(tt.expectedParams), len(function.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, testIndex, function.Parameters[i],
				ident)
		}
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`
	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParseErrors(t, p, []string{})

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n", len(function.Parameters))
	}

	testLiteralExpression(t, 0, function.Parameters[0], "x")
	testLiteralExpression(t, 0, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n", len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}

	testInfixExpression(t, 0, bodyStmt.Expression, "x", "+", "y")
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`
	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParseErrors(t, p, []string{})
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}
	if !testInfixExpression(t, 0, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(exp.Consequence.Statements))
	}
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}
	if !testIdentifier(t, 0, consequence.Expression, "x") {
		return
	}
	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}
func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`
	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParseErrors(t, p, []string{})
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}
	if !testInfixExpression(t, 0, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(exp.Consequence.Statements))
	}
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence.Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}
	if !testIdentifier(t, 0, consequence.Expression, "x") {
		return
	}
	if exp.Alternative == nil {
		t.Errorf("exp.Alternative is nil")
		return
	}
	if exp.Alternative.Statements == nil {
		t.Errorf("exp.Alternative.Statements is nil")
		return
	}
	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("alternative is not 1 statements. got=%d\n", len(exp.Alternative.Statements))
	}
	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Alternative.Statements[0] is not ast.ExpressionStatement. got=%T", exp.Alternative.Statements[0])
	}
	if !testIdentifier(t, 0, alternative.Expression, "y") {
		return
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParseErrors(t, p, []string{})

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		boolean, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
		}
		if boolean.Value != tt.expectedBoolean {
			t.Errorf("boolean.Value not %t. got=%t", tt.expectedBoolean,
				boolean.Value)
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
		{`"a" == "a"`, `("a" == "a")`},
	}
	for i, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParseErrors(t, p, []string{})
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("[test %d] expected=%q, got=%q", i, tt.expected, actual)
		}
	}
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for i, tt := range infixTests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParseErrors(t, p, []string{})

		if len(program.Statements) != 1 {
			t.Fatalf("[test %d] program has incorrect number of statements, expected 1, got %d", i, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("[test %d] program.Statement[0] is not *ast.ExpressionStatement, got %T", i, program.Statements[0])
		}

		if !testInfixExpression(t, i, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for i, tt := range prefixTests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParseErrors(t, p, []string{})

		if len(program.Statements) != 1 {
			t.Fatalf("[test %d] program has incorrect number of statements, expected 1, got %d", i, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statement[0] is not *ast.ExpressionStatement, got %T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("exp is not *ast.PrefixExpression, got %T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %q, got %q", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, i, exp.Right, tt.value) {
			return
		}
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	checkParseErrors(t, p, []string{})

	if len(program.Statements) != 1 {
		t.Fatalf("program has incorrect number of statements, expected 1, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not *ast.ExpressionStatement, got %T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.Identifier, got %T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("ident.Value not %d, got %d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("ident.TokenLiteral not %q, got %q", "5", literal.TokenLiteral())
	}

}

func TestStringLiteralExpression(t *testing.T) {
	input := `"foobar";`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	checkParseErrors(t, p, []string{})

	if len(program.Statements) != 1 {
		t.Fatalf("program has incorrect number of statements, expected 1, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not *ast.ExpressionStatement, got %T", program.Statements[0])
	}

	str, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral, got %T", stmt.Expression)
	}

	if str.Value != "foobar" {
		t.Errorf("str.Value not %q, got %q", "foobar", str.Value)
	}

	if str.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %q, got %q", "foobar", str.TokenLiteral())
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	checkParseErrors(t, p, []string{})

	if len(program.Statements) != 1 {
		t.Fatalf("program has incorrect number of statements, expected 1, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not *ast.ExpressionStatement, got %T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier, got %T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %q, got %q", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %q, got %q", "foobar", ident.TokenLiteral())
	}
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true", "y", true},
		{"let foobar = y", "foobar", "y"},
		{`let a = "foobar"`, "a", "foobar"},
	}
	for testIndex, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParseErrors(t, p, []string{})

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, testIndex, val, tt.expectedValue) {
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

func testIntegerLiteral(t *testing.T, testIndex int, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("[test %d] il not *ast.IntegerLiteral, got %T", testIndex, il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d, got %d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d, got %s", value, integ.TokenLiteral())
		return false
	}

	return true
}

func testStringLiteral(t *testing.T, testIndex int, sl ast.Expression, value string) bool {
	str, ok := sl.(*ast.StringLiteral)
	if !ok {
		t.Errorf("[test %d] sl not *ast.StringLiteral, got %T", testIndex, sl)
	}

	if str.Value != value {
		t.Errorf("str.Value not %q, got %q", value, str.Value)
		return false
	}

	if str.TokenLiteral() != fmt.Sprintf("%s", value) {
		t.Errorf("str.TokenLiteral not %q, got %q", value, str.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, testIndex int, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("[test %d] exp not *ast.Identifier, got %T", testIndex, exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("[test %d] ident.Value not %q, got %q", testIndex, value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("[test %d] ident.TokenLiteral not %q, got %q", testIndex, value, ident.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, testIndex int, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("[test %d] exp not *ast.Boolean, got %T", testIndex, exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("[test %d] bo.Value not %t, got %t", testIndex, value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("[test %d] bo.TokenLiteral not %t, got %s", testIndex, value, bo.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, testIndex int, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, testIndex, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, testIndex, exp, v)
	case string:
		if _, ok := exp.(*ast.StringLiteral); ok {
			return testStringLiteral(t, testIndex, exp, v)
		} else {
			return testIdentifier(t, testIndex, exp, v)
		}
	case bool:
		return testBooleanLiteral(t, testIndex, exp, v)
	default:
		t.Errorf("[test %d] type of exp not handled, got %T", testIndex, exp)
		return false
	}
}

func testInfixExpression(t *testing.T, testIndex int, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("[test %d] exp is not *ast.InfixExpression, got %T", testIndex, exp)
		return false
	}

	if !testLiteralExpression(t, testIndex, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("[test %d] exp.Operator is not %q, got %q", testIndex, operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, testIndex, opExp.Right, right) {
		return false
	}

	return true
}
