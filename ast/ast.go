package ast

import (
	"bytes"

	"github.com/JunNishimura/go-lisp/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Atom struct {
	Token token.Token
	Value interface{}
}

func (a *Atom) TokenLiteral() string { return a.Token.Literal }
func (a *Atom) String() string       { return a.Token.Literal }

type ConsCell struct {
	Token token.Token
	Car   Node
	Cdr   Node
}

func (cc *ConsCell) TokenLiteral() string { return cc.Token.Literal }
func (cc *ConsCell) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(cc.Car.String())
	out.WriteString(" . ")
	out.WriteString(cc.Cdr.String())
	out.WriteString(")")

	return out.String()
}

// type Statement interface {
// 	Node
// 	statementNode()
// }

// type Expression interface {
// 	Node
// 	expressionNode()
// }

// type Program struct {
// 	Statements []Statement
// }

// func (p *Program) TokenLiteral() string {
// 	if len(p.Statements) > 0 {
// 		return p.Statements[0].TokenLiteral()
// 	} else {
// 		return ""
// 	}
// }

// func (p *Program) String() string {
// 	var out bytes.Buffer

// 	for _, exp := range p.Statements {
// 		out.WriteString(exp.String())
// 	}

// 	return out.String()
// }

// type ExpressionStatement struct {
// 	Token token.Token
// 	Expression
// }

// func (es *ExpressionStatement) statementNode()       {}
// func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
// func (es *ExpressionStatement) String() string {
// 	if es.Expression != nil {
// 		return es.Expression.String()
// 	}
// 	return ""
// }

// type IntegerLiteral struct {
// 	Token token.Token
// 	Value int64
// }

// func (il *IntegerLiteral) expressionNode()      {}
// func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
// func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// type PrefixExpression struct {
// 	Token    token.Token
// 	Operator string
// 	Operands []Expression
// }

// func (pe *PrefixExpression) expressionNode()      {}
// func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
// func (pe *PrefixExpression) String() string {
// 	var out bytes.Buffer

// 	out.WriteString("(")
// 	out.WriteString(pe.Operator)
// 	for _, operand := range pe.Operands {
// 		out.WriteString(" ")
// 		out.WriteString(operand.String())
// 	}
// 	out.WriteString(")")

// 	return out.String()
// }
