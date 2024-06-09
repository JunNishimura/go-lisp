package parser

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

// func TestExpressions(t *testing.T) {
// 	input := `
// (+ 1 2)
// (- 3 4)
// (* 5 6)
// (/ 7 8)
// `
// 	l := lexer.New(input)
// 	p := New(l)

// 	program := p.ParseProgram()
// 	if program == nil {
// 		t.Fatalf("ParseProgram() returned nil")
// 	}
// 	if len(program.Statements) != 4 {
// 		t.Fatalf("program.Statements does not contain 4 statements. got=%d", len(program.Statements))
// 	}
// }

func TestIntegerAtom(t *testing.T) {
	input := "1"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

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

// func TestIntegerLiteralExpression(t *testing.T) {
// 	input := "(1)"

// 	l := lexer.New(input)
// 	p := New(l)
// 	program := p.ParseProgram()

// 	if len(program.Statements) != 1 {
// 		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
// 	}
// 	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
// 	if !ok {
// 		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
// 	}

// 	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
// 	if !ok {
// 		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
// 	}
// 	if literal.Value != 1 {
// 		t.Fatalf("literal.Value not %d. got=%d", 1, literal.Value)
// 	}
// 	if literal.TokenLiteral() != "1" {
// 		t.Fatalf("literal.TokenLiteral not %s. got=%s", "1", literal.TokenLiteral())
// 	}
// }
