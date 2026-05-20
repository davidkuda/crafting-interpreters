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

type assignExpr struct {
	name  Token
	value Expr
}

func (*assignExpr) exprNode() {}

func (a *assignExpr) String() string {
	return fmt.Sprintf("%s=%s", a.name.Lexeme, a.value)
}

type binaryExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (*binaryExpr) exprNode() {}

func (b *binaryExpr) String() string {
	return fmt.Sprintf(
		"(%s %s %s)",
		b.Operator.Lexeme,
		b.Left.String(),
		b.Right.String(),
	)
}

type groupingExpr struct {
	Expression Expr
}

func (*groupingExpr) exprNode() {}

func (g *groupingExpr) String() string {
	return fmt.Sprintf(
		"(group %v)",
		g.Expression.String(),
	)
}

type literalExpr struct {
	Value any
}

func (*literalExpr) exprNode() {}

func (l *literalExpr) String() string {
	return literalToString(l.Value)
}

type unaryExpr struct {
	Operator Token
	Right    Expr
}

func (*unaryExpr) exprNode() {}

func (u *unaryExpr) String() string {
	return fmt.Sprintf(
		"(%s %s)",
		u.Operator.Lexeme,
		u.Right,
	)
}

type variableExpr struct {
	Name Token
}

func (*variableExpr) exprNode() {}

func (v *variableExpr) String() string {
	return fmt.Sprintf("(var %s)", v.Name.Lexeme)
}

func FormatExprStringer(expr Expr) string {
	switch e := expr.(type) {
	case *binaryExpr:
		return e.String()

	case *groupingExpr:
		return e.String()

	case *literalExpr:
		return e.String()

	case *unaryExpr:
		return e.String()

	default:
		panic(fmt.Sprintf("unknown Expr type %T", expr))
	}
}

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
