package evaluator

import (
	"github.com/JunNishimura/go-lisp/object"
)

var builtinFuncs = map[string]*object.Builtin{
	"+": {
		Fn: func(args ...object.Object) object.Object {
			var sum int64
			for _, arg := range args {
				if arg.Type() != object.INTEGER_OBJ {
					return newError("argument to `+` must be INTEGER, got %s", arg.Type())
				}
				sum += arg.(*object.Integer).Value
			}
			return &object.Integer{Value: sum}
		},
	},
	"-": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) == 0 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if len(args) == 1 {
				if args[0].Type() != object.INTEGER_OBJ {
					return newError("argument to `-` must be INTEGER, got %s", args[0].Type())
				}
				return &object.Integer{Value: -args[0].(*object.Integer).Value}
			}

			var diff int64
			for i, arg := range args {
				if arg.Type() != object.INTEGER_OBJ {
					return newError("argument to `-` must be INTEGER, got %s", arg.Type())
				}
				if i == 0 {
					diff = arg.(*object.Integer).Value
				} else {
					diff -= arg.(*object.Integer).Value
				}
			}
			return &object.Integer{Value: diff}
		},
	},
	"*": {
		Fn: func(args ...object.Object) object.Object {
			var product int64 = 1
			for _, arg := range args {
				if arg.Type() != object.INTEGER_OBJ {
					return newError("argument to `*` must be INTEGER, got %s", arg.Type())
				}
				product *= arg.(*object.Integer).Value
			}
			return &object.Integer{Value: product}
		},
	},
	"/": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) == 0 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if len(args) == 1 {
				if args[0].Type() != object.INTEGER_OBJ {
					return newError("argument to `/` must be INTEGER, got %s", args[0].Type())
				}
				return &object.Integer{Value: 1 / args[0].(*object.Integer).Value}
			}

			var quotient int64
			for i, arg := range args {
				if arg.Type() != object.INTEGER_OBJ {
					return newError("argument to `/` must be INTEGER, got %s", arg.Type())
				}
				if i == 0 {
					quotient = arg.(*object.Integer).Value
				} else {
					quotient /= arg.(*object.Integer).Value
				}
			}
			return &object.Integer{Value: quotient}
		},
	},
	"=": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) == 0 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			compTo, ok := args[0].(*object.Integer)
			if !ok {
				return newError("argument to `=` must be INTEGER, got %s", args[0].Type())
			}
			for _, arg := range args[1:] {
				compFrom, ok := arg.(*object.Integer)
				if !ok {
					return newError("argument to `=` must be INTEGER, got %s", arg.Type())
				}
				if compFrom.Value != compTo.Value {
					return Nil
				}
			}
			return True
		},
	},
	"/=": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) == 0 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			compTo, ok := args[0].(*object.Integer)
			if !ok {
				return newError("argument to `/=` must be INTEGER, got %s", args[0].Type())
			}
			for _, arg := range args[1:] {
				compFrom, ok := arg.(*object.Integer)
				if !ok {
					return newError("argument to `/=` must be INTEGER, got %s", arg.Type())
				}
				if compFrom.Value == compTo.Value {
					return Nil
				}
			}
			return True
		},
	},
	"<": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) == 0 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			compTo, ok := args[0].(*object.Integer)
			if !ok {
				return newError("argument to `<` must be INTEGER, got %s", args[0].Type())
			}
			for _, arg := range args[1:] {
				compFrom, ok := arg.(*object.Integer)
				if !ok {
					return newError("argument to `<` must be INTEGER, got %s", arg.Type())
				}
				if compTo.Value >= compFrom.Value {
					return Nil
				}
				compTo = compFrom
			}
			return True
		},
	},
	"<=": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) == 0 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			compTo, ok := args[0].(*object.Integer)
			if !ok {
				return newError("argument to `<=` must be INTEGER, got %s", args[0].Type())
			}
			for _, arg := range args[1:] {
				compFrom, ok := arg.(*object.Integer)
				if !ok {
					return newError("argument to `<=` must be INTEGER, got %s", arg.Type())
				}
				if compTo.Value > compFrom.Value {
					return Nil
				}
				compTo = compFrom
			}
			return True
		},
	},
	">": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) == 0 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			compTo, ok := args[0].(*object.Integer)
			if !ok {
				return newError("argument to `>` must be INTEGER, got %s", args[0].Type())
			}
			for _, arg := range args[1:] {
				compFrom, ok := arg.(*object.Integer)
				if !ok {
					return newError("argument to `>` must be INTEGER, got %s", arg.Type())
				}
				if compTo.Value <= compFrom.Value {
					return Nil
				}
				compTo = compFrom
			}
			return True
		},
	},
	">=": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) == 0 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			compTo, ok := args[0].(*object.Integer)
			if !ok {
				return newError("argument to `>=` must be INTEGER, got %s", args[0].Type())
			}
			for _, arg := range args[1:] {
				compFrom, ok := arg.(*object.Integer)
				if !ok {
					return newError("argument to `>=` must be INTEGER, got %s", arg.Type())
				}
				if compTo.Value < compFrom.Value {
					return Nil
				}
				compTo = compFrom
			}
			return True
		},
	},
}
