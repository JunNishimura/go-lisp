package object

import (
	"bytes"
	"fmt"

	"github.com/JunNishimura/go-lisp/ast"
)

const (
	ERROR_OBJ    = "ERROR"
	NIL_OBJ      = "NIL"
	INTEGER_OBJ  = "INTEGER"
	FUNCTION_OBJ = "FUNCTION"
	SYMBOL_OBJ   = "SYMBOL"
	CONSCELL_OBJ = "CONSCELL"
	LIST_OBJ     = "LIST"
	BUILTIN_OBJ  = "BUILTIN"
)

type BuiltInFunction func(args ...Object) Object

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type Nil struct{}

func (n *Nil) Type() ObjectType { return NIL_OBJ }
func (n *Nil) Inspect() string  { return "nil" }

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Symbol struct {
	Name         string
	Value        Object
	PropertyList []Object
	Function     *Function
}

func (s *Symbol) Type() ObjectType { return SYMBOL_OBJ }
func (s *Symbol) Inspect() string  { return s.Value.Inspect() }

type Function struct {
	Parameters []*ast.Symbol
	Body       ast.SExpression
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	out.WriteString("(lambda (")
	for i, p := range f.Parameters {
		if i > 0 {
			out.WriteString(" ")
		}
		out.WriteString(p.String())
	}
	out.WriteString(") ")
	out.WriteString(f.Body.String())
	out.WriteString(")")

	return out.String()
}

type Builtin struct {
	Fn BuiltInFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

type ConsCell struct {
	Car Object
	Cdr Object
}

func (cc *ConsCell) Type() ObjectType { return CONSCELL_OBJ }
func (cc *ConsCell) Inspect() string {
	return fmt.Sprintf("(%s . %s)", cc.Car.Inspect(), cc.Cdr.Inspect())
}

type List struct {
	SExpressions []Object
}

func (l *List) Type() ObjectType { return LIST_OBJ }
func (l *List) Inspect() string {
	var out bytes.Buffer

	out.WriteString("(")
	for i, s := range l.SExpressions {
		if i > 0 {
			out.WriteString(" ")
		}
		out.WriteString(s.Inspect())
	}
	out.WriteString(")")

	return out.String()
}
