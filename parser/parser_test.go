package parser //

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
		expected string
	}{
		{
			name:     "parse positive integer",
			input:    "+1",
			expected: "+1",
		},
		{
			name:     "parse negative integer",
			input:    "-1",
			expected: "-1",
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
			if atom.String() != tt.expected {
				t.Fatalf("literal.Value not %s. got=%s", tt.expected, atom.String())
			}
		})
	}
}

func TestDottedPair(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *ast.ConsCell
	}{
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
				t.Fatalf("exp not *ast.List. got=%T", program.Expressions[0])
			}
			if cc.String() != tt.expected.String() {
				t.Fatalf("cc.String() not %s. got=%s", tt.expected.String(), cc.String())
			}
		})
	}
}
