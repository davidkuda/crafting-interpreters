package main

import "fmt"

// Walks through source code and generates tokens from it.
type Scanner struct {
	// source represents the source code, e.g. `var lang = "golox";`
	Source []byte
	Tokens []Token

	start   int
	current int
	line    int
}

func NewScanner(source []byte) Scanner {
	return Scanner{
		Source: source,
		Tokens: make([]Token, 0),
		line:   1,
	}
}

func (s *Scanner) ScanTokens() {
	for s.current < len(s.Source) {
		fmt.Println(s.current)
		s.scanToken()
	}
}

func (s *Scanner) scanToken() {
	c := s.Source[s.current]
	s.current++
	switch c {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	}
}

func (s *Scanner) addToken(tt TokenType) {
	lexeme := string(s.Source[s.start:s.current])
	t := NewToken(
		tt,
		lexeme,
		"",
		s.line,
	)
	s.Tokens = append(s.Tokens, t)
}
