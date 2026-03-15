package main

// type Expr interface {}
type Expr interface {
	Accept(Visitor)
}

type Visitor interface {
	VisitBinary(*Binary) any
	VisitGrouping(Grouping) any
	VisitLiteral(Literal) any
	VisitUnary(Unary) any
}

type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

func (b *Binary) Accept(v Visitor) any {
	return v.VisitBinary(b)
}

type Grouping struct {
	expression Expr
}

type Literal struct {
	value any
}

type Unary struct {
	operator Token
	right    Expr
}

type ASTPrinter struct{}

func (ast *ASTPrinter) print(expr Expr) string {
	return expr.Accept(ast)
}

func (ast *ASTPrinter) VisitBinary(b *Binary) string {
	return ""
}

func (ast *ASTPrinter) VisitGrouping(expr Expr) string {
	return ""
}

func (ast *ASTPrinter) VisitLiteral(expr Expr) string {
	return ""
}

func (ast *ASTPrinter) VisitUnary(expr Expr) string {
	return ""
}
