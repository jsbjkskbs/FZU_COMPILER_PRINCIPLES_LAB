package lexer

import (
	"fmt"

	. "app/utils/collections"
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

// ItemType stand for the type of token
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

type TokenSpecificType uint8

// TokenSpecificType stand for the specific type of token
// 1. basic types: int、float、string、bool、byte、int8、int16、int32、int64、uint、uint8、uint16、uint32、uint64、float32、float64
// 2. constants: int、float、char、string、bool
// 3. operators: +、-、*、/、%、=、==、!=、<、<=、>、>=、&&、||、++、--、!、&、|、^、<<、>>
// 4. delimiters: ()、{}、[]、,、;、.、:
// 5. reserved words: break、case、chan、const、continue、default、defer、do、else、false、for、func、go、goto、if、import、interface、map、package、range、return、select、struct、switch、true、type、var
// 6. identifiers: variable names, function names, etc.
// 7. unknown: unknown types, constants, operators, delimiters, reserved words
const (
	Unknown TokenSpecificType = iota
	TypeInt
	TypeInt8
	TypeInt16
	TypeInt32
	TypeInt64
	TypeUnsignedInt
	TypeUnsignedInt8
	TypeUnsignedInt16
	TypeUnsignedInt32
	TypeUnsignedInt64
	TypeFloat
	TypeFloat32
	TypeFloat64
	TypeBool
	TypeString
	TypeByte
	ConstantInt
	ConstantFloat
	ConstantChar
	ConstantStringDoubleQuote
	ConstantStringBacktick
	ConstantBoolTrue
	ConstantBoolFalse
	OperatorPlus
	OperatorMinus
	OperatorMultiply
	OperatorDivide
	OperatorModulo
	OperatorAssignment
	OperatorEqual
	OperatorNotEqual
	OperatorLessThan
	OperatorLessThanOrEqual
	OperatorGreaterThan
	OperatorGreaterThanOrEqual
	OperatorAnd
	OperatorOr
	OperatorIncrement
	OperatorDecrement
	OperatorNot
	OperatorBitwiseAnd
	OperatorBitwiseOr
	OperatorBitwiseXor
	OperatorLeftShift
	OperatorRightShift
	DelimiterLeftParenthesis
	DelimiterRightParenthesis
	DelimiterLeftBrace
	DelimiterRightBrace
	DelimiterLeftBracket
	DelimiterRightBracket
	DelimiterComma
	DelimiterSemicolon
	DelimiterDot
	DelimiterColon
	ReservedWordBreak
	ReservedWordCase
	ReservedWordChan
	ReservedWordConst
	ReservedWordContinue
	ReservedWordDefault
	ReservedWordDefer
	ReservedWordDo
	ReservedWordElse
	ReservedWordFalse
	ReservedWordFor
	ReservedWordFunc
	ReservedWordGo
	ReservedWordGoto
	ReservedWordIf
	ReservedWordImport
	ReservedWordInterface
	ReservedWordMap
	ReservedWordPackage
	ReservedWordRange
	ReservedWordReturn
	ReservedWordSelect
	ReservedWordStruct
	ReservedWordSwitch
	ReservedWordTrue
	ReservedWordType
	ReservedWordVar
	ReservedWordRune
	ReservedWordWhile
	Identifier
)

func (t TokenSpecificType) ToString() string {
	switch t {
	case TypeInt:
		return "int"
	case TypeInt8:
		return "int8"
	case TypeInt16:
		return "int16"
	case TypeInt32:
		return "int32"
	case TypeInt64:
		return "int64"
	case TypeUnsignedInt:
		return "uint"
	case TypeUnsignedInt8:
		return "uint8"
	case TypeUnsignedInt16:
		return "uint16"
	case TypeUnsignedInt32:
		return "uint32"
	case TypeUnsignedInt64:
		return "uint64"
	case TypeFloat:
		return "float"
	case TypeFloat32:
		return "float32"
	case TypeFloat64:
		return "float64"
	case TypeBool:
		return "bool"
	case TypeString:
		return "string"
	case TypeByte:
		return "byte"
	case ConstantInt:
		return "constant_int"
	case ConstantFloat:
		return "constant_float"
	case ConstantChar:
		return "constant_char"
	case ConstantStringDoubleQuote:
		return "constant_string_double_quote"
	case ConstantStringBacktick:
		return "constant_string_backtick"
	case ConstantBoolTrue:
		return "constant_bool_true"
	case ConstantBoolFalse:
		return "constant_bool_false"
	case OperatorPlus:
		return "+"
	case OperatorMinus:
		return "-"
	case OperatorMultiply:
		return "*"
	case OperatorDivide:
		return "/"
	case OperatorModulo:
		return "%"
	case OperatorAssignment:
		return "="
	case OperatorEqual:
		return "=="
	case OperatorNotEqual:
		return "!="
	case OperatorLessThan:
		return "<"
	case OperatorLessThanOrEqual:
		return "<="
	case OperatorGreaterThan:
		return ">"
	case OperatorGreaterThanOrEqual:
		return ">="
	case OperatorAnd:
		return "&&"
	case OperatorOr:
		return "||"
	case OperatorIncrement:
		return "++"
	case OperatorDecrement:
		return "--"
	case OperatorNot:
		return "!"
	case OperatorBitwiseAnd:
		return "&"
	case OperatorBitwiseOr:
		return "|"
	case OperatorBitwiseXor:
		return "^"
	case OperatorLeftShift:
		return "<<"
	case OperatorRightShift:
		return ">>"
	case DelimiterLeftParenthesis:
		return "("
	case DelimiterRightParenthesis:
		return ")"
	case DelimiterLeftBrace:
		return "{"
	case DelimiterRightBrace:
		return "}"
	case DelimiterLeftBracket:
		return "["
	case DelimiterRightBracket:
		return "]"
	case DelimiterComma:
		return ","
	case DelimiterSemicolon:
		return ";"
	case DelimiterDot:
		return "."
	case DelimiterColon:
		return ":"
	case ReservedWordBreak:
		return "break"
	case ReservedWordCase:
		return "case"
	case ReservedWordChan:
		return "chan"
	case ReservedWordConst:
		return "const"
	case ReservedWordContinue:
		return "continue"
	case ReservedWordDefault:
		return "default"
	case ReservedWordDefer:
		return "defer"
	case ReservedWordDo:
		return "do"
	case ReservedWordElse:
		return "else"
	case ReservedWordFalse:
		return "false"
	case ReservedWordFor:
		return "for"
	case ReservedWordFunc:
		return "func"
	case ReservedWordGo:
		return "go"
	case ReservedWordGoto:
		return "goto"
	case ReservedWordIf:
		return "if"
	case ReservedWordImport:
		return "import"
	case ReservedWordInterface:
		return "interface"
	case ReservedWordMap:
		return "map"
	case ReservedWordPackage:
		return "package"
	case ReservedWordRange:
		return "range"
	case ReservedWordReturn:
		return "return"
	case ReservedWordSelect:
		return "select"
	case ReservedWordStruct:
		return "struct"
	case ReservedWordSwitch:
		return "switch"
	case ReservedWordTrue:
		return "true"
	case ReservedWordType:
		return "type"
	case ReservedWordVar:
		return "var"
	case ReservedWordRune:
		return "rune"
	case Identifier:
		return "identifier"
	default:
		return "unknown"
	}
}

type Token struct {
	Type      ItemType
	Val       string
	Line, Pos int64

	_type TokenSpecificType
}

// SpecificType returns the specific type of the token
// It is used to determine the specific type of the token, such as int, float, string, etc.
func (t *Token) SpecificType() TokenSpecificType {
	return t._type
}

// TypeString returns the string representation of the token type
// It is used to determine the type of the token, such as TYPE, INTEGER, FLOAT, etc.
func (t *Token) String() string {
	return fmt.Sprintf("(%v, %s)", t.Type, t.Val)
}

// AllocSize returns the size of the token should be allocated in memory
// It is used to determine the size of the token in memory, such as int, float, string, etc.
func (t *Token) AllocSize() int {
	if t.Type != TYPE {
		return -1
	}
	switch t._type {
	case TypeInt:
		return 4
	case TypeInt8:
		return 1
	case TypeInt16:
		return 2
	case TypeInt32:
		return 4
	case TypeInt64:
		return 8
	case TypeUnsignedInt:
		return 4
	case TypeUnsignedInt8:
		return 1
	case TypeUnsignedInt16:
		return 2
	case TypeUnsignedInt32:
		return 4
	case TypeUnsignedInt64:
		return 8
	case TypeFloat:
		return 4
	case TypeFloat32:
		return 4
	case TypeFloat64:
		return 8
	case TypeBool:
		return 1
	case TypeString:
		return -1 // string is a reference type, so it doesn't have a fixed size
	case TypeByte:
		return 1
	}
	return -1
}

func (t *Token) parseType() {
	switch t.Val {
	case "int":
		t._type = TypeInt
	case "int8":
		t._type = TypeInt8
	case "int16":
		t._type = TypeInt16
	case "int32":
		t._type = TypeInt32
	case "int64":
		t._type = TypeInt64
	case "uint":
		t._type = TypeUnsignedInt
	case "uint8":
		t._type = TypeUnsignedInt8
	case "uint16":
		t._type = TypeUnsignedInt16
	case "uint32":
		t._type = TypeUnsignedInt32
	case "uint64":
		t._type = TypeUnsignedInt64
	case "float":
		t._type = TypeFloat
	case "float32":
		t._type = TypeFloat32
	case "float64":
		t._type = TypeFloat64
	case "bool":
		t._type = TypeBool
	case "string":
		t._type = TypeString
	case "byte":
		t._type = TypeByte
	default:
		t._type = Unknown
	}
}

func (t *Token) parseConstant() {
	switch t.Type {
	case INTEGER:
		t._type = ConstantInt
	case FLOAT:
		t._type = ConstantFloat
	case CHAR:
		t._type = ConstantChar
	// No need to parse string, it should be parsed in the lexer
	case STRING:
	default:
		t._type = Unknown
	}
}

func (t *Token) parseOperator() {
	switch t.Val {
	case "+":
		t._type = OperatorPlus
	case "-":
		t._type = OperatorMinus
	case "*":
		t._type = OperatorMultiply
	case "/":
		t._type = OperatorDivide
	case "%":
		t._type = OperatorModulo
	case "=":
		t._type = OperatorAssignment
	case "==":
		t._type = OperatorEqual
	case "!=":
		t._type = OperatorNotEqual
	case "<":
		t._type = OperatorLessThan
	case "<=":
		t._type = OperatorLessThanOrEqual
	case ">":
		t._type = OperatorGreaterThan
	case ">=":
		t._type = OperatorGreaterThanOrEqual
	case "&&":
		t._type = OperatorAnd
	case "||":
		t._type = OperatorOr
	case "++":
		t._type = OperatorIncrement
	case "--":
		t._type = OperatorDecrement
	case "!":
		t._type = OperatorNot
	case "&":
		t._type = OperatorBitwiseAnd
	case "|":
		t._type = OperatorBitwiseOr
	case "^":
		t._type = OperatorBitwiseXor
	case "<<":
		t._type = OperatorLeftShift
	case ">>":
		t._type = OperatorRightShift
	default:
		t._type = Unknown
	}
}

func (t *Token) parseDelimiter() {
	switch t.Val {
	case "(":
		t._type = DelimiterLeftParenthesis
	case ")":
		t._type = DelimiterRightParenthesis
	case "{":
		t._type = DelimiterLeftBrace
	case "}":
		t._type = DelimiterRightBrace
	case "[":
		t._type = DelimiterLeftBracket
	case "]":
		t._type = DelimiterRightBracket
	case ",":
		t._type = DelimiterComma
	case ";":
		t._type = DelimiterSemicolon
	case ".":
		t._type = DelimiterDot
	case ":":
		t._type = DelimiterColon
	default:
		t._type = Unknown
	}
}

func (t *Token) parseReservedWord() {
	switch t.Val {
	case "break":
		t._type = ReservedWordBreak
	case "case":
		t._type = ReservedWordCase
	case "chan":
		t._type = ReservedWordChan
	case "const":
		t._type = ReservedWordConst
	case "continue":
		t._type = ReservedWordContinue
	case "default":
		t._type = ReservedWordDefault
	case "defer":
		t._type = ReservedWordDefer
	case "do":
		t._type = ReservedWordDo
	case "else":
		t._type = ReservedWordElse
	case "false":
		t._type = ReservedWordFalse
	case "for":
		t._type = ReservedWordFor
	case "func":
		t._type = ReservedWordFunc
	case "go":
		t._type = ReservedWordGo
	case "goto":
		t._type = ReservedWordGoto
	case "if":
		t._type = ReservedWordIf
	case "import":
		t._type = ReservedWordImport
	case "interface":
		t._type = ReservedWordInterface
	case "map":
		t._type = ReservedWordMap
	case "package":
		t._type = ReservedWordPackage
	case "range":
		t._type = ReservedWordRange
	case "return":
		t._type = ReservedWordReturn
	case "select":
		t._type = ReservedWordSelect
	case "struct":
		t._type = ReservedWordStruct
	case "switch":
		t._type = ReservedWordSwitch
	case "true":
		t._type = ReservedWordTrue
	case "type":
		t._type = ReservedWordType
	case "var":
		t._type = ReservedWordVar
	case "rune":
		t._type = ReservedWordRune
	case "while":
		t._type = ReservedWordWhile
	default:
		t._type = Unknown
	}
}

func (t *Token) parseIdentifier() {
	t._type = Identifier
}

func (t *Token) parse() {
	switch t.Type {
	case TYPE:
		t.parseType()
	case INTEGER, FLOAT, CHAR, STRING:
		t.parseConstant()
	case OPERATOR:
		t.parseOperator()
	case DELIMITER:
		t.parseDelimiter()
	case RESERVED:
		t.parseReservedWord()
	case IDENTIFIER:
		t.parseIdentifier()
	default:
		t._type = Unknown
	}
}

var _BasicType = func() Set[string] {
	s := NewSet[string]()
	s.AddAll("int", "float", "string", "bool", "byte", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64")
	return s
}()

var _Operators = func() Set[string] {
	s := NewSet[string]()
	s.AddAll("+", "-", "*", "/", "%", "=", "==", "!=", "<", "<=", ">", ">=", "&&", "||", "++", "--", "!", "&", "|", "^", "<<", ">>")
	return s
}()

var _Delimiters = func() Set[string] {
	s := NewSet[string]()
	s.AddAll("(", ")", "{", "}", "[", "]", ",", ";", ".", ":")
	return s
}()

var _ReservedWords = func() Set[string] {
	s := NewSet[string]()
	s.AddAll("break", "case", "chan", "const", "continue", "default", "defer", "do", "else", "false", "for", "func", "go", "goto", "if", "import", "interface", "map", "package", "range", "return", "select", "struct", "switch", "true", "type", "var", "rune", "while")
	return s
}()
