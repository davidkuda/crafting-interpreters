package golox

type LoxCallable interface {
	// arity: https://craftinginterpreters.com/functions.html#checking-arity
	// The number of arguments a function or operation expects.
	// With functions, the arity is determined by the number of parameters it declares.
	Arity() int
	Call(interpreter *Interpreter, arguments []any) (any, error)
}

func newLoxFunction(declaration *functionStmt) *loxFunction {
	return &loxFunction{
		declaration: declaration,
	}
}

type loxFunction struct {
	declaration *functionStmt
}

func (f *loxFunction) Arity() int {
	return len(f.declaration.params)
}

func (f *loxFunction) Call(interpreter *Interpreter, arguments []any) (any, error) {
	env := NewEnvironment(interpreter.globals)

	for i := 0; i < len(f.declaration.params); i++ {
		env.Define(f.declaration.params[i].Lexeme, arguments[i])
	}

	interpreter.executeBlock(f.declaration.body, env)

	return nil, nil
}

func (f *loxFunction) String() string {
	return "<fn " + f.declaration.name.Lexeme + ">"
}
