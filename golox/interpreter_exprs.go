package golox

import (
	"errors"
	"fmt"
)

func (i *Interpreter) evaluate(expr Expr) (any, error) {
	switch e := expr.(type) {

	case *binaryExpr:
		return i.visitBinary(expr)

	case *unaryExpr:
		return i.visitUnary(expr)

	case *groupingExpr:
		return i.evaluate(e.Expression)

	case *literalExpr:
		return e.Value, nil

	case *logicalExpr:
		return i.visitLogical(expr)

	case *variableExpr:
		return i.visitVariable(expr)

	case *assignExpr:
		return i.visitAssignExpr(expr)

	case *callExpr:
		return i.visitCallExpr(expr)
	}

	return nil, errors.New("reached end of eval without evaluating anything")
}

func (i *Interpreter) visitLogical(expr Expr) (any, error) {
	logical, ok := expr.(*logicalExpr)
	if !ok {
		return nil, errors.New("not a logical expression")
	}

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

func (i *Interpreter) visitBinary(expr Expr) (any, error) {
	binary, ok := expr.(*binaryExpr)
	if !ok {
		return nil, errors.New("not a binary")
	}

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

func (i *Interpreter) visitUnary(expr Expr) (any, error) {
	unary, ok := expr.(*unaryExpr)
	if !ok {
		return nil, errors.New("not a unary")
	}

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

func (i *Interpreter) visitVariable(expr Expr) (any, error) {
	variableExpr, ok := expr.(*variableExpr)
	if !ok {
		return nil, errors.New("not a variable expression")
	}
	return i.environment.Get(variableExpr.Name)
}

func (i *Interpreter) visitAssignExpr(expr Expr) (any, error) {
	a, ok := expr.(*assignExpr)
	if !ok {
		return nil, errors.New("not a variable expression")
	}
	value, err := i.evaluate(a.value)
	if err != nil {
		return nil, err
	}

	i.environment.Assign(a.name, value)

	return value, nil
}

func (i *Interpreter) visitCallExpr(expr Expr) (any, error) {
	//
	call, ok := expr.(*callExpr)
	if !ok {
		return nil, errors.New("not a call expression")
	}

	arguments := make([]any, 0)
	for _, argument := range call.arguments {
		arg, err := i.evaluate(argument)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, arg)
	}

	function, ok := call.callee.(LoxCallable)
	if !ok {
		return nil, NewRuntimeError(call.paren, "can only call functions and classes.")
	}

	if len(arguments) != function.Arity() {
		msg := fmt.Sprintf("expected %d arguments but got %d", function.Arity(), len(arguments))
		return nil, NewRuntimeError(call.paren, msg)
	}

	return function.Call(i, arguments)
}
