package main

import (
	"testing"
)

func TestScanner(t *testing.T) {
	var tests = []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			"single characters",
			"(+++)",
			[]Token{
				{LEFT_PAREN, "(", nil, 1},
				{PLUS, "+", nil, 1},
				{PLUS, "+", nil, 1},
				{PLUS, "+", nil, 1},
				{RIGHT_PAREN, ")", nil, 1},
			},
		},
		{
			"comments",
			"(*) // this lexeme is a comment (without further tokens)\n\t( + )",
			[]Token{
				{LEFT_PAREN, "(", nil, 1},
				{STAR, "*", nil, 1},
				{RIGHT_PAREN, ")", nil, 1},
				{LEFT_PAREN, "(", nil, 2},
				{PLUS, "+", nil, 2},
				{RIGHT_PAREN, ")", nil, 2},
			},
		},
		{
			"strings",
			"\"one\" \"two\" \"three\"",
			[]Token{
				{STRING, "\"one\"", "one", 1},
				{STRING, "\"two\"", "two", 1},
				{STRING, "\"three\"", "three", 1},
			},
		},
	}

	var failed bool

	for _, test := range tests {

		s := NewScanner([]byte(test.input))
		s.ScanTokens()

		if len(test.expected) != len(s.Tokens) {
			t.Logf(
				"[%s]: failed scanning tokens: number of tokens are wrong: expected=%d received=%d\n",
				test.name,
				len(test.expected),
				len(s.Tokens),
			)
			failed = true
		}

		for i, token := range s.Tokens {
			if token != test.expected[i] {
				t.Logf("[%s]: failed scanning tokens:\n", test.name)
				t.Logf("    expected: %v\n", test.expected[i])
				t.Logf("    received: %v\n", token)
				failed = true
			}
		}

		if failed {
			t.FailNow()
		}
	}
}
