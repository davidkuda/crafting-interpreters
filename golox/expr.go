package golox

import (
	"fmt"
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

type logicalExpr struct {
	left     Expr
	operator Token
	right    Expr
}

func (*logicalExpr) exprNode() {}

func (l *logicalExpr) String() string {
	return fmt.Sprintf("%q %s %q", l.left, &l.operator, l.right)
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
