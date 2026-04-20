package golox

import (
	"testing"
)

func TestParser(t *testing.T) {
	var err error

	var tests = []struct {
		name     string
		input    string
		expected string
	}{
		{
			"math comparison",
			"21 * 2 == 42",
			"(== (* 21 2) 42)",
		},
		{
			"boolean unary comparison",
			"true == !false",
			"(== true (! false))",
		},
		{
			"simple parenthesis",
			"(42)",
			"(group 42)",
		},
		{
			"another parenthesis",
			"(21 + 21) * 2 - 42",
			"(- (* (group (+ 21 21)) 2) 42)",
		},
	}

	for _, test := range tests {

		tokens, errs := Scan([]byte(test.input))
		if errs != nil {
			t.Fatalf("could not scan input: %s: %v", test.input, err)
		}

		ast, err := Parse(tokens)
		if err != nil {
			t.Fatalf("could not parse input: %s: %v", test.input, err)
		}

		out := FormatExpr(ast)
		if out != test.expected {
			t.Fatalf("failed parsing %s: wanted %s, got %s", test.name, out, test.expected)
		}
	}
}
