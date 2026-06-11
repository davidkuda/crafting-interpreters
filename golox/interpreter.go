package golox

import (
	"errors"
	"fmt"
)

type Interpreter struct {
	environment Environment
}

func NewInterpreter() Interpreter {
	return Interpreter{
		environment: NewEnvironment(nil),
	}
}

type InterpretError struct {
	environment Environment
	Token       Token
	Msg         string
}

func (e InterpretError) Error() string {
	return fmt.Sprintf("[line %d] Error at \"%s\": %s",
		e.Token.Line, e.Token.Lexeme, e.Msg)
}

func NewInterpretError(token Token, msg string) InterpretError {
	return InterpretError{
		Token: token,
		Msg:   msg,
	}
}

func (i *Interpreter) Interpret(statements []Stmt) error {
	var err error

	for _, statement := range statements {
		err = i.execute(statement)
		if err != nil {
			return fmt.Errorf("cant execute: %v", err)
		}
	}
	return nil
}

func (i *Interpreter) execute(statement Stmt) error {
	switch statement.(type) {
	case *blockStmt:
		return i.visitBlockStatement(statement)

	case *exprStmt:
		return i.visitExpressionStmt(statement)

	case *ifStmt:
		return i.visitIfStmt(statement)

	case *printStmt:
		return i.visitPrintStmt(statement)

	case *whileStmt:
		return i.visitWhileStmt(statement)

	case *varStmt:
		return i.visitVarStmt(statement)
	}

	return errors.New("no expression in statement")
}

func (i *Interpreter) visitBlockStatement(s Stmt) error {
	block, ok := s.(*blockStmt)
	if !ok {
		return errors.New("not a blockStmt")
	}

	enclosing := i.environment
	return i.executeBlock(block.statements, NewEnvironment(&enclosing))
}

func (i *Interpreter) executeBlock(statements []Stmt, env Environment) error {
	var err error

	previous := i.environment
	i.environment = env

	for _, statement := range statements {
		err = i.execute(statement)
		if err != nil {
			i.environment = previous
			return err
		}
	}

	i.environment = previous

	return nil
}

func (i *Interpreter) visitExpressionStmt(s Stmt) error {
	stmt, ok := s.(*exprStmt)
	if !ok {
		return errors.New("not an exprStmt")
	}

	if stmt.Expression == nil {
		return errors.New("expected exprStmt.Expression")
	}

	_, err := i.evaluate(stmt.Expression)
	if err != nil {
		return fmt.Errorf("can't evaluate stmt.Expression: %v", err)
	}

	return nil
}

func (i *Interpreter) visitIfStmt(s Stmt) error {
	var err error

	stmt, ok := s.(*ifStmt)
	if !ok {
		return errors.New("not an ifStmt")
	}

	condition, err := i.evaluate(stmt.condition)
	if err != nil {
		return fmt.Errorf("could not evaluate ifStmt.condition: %v", err)
	}

	if isTruthy(condition) {
		err = i.execute(stmt.thenBranch)
		if err != nil {
			return fmt.Errorf("could not execute ifStmt.thenBranch: %v", err)
		}
	} else if stmt.elseBranch != nil {
		err = i.execute(stmt.elseBranch)
		if err != nil {
			return fmt.Errorf("could not execute ifStmt.elseBranch: %v", err)
		}
	}

	return nil
}

func (i *Interpreter) visitPrintStmt(s Stmt) error {
	stmt, ok := s.(*printStmt)
	if !ok {
		return errors.New("not a printStmt")
	}

	if stmt.Expression == nil {
		return errors.New("expected stmt.Expression")
	}

	val, err := i.evaluate(stmt.Expression)

	if err != nil {
		fmt.Errorf("can't evaluate stmt.Print: %v", err)
	}

	fmt.Println(val)

	return nil
}

func (i *Interpreter) visitWhileStmt(s Stmt) error {
	stmt, ok := s.(*whileStmt)
	if !ok {
		return errors.New("not a whileStmt")
	}

	for {
		condition, err := i.evaluate(stmt.condition)
		if err != nil {
			return err
		}

		if !isTruthy(condition) {
			break
		}

		i.execute(stmt.body)
	}

	return nil
}

func (i *Interpreter) visitVarStmt(s Stmt) error {
	var err error

	stmt, ok := s.(*varStmt)
	if !ok {
		return errors.New("not a varStmt")
	}

	var value any
	// here is again a language design choice.
	// e.g. if you wanted a default value for your types,
	// you could add it here.
	if stmt.initializer != nil {
		value, err = i.evaluate(stmt.initializer)
		if err != nil {
			return err
		}
	}

	i.environment.Define(stmt.name.Lexeme, value)

	return nil
}

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
