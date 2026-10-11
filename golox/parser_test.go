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

	g := New()

	for _, test := range tests {

		err = g.ScanSourceCode([]byte(test.input))
		if err != nil {
			t.Fatalf("could not scan input: %s: %v", test.input, err)
		}

		err = g.ParseTokens()
		if err != nil {
			t.Fatalf("could not parse input: %s: %v", test.input, err)
		}

		stmts := g.parser.Statements
		if len(stmts) != 1 {
			t.Fatalf("expected 1 statement, got %d statement", len(stmts))
		}

		s := stmts[0]

		stmt, ok := s.(*exprStmt)
		if !ok {
			t.Fatal("not an exprStmt\n")
		}

		if stmt.Expression == nil {
			t.Fatal("expected exprStmt.Expression")
		}

		expr := stmt.Expression

		out := FormatExpr(expr)
		if out != test.expected {
			t.Fatalf("failed parsing %s: wanted %s, got %s", test.name, out, test.expected)
		}
	}
}
