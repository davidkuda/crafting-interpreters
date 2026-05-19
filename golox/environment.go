package golox

type Environment struct {
	values map[string]any
}

func NewEnvironment() Environment {
	return Environment{
		values: make(map[string]any),
	}
}

// here, we make a choice: we don't check whether
// the variable exists, we either create or overwrite.
// you could design your language differently, e.g.
// return an error if the variable was already declared.
// e.g. you could say that var a = 42 defines,
// a = 21 changes the value.
// but yeah, we will just overwrite whenever we define. :)
func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Get(name Token) (any, error) {
	val, ok := e.values[name.Lexeme]
	// here again, we have to think about our language design.
	// how do we want to deal with this situation?
	// you could return a default value, a RuntimeError, a SyntaxError...
	// we return a RuntimeError.
	if !ok {
		return nil, NewRuntimeError(name, "undefined variable: " + name.Lexeme)
	}
	return val, nil
}
