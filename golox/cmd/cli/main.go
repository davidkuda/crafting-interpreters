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
	hadError bool
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
	}
}

func (g *cli) run(b []byte) {
	scanner := golox.NewScanner(b)
	scanner.ScanTokens()
	if len(scanner.Errors) > 0 {
		g.hadError = true
	}
	for _, err := range scanner.Errors {
		err.Report()
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
