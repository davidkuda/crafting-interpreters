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

func TestInterpretBinary(t *testing.T) {
	left := Literal{float64(21)}
	right := Literal{float64(21)}

	tkn := NewToken(PLUS, "+", nil, 0)

	binary := Binary{&left, tkn, &right}

	val, err := evaluate(&binary)
	if err != nil {
		t.Fatalf("could not evaluate binary: %v (expected=42 got=%v)", err, val)
	}

	t.Log(val)
}

func TestInterpretIsEqual(t *testing.T) {
	val := isEqual(42, "answer")
	if val != false {
		t.Fatal("isEqual got it wrong")
	}

	val = isEqual(42, 42)
	if val != true {
		t.Fatal("isEqual got it wrong")
	}
}
