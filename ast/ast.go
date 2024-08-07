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
	Right    SExpression
}

func (pa *PrefixAtom) TokenLiteral() string { return pa.Token.Literal }
func (pa *PrefixAtom) String() string       { return fmt.Sprintf("%s%s", pa.Operator, pa.Right.String()) }

type Symbol struct {
	Token token.Token
	Value string
}

func (s *Symbol) TokenLiteral() string { return s.Token.Literal }
func (s *Symbol) String() string       { return s.Value }

type True struct {
	Token token.Token
}

func (t *True) TokenLiteral() string { return t.Token.Literal }
func (t *True) String() string       { return "T" }

type SpecialForm struct {
	Token token.Token
	Value string
}

func (s *SpecialForm) TokenLiteral() string { return s.Token.Literal }
func (s *SpecialForm) String() string       { return s.Value }

type Nil struct {
	Token token.Token
}

func (n *Nil) TokenLiteral() string { return n.Token.Literal }
func (n *Nil) String() string       { return "NIL" }
func (n *Nil) Car() SExpression     { return n }
func (n *Nil) Cdr() SExpression     { return n }

type List interface {
	SExpression
	Car() SExpression
	Cdr() SExpression
}

type ConsCell struct {
	CarField SExpression
	CdrField SExpression
}

func (cc *ConsCell) String() string {
	var out bytes.Buffer

	out.WriteString("(")

	consCell := cc
	for {
		out.WriteString(consCell.Car().String())

		if _, ok := consCell.Cdr().(*Nil); ok {
			break
		}

		if cdr, ok := consCell.Cdr().(*ConsCell); ok {
			out.WriteString(" ")
			consCell = cdr
		} else {
			out.WriteString(" . ")
			out.WriteString(consCell.Cdr().String())
			break
		}
	}
	out.WriteString(")")

	return out.String()
}
func (cc *ConsCell) Car() SExpression { return cc.CarField }
func (cc *ConsCell) Cdr() SExpression { return cc.CdrField }

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
