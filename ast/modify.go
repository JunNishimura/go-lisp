package ast

type ModifierFun func(SExpression) SExpression

func Modify(sexp SExpression, modifier ModifierFun, symbolValue string) SExpression {
	switch st := sexp.(type) {
	case *ConsCell:
		if car, ok := st.CarField.(*Symbol); ok && car.Value == symbolValue {
			return modifier(sexp)
		}
		st.CarField = Modify(st.CarField, modifier, symbolValue)
		st.CdrField = Modify(st.CdrField, modifier, symbolValue)
	}

	return sexp
}
