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

func TestLexer02(t *testing.T) {
	testStr := `
package main
//这是一条注释
import (
    	"fmt"
)

func main() {
	fmt.Println("Hello, World!"+test(1))
}

func test(v int) bool {
	str := "\"Hello\"\nWorld\""
    fmt.Println(str)
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

// 测试标识符解析
func TestLexerL01(t *testing.T) {
	testStr := `
package main

func myFunction() {}

var myVar int

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

// 测试字符串解析
func TestLexerL02(t *testing.T) {
	testStr := `
str1 := "Hello, World!"
str2 := "This is a test with \\n escape."
str3 := "This is a test with \n escape."
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

// 测试整数和浮点数解析
func TestLexerL03(t *testing.T) {
	testStr := `
var intNum = 123
var floatNum = 123.456
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

// 测试注释解析
func TestLexerL04(t *testing.T) {
	testStr := `
// This is a single-line comment
/* This is a
multi-line comment */
var x = 42
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

// 测试运算符解析
func TestLexerL05(t *testing.T) {
	testStr := `
var result = 10 + 20 - 5 * 2 / 1
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

// 测试复杂表达式解析
func TestLexerL06(t *testing.T) {
	testStr := `
var result = (10 + 5*3) / 2 - 4
var str = "Hello, " + "World!"
var boolVal = true && false || !true
var arr = []int{1, 2, 3, 4, 5}
var mapVar = map[string]int{"one": 1, "two": 2}
var funcVar = func(x int) int {
	return x * x
}
var chanVar = make(chan int)
var interfaceVar interface{} = "This is an interface"
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
