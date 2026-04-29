package golox

import (
	"errors"
	// "fmt"
)

func Interpret(expression Expr) (any, error) {
	return evaluate(expression)
}

func evaluate(expr Expr) (any, error) {
	switch e := expr.(type) {

	case *Binary:
		return visitBinary(expr)

	case *Unary:
		return visitUnary(expr)

	case *Grouping:
		return evaluate(e.Expression)

	case *Literal:
		return e.Value, nil
	}

	return nil, errors.New("reached end of eval without evaluating anything")
}

func visitBinary(expr Expr) (any, error) {
	binary, ok := expr.(*Binary)
	if !ok {
		return nil, errors.New("not a binary")
	}

	left, err := evaluate(binary.Left)
	if err != nil {
		return nil, errors.New("something wrong")
	}

	right, err := evaluate(binary.Right)
	if err != nil {
		return nil, errors.New("something wrong")
	}

	switch binary.Operator.Type {

	// arithmetic operators:
	case PLUS:
		// add numbers:
		fLeft, okLeft := left.(float64)
		fRight, okRight := right.(float64)
		if okLeft && okRight {
			return fLeft + fRight, nil
		}

		// add strings:
		sLeft, okLeft := left.(string)
		sRight, okRight := right.(string)
		if okLeft && okRight {
			return sLeft + sRight, nil
		}

		return nil, errors.New("invalid addition")

	case MINUS:
		fLeft, ok := left.(float64)
		if !ok {
			return nil, errors.New("subtraction: minuend not a number")
		}
		fRight, ok := right.(float64)
		if !ok {
			return nil, errors.New("subtraction: subtrahend not a number")
		}
		return fLeft - fRight, nil

	case STAR:
		fLeft, ok := left.(float64)
		if !ok {
			return nil, errors.New("multiplication: factor on left not a number")
		}
		fRight, ok := right.(float64)
		if !ok {
			return nil, errors.New("multiplication: factor on right not a number")
		}
		return fLeft * fRight, nil

	case SLASH:
		fLeft, ok := left.(float64)
		if !ok {
			return nil, errors.New("division: dividend not a number")
		}
		fRight, ok := right.(float64)
		if !ok {
			return nil, errors.New("division: divisor not a number")
		}
		return fLeft / fRight, nil

	// comparisons:
	case GREATER:
		fLeft, ok := left.(float64)
		if !ok {
			return nil, errors.New("comparison: expected number")
		}
		fRight, ok := right.(float64)
		if !ok {
			return nil, errors.New("comparison: expected number")
		}
		return fLeft > fRight, nil

	case GREATER_EQUAL:
		fLeft, ok := left.(float64)
		if !ok {
			return nil, errors.New("comparison: expected number")
		}
		fRight, ok := right.(float64)
		if !ok {
			return nil, errors.New("comparison: expected number")
		}
		return fLeft >= fRight, nil

	case LESS:
		fLeft, ok := left.(float64)
		if !ok {
			return nil, errors.New("comparison: expected number")
		}
		fRight, ok := right.(float64)
		if !ok {
			return nil, errors.New("comparison: expected number")
		}
		return fLeft < fRight, nil

	case LESS_EQUAL:
		fLeft, ok := left.(float64)
		if !ok {
			return nil, errors.New("comparison: expected number")
		}
		fRight, ok := right.(float64)
		if !ok {
			return nil, errors.New("comparison: expected number")
		}
		return fLeft <= fRight, nil

	case BANG_EQUAL:
		return !isEqual(left, right), nil

	case EQUAL_EQUAL:
		return isEqual(left, right), nil

	default:
		return nil, errors.New("invalid binary")
	}

	return nil, errors.New("invalid binary")
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

func isEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil {
		return false
	}

	return a == b
}
