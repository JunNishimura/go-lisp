package evaluator

import (
	"testing"

	"github.com/JunNishimura/go-lisp/ast"
	"github.com/JunNishimura/go-lisp/lexer"
	"github.com/JunNishimura/go-lisp/object"
	"github.com/JunNishimura/go-lisp/parser"
)

func TestDefineMacros(t *testing.T) {
	input := "(defmacro myMacro (x y) '(+ x y))"

	env := object.NewEnvironment()
	program := testParseProgram(input)

	DefineMacros(program, env)

	obj, ok := env.Get("myMacro")
	if !ok {
		t.Fatalf("macro not in environment")
	}

	macro, ok := obj.(*object.Macro)
	if !ok {
		t.Fatalf("object is not Macro. got=%T (%+v)", obj, obj)
	}

	if len(macro.Parameters) != 2 {
		t.Fatalf("wrong number of identifiers. got=%d", len(macro.Parameters))
	}

	if macro.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", macro.Parameters[0])
	}
	if macro.Parameters[1].String() != "y" {
		t.Fatalf("parameter is not 'y'. got=%q", macro.Parameters[1])
	}

	expectedBody := "(quote (+ x y))"
	if macro.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, macro.Body.String())
	}
}

func testParseProgram(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func TestExpandMacros(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "expands macro which has no arguments",
			input: `
				(defmacro hoge () '1)
				(hoge)
			`,
			expected: "1",
		},
		{
			name: "expands macro which has arguments",
			input: `
				(defmacro hoge (x y) ` + "`" + `(- ,y ,x))
				(hoge (+ 2 2) (- 10 5))
			`,
			expected: "(- (- 10 5) (+ 2 2))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := testParseProgram(tt.expected)

			program := testParseProgram(tt.input)
			env := object.NewEnvironment()

			DefineMacros(program, env)

			expanded := ExpandMacros(program, env)

			if expanded.String() != expected.String() {
				t.Errorf("not equal. got=%q, want=%q", expanded.String(), expected.String())
			}
		})
	}
}
