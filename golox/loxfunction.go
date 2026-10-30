package golox

import (
	"errors"
	"fmt"
)

func newLoxFunction(declaration *functionStmt) *loxFunction {
	return &loxFunction{
		declaration: declaration,
	}
}

type loxFunction struct {
	declaration *functionStmt
}

func (f *loxFunction) Arity() int {
	return len(f.declaration.params)
}

func (f *loxFunction) Call(interpreter *Interpreter, arguments []any) (any, error) {
	env := NewEnvironment(interpreter.globals)

	for i := 0; i < len(f.declaration.params); i++ {
		env.Define(f.declaration.params[i].Lexeme, arguments[i])
	}

	err := interpreter.executeBlock(f.declaration.body, env)

	var returnValue ReturnValue
	if errors.As(err, &returnValue) {
		return returnValue.Value, nil
	}

	if err != nil {
		return nil, fmt.Errorf("Call(): %w", err)
	}

	return nil, nil
}

func (f *loxFunction) String() string {
	return "<fn " + f.declaration.name.Lexeme + ">"
}
