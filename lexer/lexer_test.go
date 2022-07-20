package lexer_test

import (
	"monkey/lexer"
	"monkey/token"
	"testing"
)

func TestSingleCharacterTokens(t *testing.T) {
	input := `=!+-*/<>(){},;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.BANG, "!"},
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.STAR, "*"},
		{token.SLASH, "/"},
		{token.LESS_THAN, "<"},
		{token.GREATER_THAN, ">"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
	}

	l := lexer.New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokenType wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestSource(t *testing.T) {
	input := `
	let five = 5;
	let ten = 10;

	let add = fn(x, y) {
		x + y;
	};

	let result = add(five, ten);

	if (5 < 10) {
		return true;
	} else {
		return false;
	}

	10==10;
	10!=9;
	"foobar"
	"foo bar"
	[1, 2];
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		// line 1
		{token.LET, "let"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		// line 2
		{token.LET, "let"},
		{token.IDENTIFIER, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		// line 4
		{token.LET, "let"},
		{token.IDENTIFIER, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},

		// line 5
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.SEMICOLON, ";"},

		// line 6
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		// line 8
		{token.LET, "let"},
		{token.IDENTIFIER, "result"},
		{token.ASSIGN, "="},
		{token.IDENTIFIER, "add"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "five"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		// line 10
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LESS_THAN, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},

		// line 11
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},

		// line 12
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},

		// line 13
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},

		// line 14
		{token.RBRACE, "}"},

		// line 16
		{token.INT, "10"},
		{token.EQUAL, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		// line 17
		{token.INT, "10"},
		{token.NOT_EQUAL, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},

		// line 21
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},

		// line 23
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},

		// EOF
		{token.EOF, ""},
	}

	l := lexer.New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokenType wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong, expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
