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
// 多层if-else嵌套测试
func TestIFLexer(t *testing.T) {
	testStr := `
	if r == '/' {
		if nextRune == '/' && nextRune != '*' {
		} else if nextRune == '*' {
			Println("")
		} else {
			Println("")
		}
	}`
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

// switch语句测试
func TestSwitchLexer(t *testing.T) {
	testStr := `
switch x {
case 1:
    fmt.Println("One")
case 2:
    fmt.Println("Two")
default:
    fmt.Println("Default")
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
