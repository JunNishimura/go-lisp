package ast

import (
	"bytes"

	"github.com/JunNishimura/go-lisp/token"
)

type Cell interface {
	TokenLiteral() string
	String() string
}

type Atom interface {
	Cell
	Atom()
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) Atom()                {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type NilLiteral struct {
	Token token.Token
}

func (nl *NilLiteral) Atom()                {}
func (nl *NilLiteral) TokenLiteral() string { return nl.Token.Literal }
func (nl *NilLiteral) String() string       { return nl.Token.Literal }

type ConsCell struct {
	Operator token.Token
	Car      Cell
	Cdr      Cell
}

func (cc *ConsCell) TokenLiteral() string { return cc.Operator.Literal }
func (cc *ConsCell) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(cc.Car.String())
	if _, ok := cc.Cdr.(*NilLiteral); !ok {
		out.WriteString(" . ")
		out.WriteString(cc.Cdr.String())
	}
	out.WriteString(")")

	return out.String()
}

type Program struct {
	SExpressions []Cell
}
