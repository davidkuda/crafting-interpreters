package golox

import (
	"fmt"
	"strconv"
)

// Lehrstück: https://github.com/golang/go/blob/6614616b7576a8011053c4b50fbb5e64d469837b/src/go/ast/ast.go#L42
type Expr interface {
	exprNode()
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (*Binary) exprNode() {}

type Grouping struct {
	Expression Expr
}

func (*Grouping) exprNode() {}

type Literal struct {
	Value any
}

func (*Literal) exprNode() {}

type Unary struct {
	Operator Token
	Right    Expr
}

func (*Unary) exprNode() {}

func FormatExpr(expr Expr) string {
	switch e := expr.(type) {
	case *Binary:
		return fmt.Sprintf(
			"(%s %s %s)",
			e.Operator.Lexeme,
			FormatExpr(e.Left),
			FormatExpr(e.Right),
		)

	case *Grouping:
		return fmt.Sprintf(
			"(group %v)",
			FormatExpr(e.Expression),
		)

	case *Literal:
		return literalToString(e.Value)

	case *Unary:
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
