package golox

import (
	"errors"
	"fmt"
)

func Interpret(expression Expr) {
	fmt.Println(expression)
}

func evaluate(expr Expr) (any, error) {
	fmt.Println(expr)

	switch e := expr.(type) {

	case *Unary:
		return visitUnary(expr)

	case *Grouping:
		return evaluate(e.Expression)

	case *Literal:
		return e.Value, nil
	}

	return nil, errors.New("reached end of eval without evaluating anything")
}

func visitUnary(expr Expr) (any, error) {
	unary, ok := expr.(*Unary)
	if !ok {
		return nil, errors.New("not a unary")
	}

	right, err := evaluate(unary.Right)
	if err != nil {
		return nil, err
	}

	switch unary.Operator.Type {
	case MINUS:
		f, ok := right.(float64)
		if ok {
			return -f, nil
		}
		return nil, errors.New("could not evaluate")

	case BANG:
		return !isTruthy(right), nil
	}

	// unreachable:
	return nil, errors.New("reached unreachable return")
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
