package parser //

import (
	"testing"

	"github.com/JunNishimura/go-lisp/ast"
	"github.com/JunNishimura/go-lisp/lexer"
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

// func TestConsCell(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		input    string
// 		expected ast.ConsCell
// 	}{
// 		{
// 			name:  "cons cell in which car and cdr are both atoms",
// 			input: "(cons 1 2)",
// 			expected: &ast.DottedPair{
// 				CarCell: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1},
// 				CdrCell: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "2"}, Value: 2},
// 			},
// 		},
// 		{
// 			name:  "cons cell in which car is an atom and cdr is nil",
// 			input: "(cons 1 nil)",
// 			expected: &ast.DottedPair{
// 				CarCell: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1},
// 				CdrCell: &ast.NilLiteral{Token: token.Token{Type: token.NIL, Literal: "NIL"}},
// 			},
// 		},
// 		{
// 			name:  "cons cell in which car is nil and cdr is an atom",
// 			input: "(cons nil 2)",
// 			expected: &ast.DottedPair{
// 				CarCell: &ast.NilLiteral{Token: token.Token{Type: token.NIL, Literal: "NIL"}},
// 				CdrCell: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "2"}, Value: 2},
// 			},
// 		},
// 		{
// 			name:  "cons cell in which operator is an math operator(+)",
// 			input: "(+ 1 2)",
// 			expected: &ast.DottedPair{
// 				CarCell: &ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "+"}, Value: "+"},
// 				CdrCell: &ast.DottedPair{
// 					CarCell: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1},
// 					CdrCell: &ast.DottedPair{
// 						CarCell: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "2"}, Value: 2},
// 						CdrCell: &ast.NilLiteral{Token: token.Token{Type: token.NIL, Literal: "NIL"}},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name:  "cons cell in which operator is an math operator(-)",
// 			input: "(- 1 2)",
// 			expected: &ast.DottedPair{
// 				CarCell: &ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "-"}, Value: "-"},
// 				CdrCell: &ast.DottedPair{
// 					CarCell: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1},
// 					CdrCell: &ast.DottedPair{
// 						CarCell: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "2"}, Value: 2},
// 						CdrCell: &ast.NilLiteral{Token: token.Token{Type: token.NIL, Literal: "NIL"}},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name:  "cons cell in which operator is an math operator(*)",
// 			input: "(* 1 2)",
// 			expected: &ast.DottedPair{
// 				CarCell: &ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "*"}, Value: "*"},
// 				CdrCell: &ast.DottedPair{
// 					CarCell: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1},
// 					CdrCell: &ast.DottedPair{
// 						CarCell: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "2"}, Value: 2},
// 						CdrCell: &ast.NilLiteral{Token: token.Token{Type: token.NIL, Literal: "NIL"}},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name:  "cons cell in which operator is an math operator(/)",
// 			input: "(/ 1 2)",
// 			expected: &ast.DottedPair{
// 				CarCell: &ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "/"}, Value: "/"},
// 				CdrCell: &ast.DottedPair{
// 					CarCell: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "1"}, Value: 1},
// 					CdrCell: &ast.DottedPair{
// 						CarCell: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "2"}, Value: 2},
// 						CdrCell: &ast.NilLiteral{Token: token.Token{Type: token.NIL, Literal: "NIL"}},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			l := lexer.New(tt.input)
// 			p := New(l)
// 			program := p.ParseProgram()
// 			checkParserErrors(t, p)

// 			if len(program.Expressions) != 1 {
// 				t.Fatalf("program.Expressions does not contain 1 expressions. got=%d", len(program.Expressions))
// 			}

// 			cc, ok := program.Expressions[0].(*ast.DottedPair)
// 			if !ok {
// 				t.Fatalf("exp not *ast.DottedPair. got=%T", program.Expressions[0])
// 			}

// 			if cc.Car().String() != tt.expected.Car().String() {
// 				t.Fatalf("cc.Car() not %s. got=%s", tt.expected.Car().String(), cc.Car().String())
// 			}
// 			if cc.Cdr().String() != tt.expected.Cdr().String() {
// 				t.Fatalf("cc.Cdr() not %s. got=%s", tt.expected.Cdr().String(), cc.Cdr().String())
// 			}
// 		})
// 	}
// }

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
