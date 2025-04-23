package lexer

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"

	"app/utils"
)

type Lexer struct {
	_reader       *bufio.Reader
	_line, _pos   int64
	_lineLengths  []int64
	_lastRuneSize int
}

func NewLexer(r io.Reader) *Lexer {
	return &Lexer{
		_reader:      bufio.NewReader(r),
		_line:        0,
		_pos:         1,
		_lineLengths: []int64{},
	}
}

func (l *Lexer) NextToken() (Token, error) {
	l.skipWhiteSpace()

	r, err := l.nextRune()
	if err == io.EOF {
		return Token{Type: EOF}, nil
	}

	if r == '/' {
		nextRune, _ := l.nextRune()
		if nextRune == '/' {
			l.skipAnnotation()
			return l.NextToken()
		} else if nextRune == '*' {
			l.skipAnnotation2()
			return l.NextToken()
		} else {
			l.retract()
		}
	}

	if r == '"' {
		s := ""
		escape := false
		for {
			r, err := l.nextRune()
			if err != nil {
				return Token{}, fmt.Errorf("string not closed, line: %d, pos: %d", l._line, l._pos)
			}
			if r == '"' && !escape {
				break
			}
			if r == '\\' {
				escape = !escape
				continue
			}
			if escape {
				switch r {
				case 'n':
					s += "\\n"
				case 't':
					s += "\\t"
				case 'r':
					s += "\\r"
				case 'b':
					s += "\\b"
				case 'f':
					s += "\\f"
				default:
					s += string(r)
				}
				escape = false
				continue
			}
			s += string(r)
		}
		return Token{Type: STRING, Val: s, Line: l._line, Pos: l._pos}, nil
	}

	if utils.IsLetter(r) {
		s := string(r)
		for {
			r, err := l.nextRune()
			if err != nil || !(utils.IsLetter(r) || utils.IsDigit(r) || r == '_') {
				l.retract()
				break
			}
			s += string(r)
		}
		if _BasicType.Contains(s) {
			return Token{Type: TYPE, Val: s, Line: l._line, Pos: l._pos}, nil
		} else if _ReservedWords.Contains(s) {
			return Token{Type: RESERVED, Val: s, Line: l._line, Pos: l._pos}, nil
		} else {
			return Token{Type: IDENTIFIER, Val: s, Line: l._line, Pos: l._pos}, nil
		}
	}

	if utils.IsDigit(r) {
		s := string(r)
		for {
			r, err := l.nextRune()
			if err != nil || !(utils.IsDigit(r) || r == '.') {
				l.retract()
				break
			}
			s += string(r)
		}
		if strings.Contains(s, ".") {
			return Token{Type: FLOAT, Val: s, Line: l._line, Pos: l._pos}, nil
		} else {
			return Token{Type: INTEGER, Val: s, Line: l._line, Pos: l._pos}, nil
		}
	}

	if _Operators.Contains(string(r)) {
		return Token{Type: OPERATOR, Val: string(r), Line: l._line, Pos: l._pos}, nil
	}

	if _Delimiters.Contains(string(r)) {
		return Token{Type: DELIMITER, Val: string(r), Line: l._line, Pos: l._pos}, nil
	}

	return Token{}, fmt.Errorf("unknown character: %c, line: %d, pos: %d", r, l._line, l._pos)
}

func (l *Lexer) nextRune() (rune, error) {
	r, s, err := l._reader.ReadRune()
	l._lastRuneSize = s
	if err != nil {
		return 0, err
	}
	if r == '\n' {
		l._line++
		l._lineLengths = append(l._lineLengths, l._pos)
		l._pos = 0
	} else {
		l._pos++
	}
	return r, nil
}

func (l *Lexer) retract() {
	_ = l._reader.UnreadRune()
	if l._pos > 0 {
		l._pos--
	} else if l._line > 1 {
		l._line--
		l._pos = l._lineLengths[l._line-1]
		l._lineLengths = l._lineLengths[:l._line]
	}
}

func (l *Lexer) skipWhiteSpace() {
	for {
		r, err := l.nextRune()
		if err != nil {
			break
		}
		if !unicode.IsSpace(r) {
			l.retract()
			break
		}
	}
}

func (l *Lexer) skipAnnotation() {
	for {
		r, err := l.nextRune()
		if err != nil || r == '\n' {
			break
		}
	}
}

func (l *Lexer) skipAnnotation2() {
	for {
		r1, err := l.nextRune()
		if err != nil {
			break
		}
		if r1 == '*' {
			r2, err := l.nextRune()
			if err != nil {
				break
			}
			if r2 == '/' {
				break
			}
		}
	}
}
