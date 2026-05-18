package golox

import (
	"fmt"
	"strconv"
)

// Lehrstück: https://github.com/golang/go/blob/6614616b7576a8011053c4b50fbb5e64d469837b/src/go/ast/ast.go#L42
type Expr interface {
	exprNode()
	fmt.Stringer
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (*Binary) exprNode() {}

func (b *Binary) String() string {
	return fmt.Sprintf(
		"(%s %s %s)",
		b.Operator.Lexeme,
		b.Left.String(),
		b.Right.String(),
	)
}

type Grouping struct {
	Expression Expr
}

func (*Grouping) exprNode() {}

func (g *Grouping) String() string {
	return fmt.Sprintf(
		"(group %v)",
		g.Expression.String(),
	)
}

type Literal struct {
	Value any
}

func (*Literal) exprNode() {}

func (l *Literal) String() string {
	return literalToString(l.Value)
}

type Unary struct {
	Operator Token
	Right    Expr
}

func (*Unary) exprNode() {}

func (u *Unary) String() string {
	return fmt.Sprintf(
		"(%s %s)",
		u.Operator.Lexeme,
		u.Right,
	)
}

func FormatExprStringer(expr Expr) string {
	switch e := expr.(type) {
	case *Binary:
		return e.String()

	case *Grouping:
		return e.String()

	case *Literal:
		return e.String()

	case *Unary:
		return e.String()

	default:
		panic(fmt.Sprintf("unknown Expr type %T", expr))
	}
}

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
