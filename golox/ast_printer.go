package golox

import (
	"fmt"
)

// ------------------------------------------------
// AST: Statements

func (f *functionStmt) String() string {
	return fmt.Sprintf(
		"(function name=%v args=%s body=%v)",
		f.name.Lexeme, f.params, f.body,
	)
}

func (s *exprStmt) String() string {
	return fmt.Sprintf(
		// "exprStmt{expression: %v}",
		"(exprStmt %v)",
		s.Expression,
	)
}


// ------------------------------------------------
// AST: Expressions

func (a *assignExpr) String() string {
	return fmt.Sprintf(
		"(assign name=%v value=%v)",
		a.name.Lexeme, a.value,
		)
}

func (*binaryExpr) String() string {
	return "binaryExpr"
}

func (call *callExpr) String() string {
	return fmt.Sprintf(
		"(call callee=%v arguments=%v)",
		call.callee,
		call.arguments,
	)
}

func (*groupingExpr) String() string {
	return "groupingExpr"
}

func (literal *literalExpr) String() string {
	return literalToString(literal.Value)
}

func (*logicalExpr) String() string {
	return "logicalExpr"
}

func (*unaryExpr) String() string {
	return "unaryExpr"
}

func (variable *variableExpr) String() string {
	return fmt.Sprintf(
		"(variable name=%v)",
		variable.Name.Lexeme,
	)
}
