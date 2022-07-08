package ast

import "monkey/token"

/**********************************************************
AST Base Interfaces
*********************************************************/

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

/**********************************************************
AST Nodes
*********************************************************/

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func NewProgram() *Program {
	return &Program{Statements: []Statement{}}
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (stmt *LetStatement) statementNode()       {}
func (stmt *LetStatement) TokenLiteral() string { return stmt.Token.Literal }

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (stmt *ReturnStatement) statementNode()       {}
func (stmt *ReturnStatement) TokenLiteral() string { return stmt.Token.Literal }

type ExpressionStatement struct {
	Token      token.Token // first token in the expression
	Expression Expression
}

func (stmt *ExpressionStatement) statementNode()       {}
func (stmt *ExpressionStatement) TokenLiteral() string { return stmt.Token.Literal }

type Identifier struct {
	Token token.Token
	Value string
}

func (expr *Identifier) expressionNode()      {}
func (expr *Identifier) TokenLiteral() string { return expr.Token.Literal }
