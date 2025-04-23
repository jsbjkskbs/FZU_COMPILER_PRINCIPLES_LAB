package lexer_test

import (
	"fmt"
	"strings"
	"testing"

	"app/lexer"
	"app/utils/log"
)

func LexerAct(str string) {
	l := lexer.NewLexer(strings.NewReader(str))
	for {
		token, err := l.NextToken()
		if err != nil {
			fmt.Println(
				log.Sprintf(log.Argument{FrontColor: log.Red, Highlight: true, Format: "Error: %s", Args: []any{err.Error()}}),
			)
		}
		if token.Type == lexer.EOF {
			break
		}
		if err == nil {
			fmt.Printf(
				"(%s, %s)\n",
				log.Sprintf(log.Argument{FrontColor: log.Green, Format: "%s", Args: []any{token.Type.ToString()}}),
				log.Sprintf(log.Argument{FrontColor: log.Yellow, Format: "%s", Args: []any{token.Val}}),
			)
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
1aaa
`
	LexerAct(testStr)
}
