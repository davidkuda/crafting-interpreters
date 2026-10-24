package golox

import (
	"errors"
	"fmt"
)

func (i *Interpreter) execute(statement Stmt) error {
	switch stmt := statement.(type) {
	case *blockStmt:
		return i.visitBlockStatement(statement)

	case *exprStmt:
		return i.visitExpressionStmt(statement)

	case *functionStmt:
		return i.visitFunctionStmt(stmt)

	case *ifStmt:
		return i.visitIfStmt(statement)

	case *printStmt:
		return i.visitPrintStmt(statement)

	case *whileStmt:
		return i.visitWhileStmt(statement)

	case *varStmt:
		return i.visitVarStmt(statement)
	}

	return errors.New("could not match statement type")
}

func (i *Interpreter) visitBlockStatement(s Stmt) error {
	i.logger.Debug("i.visitBlockStatement()")

	block, ok := s.(*blockStmt)
	if !ok {
		return errors.New("not a blockStmt")
	}

	enclosing := i.environment
	return i.executeBlock(block.statements, NewEnvironment(enclosing))
}

func (i *Interpreter) executeBlock(statements []Stmt, env *Environment) error {
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
	i.logger.Debug("i.visitExpressionStmt()")

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

func (i *Interpreter) visitFunctionStmt(s *functionStmt) error {
	i.logger.Debug("i.visitFunctionStmt()")

	function := newLoxFunction(s)
	i.environment.Define(s.name.Lexeme, function)

	return nil
}

func (i *Interpreter) visitIfStmt(s Stmt) error {
	i.logger.Debug("i.visitIfStmt()")

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
	i.logger.Debug("i.visitPrintStmt()")

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
	i.logger.Debug("i.visitWhileStmt()")

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
	i.logger.Debug("i.visitVarStmt()")

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
