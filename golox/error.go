package golox

import "fmt"

// TODO: Error might be confused with go error.
// either there is a better name, or use go errors?
type Error struct {
	line    int
	message string
	where   string
}

func NewError(line int, message string) *Error {
	return &Error{
		line:    line,
		message: message,
		where:   "",
	}
}

func (e *Error) Report() {
	fmt.Printf("[line %d] Error %s: %s\n", e.line, e.where, e.message)
}
