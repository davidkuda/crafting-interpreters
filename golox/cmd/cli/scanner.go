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
		fmt.Println("TOKEN DETECTED")
		break
	}
}

func (s *Scanner) addToken(tt TokenType) {
	t := NewToken(
		tt,
		string(s.Source[s.start:s.current]),
		"",
		s.line,
	)
	s.Tokens = append(s.Tokens, t)
}
