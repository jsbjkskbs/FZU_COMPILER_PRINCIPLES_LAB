package lexer_test

import (
	"fmt"
	"strings"
	"testing"

	"app/lexer"
)

func LexerAct(str string) {
	l := lexer.NewLexer(strings.NewReader(str))
	for {
		token, err := l.NextToken()
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		if token.Type == lexer.EOF {
			break
		}
		if err == nil {
			fmt.Println(token.String())
		}
	}
}

func TestLexer01(t *testing.T) {
	testStr := `
package main

import (
	"fmt"
)

func main() {
	var a int = 1
	fmt.Println("Hello, World!")
}
`
	LexerAct(testStr)
}
