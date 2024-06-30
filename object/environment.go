package object

import "strings"

// all the keys in the environment are case-insensitive
type envKey string

func toEnvKey(key string) envKey {
	// env keys are uppercase
	return envKey(strings.ToUpper(key))
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

type Environment struct {
	store map[envKey]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[envKey]Object)}
}

func (e *Environment) Get(key string) (Object, bool) {
	obj, ok := e.store[toEnvKey(key)]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(key)
	}

	return obj, ok
}

func (e *Environment) Set(key string, value Object) Object {
	e.store[toEnvKey(key)] = value
	return value
}
