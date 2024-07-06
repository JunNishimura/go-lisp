package evaluator

import (
	"fmt"

	"github.com/JunNishimura/go-lisp/ast"
	"github.com/JunNishimura/go-lisp/object"
)

var specialForms = map[string]*object.SpecialForm{
	"lambda": {
		Fn: func(sexp ast.SExpression, env *object.Environment) object.Object {
			consCell, ok := sexp.(*ast.ConsCell)
			if !ok {
				return newError("cdr of lambda must be a cons cell, got %T", sexp)
			}

			params, err := evalLambdaParams(consCell.Car())
			if err != nil {
				return newError(err.Error())
			}

			cdr := consCell.Cdr()
			consCell, ok = cdr.(*ast.ConsCell)
			if !ok {
				return newError("cdr of lambda must be a cons cell, got %T", cdr)
			}
			body := consCell.Car()

			return &object.Function{Parameters: params, Body: body, Env: env}
		},
	},
}

func evalLambdaParams(sexp ast.SExpression) ([]*ast.Symbol, error) {
	params := []*ast.Symbol{}

	if _, ok := sexp.(*ast.Nil); ok {
		return params, nil
	}

	consCell, ok := sexp.(*ast.ConsCell)
	if !ok {
		return nil, fmt.Errorf("parameters must be a list, got %T", sexp)
	}

	for {
		symbol, ok := consCell.Car().(*ast.Symbol)
		if !ok {
			return nil, fmt.Errorf("parameter must be a symbol, got %T", consCell.Car())
		}
		params = append(params, symbol)

		switch cdr := consCell.Cdr().(type) {
		case *ast.Nil:
			return params, nil
		case *ast.ConsCell:
			consCell = cdr
		default:
			return nil, fmt.Errorf("parameters must be a list, got %T", consCell.Cdr())
		}
	}
}
