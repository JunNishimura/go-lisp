package object

import (
	"fmt"
)

const (
	ERROR_OBJ    = "ERROR"
	NIL_OBJ      = "NIL"
	INTEGER_OBJ  = "INTEGER"
	CONSCELL_OBJ = "CONSCELL"
)

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

type ConsCell struct {
	Car Object
	Cdr Object
}

func (c *ConsCell) Type() ObjectType { return CONSCELL_OBJ }
func (c *ConsCell) Inspect() string {
	return fmt.Sprintf("(%s . %s)", c.Car.Inspect(), c.Cdr.Inspect())
}
