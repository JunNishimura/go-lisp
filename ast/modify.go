package ast

import (
	"slices"
)

type ModifierFun func(SExpression) SExpression

func Modify(sexp SExpression, modifier ModifierFun, triggers []string) SExpression {
	switch st := sexp.(type) {
	case *Program:
		for i, sexp := range st.Expressions {
			st.Expressions[i] = Modify(sexp, modifier, triggers)
		}
	case *ConsCell:
		if symbol, ok := st.CarField.(*Symbol); ok && slices.Contains(triggers, symbol.Value) {
			// return not only the car field but also the cdr field
			// since args(cdr field) are needed to modify the AST
			return modifier(sexp)
		}
		st.CarField = Modify(st.CarField, modifier, triggers)
		st.CdrField = Modify(st.CdrField, modifier, triggers)
	}

	return sexp
}
