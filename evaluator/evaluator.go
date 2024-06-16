package evaluator

import (
	"fmt"

	"github.com/JunNishimura/go-lisp/ast"
	"github.com/JunNishimura/go-lisp/object"
)

func Eval(sexp ast.SExpression) object.Object {
	switch sexp := sexp.(type) {
	case *ast.Program:
		return evalProgram(sexp)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: sexp.Value}
	default:
		return newError("unknown expression type: %T", sexp)
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
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
