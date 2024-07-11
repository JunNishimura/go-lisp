package evaluator

import (
	"fmt"

	"github.com/JunNishimura/go-lisp/ast"
	"github.com/JunNishimura/go-lisp/object"
	"github.com/JunNishimura/go-lisp/token"
)

func evalBackquote(sexp ast.SExpression, env *object.Environment) object.Object {
	consCell, ok := sexp.(*ast.ConsCell)
	if !ok {
		return newError("cdr of backquote must be a cons cell, got %T", sexp)
	}

	unquoted := evalUnquote(consCell.Car(), env)

	return &object.Quote{
		SExpression: unquoted,
	}
}

func evalUnquote(sexp ast.SExpression, env *object.Environment) ast.SExpression {
	return ast.Modify(sexp, func(sexp ast.SExpression) ast.SExpression {
		consCell, ok := sexp.(*ast.ConsCell)
		if !ok {
			return sexp
		}

		if car, ok := consCell.Car().(*ast.Symbol); !ok || car.Value != "unquote" {
			return sexp
		}

		// evaluate unquoted s-expression
		cdr := consCell.Cdr().(*ast.ConsCell)
		evaluated := Eval(cdr.Car(), env)

		return convertObjectToSExpression(evaluated)
	}, "unquote")
}

func convertObjectToSExpression(obj object.Object) ast.SExpression {
	switch obj := obj.(type) {
	case *object.Integer:
		t := token.Token{
			Type:    token.INT,
			Literal: fmt.Sprintf("%d", obj.Value),
		}
		return &ast.IntegerLiteral{Token: t, Value: obj.Value}
	default:
		return nil
	}
}
