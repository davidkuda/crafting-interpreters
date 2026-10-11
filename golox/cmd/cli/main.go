package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/davidkuda/golox"
)

type program struct {
	interpreter     *golox.Golox
	hadError        bool
	hadRuntimeError bool
}

func main() {
	program := program{}
	program.interpreter = golox.New()

	if len(os.Args) > 2 {
		fmt.Println("usage: golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		// run a script written in a file:
		filePath := os.Args[1]
		program.runFile(filePath)
	} else {
		// start a REPL:
		program.runPrompt()
	}
}

func (p *program) runFile(path string) {
	b, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("could not read file.")
		os.Exit(1)
	}
	p.run(b)
	if p.hadError {
		os.Exit(65)
	}

	if p.hadRuntimeError {
		os.Exit(70)
	}

	os.Exit(0)
}

func (p *program) runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				// ReadBytes returns io.EOF on CTRL-D
				fmt.Println("exiting program. goodbye.")
				os.Exit(0)
				return
			}
			fmt.Println("bad input:", err)
			os.Exit(1)
		}
		p.run(line)
		fmt.Println()
		p.hadError = false
		p.hadRuntimeError = false
	}
}

func (p *program) run(code []byte) {
	var err error

	err = p.interpreter.ScanSourceCode(code)
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

	err = p.interpreter.InterpretAST()
	if err != nil {
		p.hadRuntimeError = true
		fmt.Printf("failed interpretation: %v\n", err)
		fmt.Println()
		return
	}
}
