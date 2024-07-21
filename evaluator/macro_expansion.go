package evaluator

import (
	"github.com/JunNishimura/go-lisp/ast"
	"github.com/JunNishimura/go-lisp/object"
)

var macroNames = []string{}

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

	macroNames = append(macroNames, macroName)
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

	return consCell.Car(), true
}

func ExpandMacros(program ast.SExpression, env *object.Environment) ast.SExpression {
	return ast.ModifyByMacro(program, func(sexp ast.SExpression) ast.SExpression {
		consCell, ok := sexp.(*ast.ConsCell)
		if !ok {
			return sexp
		}

		macro, ok := isMacroCall(consCell, env)
		if !ok {
			return sexp
		}

		args := quoteArgs(consCell)

		evalEnv := extendMacroEnv(macro, args)

		evaluated := Eval(macro.Body, evalEnv)

		quote, ok := evaluated.(*object.Quote)
		if !ok {
			panic("we only support returning AST-nodes from macros")
		}

		return quote.SExpression
	}, macroNames)
}

func isMacroCall(consCell *ast.ConsCell, env *object.Environment) (*object.Macro, bool) {
	symbol, ok := consCell.Car().(*ast.Symbol)

	if !ok {
		return nil, false
	}

	obj, ok := env.Get(symbol.Value)
	if !ok {
		return nil, false
	}

	macro, ok := obj.(*object.Macro)
	if !ok {
		return nil, false
	}

	return macro, true
}

func quoteArgs(consCell *ast.ConsCell) []*object.Quote {
	args := []*object.Quote{}

	consCell, ok := consCell.Cdr().(*ast.ConsCell)
	if !ok {
		return args
	}

	for {
		args = append(args, &object.Quote{SExpression: consCell.Car()})

		if _, ok := consCell.Cdr().(*ast.Nil); ok {
			break
		}

		consCell, ok = consCell.Cdr().(*ast.ConsCell)
		if !ok {
			return args
		}
	}

	return args
}

func extendMacroEnv(macro *object.Macro, args []*object.Quote) *object.Environment {
	env := object.NewEnclosedEnvironment(macro.Env)

	for i, param := range macro.Parameters {
		env.Set(param.Value, args[i])
	}

	return env
}
