package lexer_test

import (
	"strings"
	"testing"

	"app/lexer"
)

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
	l := lexer.NewLexer(strings.NewReader(testStr))
	for {
		token, err := l.NextToken()
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		if token.Type == lexer.EOF {
			break
		}
		t.Log(token.String())
	}
}
