package ast

type ModifierFun func(SExpression) SExpression

func Modify(node SExpression, modifier ModifierFun) SExpression {
	switch node := node.(type) {
	case *ConsCell:
		node.CarField = Modify(node.CarField, modifier)
		node.CdrField = Modify(node.CdrField, modifier)
	case *PrefixAtom:
		node.Right = Modify(node.Right, modifier)
	}

	return modifier(node)
}
