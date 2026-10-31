package golox

import (
	"fmt"
	"log/slog"
	"os"
)

type Interpreter struct {
	globals     *Environment
	environment *Environment
	logger      slog.Logger
}

func NewInterpreter() *Interpreter {
	globals := NewEnvironment(nil)
	globals.Define("clock", Clock{})

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		// Level: slog.LevelInfo,
	}))

	return &Interpreter{
		globals:     globals,
		environment: globals,
		logger:      *logger,
	}
}

type InterpretError struct {
	environment Environment
	Token       Token
	Msg         string
}

func (e InterpretError) Error() string {
	return fmt.Sprintf("[line %d] Error at \"%s\": %s",
		e.Token.Line, e.Token.Lexeme, e.Msg)
}

func NewInterpretError(token Token, msg string) InterpretError {
	return InterpretError{
		Token: token,
		Msg:   msg,
	}
}

func (i *Interpreter) Interpret(statements []Stmt) error {
	var err error

	for _, statement := range statements {
		_, err = i.execute(statement)
		if err != nil {
			return fmt.Errorf("i.Interpret(): %w", err)
		}
	}
	return nil
}

// from page 101:
// Lox follows Ruby's simple rule: false and nil are falsey,
// and everything else is truthy.
func isTruthy(object any) bool {
	if object == nil {
		return false
	}

	b, ok := object.(bool)
	if ok {
		return b
	}

	return true
}

func isEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil {
		return false
	}

	return a == b
}
