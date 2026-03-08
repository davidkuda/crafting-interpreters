package main

import "fmt"

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any // java Object =~ go any
	Line    int
}

func NewToken(
	ttype TokenType,
	lexeme string,
	literal any,
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
	return t.Lexeme + " " + fmt.Sprintf("%v", t.Literal)
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
	FOR
	FUN
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
