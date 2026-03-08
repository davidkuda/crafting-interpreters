package main

import "strconv"

// Walks through source code and generates tokens from it.
type Scanner struct {
	// source represents the source code, e.g. `var lang = "golox";`
	Source []byte
	Tokens []Token

	start   int // start of a lexeme
	current int // current character in the loop
	line    int // lines when reading a file (incremented by \n)
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
		s.start = s.current
		s.scanToken()
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.Source)
}

func (s *Scanner) scanToken() {
	c := s.advance()

	switch c {
	// ignore whitespace:
	case ' ':
	case '\r':
	case '\t':
		break

	case '\n':
		s.line++
		break

	// single-character tokens:
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
	case '/':
		if s.nextIsSlash() {
			// comments:
			// comments are lexemes, but they are not meaningful, no further action are based
			// on comments. Therefore, we ignore them and don't create a comment token.
			// In fact, TokenType Comment doesn't even exist. :)
			for s.peek() != '\n' && !s.isAtEnd() {
				_ = s.advance()
			}
		} else {
			s.addToken(SLASH)
		}

	// one or two character tokens:
	case '!':
		if s.nextIsEqualSign() {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '=':
		if s.nextIsEqualSign() {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '<':
		if s.nextIsEqualSign() {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.nextIsEqualSign() {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}

	// longer lexemes:

	// strings:
	case '"':
		s.string()

	default:
		if s.isDigit(c) {
			s.number()
		} else {
			// TODO: report error, maybe with an error struct?
		}
	}
}

func (s *Scanner) addToken(tt TokenType) {
	lexeme := string(s.Source[s.start:s.current])
	t := NewToken(
		tt,
		lexeme,
		nil,
		s.line,
	)
	s.Tokens = append(s.Tokens, t)
}

func (s *Scanner) addTokenWithLiteral(tt TokenType, l any) {
	lexeme := string(s.Source[s.start:s.current])
	t := NewToken(
		tt,
		lexeme,
		l,
		s.line,
	)
	s.Tokens = append(s.Tokens, t)
}

func (s *Scanner) advance() byte {
	char := s.Source[s.current]
	s.current++
	return char
}

func (s *Scanner) nextIsEqualSign() bool {
	return s.match('=')
}

func (s *Scanner) nextIsSlash() bool {
	return s.match('/')
}

func (s *Scanner) match(char byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.Source[s.current] == char {
		s.current++
		return true
	} else {
		return false
	}
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\000'
	}
	return s.Source[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.Source) {
		return '\000'
	}
	return s.Source[s.current+1]
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		// TODO: error:
		// Lox.error(line, "unterminated string")
		return
	}

	s.advance()

	// +1 and -1 are necessary to not include the double quotes in the string
	start := s.start + 1
	end := s.current - 1
	str := string(s.Source[start:end])
	s.addTokenWithLiteral(STRING, str)
}

func (s *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	// fraction?
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	str := string(s.Source[s.start:s.current])
	v, err := strconv.ParseFloat(str, 64)
	if err != nil {
		// TODO: handle error
	}
	s.addTokenWithLiteral(NUMBER, v)
}
