package lexer

import (
	"testing"

	"github.com/JunNishimura/go-lisp/token"
)

func TestNextToken(t *testing.T) {
	input := `
8
+5
-10
(+ . (1 . (2 . nil)))
(+ 1 2)
(- 3 4)
(* 5 6)
(/ 7 8)
+(+ 1 2)
-(- 3 4)
(lambda (x) (+ x 1))
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "8"},
		{token.PLUS, "+"},
		{token.INT, "5"},
		{token.MINUS, "-"},
		{token.INT, "10"},
		{token.LPAREN, "("},
		{token.SYMBOL, "+"},
		{token.DOT, "."},
		{token.LPAREN, "("},
		{token.INT, "1"},
		{token.DOT, "."},
		{token.LPAREN, "("},
		{token.INT, "2"},
		{token.DOT, "."},
		{token.NIL, "nil"},
		{token.RPAREN, ")"},
		{token.RPAREN, ")"},
		{token.RPAREN, ")"},
		{token.LPAREN, "("},
		{token.SYMBOL, "+"},
		{token.INT, "1"},
		{token.INT, "2"},
		{token.RPAREN, ")"},
		{token.LPAREN, "("},
		{token.SYMBOL, "-"},
		{token.INT, "3"},
		{token.INT, "4"},
		{token.RPAREN, ")"},
		{token.LPAREN, "("},
		{token.SYMBOL, "*"},
		{token.INT, "5"},
		{token.INT, "6"},
		{token.RPAREN, ")"},
		{token.LPAREN, "("},
		{token.SYMBOL, "/"},
		{token.INT, "7"},
		{token.INT, "8"},
		{token.RPAREN, ")"},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.SYMBOL, "+"},
		{token.INT, "1"},
		{token.INT, "2"},
		{token.RPAREN, ")"},
		{token.MINUS, "-"},
		{token.LPAREN, "("},
		{token.SYMBOL, "-"},
		{token.INT, "3"},
		{token.INT, "4"},
		{token.RPAREN, ")"},
		{token.LPAREN, "("},
		{token.LAMBDA, "lambda"},
		{token.LPAREN, "("},
		{token.SYMBOL, "x"},
		{token.RPAREN, ")"},
		{token.LPAREN, "("},
		{token.SYMBOL, "+"},
		{token.SYMBOL, "x"},
		{token.INT, "1"},
		{token.RPAREN, ")"},
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
