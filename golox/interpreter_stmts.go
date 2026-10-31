package golox

import (
	"errors"
	"fmt"
)

type execResult struct {
	value        any
	shouldReturn bool
}

func (i *Interpreter) execute(statement Stmt) (execResult, error) {
	switch stmt := statement.(type) {
	case *blockStmt:
		return i.visitBlockStatement(stmt)

	case *exprStmt:
		return i.visitExpressionStmt(stmt)

	case *functionStmt:
		return i.visitFunctionStmt(stmt)

	case *ifStmt:
		return i.visitIfStmt(stmt)

	case *printStmt:
		return i.visitPrintStmt(stmt)

	case *returnStmt:
		return i.visitReturnStmt(stmt)

	case *whileStmt:
		return i.visitWhileStmt(stmt)

	case *varStmt:
		return i.visitVarStmt(stmt)
	}

	return execResult{}, errors.New("could not match statement type")
}

func (i *Interpreter) visitBlockStatement(block *blockStmt) (execResult, error) {
	i.logger.Debug("i.visitBlockStatement()")

	enclosing := i.environment
	return i.executeBlock(block.statements, NewEnvironment(enclosing))
}

func (i *Interpreter) executeBlock(statements []Stmt, env *Environment) (execResult, error) {
	i.logger.Debug("i.executeBlock()")

	previous := i.environment
	i.environment = env

	for _, statement := range statements {
		res, err := i.execute(statement)
		if err != nil {
			i.environment = previous
			return execResult{}, err
		}
		if res.shouldReturn {
			return res, nil
		}
	}

	i.environment = previous

	return execResult{}, nil
}

func (i *Interpreter) visitExpressionStmt(stmt *exprStmt) (execResult, error) {
	i.logger.Debug("i.visitExpressionStmt()")

	if stmt.Expression == nil {
		return execResult{}, errors.New("expected exprStmt.Expression")
	}

	_, err := i.evaluate(stmt.Expression)
	if err != nil {
		return execResult{}, fmt.Errorf("can't evaluate stmt.Expression: %v", err)
	}

	return execResult{}, nil
}

func (i *Interpreter) visitFunctionStmt(s *functionStmt) (execResult, error) {
	i.logger.Debug("i.visitFunctionStmt()")

	function := newLoxFunction(s)
	i.environment.Define(s.name.Lexeme, function)

	return execResult{}, nil
}

func (i *Interpreter) visitIfStmt(stmt *ifStmt) (execResult, error) {
	i.logger.Debug("i.visitIfStmt()")

	var err error

	condition, err := i.evaluate(stmt.condition)
	if err != nil {
		return execResult{}, fmt.Errorf("could not evaluate ifStmt.condition: %v", err)
	}

	if isTruthy(condition) {
		res, err := i.execute(stmt.thenBranch)
		if err != nil {
			return execResult{}, fmt.Errorf("executing then branch: %w", err)
		}
		return res, nil
	} else if stmt.elseBranch != nil {
		res, err := i.execute(stmt.elseBranch)
		if err != nil {
			return execResult{}, fmt.Errorf("executing else branch: %w", err)
		}
		return res, nil
	}

	return execResult{}, nil
}

func (i *Interpreter) visitPrintStmt(stmt *printStmt) (execResult, error) {
	i.logger.Debug("i.visitPrintStmt()")

	if stmt.Expression == nil {
		return execResult{}, errors.New("expected stmt.Expression")
	}

	val, err := i.evaluate(stmt.Expression)

	if err != nil {
		return execResult{}, fmt.Errorf("can't evaluate stmt.Print: %v", err)
	}

	fmt.Println(val)

	return execResult{}, nil
}

func (i *Interpreter) visitReturnStmt(s *returnStmt) (execResult, error) {
	i.logger.Debug("i.visitReturnStmt()")

	var err error
	var value any

	if s.value != nil {
		value, err = i.evaluate(s.value)
		if err != nil {
			return execResult{}, err
		}
	}

	res := execResult{
		value:        value,
		shouldReturn: true,
	}

	return res, nil
}

func (i *Interpreter) visitWhileStmt(stmt *whileStmt) (execResult, error) {
	i.logger.Debug("i.visitWhileStmt()")

	for {
		condition, err := i.evaluate(stmt.condition)
		if err != nil {
			return execResult{}, err
		}

		if !isTruthy(condition) {
			break
		}

		i.execute(stmt.body)
	}

	return execResult{}, nil
}

func (i *Interpreter) visitVarStmt(stmt *varStmt) (execResult, error) {
	i.logger.Debug("i.visitVarStmt()")

	var err error

	var value any
	// here is again a language design choice.
	// e.g. if you wanted a default value for your types,
	// you could add it here.
	if stmt.initializer != nil {
		value, err = i.evaluate(stmt.initializer)
		if err != nil {
			return execResult{}, err
		}
	}

	i.environment.Define(stmt.name.Lexeme, value)

	return execResult{}, nil
}
