package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

type golox struct {
	hadError bool
}

func main() {
	golox := golox{}

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

func (g *golox) runFile(path string) {
	b, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("could not read file.")
		os.Exit(1)
	}
	g.run(b)
	if g.hadError {
		os.Exit(65)
	}
}

func (g *golox) runPrompt() {
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
	}
}

func (g *golox) run(b []byte) {
	scanner := NewScanner(b)
	scanner.ScanTokens()
	if len(scanner.Errors) > 0 {
		g.hadError = true
	}
	for _, err := range scanner.Errors {
		err.report()
	}
	if len(scanner.Tokens) > 0 {
		fmt.Println("Scanned Tokens:")
		for _, token := range scanner.Tokens {
			fmt.Printf("   %+v\n", token)
		}
		fmt.Println("Scanner Errors:")
		for _, err := range scanner.Errors {
			fmt.Printf("   %+v\n", err)
		}
	}
}
