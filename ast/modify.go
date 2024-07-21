package ast

import (
	"slices"
)

type ModifierFun func(SExpression) SExpression

func ModifyByUnquote(sexp SExpression, modifier ModifierFun) SExpression {
	return modify(sexp, modifier, func(sexp SExpression) bool {
		if spForm, ok := sexp.(*SpecialForm); ok {
			return spForm.Value == "unquote"
		}
		return false
	})
}

func ModifyByMacro(sexp SExpression, modifier ModifierFun, macroNames []string) SExpression {
	return modify(sexp, modifier, func(sexp SExpression) bool {
		if symbol, ok := sexp.(*Symbol); ok {
			return slices.Contains(macroNames, symbol.Value)
		}
		return false
	})
}

func modify(sexp SExpression, modifier ModifierFun, targetCond func(sexp SExpression) bool) SExpression {
	switch st := sexp.(type) {
	case *Program:
		for i, sexp := range st.Expressions {
			st.Expressions[i] = modify(sexp, modifier, targetCond)
		}
	case *ConsCell:
		if targetCond(st.Car()) {
			// return not only the car field but also the cdr field
			// since args(cdr field) are needed to modify the AST
			return modifier(sexp)
		}
		st.CarField = modify(st.CarField, modifier, targetCond)
		st.CdrField = modify(st.CdrField, modifier, targetCond)
	}

	return sexp
}
