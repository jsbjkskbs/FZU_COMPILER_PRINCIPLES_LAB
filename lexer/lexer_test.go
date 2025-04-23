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

// 测试基本语法
func TestLexer03(t *testing.T) {
	testStr := `
	var b int = 2
	var c string = "123"
	var d bool = true
	var e []int = []int{1, 2, 3}
	f := "Hello, World!"
	f = "Hello, Go!"
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

// 测试字符串中的"注释"
func TestLexer04(t *testing.T) {
	testStr := `
s := "123 /* 456 */ 789"
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

// 测试结构体
func TestLexer05(t *testing.T) {
	testStr := `
type struct_variable_type struct {
	a int
	b string
	c bool
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

// 测试switch语句
func TestLexer06(t *testing.T) {
	testStr := `
	var a int = 3
	switch a {
	case 1:
		println("a is 1")
	case 2:
		println("a is 2")
	default:
		println("a is default")
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

// 测试循环语句
func TestLexer_LoopsAndConditions(t *testing.T) {
	testStr := `
for i := 0; i < 10; i++ {
	if i%2 == 0 {
		println("Even:", i)
	} else {
		println("Odd:", i)
	}
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

// 测试复杂的表达式
func TestLexer_ComplexExpressions(t *testing.T) {
	testStr := `
var result = (10. + 20.) * (3.0 / 5.) - 15.0
println("Result:", result)
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

// 测试字符串和转义字符
func TestLexer_StringsAndEscapes(t *testing.T) {
	testStr := `
var str1 = "Hello, World!"
var str2 = "Line1\nLine2\tTabbed"
println(str1)
println(str2)
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

// 测试注释
func TestLexer_AnnotationsAndComments(t *testing.T) {
	testStr := `
// This is a single-line comment
/*
This is a multi-line comment
spanning multiple lines
*/
var x = 42 // Inline comment
println(x)
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
