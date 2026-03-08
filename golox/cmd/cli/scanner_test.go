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
