package golox

import "errors"

// "fmt"

// Recursive Descent Parsing:
// top-down parser: start at the top grammar rule (see README)
// and work your way down to the tree leaves.
// Each expression grammar rule becomes a function.
type Parser struct {
	Tokens  []Token
	current int // used as a pointer to know to parse the next Token
}

func NewParser(t []Token) *Parser {
	return &Parser{
		Tokens: t,
	}
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// rules as functions:

// rule: expression -> equality ;
func (p *Parser) expression() Expr {
	return p.equality()
}

/* Java code from the book page 85:
private Expr equality() {
	Expr expr = comparison();
	while (match(BANG_EQUAL, EQUAL_EQUAL)) {
		Token operator = previous();
		Expr right = comparison();
		expr = new Expr.Binary(expr, operator, right);
	}
	return expr
}
*/
// rule: equality -> comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = &Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = &Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = &Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = &Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return &Unary{operator, right}
	}

	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match(FALSE) {
		return &Literal{false}, nil
	}

	if p.match(TRUE) {
		return &Literal{true}, nil
	}

	if p.match(NIL) {
		return &Literal{nil}, nil
	}

	if p.match(NUMBER, STRING) {
		return &Literal{p.previous().Literal}, nil
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return &Grouping{expr}, nil
	}

	return nil, errors.New("could not match primary expression")
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
