package golox

import (
	"errors"
	"fmt"
)

func (i *Interpreter) evaluate(expr Expr) (any, error) {
	switch e := expr.(type) {

	case *assignExpr:
		return i.visitAssignExpr(e)

	case *binaryExpr:
		return i.visitBinary(e)

	case *callExpr:
		return i.visitCallExpr(e)

	case *groupingExpr:
		return i.visitGroupingExpr(e)

	case *literalExpr:
		return i.visitLiteralExpr(e)

	case *logicalExpr:
		return i.visitLogical(e)

	case *unaryExpr:
		return i.visitUnary(e)

	case *variableExpr:
		return i.visitVariable(e)
	}

	return nil, errors.New("reached end of eval without evaluating anything")
}

func (i *Interpreter) visitAssignExpr(a *assignExpr) (any, error) {
	value, err := i.evaluate(a.value)
	if err != nil {
		return nil, err
	}

	i.environment.Assign(a.name, value)

	return value, nil
}

func (i *Interpreter) visitBinary(binary *binaryExpr) (any, error) {
	left, err := i.evaluate(binary.Left)
	if err != nil {
		return nil, err
	}

	right, err := i.evaluate(binary.Right)
	if err != nil {
		return nil, err
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

		return nil, NewInterpretError(binary.Operator, "invalid addition: operands must be two numbers or two strings.")

	case MINUS:
		fLeft, ok := left.(float64)
		if !ok {
			return nil, NewInterpretError(binary.Operator, "subtraction: minuend not a number")
		}
		fRight, ok := right.(float64)
		if !ok {
			return nil, NewInterpretError(binary.Operator, "subtraction: subtrahend not a number")
		}
		return fLeft - fRight, nil

	case STAR:
		fLeft, ok := left.(float64)
		if !ok {
			return nil, NewInterpretError(binary.Operator, "multiplication: factor on left not a number")
		}
		fRight, ok := right.(float64)
		if !ok {
			return nil, NewInterpretError(binary.Operator, "multiplication: factor on right not a number")
		}
		return fLeft * fRight, nil

	case SLASH:
		fLeft, ok := left.(float64)
		if !ok {
			return nil, NewInterpretError(binary.Operator, "invalid division: dividend not a number")
		}
		fRight, ok := right.(float64)
		if !ok {
			return nil, NewInterpretError(binary.Operator, "invalid division: divisor not a number")
		}
		if fRight == float64(0) {
			return nil, NewInterpretError(binary.Operator, "invalid division: divisor is 0")
		}
		return fLeft / fRight, nil

	// comparisons:
	case GREATER:
		fLeft, lok := left.(float64)
		fRight, rok := right.(float64)
		if !lok || !rok {
			return nil, NewInterpretError(binary.Operator, "comparison: expected number")
		}
		return fLeft > fRight, nil

	case GREATER_EQUAL:
		fLeft, lok := left.(float64)
		fRight, rok := right.(float64)
		if !lok || !rok {
			return nil, NewInterpretError(binary.Operator, "comparison: expected number")
		}
		return fLeft >= fRight, nil

	case LESS:
		fLeft, lok := left.(float64)
		fRight, rok := right.(float64)
		if !lok || !rok {
			return nil, NewInterpretError(binary.Operator, "comparison: expected number")
		}
		return fLeft < fRight, nil

	case LESS_EQUAL:
		fLeft, lok := left.(float64)
		fRight, rok := right.(float64)
		if !lok || !rok {
			return nil, NewInterpretError(binary.Operator, "comparison: expected number")
		}
		return fLeft <= fRight, nil

	case BANG_EQUAL:
		return !isEqual(left, right), nil

	case EQUAL_EQUAL:
		return isEqual(left, right), nil

	default:
		return nil, NewInterpretError(binary.Operator, "invalid binary")
	}
}

func (i *Interpreter) visitCallExpr(call *callExpr) (any, error) {
	callee, err := i.evaluate(call.callee)
	if err != nil {
		return nil, fmt.Errorf("i.evaluate(call.callee): %w", err)
	}

	arguments := make([]any, 0)
	for _, argument := range call.arguments {
		arg, err := i.evaluate(argument)
		if err != nil {
			return nil, fmt.Errorf("i.evaluate(argument=%v): %w", arg, err)
		}
		arguments = append(arguments, arg)
	}

	// type assertion: does callee implement LoxCallable?
	function, ok := callee.(LoxCallable)
	if !ok {
		return nil, NewRuntimeError(call.paren, "can only call functions and classes.")
	}

	if len(arguments) != function.Arity() {
		msg := fmt.Sprintf("expected %d arguments but got %d", function.Arity(), len(arguments))
		return nil, NewRuntimeError(call.paren, msg)
	}

	return function.Call(i, arguments)
}

func (i *Interpreter) visitGroupingExpr(grouping *groupingExpr) (any, error) {
	return i.evaluate(grouping.Expression)
}

func (i *Interpreter) visitLiteralExpr(literal *literalExpr) (any, error) {
	return literal.Value, nil
}

func (i *Interpreter) visitLogical(logical *logicalExpr) (any, error) {
	left, err := i.evaluate(logical.left)
	if err != nil {
		return nil, err
	}

	if logical.operator.Type == OR {
		if isTruthy(left) {
			return left, nil
		}
	} else {
		if !isTruthy(left) {
			return left, nil
		}
	}

	return i.evaluate(logical.right)
}

func (i *Interpreter) visitUnary(unary *unaryExpr) (any, error) {
	right, err := i.evaluate(unary.Right)
	if err != nil {
		return nil, err
	}

	switch unary.Operator.Type {
	case MINUS:
		f, ok := right.(float64)
		if ok {
			return -f, nil
		}
		return nil, NewInterpretError(unary.Operator, "minus unary: expected number")

	case BANG:
		return !isTruthy(right), nil

	default:
		return nil, NewInterpretError(unary.Operator, "invalid unary expression")
	}
}

func (i *Interpreter) visitVariable(vExpr *variableExpr) (any, error) {
	return i.environment.Get(vExpr.Name)
}
