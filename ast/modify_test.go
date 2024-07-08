package ast

import (
	"reflect"
	"testing"
)

func TestModify(t *testing.T) {
	turnOneIntoTwo := func(sexp SExpression) SExpression {
		integer, ok := sexp.(*IntegerLiteral)
		if !ok {
			return sexp
		}

		if integer.Value == 1 {
			integer.Value = 2
		}

		return integer
	}

	tests := []struct {
		name     string
		input    SExpression
		expected SExpression
	}{
		{
			name:     "integer literal",
			input:    &IntegerLiteral{Value: 1},
			expected: &IntegerLiteral{Value: 2},
		},
		{
			name: "prefix atom",
			input: &PrefixAtom{
				Operator: "+",
				Right:    &IntegerLiteral{Value: 1},
			},
			expected: &PrefixAtom{
				Operator: "+",
				Right:    &IntegerLiteral{Value: 2},
			},
		},
		{
			name: "cons cell",
			input: &ConsCell{
				CarField: &IntegerLiteral{Value: 1},
				CdrField: &IntegerLiteral{Value: 1},
			},
			expected: &ConsCell{
				CarField: &IntegerLiteral{Value: 2},
				CdrField: &IntegerLiteral{Value: 2},
			},
		},
		{
			name: "nested cons cell",
			input: &ConsCell{
				CarField: &IntegerLiteral{Value: 1},
				CdrField: &ConsCell{
					CarField: &IntegerLiteral{Value: 1},
					CdrField: &Nil{},
				},
			},
			expected: &ConsCell{
				CarField: &IntegerLiteral{Value: 2},
				CdrField: &ConsCell{
					CarField: &IntegerLiteral{Value: 2},
					CdrField: &Nil{},
				},
			},
		},
	}

	for _, tt := range tests {
		modified := Modify(tt.input, turnOneIntoTwo)
		if isEqual := reflect.DeepEqual(modified, tt.expected); !isEqual {
			t.Errorf("tests[%s] - not equal. expected=%+v, got=%+v",
				tt.name, tt.expected, modified)
		}
	}
}
