package ast

import (
	"bytes"
	"strings"

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

type NilLiteral struct {
	Token token.Token
}

func (nl *NilLiteral) TokenLiteral() string { return nl.Token.Literal }
func (nl *NilLiteral) String() string       { return strings.ToUpper(nl.Token.Literal) }

type ConsLiteral struct {
	Token token.Token
}

func (cl *ConsLiteral) TokenLiteral() string { return cl.Token.Literal }
func (cl *ConsLiteral) String() string       { return "" }

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type ConsCell interface {
	SExpression
	Car() SExpression
	Cdr() SExpression
}

type DottedPair struct {
	CarCell SExpression
	CdrCell SExpression
}

func (d *DottedPair) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(d.CarCell.String())
	if _, ok := d.CdrCell.(*NilLiteral); !ok {
		out.WriteString(" . ")
		out.WriteString(d.CdrCell.String())
	}
	out.WriteString(")")

	return out.String()
}

func (d *DottedPair) Car() SExpression { return d.CarCell }
func (d *DottedPair) Cdr() SExpression { return d.CdrCell }

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
