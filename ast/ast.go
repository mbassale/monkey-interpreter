package ast

import (
	"bytes"
	"monkey/token"
	"strings"
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

type BlockStatement struct {
	Token      token.Token // the '{' token
	Statements []Statement
}

func (stmt *BlockStatement) statementNode()       {}
func (stmt *BlockStatement) TokenLiteral() string { return stmt.Token.Literal }
func (stmt *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range stmt.Statements {
		out.WriteString(s.String())
	}
	return out.String()
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

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (expr *PrefixExpression) expressionNode()      {}
func (expr *PrefixExpression) TokenLiteral() string { return expr.Token.Literal }
func (expr *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(expr.Operator)
	out.WriteString(expr.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Token    token.Token // operator
	Left     Expression
	Operator string
	Right    Expression
}

func (expr *InfixExpression) expressionNode()      {}
func (expr *InfixExpression) TokenLiteral() string { return expr.Token.Literal }
func (expr *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(expr.Left.String())
	out.WriteString(" " + expr.Operator + " ")
	out.WriteString(expr.Right.String())
	out.WriteString(")")
	return out.String()
}

type IfExpression struct {
	Token      token.Token // the 'if' token
	Condition  Expression
	ThenBranch *BlockStatement
	ElseBranch *BlockStatement
}

func (expr *IfExpression) expressionNode()      {}
func (expr *IfExpression) TokenLiteral() string { return expr.Token.Literal }
func (expr *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(expr.Condition.String())
	out.WriteString(" ")
	out.WriteString(expr.ThenBranch.String())
	if expr.ElseBranch != nil {
		out.WriteString("else ")
		out.WriteString(expr.ElseBranch.String())
	}
	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token // the 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (expr *FunctionLiteral) expressionNode()      {}
func (expr *FunctionLiteral) TokenLiteral() string { return expr.Token.Literal }
func (expr *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, param := range expr.Parameters {
		params = append(params, param.String())
	}

	out.WriteString(expr.Token.Literal)
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(expr.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     token.Token // the '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (expr *CallExpression) expressionNode()      {}
func (expr *CallExpression) TokenLiteral() string { return expr.Token.Literal }
func (expr *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, arg := range expr.Arguments {
		args = append(args, arg.String())
	}

	out.WriteString(expr.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
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

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (expr *BooleanLiteral) expressionNode()      {}
func (expr *BooleanLiteral) TokenLiteral() string { return expr.Token.Literal }
func (expr *BooleanLiteral) String() string       { return expr.Token.Literal }

type StringLiteral struct {
	Token token.Token
	Value string
}

func (expr *StringLiteral) expressionNode()      {}
func (expr *StringLiteral) TokenLiteral() string { return expr.Token.Literal }
func (expr *StringLiteral) String() string       { return expr.Token.Literal }

type ArrayLiteral struct {
	Token    token.Token // the '[' token
	Elements []Expression
}

func (expr *ArrayLiteral) expressionNode()      {}
func (expr *ArrayLiteral) TokenLiteral() string { return expr.Token.Literal }

func (expr *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, elem := range expr.Elements {
		elements = append(elements, elem.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
