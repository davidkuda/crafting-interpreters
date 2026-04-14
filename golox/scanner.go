package golox

import (
	"fmt"
	"strconv"
)

func Scan(b []byte) ([]Token, []*Error) {
	s := NewScanner(b)
	s.ScanTokens()
	return s.Tokens, s.Errors
}

// Walks through source code and generates tokens from it.
type Scanner struct {
	// source represents the source code, e.g. `var lang = "golox";`
	Source []byte
	Tokens []Token
	Errors []*Error

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
	for !s.isAtEnd() {
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
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			s.Errors = append(s.Errors, NewError(s.line, "unexpected character"))
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
		err := NewError(s.line, "unterminated string")
		s.Errors = append(s.Errors, err)
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
		s.Errors = append(s.Errors, NewError(s.line, fmt.Sprintf("could not parse %s as float", str)))
	}
	s.addTokenWithLiteral(NUMBER, v)
}

func (s *Scanner) isAlpha(c byte) bool {
	isSmallLetter := c >= 'a' && c <= 'z'
	isCapitalLetter := c >= 'A' && c <= 'Z'

	return isSmallLetter || isCapitalLetter || c == '_'
}

func (s *Scanner) isAlphaNumeric(c byte) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := string(s.Source[s.start:s.current])
	isKeyword, keyword := s.isKeyword(text)
	if isKeyword {
		//
		s.addToken(keyword)
	} else {
		s.addToken(IDENTIFIER)
	}
}

func (s *Scanner) isKeyword(text string) (bool, TokenType) {
	switch text {
	case "and":
		return true, AND
	case "class":
		return true, CLASS
	case "else":
		return true, ELSE
	case "false":
		return true, FALSE
	case "for":
		return true, FOR
	case "fun":
		return true, FUN
	case "if":
		return true, IF
	case "nil":
		return true, NIL
	case "or":
		return true, OR
	case "print":
		return true, PRINT
	case "return":
		return true, RETURN
	case "super":
		return true, SUPER
	case "this":
		return true, THIS
	case "true":
		return true, TRUE
	case "var":
		return true, VAR
	case "while":
		return true, WHILE
	}

	return false, 0
}
