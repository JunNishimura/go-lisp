package lexer

import (
	"testing"

	"github.com/JunNishimura/go-lisp/token"
)

func TestNextToken(t *testing.T) {
	input := `
		(+ 1 2)
		(- 3 4)
		(* 5 6)
		(/ 7 8)
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LPAREN, "("},
		{token.PLUS, "+"},
		{token.INT, "1"},
		{token.INT, "2"},
		{token.RPAREN, ")"},
		{token.LPAREN, "("},
		{token.MINUS, "-"},
		{token.INT, "3"},
		{token.INT, "4"},
		{token.RPAREN, ")"},
		{token.LPAREN, "("},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.INT, "6"},
		{token.RPAREN, ")"},
		{token.LPAREN, "("},
		{token.SLASH, "/"},
		{token.INT, "7"},
		{token.INT, "8"},
		{token.RPAREN, ")"},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
