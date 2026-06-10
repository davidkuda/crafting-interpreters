package golox

import "fmt"

type GoloxError struct {
	line    int
	message string
	where   string
}

func NewError(line int, message string, where string) *GoloxError {
	return &GoloxError{
		line:    line,
		message: message,
		where:   where,
	}
}

func (e *GoloxError) Error() string {
	return fmt.Sprintf("[line %d] Error %q: %s", e.line, e.where, e.message)
}

func ReportGoloxError(token Token, message string) {
	if token.Type == EOF {
		fmt.Printf(
			"[line %d] Error %s: %s\n",
			token.Line,
			"at end",
			message,
		)
	} else {
		fmt.Printf(
			"[line %d] Error at %s: %s\n",
			token.Line,
			token.Lexeme,
			message,
		)
	}
}

func (e *GoloxError) Report() {
	fmt.Printf("[line %d] Error %s: %s\n", e.line, e.where, e.message)
}

type RuntimeError struct {
	Token   Token
	Message string
}

func NewRuntimeError(tkn Token, msg string) *RuntimeError {
	return &RuntimeError{
		Token:   tkn,
		Message: msg,
	}
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("runtime error: line %d: %s: %s",
		e.Token.Line, e.Token.Lexeme, e.Message)
}
