package golox

import (
	"errors"
	"fmt"
)

type InterpretError struct {
	Token Token
	Msg   string
}

func (e InterpretError) Error() string {
	return fmt.Sprintf("[line %d] Error at \"%s\": %s",
		e.Token.Line, e.Token.Lexeme, e.Msg)
}

func NewInterpretError(token Token, msg string) InterpretError {
	return InterpretError{token, msg}
}

func Interpret(statements []Stmt) error {
	var err error

	for _, statement := range statements {
		err = execute(statement)
		if err != nil {
			return fmt.Errorf("cant execute: %v", err)
		}
	}
	return nil
}

func execute(statement Stmt) error {
	if statement.Print != nil {
		return visitPrintStmt(statement)
	}

	if statement.Expression != nil {
		return visitExpressionStmt(statement)
	}

	return errors.New("no expression in statement")
}

func visitExpressionStmt(stmt Stmt) error {
	if stmt.Expression == nil {
		return errors.New("expected stmt.Expression")
	}

	_, err := evaluate(stmt.Expression)
	if err != nil {
		return fmt.Errorf("can't evaluate stmt.Expression: %v", err)
	}

	return nil
}

func visitPrintStmt(stmt Stmt) error {
	val, err := evaluate(stmt.Print)
	if stmt.Print == nil {
		return errors.New("expected stmt.Print")
	}

	if err != nil {
		fmt.Errorf("can't evaluate stmt.Print: %v", err)
	}

	fmt.Println(val)

	return nil
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
		return nil, err
	}

	right, err := evaluate(binary.Right)
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
		return nil, NewInterpretError(unary.Operator, "minus unary: expected number")

	case BANG:
		return !isTruthy(right), nil

	default:
		return nil, NewInterpretError(unary.Operator, "invalid unary expression")
	}
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
