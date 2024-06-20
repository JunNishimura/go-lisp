package evaluator

import (
	"fmt"

	"github.com/JunNishimura/go-lisp/ast"
	"github.com/JunNishimura/go-lisp/object"
)

var (
	Nil = &object.Nil{}
)

func Eval(sexp ast.SExpression) object.Object {
	switch sexp := sexp.(type) {
	case *ast.Program:
		return evalProgram(sexp)
	case *ast.NilLiteral:
		return Nil
	case *ast.IntegerLiteral:
		return &object.Integer{Value: sexp.Value}
	case *ast.DottedPair:
		car := Eval(sexp.CarCell)
		if isError(car) {
			return car
		}
		cdr := Eval(sexp.CdrCell)
		if isError(cdr) {
			return cdr
		}
		return evalDottedPair(sexp)
	default:
		return newError("unknown expression type: %T", sexp)
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, exp := range program.Expressions {
		result = Eval(exp)

		switch result := result.(type) {
		case *object.Error:
			return result
		}
	}

	return result
}
