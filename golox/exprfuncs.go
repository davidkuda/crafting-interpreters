package golox

import (
	"fmt"
	"strconv"
)

func FormatExpr(expr Expr) string {
	switch e := expr.(type) {
	case *binaryExpr:
		return fmt.Sprintf(
			"(%s %s %s)",
			e.Operator.Lexeme,
			FormatExpr(e.Left),
			FormatExpr(e.Right),
		)

	case *groupingExpr:
		return fmt.Sprintf(
			"(group %v)",
			FormatExpr(e.Expression),
		)

	case *literalExpr:
		return literalToString(e.Value)

	case *unaryExpr:
		return fmt.Sprintf(
			"(%s %s)",
			e.Operator.Lexeme,
			FormatExpr(e.Right),
		)

	default:
		panic(fmt.Sprintf("unknown Expr type %T", expr))
	}
}

func literalToString(v any) string {
	switch x := v.(type) {
	case nil:
		return "nil"
	case string:
		return x
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case bool:
		if x {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprint(x)
	}
}
