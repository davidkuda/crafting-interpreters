package golox

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
				{EOF, "", nil, 1},
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
				{EOF, "", nil, 2},
			},
		},
		{
			"strings",
			"\"one\" \"two\" \"three\"",
			[]Token{
				{STRING, "\"one\"", "one", 1},
				{STRING, "\"two\"", "two", 1},
				{STRING, "\"three\"", "three", 1},
				{EOF, "", nil, 1},
			},
		},
		{
			"numbers",
			"\"one\" (42) \"three\"\n108\n42.42",
			[]Token{
				{STRING, "\"one\"", "one", 1},
				{LEFT_PAREN, "(", nil, 1},
				{NUMBER, "42", float64(42), 1},
				{RIGHT_PAREN, ")", nil, 1},
				{STRING, "\"three\"", "three", 1},
				{NUMBER, "108", float64(108), 2},
				{NUMBER, "42.42", float64(42.42), 3},
				{EOF, "", nil, 3},
			},
		},
		{
			"keywords and identifiers",
			"var answer = 42",
			[]Token{
				{VAR, "var", nil, 1},
				{IDENTIFIER, "answer", nil, 1},
				{EQUAL, "=", nil, 1},
				{NUMBER, "42", float64(42), 1},
				{EOF, "", nil, 1},
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
				t.Logf("    expected: %+v\n", test.expected[i])
				t.Logf("    received: %+v\n", token)
				failed = true
			}
		}
	}

	if failed {
		t.FailNow()
	}
}
