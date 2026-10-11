package golox

import (
	"errors"
	"fmt"
)

func New() *Golox {
	return &Golox{
		interpreter: NewInterpreter(),
	}
}

type Golox struct {
	scanner     Scanner
	parser      *Parser
	interpreter *Interpreter
}

func (g *Golox) ScanSourceCode(code []byte) error {
	g.scanner = NewScanner(code)
	g.scanner.ScanTokens()
	if len(g.scanner.Errors) > 0 {
		return errors.New("failed scanning golox code")
	}

	return nil
}

func (g *Golox) ReportScanErrors() {
	fmt.Println("Scan Errors:")
	for _, err := range g.scanner.Errors {
		fmt.Println("  - ", err)
	}
}

// use for debugging:
func (g *Golox) PrintTokens() {
	fmt.Println("Tokens:")
	for _, t := range g.scanner.Tokens {
		fmt.Println("  ", t)
	}
}

func (g *Golox) ParseTokens() error {
	g.parser = NewParser(g.scanner.Tokens)
	err := g.parser.Parse()
	if err != nil {
		return fmt.Errorf("failed parsing golox code: %s", err)
	}
	return nil
}

func (g *Golox) ReportParseErrors() {
	fmt.Println("Parse Errors:")
	for _, err := range g.parser.errors {
		fmt.Println("  - ", err)
	}
}

func (g *Golox) InterpretAST() error {
	var err error

	for _, statement := range g.parser.Statements {
		err = g.interpreter.execute(statement)
		if err != nil {
			return fmt.Errorf("cant execute: %v", err)
		}
	}
	return nil
}

func (g *Golox) ReportInterpretErrors() {}
