package golox

import (
	"testing"
)

func TestAST(t *testing.T) {
	expr := &binaryExpr{
		Left: &unaryExpr{
			Operator: NewToken(MINUS, "-", nil, 1),
			Right:    &literalExpr{123},
		},
		Operator: NewToken(STAR, "*", nil, 1),
		Right: &groupingExpr{
			Expression: &literalExpr{45.67},
		},
	}

	s := FormatExpr(expr)

	expected := "(* (- 123) (group 45.67))"

	if s != expected {
		t.Fatalf("ASTPrinter got it wrong: got=%q expected=%q", s, expected)
	}
}
