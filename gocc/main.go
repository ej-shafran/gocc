package main

import (
	"fmt"
	"gocc/lexer"
	"os"
)

func main() {
	file, err := os.Open("return_2.c")
	if err != nil {
		panic(err)
	}

	l := lexer.NewLexer(file)
	for {
		tok, err := l.NextToken()
		if err != nil {
			panic(err)
		}

		if tok.Kind == lexer.EOF {
			break
		}

		fmt.Printf("return_2.c:%d:%d: %s\n", tok.Pos.Line, tok.Pos.Col, tok.Value)
	}
}
