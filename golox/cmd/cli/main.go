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
	interpreter     golox.Interpreter
	hadError        bool
	hadRuntimeError bool
}

func main() {
	program := program{}
	program.interpreter = golox.NewInterpreter()

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
	tokens, errs := golox.Scan(code)
	if len(errs) > 0 {
		p.hadError = true
		for _, err := range errs {
			fmt.Println(err)
		}
		return
	}

	ast, err := golox.Parse(tokens)
	if err != nil {
		p.hadError = true
		fmt.Printf("failed parsing tokens: %v\n", err)
		fmt.Println()
		return
	}
	// for _, declaration := range ast {
	// 	fmt.Printf("ast:    %+v\n", declaration)
	// }

	err = p.interpreter.Interpret(ast)
	if err != nil {
		p.hadRuntimeError = true
		fmt.Printf("failed interpretation: %v\n", err)
		fmt.Println()
		return
	}
}
