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

	// TODO: what if !p.isAtEnd() ?
	for !p.isAtEnd() {
		statement, err := p.statement()
		if err != nil {
			// 
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

// implement error interface (see https://go.dev/tour/methods/19 as reference):
func (e *ParseError) Error() string {
	return fmt.Sprintf("[line %d] Error at \"%s\": %s",
		e.Token.Line, e.Token.Lexeme, e.Msg)
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// rules as functions:

func (p *Parser) statement() (Stmt, error) {
	if p.match(PRINT) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement() (Stmt, error) {
	val, err := p.expression()
	if err != nil {
		return Stmt{}, err
	}

	p.consume(SEMICOLON, "Expect ';' after value.")

	return Stmt{Print: val}, nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	val, err := p.expression()
	if err != nil {
		return Stmt{}, err
	}

	p.consume(SEMICOLON, "Expect ';' after value.")

	return Stmt{Expression: val}, nil
}

// rule: expression -> equality ;
func (p *Parser) expression() (Expr, error) {
	return p.equality()
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
		expr = &Binary{expr, operator, right}
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
		expr = &Binary{expr, operator, right}
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
		expr = &Binary{expr, operator, right}
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
		expr = &Binary{expr, operator, right}
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
		return &Unary{operator, right}, nil
	}

	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	return expr, nil
}

// primary -> NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
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
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}

		return &Grouping{expr}, nil
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
