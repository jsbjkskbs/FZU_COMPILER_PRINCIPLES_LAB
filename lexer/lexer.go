package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync/atomic"
	"unicode"

	"app/config"
	"app/utils"
)

var _PrintNoBufferReaderOnce = atomic.Bool{}

type Lexer struct {
	_reader      io.RuneScanner
	_line, _pos  int64
	_lineLengths []int64
}

func NewLexer(r io.Reader) *Lexer {
	var reader io.RuneScanner
	if _PrintNoBufferReaderOnce.CompareAndSwap(false, true) {
		if config.Config.Lexer.UsingNoBufferedReader {
			println("Lexer: using no buffered reader")
		} else {
			println("Lexer: using buffered reader")
		}
	}
	if config.Config.Lexer.UsingNoBufferedReader {
		reader = bufio.NewReaderSize(r, 16)
	}
	reader = bufio.NewReader(r)
	return &Lexer{
		_reader:      reader,
		_line:        0,
		_pos:         1,
		_lineLengths: []int64{},
	}
}

func (l *Lexer) NextToken() (Token, error) {
	if l._reader == nil {
		return Token{}, fmt.Errorf("lexer is not initialized")
	}
	token, err := l.nextToken()
	token.parse()
	return token, err
}

func (l *Lexer) nextToken() (Token, error) {
	err := l.skipWhiteSpace()
	if errors.Is(err, io.EOF) {
		return Token{Type: EOF}, nil
	}

	r, err := l.nextRune()
	if errors.Is(err, io.EOF) {
		return Token{Type: EOF}, nil
	}

	if r == '/' {
		nextRune, err := l.nextRune()
		if errors.Is(err, io.EOF) {
			l.retract()
		} else {
			if nextRune == '/' {
				err := l.skipAnnotation()
				if err != nil {
					if errors.Is(err, io.EOF) {
						return Token{Type: EOF}, nil
					}
					return Token{}, err
				}
				return l.NextToken()
			} else if nextRune == '*' {
				err := l.skipAnnotation2()
				if err != nil {
					if errors.Is(err, io.EOF) {
						return Token{Type: EOF}, nil
					}
					return Token{}, err
				}
				return l.NextToken()
			} else {
				l.retract()
			}
		}
	}

	if r == '"' {
		return l.ReadString()
	}

	if r == '\'' {
		return l.ReadChar()
	}

	if r == '`' {
		return l.ReadString2()
	}

	if utils.IsLetter(r) || r == '_' {
		return l.ReadWord(r)
	}

	if utils.IsDigit(r) {
		return l.ReadNumber(r)
	}

	if _Operators.ContainsFunc(func(s string) bool {
		return strings.HasPrefix(s, string(r))
	}) {
		return l.ReadOperator(r)
	}

	if _Delimiters.Contains(string(r)) {
		return Token{Type: DELIMITER, Val: string(r), Line: l._line, Pos: l._pos}, nil
	}

	return Token{}, fmt.Errorf("unknown character: %c, at line %d, pos %d", r, l._line, l._pos)
}

func (l *Lexer) nextRune() (rune, error) {
	r, _, err := l._reader.ReadRune()
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

func (l *Lexer) skipWhiteSpace() error {
	for {
		r, err := l.nextRune()
		if err != nil {
			return err
		}
		if !unicode.IsSpace(r) {
			l.retract()
			return nil
		}
	}
}

func (l *Lexer) skipAnnotation() error {
	for {
		r, err := l.nextRune()
		if err != nil {
			return err
		}
		if r == '\n' {
			return nil
		}
	}
}

func (l *Lexer) skipAnnotation2() error {
	for {
		r1, err := l.nextRune()
		if err != nil {
			return err
		}
		if r1 == '*' {
			r2, err := l.nextRune()
			if err != nil {
				return err
			}
			if r2 == '/' {
				return nil
			}
		}
	}
}

func (l *Lexer) ReadString() (Token, error) {
	s := ""
	escape := false
	u := ""
	o := ""
	escapeAsUnicodeUpper := false
	escapeAsUnicodeLower := false
	escapeAsOctal := false
	widthOfUnicode := 0
	widthOfOctal := 0
	for {
		r, err := l.nextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return Token{Type: EOF}, fmt.Errorf("string not closed, line %d, pos %d", l._line, l._pos)
			} else {
				return Token{}, err
			}
		}
		if r == '"' && !escape {
			break
		}
		if r == '\\' {
			if escape {
				s += "\\"
			}
			escape = !escape
			continue
		}
		if escape {
			if escapeAsUnicodeLower {
				if !utils.IsHex(r) {
					return Token{}, fmt.Errorf("illegal hex for unicode[lower] %s, at line %d, pos %d", u, l._line, l._pos)
				} else {
					widthOfUnicode++
					u += string(r)
					if widthOfUnicode == 4 {
						s += string(utils.HexToRune(u))
						u = ""
						widthOfUnicode = 0
						escapeAsUnicodeLower = false
						escape = false
					}
				}
			} else if escapeAsUnicodeUpper {
				if !utils.IsHex(r) {
					return Token{}, fmt.Errorf("illegal hex for unicode[upper] %s, at line %d, pos %d", u, l._line, l._pos)
				} else {
					widthOfUnicode++
					u += string(r)
					if widthOfUnicode == 8 {
						s += string(utils.HexToRune(u))
						u = ""
						widthOfUnicode = 0
						escapeAsUnicodeUpper = false
						escape = false
					}
				}
			} else if escapeAsOctal {
				if '0' > r || r > '7' {
					return Token{}, fmt.Errorf("illegal octal %s, at line %d, pos %d", o, l._line, l._pos)
				} else {
					widthOfOctal++
					o += string(r)
					if widthOfOctal == 2 {
						s += string(utils.OctalToRune(o))
						o = ""
						widthOfOctal = 0
						escapeAsOctal = false
						escape = false
					}
				}
			} else {
				switch r {
				case 'n':
					s += "\n"
				case 't':
					s += "\t"
				case 'r':
					s += "\r"
				case 'b':
					s += "\b"
				case 'f':
					s += "\f"
				case 'a':
					s += "\a"
				case 'v':
					s += "\v"
				case 'u': // escape unicode
					escapeAsUnicodeLower = true
				case 'U': // escape unicode
					escapeAsUnicodeUpper = true
				case '0': // escape octal
					escapeAsOctal = true
				case '"':
					s += "\""
				default:
					return Token{}, fmt.Errorf("illegal escape \\%s, at line %d, pos %d", string(r), l._line, l._pos)
				}
			}
			if !escapeAsUnicodeLower && !escapeAsUnicodeUpper && !escapeAsOctal {
				escape = false
			}
			continue
		}
		if r == '\n' {
			if errors.Is(err, io.EOF) {
				return Token{Type: EOF}, fmt.Errorf("string not closed, line %d, pos %d", l._line-1, l._lineLengths[l._line-1])
			} else {
				return Token{}, fmt.Errorf("string not closed, line %d, pos %d", l._line-1, l._lineLengths[l._line-1])
			}
		}
		s += string(r)
	}
	if escapeAsUnicodeLower {
		return Token{}, fmt.Errorf("illegal unicode[lower] %s, at line %d, pos %d", u, l._line, l._pos)
	}
	if escapeAsUnicodeUpper {
		return Token{}, fmt.Errorf("illegal unicode[upper] %s, at line %d, pos %d", u, l._line, l._pos)
	}
	if escapeAsOctal {
		return Token{}, fmt.Errorf("illegal octal %s, at line %d, pos %d", o, l._line, l._pos)
	}
	return Token{Type: STRING, Val: s, Line: l._line, Pos: l._pos, _type: ConstantStringDoubleQuote}, nil
}

func (l *Lexer) ReadString2() (Token, error) {
	s := ""
	for {
		r, err := l.nextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return Token{Type: EOF}, fmt.Errorf("string not closed, line %d, pos %d", l._line, l._pos)
			} else {
				return Token{}, err
			}
		}
		if r == '`' {
			break
		}
		s += string(r)
	}
	return Token{Type: STRING, Val: s, Line: l._line, Pos: l._pos, _type: ConstantStringBacktick}, nil
}

func (l *Lexer) ReadChar() (Token, error) {
	s := ""
	escape := false
	escapeAsUnicodeUpper := false
	escapeAsUnicodeLower := false
	illegalUnicode := false
	width := 0
	for {
		r, err := l.nextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return Token{Type: EOF}, fmt.Errorf("char not closed, line %d, pos %d", l._line, l._pos)
			} else {
				return Token{}, err
			}
		}
		if r == '\'' && !escape {
			break
		}
		if r == '\\' {
			if escape {
				s += "\\"
				width++
			}
			escape = !escape
			continue
		}
		if escape {
			switch r {
			case 'n':
				s += "\n"
			case 't':
				s += "\t"
			case 'r':
				s += "\r"
			case 'b':
				s += "\b"
			case 'f':
				s += "\f"
			case 'u': // escapeAsUnicode
				if escapeAsUnicodeLower || escapeAsUnicodeUpper {
					illegalUnicode = true
				}
				escapeAsUnicodeLower = true
				s += string(r)
			case 'U': // escapeAsUnicode
				if escapeAsUnicodeLower || escapeAsUnicodeUpper {
					illegalUnicode = true
				}
				escapeAsUnicodeUpper = true
				s += string(r)
			default:
				if (escapeAsUnicodeLower || escapeAsUnicodeUpper) && !utils.IsHex(r) {
					illegalUnicode = true
				}
				s += string(r)
			}
			width++
			escape = false
			continue
		}
		s += string(r)
		width++
	}
	// check if the char is valid[not starting with \ and too long]
	if width > 1 && (!escapeAsUnicodeLower && !escapeAsUnicodeUpper) {
		return Token{}, fmt.Errorf("illegal char[too long] %s, at line %d, pos %d", s, l._line, l._pos)
	}
	if escapeAsUnicodeLower {
		if width != 5 {
			return Token{}, fmt.Errorf("illegal char[unmatched unicode length] %s, at line %d, pos %d", s, l._line, l._pos)
		}
		return Token{Type: CHAR, Val: string(utils.HexToRune(s[1:])), Line: l._line, Pos: l._pos}, nil
	}
	if escapeAsUnicodeUpper {
		if width != 9 {
			return Token{}, fmt.Errorf("illegal char[unmatched unicode length] %s, at line %d, pos %d", s, l._line, l._pos)
		}
		return Token{Type: CHAR, Val: string(utils.HexToRune(s[1:])), Line: l._line, Pos: l._pos}, nil
	}
	if (escapeAsUnicodeLower || escapeAsUnicodeUpper) && illegalUnicode {
		return Token{}, fmt.Errorf("illegal char[escapeAsUnicode] %s, at line %d, pos %d", s, l._line, l._pos)
	}
	return Token{Type: CHAR, Val: s, Line: l._line, Pos: l._pos}, nil
}

func (l *Lexer) ReadWord(r rune) (Token, error) {
	s := string(r)
	var errWhenPassed error
	for {
		r, err := l.nextRune()
		if err != nil || !(utils.IsLetter(r) || utils.IsDigit(r) || r == '_') {
			if errors.Is(err, io.EOF) {
				errWhenPassed = io.EOF
			}
			l.retract()
			break
		}
		s += string(r)
	}
	if _BasicType.Contains(s) {
		return Token{Type: TYPE, Val: s, Line: l._line, Pos: l._pos}, errWhenPassed
	} else if _ReservedWords.Contains(s) {
		return Token{Type: RESERVED, Val: s, Line: l._line, Pos: l._pos}, errWhenPassed
	} else {
		return Token{Type: IDENTIFIER, Val: s, Line: l._line, Pos: l._pos}, errWhenPassed
	}
}

func (l *Lexer) ReadNumber(r rune) (Token, error) {
	s := string(r)
	illegalSuffix := false
	tokenWhenWrong := Token{}
	var errWhenPassed error
	for {
		nr, err := l.nextRune()
		if err != nil || !(utils.IsDigit(nr) || utils.IsLetter(nr) || nr == '_' || nr == '.') {
			if errors.Is(err, io.EOF) {
				tokenWhenWrong.Type = EOF
				errWhenPassed = io.EOF
			}
			l.retract()
			break
		}
		if utils.IsLetter(nr) || nr == '_' {
			illegalSuffix = true
		}
		s += string(nr)
	}
	if illegalSuffix && !strings.HasPrefix(s, "0x") && !strings.HasPrefix(s, "0X") {
		return tokenWhenWrong, fmt.Errorf("illegal number[suffix] %s, at line %d, pos %d", s, l._line, l._pos)
	}
	dotCount := strings.Count(s, ".")
	if dotCount == 1 {
		if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
			return tokenWhenWrong, fmt.Errorf("illegal number[hex] %s, at line %d, pos %d", s, l._line, l._pos)
		}
		if strings.HasPrefix(s, "00") {
			parts := strings.Split(s, ".")
			l := utils.RemoveLeadingZeros(parts[0])
			if l == "" {
				l = "0"
			}
			s = l + "." + parts[1]
		}
		return Token{Type: FLOAT, Val: s, Line: l._line, Pos: l._pos}, errWhenPassed
	} else if dotCount == 0 {
		if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
			if len(s) < 3 {
				return tokenWhenWrong, fmt.Errorf("illegal number[hex] %s, at line %d, pos %d", s, l._line, l._pos)
			}
			if strings.ContainsFunc(s[2:], func(r rune) bool {
				return !utils.IsHex(r)
			}) {
				return tokenWhenWrong, fmt.Errorf("illegal number[hex] %s, at line %d, pos %d", s, l._line, l._pos)
			} else {
				return Token{Type: INTEGER, Val: s, Line: l._line, Pos: l._pos}, errWhenPassed
			}
		} else if strings.HasPrefix(s, "0") {
			if len(s) > 1 {
				return tokenWhenWrong, fmt.Errorf("illegal number[integer] %s, at line %d, pos %d", s, l._line, l._pos)
			}
		}
		return Token{Type: INTEGER, Val: s, Line: l._line, Pos: l._pos}, errWhenPassed
	} else {
		return tokenWhenWrong, fmt.Errorf("illegal number[too many dots] %s, at line %d, pos %d", s, l._line, l._pos)
	}
}

func (l *Lexer) ReadOperator(r rune) (Token, error) {
	prefix := string(r)
	previousSet := _Operators
	currentSet := previousSet.Filter(func(s string) bool {
		return strings.HasPrefix(s, prefix)
	})
	bestMatch := ""
	var errWhenPassed error
	tokenWhenError := Token{}
	for {
		if currentSet.Contains(prefix) {
			bestMatch = prefix
		}

		r, err := l.nextRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				errWhenPassed = io.EOF
				tokenWhenError.Type = EOF
			}
			l.retract()
			break
		}
		prefix += string(r)
		currentSet = currentSet.Filter(func(s string) bool {
			return strings.HasPrefix(s, prefix)
		})

		if currentSet.Size() == 0 {
			l.retract()
			break
		}
	}

	if bestMatch == "" {
		// impossible to reach here
		return tokenWhenError, fmt.Errorf("illegal operator %s, at line %d, pos %d", prefix, l._line, l._pos)
	}

	// retract to the best match
	retractStep := len(prefix) - len(bestMatch)
	for i := 0; i < retractStep; i++ {
		l.retract()
	}

	return Token{Type: OPERATOR, Val: bestMatch, Line: l._line, Pos: l._pos}, errWhenPassed
}
