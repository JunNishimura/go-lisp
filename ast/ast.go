package ast

import (
	"bytes"
	"fmt"

	"github.com/JunNishimura/go-lisp/token"
)

type SExpression interface {
	String() string
}

type Atom interface {
	SExpression
	TokenLiteral() string
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type PrefixAtom struct {
	Token    token.Token
	Operator string
	Right    Atom
}

func (pa *PrefixAtom) TokenLiteral() string { return pa.Token.Literal }
func (pa *PrefixAtom) String() string       { return fmt.Sprintf("%s%s", pa.Operator, pa.Right.String()) }

type Symbol struct {
	Token token.Token
	Value string
}

func (s *Symbol) TokenLiteral() string { return s.Token.Literal }
func (s *Symbol) String() string       { return s.Value }

type List struct {
	Car SExpression
	Cdr []SExpression
}

func (l *List) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(l.Car.String())

	for _, cdr := range l.Cdr {
		out.WriteString(" ")
		out.WriteString(cdr.String())
	}

	out.WriteString(")")
	return out.String()
}

type Program struct {
	Expressions []SExpression
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, exp := range p.Expressions {
		out.WriteString(exp.String())
	}

	return out.String()
}
