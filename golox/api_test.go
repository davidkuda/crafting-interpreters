package golox

import (
	"fmt"
	"testing"
)

type testProgram struct {
	interpreter     *Golox
	hadError        bool
	hadRuntimeError bool
}

func TestGoloxFunctions(t *testing.T) {
	p := testProgram{}
	p.interpreter = New()

	// srcCode := `
	// fun add(a, b) {
	// 	a+b
	// }
	//
	// add(1, 2)
	// `
	srcCode := "add(1, 2)"

	var err error

	err = p.interpreter.ScanSourceCode([]byte(srcCode))
	if err != nil {
		p.hadError = true
		fmt.Println(err)
		p.interpreter.ReportScanErrors()
		return
	}

	err = p.interpreter.ParseTokens()
	if err != nil {
		p.hadError = true
		fmt.Println(err)
		p.interpreter.ReportParseErrors()
		return
	}

	stmts := p.interpreter.parser.Statements
	fmt.Println("len(stmts):", len(stmts))
	fmt.Println(stmts[0])

	stmt, ok := stmts[0].(*exprStmt)
	if !ok {
		t.Fatalf("expected *exprStmt, got %T", stmts[0])
	}

	call, ok := stmt.Expression.(*callExpr)
	if !ok {
		t.Fatalf("expected *callExpr, got %T", stmt.Expression)
	}

	if len(call.arguments) != 2 {
		t.Fatalf("expected 2 argument, got %d", len(call.arguments))
	}

	if call.paren.Lexeme != ")" {
		t.Errorf("expected closing parenthesis, got %q", call.paren.Lexeme)
	}
}
