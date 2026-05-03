package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/davidkuda/golox"
)

type cli struct {
	hadError        bool
	hadRuntimeError bool
}

func main() {
	golox := cli{}

	if len(os.Args) > 2 {
		fmt.Println("usage: golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		// run a script written in a file:
		filePath := os.Args[1]
		golox.runFile(filePath)
	} else {
		// start a REPL:
		golox.runPrompt()
	}
}

func (g *cli) runFile(path string) {
	b, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("could not read file.")
		os.Exit(1)
	}
	g.run(b)
	if g.hadError {
		os.Exit(65)
	}

	if g.hadRuntimeError {
		os.Exit(70)
	}

	os.Exit(0)
}

func (g *cli) runPrompt() {
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
		g.run(line)
		g.hadError = false
		g.hadRuntimeError = false
	}
}

func (g *cli) run(code []byte) {
	tokens, errs := golox.Scan(code)
	if len(errs) > 0 {
		g.hadError = true
		for _, err := range errs {
			fmt.Println(err)
		}
		return
	}

	ast, err := golox.Parse(tokens)
	if err != nil {
		g.hadError = true
		fmt.Printf("failed parsing tokens: %v\n", err)
		return
	}
	
	fmt.Println(ast)

	val, err := golox.Interpret(ast)
	if err != nil {
		g.hadError = true
		fmt.Printf("failed interpretation: %v\n", err)
		return
	}

	fmt.Println(val)
}
