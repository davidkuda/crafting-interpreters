package golox

import (
	"testing"
)

func TestInterpretLiteral(t *testing.T) {
	expr := Literal{42}
	val, err := evaluate(&expr)
	if err != nil {
		t.Fatalf("could not evaluate expr=%s: %v", expr, err)
	}
	if val != 42 {
		t.Fatalf("failed evaluating Literal{42}: expected=42, got=%v", val)
	}
}

func TestInterpretUnary(t *testing.T) {
	var tests = []struct {
		input    Unary
		expected any
	}{
		{
			Unary{NewToken(BANG, "!", nil, 0), &Literal{true}},
			false,
		},
		{
			Unary{NewToken(BANG, "!", nil, 0), &Literal{false}},
			true,
		},
		{
			Unary{
				NewToken(MINUS, "-", nil, 0),
				&Literal{float64(42)},
			},
			-float64(42),
		},
	}

	for _, test := range tests {

		expr := test.input

		val, err := evaluate(&expr)
		if err != nil {
			t.Fatalf("could not evaluate expr=%s: %v", expr.String(), err)
		}

		if val != test.expected {
			t.Fatalf("failed evaluating !true: expected=false, got=%v", val)
		}

	}

}
