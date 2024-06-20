package evaluator

import (
	"testing"

	"github.com/JunNishimura/go-lisp/lexer"
	"github.com/JunNishimura/go-lisp/object"
	"github.com/JunNishimura/go-lisp/parser"
)

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"+5", 5},
		{"-5", -5},
		// {"(+ 5 5)", 10},
		// {"(- 5 5)", 0},
		// {"(* 5 5)", 25},
		// {"(/ 5 5)", 1},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

// func TestEvalNilExpression(t *testing.T) {
// 	input := "nil"

// 	evaluated := testEval(input)
// 	if evaluated != Nil {
// 		t.Errorf("object is not Nil. got=%T (%+v)", evaluated, evaluated)
// 	}
// }
