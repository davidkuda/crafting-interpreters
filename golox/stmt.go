package golox

type Stmt interface {
	stmt()
}

type blockStmt struct {
	statements []Stmt
}

func (*blockStmt) stmt() {}

type exprStmt struct {
	Expression Expr
}

func (*exprStmt) stmt() {}

type ifStmt struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func (*ifStmt) stmt() {}

type printStmt struct {
	Expression Expr
}

func (*printStmt) stmt() {}

type whileStmt struct {
	condition Expr
	body      Stmt
}

func (*whileStmt) stmt() {}

type varStmt struct {
	name        Token
	initializer Expr
}

func (*varStmt) stmt() {}
