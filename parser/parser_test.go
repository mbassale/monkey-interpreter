package parser_test

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/parser"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		lexer := lexer.New(tt.input)
		parser := parser.New(lexer)
		program := parser.ParseProgram()
		checkParseErrors(t, parser)

		assert.Len(t, program.Statements, 1)

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 424242;
	`

	lexer := lexer.New(input)
	parser := parser.New(lexer)

	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements doesn not contain 3 statements. get=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}
		assert.Equal(t, returnStmt.TokenLiteral(), "return")
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	lexer := lexer.New(input)
	parser := parser.New(lexer)
	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expression not *ast.Identifier. got=%T", stmt.Expression)
	}

	assert.Equal(t, ident.Value, "foobar")
	assert.Equal(t, ident.TokenLiteral(), "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	lexer := lexer.New(input)
	parser := parser.New(lexer)
	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	assert.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	testLiteralExpression(t, stmt.Expression, int64(5))
}

func TestBooleanLiteralExpression(t *testing.T) {
	input := "true;false;"
	boolValues := []bool{true, false}

	lexer := lexer.New(input)
	parser := parser.New(lexer)
	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	assert.Len(t, program.Statements, 2)

	for idx, stmt := range program.Statements {
		exprStmt, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", stmt)
		}
		if !testLiteralExpression(t, exprStmt.Expression, boolValues[idx]) {
			return
		}
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world"`

	lexer := lexer.New(input)
	parser := parser.New(lexer)
	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("stmt.Expression not *ast.StringLiteral. got=%T", stmt.Expression)
	}

	assert.Equal(t, literal.Value, "hello world")
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		lexer := lexer.New(tt.input)
		parser := parser.New(lexer)
		program := parser.ParseProgram()
		checkParseErrors(t, parser)

		assert.Len(t, program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		expr, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		assert.Equal(t, expr.Operator, tt.operator)
		if !testLiteralExpression(t, expr.Right, tt.value) {
			return
		}
	}
}

func TestParserInfixExpressions(t *testing.T) {
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
		{"true == true;", true, "==", true},
		{"true != false;", true, "!=", false},
		{"false == false;", false, "==", false},
	}

	for _, tt := range infixTests {
		lexer := lexer.New(tt.input)
		parser := parser.New(lexer)
		program := parser.ParseProgram()
		checkParseErrors(t, parser)

		assert.Len(t, program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}

	for _, tt := range tests {
		lexer := lexer.New(tt.input)
		parser := parser.New(lexer)
		program := parser.ParseProgram()
		checkParseErrors(t, parser)

		actual := program.String()
		assert.Equal(t, actual, tt.expected)
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	lexer := lexer.New(input)
	parser := parser.New(lexer)
	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	assert.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	expr, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, expr.Condition, "x", "<", "y") {
		return
	}

	assert.Len(t, expr.ThenBranch.Statements, 1)

	thenBranch, ok := expr.ThenBranch.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expr.ThenBranch.Statements[0] is not ast.ExpressionStatement. got=%T", expr.ThenBranch.Statements[0])
	}

	if !testIdentifier(t, thenBranch.Expression, "x") {
		return
	}

	assert.Nil(t, expr.ElseBranch)
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	lexer := lexer.New(input)
	parser := parser.New(lexer)
	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	assert.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	expr, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, expr.Condition, "x", "<", "y") {
		return
	}

	assert.Len(t, expr.ThenBranch.Statements, 1)

	thenBranch, ok := expr.ThenBranch.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expr.ThenBranch.Statements[0] is not ast.ExpressionStatement. got=%T", expr.ThenBranch.Statements[0])
	}

	if !testIdentifier(t, thenBranch.Expression, "x") {
		return
	}

	assert.NotNil(t, expr.ElseBranch)

	elseBranch, ok := expr.ElseBranch.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expr.ElseBranch.Statements[0] is not ast.ExpressionStatement. got=%T", expr.ThenBranch.Statements[0])
	}

	if !testIdentifier(t, elseBranch.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	lexer := lexer.New(input)
	parser := parser.New(lexer)
	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	assert.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	assert.Len(t, function.Parameters, 2)
	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	assert.Len(t, function.Body.Statements, 1)

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParser(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		lexer := lexer.New(tt.input)
		parser := parser.New(lexer)
		program := parser.ParseProgram()
		checkParseErrors(t, parser)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		assert.Len(t, function.Parameters, len(tt.expectedParams))
		for i, identifier := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], identifier)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5)"

	lexer := lexer.New(input)
	parser := parser.New(lexer)
	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	assert.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	expr, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, expr.Function, "add") {
		return
	}

	assert.Len(t, expr.Arguments, 3)

	testLiteralExpression(t, expr.Arguments[0], 1)
	testInfixExpression(t, expr.Arguments[1], 2, "*", 3)
	testInfixExpression(t, expr.Arguments[2], 4, "+", 5)
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	lexer := lexer.New(input)
	parser := parser.New(lexer)
	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("expr not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	assert.Equal(t, len(array.Elements), 3)
	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func testIntegerLiteral(t *testing.T, literal ast.Expression, value int64) bool {
	integerLiteral, ok := literal.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("literal not *ast.IntegerLiteral. got=%T", literal)
		return false
	}
	assert.Equal(t, integerLiteral.Value, value)
	assert.Equal(t, integerLiteral.TokenLiteral(), fmt.Sprintf("%d", value))
	return true
}

func testBooleanLiteral(t *testing.T, literal ast.Expression, value bool) bool {
	booleanLiteral, ok := literal.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("literal not *ast.BooleanLiteral. got=%T", literal)
		return false
	}
	assert.Equal(t, booleanLiteral.Value, value)
	assert.Equal(t, booleanLiteral.TokenLiteral(), fmt.Sprintf("%t", value))
	return true
}

func testIdentifier(t *testing.T, expr ast.Expression, value string) bool {
	identifier, ok := expr.(*ast.Identifier)
	if !ok {
		t.Errorf("expr not *ast.Identifier. got=%T", expr)
		return false
	}
	assert.Equal(t, identifier.Value, value)
	assert.Equal(t, identifier.TokenLiteral(), value)
	return true
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool {
	assert.Equal(t, stmt.TokenLiteral(), "let")
	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("stmt not *ast.LetStatement. got=%T", stmt)
		return false
	}
	assert.Equal(t, letStmt.Name.Value, name)
	assert.Equal(t, letStmt.Name.TokenLiteral(), name)
	return true
}

func testLiteralExpression(t *testing.T, expr ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case bool:
		return testBooleanLiteral(t, expr, v)
	case int:
		return testIntegerLiteral(t, expr, int64(v))
	case int64:
		return testIntegerLiteral(t, expr, v)
	case string:
		return testIdentifier(t, expr, v)
	}
	t.Errorf("type of expr not handled. got=%T", expr)
	return false
}

func testInfixExpression(t *testing.T, expr ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExpr, ok := expr.(*ast.InfixExpression)
	if !ok {
		t.Errorf("expr is not ast.InfixExpression. got=%T(%s)", expr, expr)
		return false
	}

	if !testLiteralExpression(t, opExpr.Left, left) {
		return false
	}
	assert.Equal(t, opExpr.Operator, operator)
	return testLiteralExpression(t, opExpr.Right, right)
}

func checkParseErrors(t *testing.T, parser *parser.Parser) {
	errors := parser.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("Parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("Error: %q", msg)
	}
	t.FailNow()
}
