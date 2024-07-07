package evaluator

import "testing"

func TestLambdaExpression(t *testing.T) {
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
