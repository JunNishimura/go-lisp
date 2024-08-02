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
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if ok {
		if result.Value != expected {
			t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
			return false
		}
		return true
	}

	symbol, ok := obj.(*object.Symbol)
	if ok {
		integer, ok := symbol.Value.(*object.Integer)
		if !ok {
			t.Errorf("object is not Integer. got=%T (%+v)", result.Value, result.Value)
			return false
		}
		if integer.Value != expected {
			t.Errorf("object has wrong value. got=%d, want=%d", integer.Value, expected)
			return false
		}
		return true
	}

	t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
	return false
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
		{"(+ . (1 . (2 . nil)))", 3},
		{"(+ . (1 2))", 3},
		{"(+ 5 5)", 10},
		{"(- 5 5)", 0},
		{"(* 5 5)", 25},
		{"(/ 5 5)", 1},
		{"(+ -5 5)", 0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestLambda(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"((lambda () 5))", 5},
		{"((lambda (x) x) 5)", 5},
		{"((lambda (x y) (+ x y)) 5 5)", 10},
		{"(+ ((lambda () 1)) ((lambda (x y) (+ x y)) 1 2))", 4},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestQuote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"'5", "5"},
		{"'-5", "-5"},
		{"'(+ 1 2)", "(+ 1 2)"},
		{"'(+ . (1 . (2 . nil)))", "(+ 1 2)"},
		{"(quote 5)", "5"},
		{"(quote -5)", "-5"},
		{"(quote (+ 1 2))", "(+ 1 2)"},
		{"(quote (+ . (1 . (2 . nil))))", "(+ 1 2)"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if evaluated.Inspect() != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, evaluated.Inspect())
		}
	}
}

func TestBackQuote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"`5", "5"},
		{"`(+ 1 2)", "(+ 1 2)"},
		{"`(+ 1 ,(+ 1 1))", "(+ 1 2)"},
		{"`(+ ,((lambda () 1)) 2)", "(+ 1 2)"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if evaluated.Inspect() != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, evaluated.Inspect())
		}
	}
}

func TestTrueExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"t", "T"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if evaluated.Inspect() != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, evaluated.Inspect())
		}
	}
}

func testComparisonObject(t *testing.T, obj object.Object, expected string) {
	if expected == "T" {
		if obj != True {
			t.Errorf("object is not TRUE. got=%T (%+v)", obj, obj)
		}
	} else {
		if obj != Nil {
			t.Errorf("object is not NIL. got=%T (%+v)", obj, obj)
		}
	}
}

func TestComparisonExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"(= 1)", "T"},
		{"(= 1 1)", "T"},
		{"(= 1 2)", "NIL"},
		{"(/= 1)", "T"},
		{"(/= 1 1)", "NIL"},
		{"(/= 1 2)", "T"},
		{"(< 1)", "T"},
		{"(< 1 1)", "NIL"},
		{"(< 1 2)", "T"},
		{"(< 2 1)", "NIL"},
		{"(> 1)", "T"},
		{"(> 1 1)", "NIL"},
		{"(> 1 2)", "NIL"},
		{"(> 2 1)", "T"},
		{"(<= 1)", "T"},
		{"(<= 1 0)", "NIL"},
		{"(<= 1 1)", "T"},
		{"(<= 1 2)", "T"},
		{"(<= 1)", "T"},
		{"(<= 2 1)", "NIL"},
		{"(<= 2 2)", "T"},
		{"(<= 2 3)", "T"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testComparisonObject(t, evaluated, tt.expected)
	}
}

func TestIfExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"(if t 10)", "10"},
		{"(if 0 10)", "10"},
		{"(if 1 10)", "10"},
		{"(if nil 10)", "nil"},
		{"(if t 10 20)", "10"},
		{"(if nil 10 20)", "20"},
		{"(if (= 1 1) 10 20)", "10"},
		{"(if (= 1 2) 10 20)", "20"},
		{"(if t (+ 1 1) (+ 2 2))", "2"},
		{"(if nil (+ 1 1) (+ 2 2))", "4"},
		{"(if (= ((lambda (x y) (+ x y)) 1 1) (* 1 2)) 10 20)", "10"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if evaluated.Inspect() != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, evaluated.Inspect())
		}
	}
}

func TestSetqExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"(setq x 10) x", 10},
		{"(setq x 10) (setq x 20) x", 20},
		{"(setq x 10) (setq y 20) (+ x y)", 30},
		{"(setq x (+ 1 1)) x", 2},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestApplyExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"(setq f (lambda () (+ 1 1))) (apply f ())", 2},
		{"(setq f (lambda (x) (+ x x))) (apply f '(1))", 2},
		{"(setq f (lambda (x y) (+ x y))) (apply f '(1 1))", 2},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}
