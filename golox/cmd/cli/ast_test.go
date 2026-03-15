package main

import (
	"testing"
)

func TestAST(t *testing.T) {
	expr := &Binary{
		Left: &Unary{
			Operator: NewToken(MINUS, "-", nil, 1),
			Right:    &Literal{123},
		},
		Operator: NewToken(STAR, "*", nil, 1),
		Right: &Grouping{
			Expression: &Literal{45.67},
		},
	}

	s := FormatExpr(expr)

	expected := "(* (- 123) (group 45.67))"

	if s != expected {
		t.Fatalf("ASTPrinter got it wrong: got=%q expected=%q", s, expected)
	}
}
