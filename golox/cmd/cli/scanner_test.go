package main

import (
	"testing"
)

func TestScanSingleCharacterToken(t *testing.T) {
	var input string
	var s Scanner
	var failed bool

	input = "(+++)\n"
	expected := []Token{
		{LEFT_PAREN, "(", "", 1},
		{PLUS, "+", "", 1},
		{PLUS, "+", "", 1},
		{PLUS, "+", "", 1},
		{RIGHT_PAREN, ")", "", 1},
	}

	s = NewScanner([]byte(input))
	s.ScanTokens()
	for i, token := range s.Tokens {
		if token != expected[i] {
			t.Logf("failed scanning tokens: expected=%s received=%s\n", token.Lexeme, expected[i].Lexeme)
			failed = true
		}
	}

	if failed {
		t.FailNow()
	}
}

func TestScanComment(t *testing.T) {
	var input string
	var s Scanner
	var failed bool

	input = "(*) // this lexeme is a comment (without further tokens)\n\t( + )"
	s = NewScanner([]byte(input))
	s.ScanTokens()
	if s.Tokens[0].Type != SLASH {}
	if s.Tokens[0].Literal != input[3:] {}

	expected := []Token{
		{LEFT_PAREN, "(", "", 1},
		{STAR, "*", "", 1},
		{RIGHT_PAREN, ")", "", 1},
		{LEFT_PAREN, "(", "", 2},
		{PLUS, "+", "", 2},
		{RIGHT_PAREN, ")", "", 2},
	}

	if len(expected) != len(s.Tokens) {
		t.Logf("failed scanning tokens: number of tokens are wrongt: expected=%d received=%d\n", len(expected), len(s.Tokens))
		failed = true
	}

	for i, token := range s.Tokens {
		if token != expected[i] {
			t.Logf("failed scanning tokens: expected=%v received=%v\n", token, expected[i])
			failed = true
		}
	}

	if failed {
		t.FailNow()
	}
}

func TestScanString(t *testing.T) {
	var input string
	var s Scanner
	var failed bool

	input = "\"one\" \"two\" \"three\""
	expected := []Token{
		{STRING, "\"one\"", "one", 1},
		{STRING, "\"two\"", "two", 1},
		{STRING, "\"three\"", "three", 1},
	}

	s = NewScanner([]byte(input))
	s.ScanTokens()

	if len(expected) != len(s.Tokens) {
		t.Logf("failed scanning tokens: number of tokens are wrong: expected=%d received=%d\n", len(expected), len(s.Tokens))
		failed = true
	}

	for i, token := range s.Tokens {
		if token != expected[i] {
			t.Logf("failed scanning tokens:\n")
			t.Logf("    expected: %v\n", expected[i])
			t.Logf("    received: %v\n", token)
			failed = true
		}
	}

	if failed {
		t.FailNow()
	}
}
