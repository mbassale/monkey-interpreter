package lexer

import "monkey/token"

type Lexer struct {
	input        string // source
	position     int    // current char position
	nextPosition int    // next char position
	currentChar  byte   // current char
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.currentChar {
	case '=':
		if l.peekChar() == '=' {
			ch := l.currentChar
			l.readChar()
			literal := string(ch) + string(l.currentChar)
			tok = token.NewTokenWithLiteral(token.EQUAL, literal)
		} else {
			tok = token.NewToken(token.ASSIGN, l.currentChar)
		}
	case '+':
		tok = token.NewToken(token.PLUS, l.currentChar)
	case '-':
		tok = token.NewToken(token.MINUS, l.currentChar)
	case '*':
		tok = token.NewToken(token.STAR, l.currentChar)
	case '/':
		tok = token.NewToken(token.SLASH, l.currentChar)
	case '!':
		if l.peekChar() == '=' {
			ch := l.currentChar
			l.readChar()
			literal := string(ch) + string(l.currentChar)
			tok = token.NewTokenWithLiteral(token.NOT_EQUAL, literal)
		} else {
			tok = token.NewToken(token.BANG, l.currentChar)
		}
	case '<':
		tok = token.NewToken(token.LESS_THAN, l.currentChar)
	case '>':
		tok = token.NewToken(token.GREATER_THAN, l.currentChar)
	case ';':
		tok = token.NewToken(token.SEMICOLON, l.currentChar)
	case '(':
		tok = token.NewToken(token.LPAREN, l.currentChar)
	case ')':
		tok = token.NewToken(token.RPAREN, l.currentChar)
	case ',':
		tok = token.NewToken(token.COMMA, l.currentChar)
	case '{':
		tok = token.NewToken(token.LBRACE, l.currentChar)
	case '}':
		tok = token.NewToken(token.RBRACE, l.currentChar)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.currentChar) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdentifierTokenType(tok.Literal)
			return tok
		} else if isDigit(l.currentChar) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = token.NewToken(token.ILLEGAL, l.currentChar)
		}
	}

	l.readChar()

	return tok
}

func (l *Lexer) readChar() {
	if l.nextPosition >= len(l.input) {
		l.currentChar = 0
	} else {
		l.currentChar = l.input[l.nextPosition]
	}
	l.position = l.nextPosition
	l.nextPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.nextPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.nextPosition]
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.currentChar) || (l.position-position > 0 && isDigit(l.currentChar)) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.currentChar) || l.currentChar == '.' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.currentChar == ' ' || l.currentChar == '\t' || l.currentChar == '\n' || l.currentChar == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
