package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
	errors       []string
}

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{lexer: lexer, errors: []string{}}

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
		return nil
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
