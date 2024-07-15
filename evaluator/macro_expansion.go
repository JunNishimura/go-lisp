package evaluator

import (
	"github.com/JunNishimura/go-lisp/ast"
	"github.com/JunNishimura/go-lisp/object"
)

func DefineMacros(program *ast.Program, env *object.Environment) {
	definitions := []int{}

	for i, exp := range program.Expressions {
		if isMacroDefinition(exp) {
			addMacro(exp, env)
			definitions = append(definitions, i)
		}
	}

	for i := len(definitions) - 1; i >= 0; i-- {
		definitionIndex := definitions[i]
		program.Expressions = append(program.Expressions[:definitionIndex], program.Expressions[definitionIndex+1:]...)
	}
}

func isMacroDefinition(sexp ast.SExpression) bool {
	consCell, ok := sexp.(*ast.ConsCell)
	if !ok {
		return false
	}

	car, ok := consCell.Car().(*ast.Symbol)
	if !ok {
		return false
	}

	return car.Value == "defmacro"
}

func addMacro(sexp ast.SExpression, env *object.Environment) {
	macroName, ok := getMacroName(sexp)
	if !ok {
		return
	}

	params, ok := getMacroParams(sexp)
	if !ok {
		return
	}

	body, ok := getMacroBody(sexp)
	if !ok {
		return
	}

	macro := &object.Macro{
		Parameters: params,
		Body:       body,
		Env:        env,
	}

	env.Set(macroName, macro)
}

func getMacroName(sexp ast.SExpression) (string, bool) {
	consCell, ok := sexp.(*ast.ConsCell)
	if !ok {
		return "", false
	}

	consCell, ok = consCell.Cdr().(*ast.ConsCell)
	if !ok {
		return "", false
	}

	symbol, ok := consCell.Car().(*ast.Symbol)
	if !ok {
		return "", false
	}

	return symbol.Value, true
}

func getMacroParams(sexp ast.SExpression) ([]*ast.Symbol, bool) {
	consCell, ok := sexp.(*ast.ConsCell)
	if !ok {
		return nil, false
	}

	consCell, ok = consCell.Cdr().(*ast.ConsCell)
	if !ok {
		return nil, false
	}

	consCell, ok = consCell.Cdr().(*ast.ConsCell)
	if !ok {
		return nil, false
	}

	paramConsCell, ok := consCell.Car().(*ast.ConsCell)
	if !ok {
		if _, ok := consCell.Car().(*ast.Nil); ok {
			// No parameters
			return []*ast.Symbol{}, true
		}
		return nil, false
	}

	params := []*ast.Symbol{}
	for {
		symbol, ok := paramConsCell.Car().(*ast.Symbol)
		if !ok {
			return nil, false
		}

		params = append(params, symbol)

		if _, ok := paramConsCell.Cdr().(*ast.Nil); ok {
			break
		}

		paramConsCell, ok = paramConsCell.Cdr().(*ast.ConsCell)
		if !ok {
			return nil, false
		}
	}

	return params, true
}

func getMacroBody(sexp ast.SExpression) (ast.SExpression, bool) {
	consCell, ok := sexp.(*ast.ConsCell)
	if !ok {
		return nil, false
	}

	consCell, ok = consCell.Cdr().(*ast.ConsCell)
	if !ok {
		return nil, false
	}

	consCell, ok = consCell.Cdr().(*ast.ConsCell)
	if !ok {
		return nil, false
	}

	consCell, ok = consCell.Cdr().(*ast.ConsCell)
	if !ok {
		return nil, false
	}

	body, ok := consCell.Car().(ast.SExpression)
	if !ok {
		return nil, false
	}

	return body, true
}
