package ast

import (
	"bytes"
	"monkey/token"
)

/**********************************************************
AST Base Interfaces
*********************************************************/

type Node interface {
	TokenLiteral() string
	String() string
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

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
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

func (stmt *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(stmt.TokenLiteral() + " ")
	out.WriteString(stmt.Name.String())
	out.WriteString(" = ")
	if stmt.Value != nil {
		out.WriteString(stmt.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (stmt *ReturnStatement) statementNode()       {}
func (stmt *ReturnStatement) TokenLiteral() string { return stmt.Token.Literal }

func (stmt *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(stmt.TokenLiteral() + " ")
	if stmt.ReturnValue != nil {
		out.WriteString(stmt.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token // first token in the expression
	Expression Expression
}

func (stmt *ExpressionStatement) statementNode()       {}
func (stmt *ExpressionStatement) TokenLiteral() string { return stmt.Token.Literal }

func (stmt *ExpressionStatement) String() string {
	if stmt.Expression != nil {
		return stmt.Expression.String()
	}
	return ""
}

type Identifier struct {
	Token token.Token
	Value string
}

func (expr *Identifier) expressionNode()      {}
func (expr *Identifier) TokenLiteral() string { return expr.Token.Literal }

func (expr *Identifier) String() string {
	return expr.Value
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (expr *IntegerLiteral) expressionNode()      {}
func (expr *IntegerLiteral) TokenLiteral() string { return expr.Token.Literal }
func (expr *IntegerLiteral) String() string       { return expr.Token.Literal }
