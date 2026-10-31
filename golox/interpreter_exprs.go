package golox

import (
	"errors"
	"fmt"
)

type execResult struct {
	value        any
	shouldReturn bool
}

// we don't return a pointer here, because execResult is "plain old data".
// see https://en.wikipedia.org/wiki/Passive_data_structure
func (i *Interpreter) evaluate(expr Expr) (execResult, error) {
	i.logger.Debug("i.evaluate(expr)", "expr", fmt.Sprint(expr))
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

	return execResult{}, errors.New("reached end of eval without evaluating anything")
}

func (i *Interpreter) visitAssignExpr(a *assignExpr) (execResult, error) {
	i.logger.Debug("i.visitAssignExpr()")

	execres, err := i.evaluate(a.value)
	if err != nil {
		return execResult{}, err
	}

	i.environment.Assign(a.name, execres.value)

	return execres, nil
}

func (i *Interpreter) visitBinary(binary *binaryExpr) (execResult, error) {
	i.logger.Debug("i.visitBinary()")

	execResLeft, err := i.evaluate(binary.Left)
	if err != nil {
		return execResult{}, err
	}
	left := execResLeft.value

	execResRight, err := i.evaluate(binary.Right)
	if err != nil {
		return execResult{}, err
	}
	right := execResRight.value

	switch binary.Operator.Type {

	// arithmetic operators:
	case PLUS:
		// add numbers:
		fLeft, okLeft := left.(float64)
		fRight, okRight := right.(float64)
		if okLeft && okRight {
			return execResult{value: fLeft + fRight}, nil
		}

		// add strings:
		sLeft, okLeft := left.(string)
		sRight, okRight := right.(string)
		if okLeft && okRight {
			return execResult{value: sLeft + sRight}, nil
		}

		return execResult{}, NewInterpretError(binary.Operator, "invalid addition: operands must be two numbers or two strings.")

	case MINUS:
		fLeft, ok := left.(float64)
		if !ok {
			return execResult{}, NewInterpretError(binary.Operator, "subtraction: minuend not a number")
		}
		fRight, ok := right.(float64)
		if !ok {
			return execResult{}, NewInterpretError(binary.Operator, "subtraction: subtrahend not a number")
		}
		return execResult{value: fLeft - fRight}, nil

	case STAR:
		fLeft, ok := left.(float64)
		if !ok {
			return execResult{}, NewInterpretError(binary.Operator, "multiplication: factor on left not a number")
		}
		fRight, ok := right.(float64)
		if !ok {
			return execResult{}, NewInterpretError(binary.Operator, "multiplication: factor on right not a number")
		}
		return execResult{value: fLeft * fRight}, nil

	case SLASH:
		fLeft, ok := left.(float64)
		if !ok {
			return execResult{}, NewInterpretError(binary.Operator, "invalid division: dividend not a number")
		}
		fRight, ok := right.(float64)
		if !ok {
			return execResult{}, NewInterpretError(binary.Operator, "invalid division: divisor not a number")
		}
		if fRight == float64(0) {
			return execResult{}, NewInterpretError(binary.Operator, "invalid division: divisor is 0")
		}
		return execResult{value: fLeft / fRight}, nil

	// comparisons:
	case GREATER:
		fLeft, lok := left.(float64)
		fRight, rok := right.(float64)
		if !lok || !rok {
			return execResult{}, NewInterpretError(binary.Operator, "comparison: expected number")
		}
		return execResult{value: fLeft > fRight}, nil

	case GREATER_EQUAL:
		fLeft, lok := left.(float64)
		fRight, rok := right.(float64)
		if !lok || !rok {
			return execResult{}, NewInterpretError(binary.Operator, "comparison: expected number")
		}
		return execResult{value: fLeft >= fRight}, nil

	case LESS:
		fLeft, lok := left.(float64)
		fRight, rok := right.(float64)
		if !lok || !rok {
			return execResult{}, NewInterpretError(binary.Operator, "comparison: expected number")
		}
		return execResult{value: fLeft < fRight}, nil

	case LESS_EQUAL:
		fLeft, lok := left.(float64)
		fRight, rok := right.(float64)
		if !lok || !rok {
			return execResult{}, NewInterpretError(binary.Operator, "comparison: expected number")
		}
		return execResult{value: fLeft <= fRight}, nil

	case BANG_EQUAL:
		return execResult{value: !isEqual(left, right)}, nil

	case EQUAL_EQUAL:
		return execResult{value: isEqual(left, right)}, nil

	default:
		return execResult{}, NewInterpretError(binary.Operator, "invalid binary")
	}
}

func (i *Interpreter) visitCallExpr(call *callExpr) (execResult, error) {
	i.logger.Debug("i.visitCallExpr()")

	execResCallee, err := i.evaluate(call.callee)
	if err != nil {
		return execResult{}, fmt.Errorf("i.evaluate(call.callee): %w", err)
	}
	callee := execResCallee.value

	arguments := make([]any, 0)
	for _, argument := range call.arguments {
		arg, err := i.evaluate(argument)
		if err != nil {
			return execResult{}, fmt.Errorf("i.evaluate(argument=%v): %w", arg, err)
		}
		arguments = append(arguments, arg)
	}

	// type assertion: does callee implement LoxCallable?
	function, ok := callee.(LoxCallable)
	if !ok {
		return execResult{}, NewRuntimeError(call.paren, "can only call functions and classes.")
	}

	if len(arguments) != function.Arity() {
		msg := fmt.Sprintf("expected %d arguments but got %d", function.Arity(), len(arguments))
		return execResult{}, NewRuntimeError(call.paren, msg)
	}

	return function.Call(i, arguments)
}

func (i *Interpreter) visitGroupingExpr(grouping *groupingExpr) (execResult, error) {
	i.logger.Debug("i.visitGroupingExpr()")

	return i.evaluate(grouping.Expression)
}

func (i *Interpreter) visitLiteralExpr(literal *literalExpr) (execResult, error) {
	i.logger.Debug("i.visitLiteralExpr()")

	return execResult{value: literal.Value}, nil
}

func (i *Interpreter) visitLogical(logical *logicalExpr) (execResult, error) {
	i.logger.Debug("i.visitLogical()")

	execResLeft, err := i.evaluate(logical.left)
	if err != nil {
		return execResult{}, err
	}
	left := execResLeft.value

	if logical.operator.Type == OR {
		if isTruthy(left) {
			return execResult{value: left}, nil
		}
	} else {
		if !isTruthy(left) {
			return execResult{value: left}, nil
		}
	}

	return i.evaluate(logical.right)
}

func (i *Interpreter) visitUnary(unary *unaryExpr) (execResult, error) {
	i.logger.Debug("i.visitUnary()")

	execResRight, err := i.evaluate(unary.Right)
	if err != nil {
		return execResult{}, err
	}
	right := execResRight.value

	switch unary.Operator.Type {
	case MINUS:
		f, ok := right.(float64)
		if ok {
			return execResult{value: -f}, nil
		}
		return execResult{}, NewInterpretError(unary.Operator, "minus unary: expected number")

	case BANG:
		return execResult{value: !isTruthy(right)}, nil

	default:
		return execResult{}, NewInterpretError(unary.Operator, "invalid unary expression")
	}
}

func (i *Interpreter) visitVariable(vExpr *variableExpr) (execResult, error) {
	i.logger.Debug("i.visitVariable()")

	val, err := i.environment.Get(vExpr.Name)
	if err != nil {
		return execResult{}, fmt.Errorf("i.visitVariable(%v): %w", vExpr.Name, err)
	}

	return execResult{value: val}, nil
}
