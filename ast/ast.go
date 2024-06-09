package ast

import (
	"bytes"

	"github.com/JunNishimura/go-lisp/token"
)

type Cell interface {
	TokenLiteral() string
	String() string
}

type AtomValue interface {
	int64
}

type Atom[T AtomValue] struct {
	Token token.Token
	Value T
}

func (a *Atom[T]) TokenLiteral() string { return a.Token.Literal }
func (a *Atom[T]) String() string       { return a.Token.Literal }

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
	out.WriteString(" . ")
	out.WriteString(cc.Cdr.String())
	out.WriteString(")")

	return out.String()
}

type Program struct {
	SExpressions []Cell
}
