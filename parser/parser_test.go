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
	input := `
	let x = 5;
	let y = 10;
	let foobar = 424242;
	`

	lexer := lexer.New(input)
	parser := parser.New(lexer)
	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got %d", len(program.Statements))
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

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
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
		if !testIntegerLiteral(t, expr.Right, tt.integerValue) {
			return
		}
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

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("stmt.Expression not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	assert.Equal(t, literal.Value, int64(5))
	assert.Equal(t, literal.TokenLiteral(), "5")
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
