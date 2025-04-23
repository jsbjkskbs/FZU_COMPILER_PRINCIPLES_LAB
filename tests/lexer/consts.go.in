package lexer

import (
	"fmt"

	. "app/utils/set"
)

type ItemType uint8

const (
	EOF ItemType = iota + 1
	TYPE
	INTEGER
	FLOAT
	STRING
	CHAR
	OPERATOR
	DELIMITER
	RESERVED
	IMPORT
	PACKAGE
	IDENTIFIER
)

func (i ItemType) ToString() string {
	switch i {
	case EOF:
		return "文件结束符"
	case TYPE:
		return "类型"
	case INTEGER:
		return "整数"
	case FLOAT:
		return "浮点数"
	case STRING:
		return "字符串"
	case CHAR:
		return "字符"
	case OPERATOR:
		return "运算符"
	case DELIMITER:
		return "分隔符"
	case RESERVED:
		return "保留字"
	case IMPORT:
		return "导入"
	case PACKAGE:
		return "包"
	case IDENTIFIER:
		return "标识符"
	default:
		return "未知类型"
	}
}

type Token struct {
	Type      ItemType
	Val       string
	Line, Pos int64
}

func (t *Token) String() string {
	// return fmt.Sprintf("(%s, %s, CharAt{Line: %d, Pos: %d})", t.Type.ToString(), t.Val, t.Line, t.Pos)
	return fmt.Sprintf("(%s, %s)", t.Type.ToString(), t.Val)
}

var _BasicType = func() Set[string] {
	s := New[string]()
	s.AddAll("int", "float", "string", "bool", "byte")
	return s
}()

var _Operators = func() Set[string] {
	s := New[string]()
	s.AddAll("+", "-", "*", "/", "%", "=", "==", "!=", "<", "<=", ">", ">=", "&&", "||", "++", "--", "!", "&", "|", "^", "<<", ">>")
	return s
}()

var _Delimiters = func() Set[string] {
	s := New[string]()
	s.AddAll("(", ")", "{", "}", "[", "]", ",", ";", ".", ":")
	return s
}()

var _ReservedWords = func() Set[string] {
	s := New[string]()
	s.AddAll("break", "case", "chan", "const", "continue", "default", "defer", "do", "else", "false", "for", "func", "go", "goto", "if", "import", "interface", "map", "package", "range", "return", "select", "struct", "switch", "true", "type", "var", "rune")
	return s
}()
