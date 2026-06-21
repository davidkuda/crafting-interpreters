package golox

type LoxCallable interface {
	// arity: https://craftinginterpreters.com/functions.html#checking-arity
	// The number of arguments a function or operation expects.
	// With functions, the arity is determined by the number of parameters it declares.
	Arity() int
	Call(interpreter *Interpreter, arguments []any) (any, error)
}
