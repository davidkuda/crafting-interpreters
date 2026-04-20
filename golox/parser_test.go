package golox

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	var err error

	in := "21 * 2 == 42"

	tokens, errs := Scan([]byte(in))
	if errs != nil {
		t.Fatalf("could not scan input: %s: %v", in, err)
	}

	ast, err := Parse(tokens)
	if err != nil {
		t.Fatalf("could not parse input: %s: %v", in, err)
	}

	out := FormatExpr(ast)
	fmt.Printf("out: %v\n", out)

	expected := "(== (* 21 2) 42)"
	if out != expected {
		t.Fatalf("failed parsing: wanted %s, got %s", out, expected)
	}
}

func TestParserParenthesis(t *testing.T) {
	var err error

	var tests = []struct {
		name     string
		input    string
		expected string
	}{
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
			t.Fatalf("failed parsing: wanted %s, got %s", out, test.expected)
		}
	}
}

func TestParserBooleanComparison(t *testing.T) {
	var err error

	in := "true == !false"

	tokens, errs := Scan([]byte(in))
	if errs != nil {
		t.Fatalf("could not scan input: %s: %v", in, err)
	}

	ast, err := Parse(tokens)
	if err != nil {
		t.Fatalf("could not parse input: %s: %v", in, err)
	}

	out := FormatExpr(ast)
	expected := "(== true (! false))"
	if out != expected {
		t.Fatalf("failed parsing: wanted %s, got %s", out, expected)
	}
}
