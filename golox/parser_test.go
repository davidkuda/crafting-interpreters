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

	in := "(21 + 21) * 2 - 42"

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
}
