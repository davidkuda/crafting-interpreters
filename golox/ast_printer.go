package golox

import (
	"fmt"
)

func (f *functionStmt) String() string {
	return fmt.Sprintf(
		"(function name=%v args=%v body=%v)",
		f.name, f.params, f.body,
	)
}

func (s *exprStmt) String() string {
	return fmt.Sprintf(
		// "exprStmt{expression: %v}",
		"(exprStmt %v)",
		s.Expression,
	)
}

func (call *callExpr) String() string {
	return fmt.Sprintf(
		"(call callee=%v arguments=%v)",
		call.callee,
		call.arguments,
	)
}

func (literal *literalExpr) String() string {
	return literalToString(literal.Value)
}

func (variable *variableExpr) String() string {
	return fmt.Sprintf(
		"(variable name=%v)",
		variable.Name.Lexeme,
	)
}
