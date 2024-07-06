package evaluator

import (
	"fmt"

	"github.com/JunNishimura/go-lisp/ast"
	"github.com/JunNishimura/go-lisp/object"
)

var (
	Nil = &object.Nil{}
)

func Eval(sexp ast.SExpression, env *object.Environment) object.Object {
	switch sexp := sexp.(type) {
	case *ast.Program:
		return evalProgram(sexp, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: sexp.Value}
	case *ast.PrefixAtom:
		right := Eval(sexp.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixAtom(sexp.Operator, right)
	case *ast.Nil:
		return Nil
	case *ast.Symbol:
		return evalSymbol(sexp, env)
	case *ast.ConsCell:
		return evalList(sexp, env)
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

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, exp := range program.Expressions {
		result = Eval(exp, env)

		switch result := result.(type) {
		case *object.Error:
			return result
		}
	}

	return result
}

func evalPrefixAtom(operator string, right object.Object) object.Object {
	switch operator {
	case "+":
		return evalPlusPrefix(right)
	case "-":
		return evalMinusPrefix(right)
	default:
		return newError("unknown operator: %s", operator)
	}
}

func evalPlusPrefix(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: + %s", right.Type())
	}

	return right
}

func evalMinusPrefix(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: - %s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalSymbol(symbol *ast.Symbol, env *object.Environment) object.Object {
	if val, ok := env.Get(symbol.Value); ok {
		return val
	}

	if builtin, ok := builtinFuncs[symbol.Value]; ok {
		return builtin
	}

	return newError("symbol not found: %s", symbol.Value)
}

// evaluate cdr of the cons cell as arguments to the command car
func evalList(sexp *ast.ConsCell, env *object.Environment) object.Object {
	// check if the car is a special form
	if specialForm, ok := sexp.Car().(*ast.SpecialForm); ok {
		return specialForms[specialForm.TokenLiteral()].Fn(sexp.Cdr(), env)
	}

	// Evaluate the car of the cons cell
	car := Eval(sexp.Car(), env)
	if isError(car) {
		return car
	}

	// Evaluate the arguments
	args := evalArgs(sexp.Cdr(), env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}
	return applyFunction(car, args)
}

func evalArgs(sexp ast.SExpression, env *object.Environment) []object.Object {
	list := []object.Object{}

	switch sexp := sexp.(type) {
	case *ast.Nil:
		return list
	case *ast.ConsCell:
		return evalValueList(sexp, env)
	default:
		return []object.Object{newError("arguments must be a list, got %T", sexp)}
	}
}

func evalValueList(consCell *ast.ConsCell, env *object.Environment) []object.Object {
	list := []object.Object{}

	for {
		// Evaluate the car of the cons cell
		car := Eval(consCell.Car(), env)
		if isError(car) {
			return []object.Object{car}
		}
		list = append(list, car)

		// move to the next cons cell or return the list if the cdr is nil
		switch cdr := consCell.Cdr().(type) {
		case *ast.Nil:
			return list
		case *ast.ConsCell:
			consCell = cdr
		default:
			err := newError("arguments must be a list, got %T", consCell.Cdr())
			return []object.Object{err}
		}
	}
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv, err := extendFunctionEnv(fn, args)
		if err != nil {
			return newError(err.Error())
		}
		return Eval(fn.Body, extendedEnv)
	case *object.Symbol:
		extendedEnv, err := extendFunctionEnv(fn.Function, args)
		if err != nil {
			return newError(err.Error())
		}
		symbolFunc := fn.Function
		return Eval(symbolFunc.Body, extendedEnv)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) (*object.Environment, error) {
	if len(fn.Parameters) != len(args) {
		return nil, fmt.Errorf("function expects %d arguments, but got %d", len(fn.Parameters), len(args))
	}

	env := object.NewEnclosedEnvironment(fn.Env)

	for i, param := range fn.Parameters {
		env.Set(param.Value, args[i])
	}

	return env, nil
}
