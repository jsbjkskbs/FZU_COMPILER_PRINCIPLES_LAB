package main

import (
	"fmt"
	"os"

	"app/lexer"
)

var (
	l *lexer.Lexer
)

func main() {
	file, err := os.Open(`tests/1.go`)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	l = lexer.NewLexer(file)

	for {
		token, _ := l.NextToken()
		if token.Type == lexer.EOF {
			break
		}
		fmt.Println(token.String())
	}
}
