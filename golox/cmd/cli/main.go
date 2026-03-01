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
	run(b)
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
		run(line)
	}
}

func run(b []byte) {
	fmt.Println(string(b))
}
