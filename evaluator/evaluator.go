package evaluator

import (
	"fmt"

	"github.com/JunNishimura/go-lisp/ast"
	"github.com/JunNishimura/go-lisp/object"
	"github.com/JunNishimura/go-lisp/token"
)

var (
	Nil  = &object.Nil{}
	True = &object.True{}
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
	case *ast.True:
		return True
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
	switch car := sexp.Car().(type) {
	case *ast.Symbol:
		return evalNormalForm(sexp, env)
	case *ast.ConsCell:
		if isLambdaExpression(car) {
			return evalNormalForm(sexp, env)
		}
	case *ast.SpecialForm:
		return evalSpecialForm(sexp, env)
	}

	return newError("unknown operator type: %T", sexp.Car())
}

func isLambdaExpression(consCell *ast.ConsCell) bool {
	spForm, ok := consCell.Car().(*ast.SpecialForm)
	if !ok {
		return false
	}

	return spForm.Token.Type == token.LAMBDA
}

func evalNormalForm(consCell *ast.ConsCell, env *object.Environment) object.Object {
	// Evaluate the car of the cons cell
	car := Eval(consCell.Car(), env)
	if isError(car) {
		return car
	}

	// Evaluate the arguments
	args := evalArgs(consCell.Cdr(), env)
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
		if symbol, ok := car.(*object.Symbol); ok {
			list = append(list, symbol.Value)
		} else {
			list = append(list, car)
		}

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

func evalSpecialForm(sexp *ast.ConsCell, env *object.Environment) object.Object {
	spForm, ok := sexp.Car().(*ast.SpecialForm)
	if !ok {
		return newError("expect special form, got %T", sexp.Car())
	}

	switch spForm.Value {
	case "lambda":
		return evalLambda(sexp, env)
	case "quote":
		return evalQuote(sexp)
	case "backquote":
		return evalBackquote(sexp, env)
	case "if":
		return evalIf(sexp, env)
	case "setq":
		return evalSetq(sexp, env)
	}

	return newError("unknown special form: %s", spForm.Value)
}

func evalLambda(sexp *ast.ConsCell, env *object.Environment) object.Object {
	spForm, ok := sexp.Car().(*ast.SpecialForm)
	if !ok {
		return newError("expect special form, got %T", sexp.Car())
	}
	if spForm.Token.Type != token.LAMBDA {
		return newError("expect special form lambda, got %s", spForm.Token.Type)
	}

	cdr, ok := sexp.Cdr().(*ast.ConsCell)
	if !ok {
		return newError("not defined lambda parameters")
	}

	params, err := evalLambdaParams(cdr.Car())
	if err != nil {
		return newError(err.Error())
	}

	cddr, ok := cdr.Cdr().(*ast.ConsCell)
	if !ok {
		return newError("not defined lambda body")
	}

	return &object.Function{
		Parameters: params,
		Body:       cddr.Car(),
		Env:        env,
	}
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

func evalQuote(sexp *ast.ConsCell) object.Object {
	spForm, ok := sexp.Car().(*ast.SpecialForm)
	if !ok {
		return newError("expect special form, got %T", sexp.Car())
	}
	if spForm.Token.Type != token.QUOTE {
		return newError("expect special form quote, got %s", spForm.Token.Type)
	}

	cdr, ok := sexp.Cdr().(*ast.ConsCell)
	if !ok {
		return newError("not defined quote expression")
	}

	return &object.Quote{
		SExpression: cdr.Car(),
	}
}

func evalBackquote(sexp *ast.ConsCell, env *object.Environment) object.Object {
	spForm, ok := sexp.Car().(*ast.SpecialForm)
	if !ok {
		return newError("expect special form, got %T", sexp.Car())
	}
	if spForm.Token.Type != token.BACKQUOTE {
		return newError("expect special form backquote, got %s", spForm.Token.Type)
	}

	cdr, ok := sexp.Cdr().(*ast.ConsCell)
	if !ok {
		return newError("not defined backquote expression")
	}

	unquoted := evalUnquote(cdr.Car(), env)

	return &object.Quote{
		SExpression: unquoted,
	}
}

func evalUnquote(sexp ast.SExpression, env *object.Environment) ast.SExpression {
	return ast.ModifyByUnquote(sexp, func(sexp ast.SExpression) ast.SExpression {
		consCell, ok := sexp.(*ast.ConsCell)
		if !ok {
			return sexp
		}

		if car, ok := consCell.Car().(*ast.SpecialForm); !ok || car.Value != "unquote" {
			return sexp
		}

		// evaluate unquoted s-expression
		cdr := consCell.Cdr().(*ast.ConsCell)
		evaluated := Eval(cdr.Car(), env)

		return convertObjectToSExpression(evaluated)
	})
}

func convertObjectToSExpression(obj object.Object) ast.SExpression {
	switch obj := obj.(type) {
	case *object.Integer:
		t := token.Token{
			Type:    token.INT,
			Literal: fmt.Sprintf("%d", obj.Value),
		}
		return &ast.IntegerLiteral{Token: t, Value: obj.Value}
	case *object.Quote:
		return obj.SExpression
	default:
		return nil
	}
}

func evalIf(consCell *ast.ConsCell, env *object.Environment) object.Object {
	spForm, ok := consCell.Car().(*ast.SpecialForm)
	if !ok {
		return newError("expect special form, got %T", consCell.Car())
	}
	if spForm.Token.Type != token.IF {
		return newError("expect special form if, got %s", spForm.Token.Type)
	}

	cdr, ok := consCell.Cdr().(*ast.ConsCell)
	if !ok {
		return newError("not defined if condition")
	}

	// evaluate the condition
	cadr := cdr.Car()
	condition := Eval(cadr, env)
	if isError(condition) {
		return condition
	}

	cddr, ok := cdr.Cdr().(*ast.ConsCell)
	if !ok {
		return newError("not defined if consequent")
	}

	// if condition is true, evaluate the consequent
	if isTruthy(condition) {
		caddr := cddr.Car()
		return Eval(caddr, env)
	}

	// if alternative is not defined, return nil
	cdddr, ok := cddr.Cdr().(*ast.ConsCell)
	if !ok {
		if _, ok := cddr.Cdr().(*ast.Nil); ok {
			return Nil
		}
		return newError("invalid if alternative")
	}

	// evaluate the alternative
	cadddr := cdddr.Car()
	return Eval(cadddr, env)
}

func isTruthy(obj object.Object) bool {
	switch obj.(type) {
	case *object.True:
		return true
	case *object.Nil:
		return false
	default:
		return true
	}
}

func evalSetq(consCell *ast.ConsCell, env *object.Environment) object.Object {
	spForm, ok := consCell.Car().(*ast.SpecialForm)
	if !ok {
		return newError("expect special form, got %T", consCell.Car())
	}
	if spForm.Token.Type != token.SETQ {
		return newError("expect special form setq, got %s", spForm.Token.Type)
	}

	cdr, ok := consCell.Cdr().(*ast.ConsCell)
	if !ok {
		return newError("not defined name of symbol")
	}

	symbolName, ok := cdr.Car().(*ast.Symbol)
	if !ok {
		return newError("expect symbol, got %T", cdr.Car())
	}

	cddr, ok := cdr.Cdr().(*ast.ConsCell)
	if !ok {
		return newError("not defined value of symbol")
	}

	value := Eval(cddr.Car(), env)
	if isError(value) {
		return value
	}

	symbol := &object.Symbol{
		Name:  symbolName.Value,
		Value: value,
	}

	env.Set(symbolName.Value, symbol)

	return symbol
}
