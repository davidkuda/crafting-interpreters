package golox

type Stmt interface {
	stmt()
}

type exprStmt struct {
	Expression Expr
}

func (*exprStmt) stmt() {}

type printStmt struct {
	Expression Expr
}

func (*printStmt) stmt() {}

type varStmt struct {
	name        Token
	initializer Expr
}

func (*varStmt) stmt() {}
