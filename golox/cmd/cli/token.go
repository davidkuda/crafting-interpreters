package main

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal string
	Line    int
}

func NewToken(
	ttype TokenType,
	lexeme string,
	literal string,
	line int,
) Token {
	return Token{
		Type:    ttype,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func (t *Token) String() string {
	return t.Lexeme + " " + t.Literal
}

type TokenType int

// TokenType ENUM
const (
	// Single-character tokens.
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// One or two character tokens.
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals.
	IDENTIFIER
	STRING
	NUMBER

	// Keywords.
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE
)
