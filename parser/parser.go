package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS       // ==
	LESS_GREATER // < or >
	SUM          // +
	PRODUCT      // *
	PREFIX       // -X or !X
	CALL         // myFunction()
)

type (
	PrefixParseFunction func() ast.Expression
	InfixParseFunction  func(ast.Expression) ast.Expression
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
	errors       []string

	prefixParseFunctions map[token.TokenType]PrefixParseFunction
	infixParseFunctions  map[token.TokenType]InfixParseFunction
}

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{lexer: lexer, errors: []string{}}

	parser.prefixParseFunctions = make(map[token.TokenType]PrefixParseFunction)
	parser.registerPrefix(token.IDENTIFIER, parser.parseIdentifier)
	parser.registerPrefix(token.INT, parser.parseIntegerLiteral)
	parser.registerPrefix(token.BANG, parser.parsePrefixExpression)
	parser.registerPrefix(token.MINUS, parser.parsePrefixExpression)

	// read two tokens to currentToken and nextToken are both set
	parser.nextToken()
	parser.nextToken()

	return parser
}

func (parser *Parser) ParseProgram() *ast.Program {
	program := ast.NewProgram()
	for parser.currentToken.Type != token.EOF {
		stmt := parser.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		parser.nextToken()
	}
	return program
}

func (parser *Parser) Errors() []string {
	return parser.errors
}

func (parser *Parser) nextToken() {
	parser.currentToken = parser.peekToken
	parser.peekToken = parser.lexer.NextToken()
}

func (parser *Parser) parseStatement() ast.Statement {
	switch parser.currentToken.Type {
	case token.LET:
		return parser.parseLetStatement()
	case token.RETURN:
		return parser.parseReturnStatement()
	default:
		return parser.parseExpressionStatement()
	}
}

func (parser *Parser) parseLetStatement() *ast.LetStatement {

	stmt := &ast.LetStatement{Token: parser.currentToken}

	if !parser.expectPeek(token.IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}

	if !parser.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: add expression parsing
	for !parser.currentTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return stmt
}

func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: parser.currentToken}

	parser.nextToken()

	// TODO: add expression parsing
	for !parser.currentTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return stmt
}

func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: parser.currentToken}

	stmt.Expression = parser.parseExpression(LOWEST)

	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return stmt
}

func (parser *Parser) parseExpression(precedence int) ast.Expression {
	prefix := parser.prefixParseFunctions[parser.currentToken.Type]
	if prefix == nil {
		parser.noPrefixParseFunctionError(parser.currentToken.Type)
		return nil
	}
	leftExpr := prefix()

	return leftExpr
}

func (parser *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}
}

func (parser *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: parser.currentToken}

	value, err := strconv.ParseInt(parser.currentToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", parser.currentToken.Literal)
		parser.errors = append(parser.errors, msg)
		return nil
	}

	literal.Value = value
	return literal
}

func (parser *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    parser.currentToken,
		Operator: parser.currentToken.Literal,
	}
	parser.nextToken()
	expression.Right = parser.parseExpression(PREFIX)
	return expression
}

func (parser *Parser) currentTokenIs(tokenType token.TokenType) bool {
	return parser.currentToken.Type == tokenType
}

func (parser *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return parser.peekToken.Type == tokenType
}

func (parser *Parser) expectPeek(tokenType token.TokenType) bool {
	if parser.peekTokenIs(tokenType) {
		parser.nextToken()
		return true
	} else {
		parser.peekError(tokenType)
		return false
	}
}

func (parser *Parser) peekError(tokenType token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, but got %s instead", tokenType, parser.peekToken.Type)
	parser.errors = append(parser.errors, msg)
}

func (parser *Parser) noPrefixParseFunctionError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	parser.errors = append(parser.errors, msg)
}

func (parser *Parser) registerPrefix(tokenType token.TokenType, fn PrefixParseFunction) {
	parser.prefixParseFunctions[tokenType] = fn
}

func (parser *Parser) registerInfix(tokenType token.TokenType, fn InfixParseFunction) {
	parser.infixParseFunctions[tokenType] = fn
}
