package golox

// Lehrstück: https://github.com/golang/go/blob/6614616b7576a8011053c4b50fbb5e64d469837b/src/go/ast/ast.go#L42
type Expr interface {
	exprNode()
}

type assignExpr struct {
	name  Token
	value Expr
}

func (*assignExpr) exprNode() {}

type binaryExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (*binaryExpr) exprNode() {}

type callExpr struct {
	callee    Expr
	paren     Token
	arguments []Expr
}

func (*callExpr) exprNode() {}

type groupingExpr struct {
	Expression Expr
}

func (*groupingExpr) exprNode() {}

type literalExpr struct {
	Value any
}

func (*literalExpr) exprNode() {}

type logicalExpr struct {
	left     Expr
	operator Token
	right    Expr
}

func (*logicalExpr) exprNode() {}

type unaryExpr struct {
	Operator Token
	Right    Expr
}

func (*unaryExpr) exprNode() {}

type variableExpr struct {
	Name Token
}

func (*variableExpr) exprNode() {}
