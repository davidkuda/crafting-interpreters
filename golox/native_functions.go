package golox

import (
    "time"
)

type Clock struct{}

func (Clock) Arity() int {
	return 0
}

func (Clock) Call(interpreter *Interpreter, arguments []any) any {
	return float64(time.Now().UnixMilli()) / 1000.0
}

func (Clock) String() string {
	return "<native fn>"
}
