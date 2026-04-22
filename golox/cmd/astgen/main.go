package main

type Expr interface{}

type Token struct{}

type Visitor interface{
	visit()
}

type Binary struct {
	left     Expr
	operator Token
	right    Expr
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

func main() {
	expr := Binary{
		left: Unary{

		}
	}
}
