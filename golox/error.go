package golox

import "fmt"

type GoloxError struct {
	line    int
	message string
	where   string
}

func NewError(line int, message string) *GoloxError {
	return &GoloxError{
		line:    line,
		message: message,
		where:   "",
	}
}

func (e *GoloxError) Error() string {
	return fmt.Sprintf("[line %d] Error %s: %s\n", e.line, e.where, e.message)
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
