package golox

import (
	"fmt"
)

// Recursive Descent Parsing:
// top-down parser: start at the top grammar rule (see README)
// and work your way down to the tree leaves.
// Each expression grammar rule becomes a function.

func Parse(tokens []Token) ([]Stmt, error) {
	statements := make([]Stmt, 0)
	p := NewParser(tokens)

	for !p.isAtEnd() {
		statement, err := p.declaration()
		if err != nil {
			return nil, fmt.Errorf("could not parse: %v", err)
		}
		statements = append(statements, statement)
	}

	return statements, nil
}

type Parser struct {
	Tokens  []Token
	current int // used as a pointer to know to parse the next Token
}

func NewParser(t []Token) *Parser {
	return &Parser{
		Tokens: t,
	}
}

type ParseError struct {
	Token Token
	Msg   string
}

func NewParseError(token Token, msg string) *ParseError {
	return &ParseError{
		Token: token,
		Msg:   msg,
	}
}

// implement error interface (see https://go.dev/tour/methods/19 as reference):
func (e *ParseError) Error() string {
	return fmt.Sprintf("[line %d] Error at \"%s\": %s",
		e.Token.Line, e.Token.Lexeme, e.Msg)
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// rules as functions:

func (p *Parser) declaration() (Stmt, error) {
	if p.match(FUN) {
		fn, err := p.functionDeclaration("function")
		if err != nil {
			return nil, err
		}
		return fn, nil
	}

	if p.match(VAR) {
		vd, err := p.varDeclaration()
		if err != nil {
			return nil, err
		}
		return vd, nil
	}

	stmt, err := p.statement()
	if err != nil {
		// p.decleration gets called repeatedly in Parse.
		// when entering panic mode, go to the beginning
		// of the next stmt or declaration and continue
		// parsing.
		p.synchronize()
	}

	return stmt, nil
}

// can be used for functions (kind="function") or methods (kind="method").
// kind is only relevant for the error message.
func (p *Parser) functionDeclaration(kind string) (Stmt, error) {
	var err error

	name, err := p.consume(IDENTIFIER, "expect "+kind+" name")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(LEFT_PAREN, "expect '(' after "+kind+" name")
	if err != nil {
		return nil, err
	}

	parameters := make([]Token, 0)

	// we are just after the opening LEFT_PAREN. If not closing RIGHT_PAREN, we have args.
	for {
		if len(parameters) > 255 {
			return nil, NewParseError(p.peek(), "can't have more than 255 parameters")
		}

		param, err := p.consume(IDENTIFIER, "expect parameter name")
		if err != nil {
			return nil, err
		}
		parameters = append(parameters, param)

		if p.check(COMMA) {
			p.advance()
			continue
		} else if p.check(RIGHT_PAREN) {
			p.advance()
			break
		} else {
			return nil, NewParseError(p.peek(), "expect COMMA or RIGHT_PAREN after parameter")
		}
	}

	_, err = p.consume(LEFT_BRACE, "expect '{' before "+kind+" body")
	if err != nil {
		return nil, err
	}
	body, err := p.block()
	if err != nil {
		return nil, err
	}

	fn := functionStmt{
		name: name,
		params: parameters,
		body: body,
	}
	return &fn, nil
}

func (p *Parser) varDeclaration() (Stmt, error) {
	name, err := p.consume(IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer Expr
	if p.match(EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	p.consume(SEMICOLON, "Expect ';' after variable declaration.")

	return &varStmt{name, initializer}, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(FOR) {
		return p.forStatement()
	}

	if p.match(IF) {
		return p.ifStatement()
	}

	if p.match(PRINT) {
		return p.printStatement()
	}

	if p.match(WHILE) {
		return p.whileStatement()
	}

	if p.match(LEFT_BRACE) {
		statements, err := p.block()
		if err != nil {
			return nil, err
		}
		return &blockStmt{statements: statements}, nil
	}

	return p.expressionStatement()
}

func (p *Parser) forStatement() (Stmt, error) {
	// here, "desugaring" (page 145-146) happens.

	var err error

	p.consume(LEFT_PAREN, "Expect '(' after 'for'.")

	var initializer Stmt

	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	var condition Expr

	if !p.check(SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	p.consume(SEMICOLON, "Expect ';' after loop condition.")

	var increment Expr
	if !p.check(RIGHT_PAREN) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	p.consume(RIGHT_PAREN, "Expect ')' after for clauses.")

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = &blockStmt{
			statements: []Stmt{
				body,
				&exprStmt{increment},
			},
		}
	}

	if condition == nil {
		condition = &literalExpr{true}
	}

	body = &whileStmt{condition, body}

	if initializer != nil {
		body = &blockStmt{statements: []Stmt{initializer, body}}
	}

	return body, nil
}

func (p *Parser) ifStatement() (Stmt, error) {
	var err error

	_, err = p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}

	var condition Expr
	condition, err = p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(RIGHT_PAREN, "Expect ')' after if condition.")
	if err != nil {
		return nil, err
	}

	var thenBranch, elseBranch Stmt
	thenBranch, err = p.statement()
	if err != nil {
		return nil, err
	}
	if p.match(ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return &ifStmt{
		condition:  condition,
		thenBranch: thenBranch,
		elseBranch: elseBranch,
	}, nil
}

func (p *Parser) printStatement() (Stmt, error) {
	val, err := p.expression()
	if err != nil {
		return nil, err
	}

	p.consume(SEMICOLON, "Expect ';' after value.")

	return &printStmt{Expression: val}, nil
}

func (p *Parser) whileStatement() (Stmt, error) {
	p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(RIGHT_PAREN, "Expect ')' after condition.")
	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &whileStmt{condition, body}, nil
}

func (p *Parser) block() ([]Stmt, error) {
	var err error

	statements := make([]Stmt, 0)

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		statement, err := p.declaration()
		// TODO: how does this play together with synchronize?
		// eventually find out with test cases...
		if err != nil {
			return nil, err
		}
		statements = append(statements, statement)
	}

	_, err = p.consume(RIGHT_BRACE, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}
	return statements, nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	val, err := p.expression()
	if err != nil {
		return nil, err
	}

	p.consume(SEMICOLON, "Expect ';' after value.")

	return &exprStmt{Expression: val}, nil
}

// rule: expression -> assignment ;
func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

// rule: assignment -> IDENTIFIER "=" assignment | equality ;
func (p *Parser) assignment() (Expr, error) {
	var err error

	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if v, ok := expr.(*variableExpr); ok {
			name := v.Name
			return &assignExpr{name, value}, nil
		}

		// TODO: the book remarks:
		// "we report an error but we don't throw it
		// because the parser isn't in a confused state
		// where we need to go into panic mode and sync."
		// well, I do "throw" the error here...
		// how to solve this? Ideas: either just report here and
		// continue, or report at the top of the stack.
		// let's see how to deal with this once I see the problem
		// during experimentation / testing.
		return nil, NewParseError(equals, "invalid assignment target")
	}

	return expr, nil
}

func (p *Parser) or() (Expr, error) {
	var err error
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(OR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = &logicalExpr{expr, operator, right}
	}

	return expr, nil
}

func (p *Parser) and() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = &logicalExpr{expr, operator, right}
	}

	return expr, nil
}

// rule: equality -> comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &binaryExpr{expr, operator, right}
	}

	return expr, nil
}

// comparison -> term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &binaryExpr{expr, operator, right}
	}

	return expr, nil
}

// term -> factor ( ( "-" | "+" ) factor )* ;
func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &binaryExpr{expr, operator, right}
	}

	return expr, nil
}

// factor -> unary ( ( "/" | "*" ) unary )* ;
func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &binaryExpr{expr, operator, right}
	}

	return expr, nil
}

// unary -> ( "-" | "!" ) unary | primary ;
func (p *Parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &unaryExpr{operator, right}, nil
	}

	expr, err := p.call()
	if err != nil {
		return nil, err
	}

	return expr, nil
}

func (p *Parser) call() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(LEFT_PAREN) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) finishCall(callee Expr) (Expr, error) {
	arguments := make([]Expr, 0)

	for {
		if len(arguments) > 255 {
			return nil, NewParseError(p.peek(), "Can't have more than 255 arguments.")
		}

		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, expr)

		if !p.match(COMMA) {
			break
		}
	}

	paren, err := p.consume(RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	call := callExpr{
		callee:    callee,
		paren:     paren,
		arguments: arguments,
	}

	return &call, nil
}

// primary -> NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
func (p *Parser) primary() (Expr, error) {
	if p.match(FALSE) {
		return &literalExpr{false}, nil
	}

	if p.match(TRUE) {
		return &literalExpr{true}, nil
	}

	if p.match(NIL) {
		return &literalExpr{nil}, nil
	}

	if p.match(NUMBER, STRING) {
		return &literalExpr{p.previous().Literal}, nil
	}

	if p.match(IDENTIFIER) {
		return &variableExpr{p.previous()}, nil
	}

	if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}

		return &groupingExpr{expr}, nil
	}

	return nil, &ParseError{p.peek(), "expect expression"}
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// helpers

func (p *Parser) match(tts ...TokenType) bool {
	for _, tt := range tts {
		if p.check(tt) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) consume(tt TokenType, message string) (Token, error) {
	if p.check(tt) {
		return p.advance(), nil
	}

	return Token{}, &ParseError{p.peek(), message}
}

func (p *Parser) check(tt TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tt
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) peek() Token {
	return p.Tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.Tokens[p.current-1]
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().Type == SEMICOLON {
			return
		}

		tt := p.peek().Type

		switch tt {
		case CLASS, FOR, FUN, IF, PRINT, RETURN, VAR, WHILE:
			return
		}

		p.advance()
	}
}
