package lexer

import (
	"testing"

	"github.com/JunNishimura/go-lisp/token"
)

func TestProgram(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []token.Token
	}{
		{
			name: "multiple atoms in multiple lines",
			input: `
				1
				hoge
				-10
				nil
			`,
			expected: []token.Token{
				{Type: token.INT, Literal: "1"},
				{Type: token.SYMBOL, Literal: "hoge"},
				{Type: token.MINUS, Literal: "-"},
				{Type: token.INT, Literal: "10"},
				{Type: token.NIL, Literal: "nil"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name: "multiple lists in multiple lines",
			input: `
				(+ 1 2)
				(hoge)
				((lambda (x) (+ x 1)) 2)
			`,
			expected: []token.Token{
				{Type: token.LPAREN, Literal: "("},
				{Type: token.SYMBOL, Literal: "+"},
				{Type: token.INT, Literal: "1"},
				{Type: token.INT, Literal: "2"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.SYMBOL, Literal: "hoge"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.LAMBDA, Literal: "lambda"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.SYMBOL, Literal: "x"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.SYMBOL, Literal: "+"},
				{Type: token.SYMBOL, Literal: "x"},
				{Type: token.INT, Literal: "1"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.INT, Literal: "2"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name: "multiple atoms and lists in multiple lines",
			input: `
				1
				(+ 1 2)
			`,
			expected: []token.Token{
				{Type: token.INT, Literal: "1"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.SYMBOL, Literal: "+"},
				{Type: token.INT, Literal: "1"},
				{Type: token.INT, Literal: "2"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.EOF, Literal: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			for i, expected := range tt.expected {
				tok := l.NextToken()
				if tok.Type != expected.Type {
					t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
						i, expected.Type, tok.Type)
				}

				if tok.Literal != expected.Literal {
					t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
						i, expected.Literal, tok.Literal)
				}
			}
		})
	}
}

func TestAtom(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []token.Token
	}{
		{
			name:  "integer",
			input: "123",
			expected: []token.Token{
				{Type: token.INT, Literal: "123"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "prefix + integer",
			input: "+123",
			expected: []token.Token{
				{Type: token.PLUS, Literal: "+"},
				{Type: token.INT, Literal: "123"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "prefix - integer",
			input: "-123",
			expected: []token.Token{
				{Type: token.MINUS, Literal: "-"},
				{Type: token.INT, Literal: "123"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "T",
			input: "t",
			expected: []token.Token{
				{Type: token.TRUE, Literal: "t"},
				{Type: token.EOF, Literal: ""},
			},
		},
	}

	for _, tt := range tests {
		l := New(tt.input)
		for i, expected := range tt.expected {
			tok := l.NextToken()
			if tok.Type != expected.Type {
				t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
					i, expected.Type, tok.Type)
			}

			if tok.Literal != expected.Literal {
				t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
					i, expected.Literal, tok.Literal)
			}
		}
	}
}

func TestList(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []token.Token
	}{
		{
			name:  "empty list",
			input: "()",
			expected: []token.Token{
				{Type: token.LPAREN, Literal: "("},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "list with integer",
			input: "(123)",
			expected: []token.Token{
				{Type: token.LPAREN, Literal: "("},
				{Type: token.INT, Literal: "123"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "dotted pair",
			input: "(+ . (1 . (2 . nil)))",
			expected: []token.Token{
				{Type: token.LPAREN, Literal: "("},
				{Type: token.SYMBOL, Literal: "+"},
				{Type: token.DOT, Literal: "."},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.INT, Literal: "1"},
				{Type: token.DOT, Literal: "."},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.INT, Literal: "2"},
				{Type: token.DOT, Literal: "."},
				{Type: token.NIL, Literal: "nil"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "plus",
			input: "(+ 1 2)",
			expected: []token.Token{
				{Type: token.LPAREN, Literal: "("},
				{Type: token.SYMBOL, Literal: "+"},
				{Type: token.INT, Literal: "1"},
				{Type: token.INT, Literal: "2"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "minus",
			input: "(- 3 4)",
			expected: []token.Token{
				{Type: token.LPAREN, Literal: "("},
				{Type: token.SYMBOL, Literal: "-"},
				{Type: token.INT, Literal: "3"},
				{Type: token.INT, Literal: "4"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "multiply",
			input: "(* 5 6)",
			expected: []token.Token{
				{Type: token.LPAREN, Literal: "("},
				{Type: token.SYMBOL, Literal: "*"},
				{Type: token.INT, Literal: "5"},
				{Type: token.INT, Literal: "6"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "divide",
			input: "(/ 7 8)",
			expected: []token.Token{
				{Type: token.LPAREN, Literal: "("},
				{Type: token.SYMBOL, Literal: "/"},
				{Type: token.INT, Literal: "7"},
				{Type: token.INT, Literal: "8"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "equal",
			input: "(= 9 10)",
			expected: []token.Token{
				{Type: token.LPAREN, Literal: "("},
				{Type: token.SYMBOL, Literal: "="},
				{Type: token.INT, Literal: "9"},
				{Type: token.INT, Literal: "10"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "less than",
			input: "(< 11 12)",
			expected: []token.Token{
				{Type: token.LPAREN, Literal: "("},
				{Type: token.SYMBOL, Literal: "<"},
				{Type: token.INT, Literal: "11"},
				{Type: token.INT, Literal: "12"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "greater than",
			input: "(> 13 14)",
			expected: []token.Token{
				{Type: token.LPAREN, Literal: "("},
				{Type: token.SYMBOL, Literal: ">"},
				{Type: token.INT, Literal: "13"},
				{Type: token.INT, Literal: "14"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "less than or equal",
			input: "(<= 15 16)",
			expected: []token.Token{
				{Type: token.LPAREN, Literal: "("},
				{Type: token.SYMBOL, Literal: "<="},
				{Type: token.INT, Literal: "15"},
				{Type: token.INT, Literal: "16"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "greater than or equal",
			input: "(>= 17 18)",
			expected: []token.Token{
				{Type: token.LPAREN, Literal: "("},
				{Type: token.SYMBOL, Literal: ">="},
				{Type: token.INT, Literal: "17"},
				{Type: token.INT, Literal: "18"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "not equal",
			input: "(/= 19 20)",
			expected: []token.Token{
				{Type: token.LPAREN, Literal: "("},
				{Type: token.SYMBOL, Literal: "/="},
				{Type: token.INT, Literal: "19"},
				{Type: token.INT, Literal: "20"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "list with prefix +",
			input: "+(+ 1 2)",
			expected: []token.Token{
				{Type: token.PLUS, Literal: "+"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.SYMBOL, Literal: "+"},
				{Type: token.INT, Literal: "1"},
				{Type: token.INT, Literal: "2"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "list with prefix -",
			input: "-(- 3 4)",
			expected: []token.Token{
				{Type: token.MINUS, Literal: "-"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.SYMBOL, Literal: "-"},
				{Type: token.INT, Literal: "3"},
				{Type: token.INT, Literal: "4"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "lambda",
			input: "(lambda (x) (+ x 1))",
			expected: []token.Token{
				{Type: token.LPAREN, Literal: "("},
				{Type: token.LAMBDA, Literal: "lambda"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.SYMBOL, Literal: "x"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.SYMBOL, Literal: "+"},
				{Type: token.SYMBOL, Literal: "x"},
				{Type: token.INT, Literal: "1"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "defmacro",
			input: "(defmacro unless (cond body) `(if (not ,cond) ,body nil))",
			expected: []token.Token{
				{Type: token.LPAREN, Literal: "("},
				{Type: token.SYMBOL, Literal: "defmacro"},
				{Type: token.SYMBOL, Literal: "unless"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.SYMBOL, Literal: "cond"},
				{Type: token.SYMBOL, Literal: "body"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.BACKQUOTE, Literal: "`"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.IF, Literal: "if"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.SYMBOL, Literal: "not"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.SYMBOL, Literal: "cond"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.SYMBOL, Literal: "body"},
				{Type: token.NIL, Literal: "nil"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "quote",
			input: "'(1 2 3)",
			expected: []token.Token{
				{Type: token.QUOTE, Literal: "'"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.INT, Literal: "1"},
				{Type: token.INT, Literal: "2"},
				{Type: token.INT, Literal: "3"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "quote(symbol)",
			input: "(quote 1 2 3)",
			expected: []token.Token{
				{Type: token.LPAREN, Literal: "("},
				{Type: token.QUOTE, Literal: "quote"},
				{Type: token.INT, Literal: "1"},
				{Type: token.INT, Literal: "2"},
				{Type: token.INT, Literal: "3"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "backquote",
			input: "`(1 2 3)",
			expected: []token.Token{
				{Type: token.BACKQUOTE, Literal: "`"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.INT, Literal: "1"},
				{Type: token.INT, Literal: "2"},
				{Type: token.INT, Literal: "3"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "comma in backquotted list",
			input: "`(1 2 ,(+ 1 2))",
			expected: []token.Token{
				{Type: token.BACKQUOTE, Literal: "`"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.INT, Literal: "1"},
				{Type: token.INT, Literal: "2"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.SYMBOL, Literal: "+"},
				{Type: token.INT, Literal: "1"},
				{Type: token.INT, Literal: "2"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "if expression",
			input: "(if (= 1 2) 3 4)",
			expected: []token.Token{
				{Type: token.LPAREN, Literal: "("},
				{Type: token.IF, Literal: "if"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.SYMBOL, Literal: "="},
				{Type: token.INT, Literal: "1"},
				{Type: token.INT, Literal: "2"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.INT, Literal: "3"},
				{Type: token.INT, Literal: "4"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.EOF, Literal: ""},
			},
		},
	}

	for _, tt := range tests {
		l := New(tt.input)
		for i, expected := range tt.expected {
			tok := l.NextToken()
			if tok.Type != expected.Type {
				t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
					i, expected.Type, tok.Type)
			}

			if tok.Literal != expected.Literal {
				t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
					i, expected.Literal, tok.Literal)
			}
		}
	}
}
