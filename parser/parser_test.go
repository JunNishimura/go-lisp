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

func TestConsCell(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *ast.ConsCell
	}{
		{
			name:  "cons cell in which car and cdr are both atoms",
			input: "(cons 1 2)",
			expected: &ast.ConsCell{
				Operator: token.Token{Type: token.IDENT, Literal: "cons"},
				Car:      &ast.Atom[int64]{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1},
				Cdr:      &ast.Atom[int64]{Token: token.Token{Type: token.INT, Literal: "2"}, Value: 2},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.SExpressions) != 1 {
				t.Fatalf("program.SExpressions does not contain 1 expressions. got=%d", len(program.SExpressions))
			}

			cc, ok := program.SExpressions[0].(*ast.ConsCell)
			if !ok {
				t.Fatalf("exp not *ast.ConsCell. got=%T", program.SExpressions[0])
			}

			if cc.Operator.Literal != tt.expected.Operator.Literal {
				t.Fatalf("cc.Operator.Literal not %s. got=%s", tt.expected.Operator.Literal, cc.Operator.Literal)
			}

			if cc.String() != tt.expected.String() {
				t.Fatalf("cc.String() not %s. got=%s", tt.expected.String(), cc.String())
			}
		})
	}
}

func TestIntegerAtom(t *testing.T) {
	input := "1"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.SExpressions) != 1 {
		t.Fatalf("program.SExpressions does not contain 1 expressions. got=%d", len(program.SExpressions))
	}
	atom, ok := program.SExpressions[0].(*ast.Atom[int64])
	if !ok {
		t.Fatalf("exp not *ast.Atom[int64]. got=%T", program.SExpressions[0])
	}
	if atom.Value != 1 {
		t.Fatalf("literal.Value not %d. got=%d", 1, atom.Value)
	}
	if atom.TokenLiteral() != "1" {
		t.Fatalf("literal.TokenLiteral not %s. got=%s", "1", atom.TokenLiteral())
	}
}
