package parser

import (
	"testing"

	"github.com/JunNishimura/go-lisp/ast"
	"github.com/JunNishimura/go-lisp/lexer"
	"github.com/JunNishimura/go-lisp/token"
)

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestIntegerAtom(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			name:     "parse positive integer",
			input:    "1",
			expected: 1,
		},
		{
			name:     "parse multiple digit integer",
			input:    "1234567890",
			expected: 1234567890,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Expressions) != 1 {
				t.Fatalf("program.Expressions does not contain 1 expressions. got=%d", len(program.Expressions))
			}
			atom, ok := program.Expressions[0].(*ast.IntegerLiteral)
			if !ok {
				t.Fatalf("exp not *ast.IntegerLiteral. got=%T", program.Expressions[0])
			}
			if atom.Value != tt.expected {
				t.Fatalf("literal.Value not %d. got=%d", tt.expected, atom.Value)
			}
		})
	}
}

func TestPrefixAtom(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *ast.PrefixAtom
	}{
		{
			name:  "parse positive integer",
			input: "+1",
			expected: &ast.PrefixAtom{
				Token:    token.Token{Type: token.PLUS, Literal: "+"},
				Operator: "+",
				Right: &ast.IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: "1"},
					Value: 1,
				},
			},
		},
		{
			name:  "parse negative integer",
			input: "-1",
			expected: &ast.PrefixAtom{
				Token:    token.Token{Type: token.MINUS, Literal: "-"},
				Operator: "-",
				Right: &ast.IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: "1"},
					Value: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Expressions) != 1 {
				t.Fatalf("program.Expressions does not contain 1 expressions. got=%d", len(program.Expressions))
			}
			atom, ok := program.Expressions[0].(*ast.PrefixAtom)
			if !ok {
				t.Fatalf("exp not *ast.PrefixExpression. got=%T", program.Expressions[0])
			}
			if atom.String() != tt.expected.String() {
				t.Fatalf("literal.Value not %s. got=%s", tt.expected, atom.String())
			}
		})
	}
}

func TestNil(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *ast.Nil
	}{
		{
			name:     "literal nil",
			input:    "nil",
			expected: &ast.Nil{Token: token.Token{Type: token.NIL, Literal: "nil"}},
		},
		{
			name:     "empty list",
			input:    "()",
			expected: &ast.Nil{Token: token.Token{Type: token.NIL, Literal: "nil"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Expressions) != 1 {
				t.Fatalf("program.Expressions does not contain 1 expressions. got=%d", len(program.Expressions))
			}
			n, ok := program.Expressions[0].(*ast.Nil)
			if !ok {
				t.Fatalf("exp not *ast.Nil. got=%T", program.Expressions[0])
			}
			if n.String() != tt.expected.String() {
				t.Fatalf("literal.Value not %s. got=%s", tt.expected.String(), n.String())
			}
		})
	}
}

func TestArithmeticOperations(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *ast.ConsCell
	}{
		{
			name:  "single element list",
			input: "(1)",
			expected: &ast.ConsCell{
				CarField: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1},
				CdrField: &ast.Nil{Token: token.Token{Type: token.NIL, Literal: "nil"}},
			},
		},
		{
			name:  "simple addition",
			input: "(+ . (1 . (2 . nil)))",
			expected: &ast.ConsCell{
				CarField: &ast.Symbol{Token: token.Token{Type: token.SYMBOL, Literal: "+"}, Value: "+"},
				CdrField: &ast.ConsCell{
					CarField: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1},
					CdrField: &ast.ConsCell{
						CarField: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "2"}, Value: 2},
						CdrField: &ast.Nil{Token: token.Token{Type: token.NIL, Literal: "nil"}},
					},
				},
			},
		},
		{
			name:  "simple subtraction",
			input: "(- . (1 . (2 . nil)))",
			expected: &ast.ConsCell{
				CarField: &ast.Symbol{Token: token.Token{Type: token.SYMBOL, Literal: "-"}, Value: "-"},
				CdrField: &ast.ConsCell{
					CarField: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1},
					CdrField: &ast.ConsCell{
						CarField: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "2"}, Value: 2},
						CdrField: &ast.Nil{Token: token.Token{Type: token.NIL, Literal: "nil"}},
					},
				},
			},
		},
		{
			name:  "simple multiplication",
			input: "(* 1 2)",
			expected: &ast.ConsCell{
				CarField: &ast.Symbol{Token: token.Token{Type: token.SYMBOL, Literal: "*"}, Value: "*"},
				CdrField: &ast.ConsCell{
					CarField: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1},
					CdrField: &ast.ConsCell{
						CarField: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "2"}, Value: 2},
						CdrField: &ast.Nil{Token: token.Token{Type: token.NIL, Literal: "nil"}},
					},
				},
			},
		},
		{
			name:  "simple division",
			input: "(/ 1 2)",
			expected: &ast.ConsCell{
				CarField: &ast.Symbol{Token: token.Token{Type: token.SYMBOL, Literal: "/"}, Value: "/"},
				CdrField: &ast.ConsCell{
					CarField: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1},
					CdrField: &ast.ConsCell{
						CarField: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "2"}, Value: 2},
						CdrField: &ast.Nil{Token: token.Token{Type: token.NIL, Literal: "nil"}},
					},
				},
			},
		},
		{
			name:  "nested cons cell",
			input: "(* (+ 1 2) (- 4 3))",
			expected: &ast.ConsCell{
				CarField: &ast.Symbol{Token: token.Token{Type: token.SYMBOL, Literal: "*"}, Value: "*"},
				CdrField: &ast.ConsCell{
					CarField: &ast.ConsCell{
						CarField: &ast.Symbol{Token: token.Token{Type: token.SYMBOL, Literal: "+"}, Value: "+"},
						CdrField: &ast.ConsCell{
							CarField: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1},
							CdrField: &ast.ConsCell{
								CarField: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "2"}, Value: 2},
								CdrField: &ast.Nil{Token: token.Token{Type: token.NIL, Literal: "nil"}},
							},
						},
					},
					CdrField: &ast.ConsCell{
						CarField: &ast.ConsCell{
							CarField: &ast.Symbol{Token: token.Token{Type: token.SYMBOL, Literal: "-"}, Value: "-"},
							CdrField: &ast.ConsCell{
								CarField: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "4"}, Value: 4},
								CdrField: &ast.ConsCell{
									CarField: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "3"}, Value: 3},
									CdrField: &ast.Nil{Token: token.Token{Type: token.NIL, Literal: "nil"}},
								},
							},
						},
						CdrField: &ast.Nil{Token: token.Token{Type: token.NIL, Literal: "nil"}},
					},
				},
			},
		},
		{
			name:  "mix of dotted pair and list",
			input: "(+ . (1 2))",
			expected: &ast.ConsCell{
				CarField: &ast.Symbol{Token: token.Token{Type: token.SYMBOL, Literal: "+"}, Value: "+"},
				CdrField: &ast.ConsCell{
					CarField: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1},
					CdrField: &ast.ConsCell{
						CarField: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "2"}, Value: 2},
						CdrField: &ast.Nil{Token: token.Token{Type: token.NIL, Literal: "nil"}},
					},
				},
			},
		},
		{
			name:  "include prefix atom",
			input: "(+ -5 5)",
			expected: &ast.ConsCell{
				CarField: &ast.Symbol{Token: token.Token{Type: token.SYMBOL, Literal: "+"}, Value: "+"},
				CdrField: &ast.ConsCell{
					CarField: &ast.PrefixAtom{
						Token:    token.Token{Type: token.MINUS, Literal: "-"},
						Operator: "-",
						Right: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
					},
					CdrField: &ast.ConsCell{
						CarField: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "5"}, Value: 5},
						CdrField: &ast.Nil{Token: token.Token{Type: token.NIL, Literal: "nil"}},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Expressions) != 1 {
				t.Fatalf("program.Expressions does not contain 1 expressions. got=%d", len(program.Expressions))
			}
			cc, ok := program.Expressions[0].(*ast.ConsCell)
			if !ok {
				t.Fatalf("exp not *ast.ConsCell. got=%T", program.Expressions[0])
			}
			if cc.String() != tt.expected.String() {
				t.Fatalf("cc.String() not %s. got=%s", tt.expected.String(), cc.String())
			}
		})
	}
}

func TestLambdaExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *ast.ConsCell
	}{
		{
			name:  "lambda expression with no parameter",
			input: "(lambda () 1)",
			expected: &ast.ConsCell{
				CarField: &ast.Symbol{Token: token.Token{Type: token.SYMBOL, Literal: "lambda"}, Value: "lambda"},
				CdrField: &ast.ConsCell{
					CarField: &ast.Nil{Token: token.Token{Type: token.NIL, Literal: "nil"}},
					CdrField: &ast.ConsCell{
						CarField: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1},
						CdrField: &ast.Nil{Token: token.Token{Type: token.NIL, Literal: "nil"}},
					},
				},
			},
		},
		{
			name:  "lambda expression with one parameter",
			input: "(lambda (x) (+ x 1))",
			expected: &ast.ConsCell{
				CarField: &ast.Symbol{Token: token.Token{Type: token.SYMBOL, Literal: "lambda"}, Value: "lambda"},
				CdrField: &ast.ConsCell{
					CarField: &ast.ConsCell{
						CarField: &ast.Symbol{Token: token.Token{Type: token.SYMBOL, Literal: "x"}, Value: "x"},
						CdrField: &ast.Nil{Token: token.Token{Type: token.NIL, Literal: "nil"}},
					},
					CdrField: &ast.ConsCell{
						CarField: &ast.ConsCell{
							CarField: &ast.Symbol{Token: token.Token{Type: token.SYMBOL, Literal: "+"}, Value: "+"},
							CdrField: &ast.ConsCell{
								CarField: &ast.Symbol{Token: token.Token{Type: token.SYMBOL, Literal: "x"}, Value: "x"},
								CdrField: &ast.ConsCell{
									CarField: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1},
									CdrField: &ast.Nil{Token: token.Token{Type: token.NIL, Literal: "nil"}},
								},
							},
						},
						CdrField: &ast.Nil{Token: token.Token{Type: token.NIL, Literal: "nil"}},
					},
				},
			},
		},
		{
			name:  "lambda expression with multiple parameters",
			input: "(lambda (x y) (+ x y))",
			expected: &ast.ConsCell{
				CarField: &ast.Symbol{Token: token.Token{Type: token.SYMBOL, Literal: "lambda"}, Value: "lambda"},
				CdrField: &ast.ConsCell{
					CarField: &ast.ConsCell{
						CarField: &ast.Symbol{Token: token.Token{Type: token.SYMBOL, Literal: "x"}, Value: "x"},
						CdrField: &ast.ConsCell{
							CarField: &ast.Symbol{Token: token.Token{Type: token.SYMBOL, Literal: "y"}, Value: "y"},
							CdrField: &ast.Nil{Token: token.Token{Type: token.NIL, Literal: "nil"}},
						},
					},
					CdrField: &ast.ConsCell{
						CarField: &ast.ConsCell{
							CarField: &ast.Symbol{Token: token.Token{Type: token.SYMBOL, Literal: "+"}, Value: "+"},
							CdrField: &ast.ConsCell{
								CarField: &ast.Symbol{Token: token.Token{Type: token.SYMBOL, Literal: "x"}, Value: "x"},
								CdrField: &ast.ConsCell{
									CarField: &ast.Symbol{Token: token.Token{Type: token.SYMBOL, Literal: "y"}, Value: "y"},
									CdrField: &ast.Nil{Token: token.Token{Type: token.NIL, Literal: "nil"}},
								},
							},
						},
						CdrField: &ast.Nil{Token: token.Token{Type: token.NIL, Literal: "nil"}},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Expressions) != 1 {
				t.Fatalf("program.Expressions does not contain 1 expressions. got=%d", len(program.Expressions))
			}
			cc, ok := program.Expressions[0].(*ast.ConsCell)
			if !ok {
				t.Fatalf("exp not *ast.ConsCell. got=%T", program.Expressions[0])
			}
			if cc.String() != tt.expected.String() {
				t.Fatalf("cc.String() not %s. got=%s", tt.expected.String(), cc.String())
			}
		})
	}
}

func TestQuoteExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ast.SExpression
	}{
		{
			name:  "atom with quote",
			input: "'1",
			expected: &ast.ConsCell{
				CarField: &ast.Symbol{Token: token.Token{Type: token.QUOTE, Literal: "'"}, Value: "quote"},
				CdrField: &ast.ConsCell{
					CarField: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1},
					CdrField: &ast.Nil{Token: token.Token{Type: token.NIL, Literal: "nil"}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Expressions) != 1 {
				t.Fatalf("program.Expressions does not contain 1 expressions. got=%d", len(program.Expressions))
			}
			cc, ok := program.Expressions[0].(*ast.ConsCell)
			if !ok {
				t.Fatalf("exp not *ast.ConsCell. got=%T", program.Expressions[0])
			}
			if cc.String() != tt.expected.String() {
				t.Fatalf("cc.String() not %s. got=%s", tt.expected.String(), cc.String())
			}
		})
	}
}
