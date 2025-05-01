# Lexer

## Purpose of this Lab
This lab requires students to write lexical analysis functions to identify individual meaningful words from the input source code, including basic reserved words, identifiers, constants, operators, and delimiters. The goal of lexical analysis is to divide the input source code into a sequence of basic symbols (tokens) and filter out comments and whitespace such as spaces and line breaks. Lexical analysis provides the token sequence for the next phase, syntax analysis.

## Tasks
Implement the following functionalities for the compiler:
1. Tokenize words and output them in tuple format (word name, word category);
2. Remove comments;
3. Remove whitespace (spaces, line breaks, tabs);
4. Detect and locate errors;
5. Build a symbol table and print it in order.

## Implementation Details

### 1. File Reading

The lab utilizes both `mmap` and `bufio.Reader` for reading file content.

- `mmap` maps the file directly into memory, avoiding multiple system calls typical of traditional file I/O. This significantly improves read efficiency and allows direct memory access to file content, making it suitable for handling large files.
- `bufio.Reader` provides a buffering mechanism that optimizes the performance of reading small chunks of data by reducing the number of direct reads from the underlying file.

By combining the strengths of both approaches, `mmap` offers memory-like access, simplifying file content processing, while `bufio.Reader` provides flexible reading interfaces (e.g., `ReadRune` and `UnreadRune`) for handling `UTF-8` characters. Additionally, `mmap` supports random access to file content, making it ideal for scenarios requiring frequent jumps, whereas `bufio.Reader` allows configurable buffer sizes, balancing memory usage and performance.

> However, in this lab, the performance difference between `mmap` and traditional I/O methods is negligible.
> For a 100,000+ lines of code file (created by duplicating 1,000+ lines of code), the time difference between `mmap` and traditional I/O methods is minimal, regardless of whether a 4KB or 16B buffer size is used.
> ```
> $ ./bin/mmap -lexer--no-buffered -b -s
> mmap: using windows implementation
> Benchmark mode enabled
> ====================
> ***  Lexer Test  ***
> *** Got 1 Files ***
> ====================
> << tests\benchmark\lexer\bits_test.go.in
> ====================
> !!! Starting tests... !!!
> ====================
> Lexer: using no buffered reader
> >> Test for tests\benchmark\lexer\bits_test.go.in finished, consume 275 ms
> ====================
> !!! All tests finished !!!
> ====================
> ```
> ```
> .\bin\traditional.exe -lexer--no-buffered -b -s
> mmap: using default implementation[os.File]
> Benchmark mode enabled
> ====================
> ***  Lexer Test  ***
> *** Got 1 Files ***
> ====================
> << tests\benchmark\lexer\bits_test.go.in
> ====================
> !!! Starting tests... !!!
> ====================
> Lexer: using no buffered reader
> >> Test for tests\benchmark\lexer\bits_test.go.in finished, consume 269 ms
> ====================
> !!! All tests finished !!!
> ====================
> ```
> 
> - In the lab, file reading is performed sequentially. In such cases, the operating system's file caching mechanisms (e.g., page cache) significantly optimize the performance of traditional I/O. The advantage of `mmap` is more evident in random access scenarios, as it allows direct access to any part of the file via memory addresses without frequent system calls. However, this advantage is not prominent in sequential reading.
> - The `UnreadRune` operation of `bufio.Reader` allows previously read characters to be pushed **back into** the buffer. However, this operation typically only involves data within the buffer and does not trigger random access to the underlying file. Therefore, the use of `UnreadRune` does not significantly increase the frequency of random access, and file reading remains primarily sequential.
> - With larger buffer sizes (e.g., 4KB), the performance of traditional I/O is already very close to that of `mmap`, as fewer system calls and larger data block reads reduce performance overhead. Even with smaller buffer sizes (e.g., 16B), the sequential reading characteristic ensures that the performance of traditional I/O does not degrade significantly.

### 2. Code File Tokenization

#### 2.1 Reference Specification
In this lab, the design of the lexical analyzer is based on the syntax rules of the Go programming language. Specifically, the definitions and handling of basic symbols (tokens) such as keywords and strings follow the Go language specification. This design choice leverages Go's clear syntax structure and wide application scenarios, providing a reliable reference framework for the lab (using Go as the reference is more direct since the implementation is in Go).

#### 2.2 Implementation Approach
This program does not directly implement a finite state machine (FSM) but instead simulates similar functionality through simple branch conditions. This implementation approach is primarily due to the static nature of lexical rules in high-level languages, which are well-defined at compile time. Branch conditions provide a more intuitive way to implement these rules without introducing the complexity of FSMs. FSMs are better suited for dynamic, runtime requirements, especially in scenarios requiring frequent extensions or modifications of states and transition rules. However, the lexical analysis tasks in this lab are relatively fixed, and branch-based implementation is sufficient to meet the requirements. Additionally, branch logic is clear, easy to understand, and debug, enabling quick implementation of lexical analysis functionality while avoiding the extra overhead of FSM state design and transition table maintenance.

In simple terms, the program logic is as follows:

1. **Branch Conditions Replace State Transitions**  
    FSMs typically handle input characters through states and transition tables, whereas this program directly uses conditional branches (e.g., `if`, `switch`) to determine the type of the current character and its processing logic. Each character type (e.g., letters, numbers, operators) corresponds to a branch that directly invokes the appropriate processing function (e.g., `ReadWord`, `ReadNumber`).

2. **Static Rules Replace Dynamic States**  
    Lexical rules in high-level languages are usually static and well-defined at compile time. This program uses hardcoded rules (e.g., `Operators`, `Delimiters` collections) to directly match input characters without dynamically maintaining states and transition tables. For example, operator matching is implemented through prefix matching in the `Operators` collection rather than state transitions.

3. **Character-by-Character Reading Replaces State Flow**  
    FSMs rely on state transitions, while this program reads the input stream character by character and immediately determines the type and processing logic for each character. For instance, when encountering a letter or an underscore (`_`), the program calls `ReadWord` to read the complete word; when encountering a digit, it calls `ReadNumber` to parse the full number.

4. **Error Handling Embedded in Branch Logic**  
    FSMs typically handle illegal input by transitioning to specific error states, whereas this program directly detects illegal characters or formats within the branch logic and returns errors.

![lexer logic](/docs/img/lexer/1.png)

#### 2.3 Implementation Details

Before entering the main decision logic, the `Lexer` skips comments and whitespace. It is important to note that some logical sequences may require reading multiple characters to determine the entry condition. If the condition is not met, the read characters must be pushed back into the buffer.

The `Lexer` maintains the state of `Line` and `Column` to facilitate error localization.

```go
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
```

##### 2.3.1 Whitespace

Whitespace includes spaces, tabs, and newlines. The program skips these characters until the next valid character is encountered.

```go
// Skip whitespace before entering decision branches
// Whitespace includes spaces, tabs, and newlines
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
```

##### 2.3.2 Annotations

According to Go's syntax rules, annotations are divided into two types:
1. Single-line annotations starting with `//`
2. Multi-line annotations starting with `/*` and ending with `*/`

For case `1.`, the program reads the entire line until the end.

> Most modern operating systems use CRLF (`\r\n`) or LF (`\n`) as line-ending sequences, while a few use CR (`\r`). The program uses `\n` as the line terminator to ensure compatibility with most systems.

```go
// The decision branch has already read `//`, now skip the comment
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
```

For case `2.`, the program checks for the multi-line annotations terminator `*/`. If not found, it continues reading subsequent lines until the terminator is encountered.

```go
// The decision branch has already read `/*`, now skip the comment
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
```

##### 2.3.3 Strings

The program supports escape characters in strings, including `\n`, `\t`, `\r`, `\b`, `\f`, `\a`, `\v`, and `\"`. It also supports Unicode escape sequences (`\u` and `\U`) and octal escape sequences (`\0`). These escape sequences are converted into their corresponding characters, which are reflected in the `Val` field of the `Token`.

In Go, strings are enclosed by either double quotes (`"`) or backticks (`` ` ``).

Double-quoted strings can contain escape characters (e.g., `\n`, `\t`). When reading double-quoted strings, the program processes escape characters and converts them into their corresponding characters.

Since strings may contain escape sequences, using a finite state machine (FSM) to handle string parsing can be complex. Instead, extracting the escape state into a separate logic makes the implementation clearer. While reading a string, the program checks if the current character is the escape character `\`. If so, it enters the escape state and determines how to handle subsequent characters based on the type of escape sequence.

The program supports the following string formats:
- Plain strings: `"hello"`, `"world"`, `"hello world"`
- Escaped strings: `"hello\nworld"`, `"hello\tworld"`, `"hello\rworld"`, `"hello\bworld"`, `"hello\fworld"`, `"hello\aworld"`, `"hello\vworld"`, `"hello\"world"`
- Unicode strings: `"好"`
- Unicode escaped strings: `"\u4e2d"`, `"\U00004e2d"`
- Octal escaped strings: `"\040"`
- Backtick-enclosed strings: `` `hello` ``, `` `hello\nworld` ``, `` `hello\tworld` ``, `` `hello\rworld` ``, `` `hello\bworld` ``, `` `hello\fworld` ``, `` `hello\aworld` ``, `` `hello\vworld` ``, `` `hello\"world` ``

```go
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
				if !utils.IsOctal(r) {
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
				case 'n', 't', 'r', 'b', 'f', 'a', 'v', '"':
					s += utils.AppendEscape(r)
				case 'u': // escape unicode
					escapeAsUnicodeLower = true
				case 'U': // escape unicode
					escapeAsUnicodeUpper = true
				case '0': // escape octal
					escapeAsOctal = true
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
```

For backtick-enclosed strings, the program reads directly until the next backtick is encountered.

```go
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
```

##### 2.3.4 Characters

Character constants are enclosed in single quotes `''` and support escape characters (e.g., `\'`, `\n`). When reading a character, the program needs to determine whether the current character is an escape character `\`. If it is, the program enters the escape state. Based on the type of escape character, the program decides how to handle subsequent characters.

The program can validate whether a character constant is legal, i.e., whether it contains only one character. If the character constant contains multiple characters, an error is returned.

It is important to note that modern programming languages provide support for Unicode characters, and character constants can also include Unicode characters. Since the length of Unicode characters is often greater than one byte, the program must account for the length of Unicode characters when reading character constants.

The program supports Unicode characters through the escape sequences `\u` and `\U`. For the `\u` escape sequence, the program reads 4 hexadecimal digits and converts them into the corresponding Unicode character. For the `\U` escape sequence, the program reads 8 hexadecimal digits and converts them into the corresponding Unicode character (the details of the corresponding Unicode character can be seen in the Token output).

The program supports the following character formats:
- Single characters: `'a'`, `'_'`, `'1'`
- Escape characters: `'\n'`, `'\t'`, `'\r'`, `'\b'`, `'\f'`, `'\a'`, `'\v'`, `'\\'`
- Unicode characters: `'好'`
- Unicode escape sequences: `'\u4e2d'`, `'\U00004e2d'`

```go
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
			case 'n', 't', 'r', 'b', 'f', 'a', 'v':
				s += utils.AppendEscape(r)
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
```

##### 2.3.5 Words

Words consist of letters, digits, and underscores `_`, and must start with a letter or an underscore `_`. The program needs to determine whether the current character is a letter or an underscore `_`. If so, it enters the word-reading state. Based on the type of the word, the program decides how to handle subsequent characters.

Words may be reserved keywords. Reserved keywords are predefined in programming languages and have special syntactic meanings. The program needs to check whether the current word is a reserved keyword. If it is, it is marked as a reserved keyword; otherwise, it is marked as an identifier.

It is evident that the logic for determining identifiers and keywords is logically independent but implementation-wise interdependent.

Unfortunately, the program does not support characters other than letters, digits, and underscores `_` as part of a word (e.g., Chinese characters in Unicode). Readers can extend this functionality as needed.

The program supports the following word formats:
- Identifiers: `a`, `_a`, `a1`, `_a1`, `a_1`, `_a_1`
- Reserved keywords: `int`, `float`, `string`, `if`, `else`, `for`, `while`, `return`, `break`, `continue`

```go
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
```

##### 2.3.6 Numbers

Numeric constants consist of digits and support both integers and floating-point numbers. The program needs to determine whether the current character is a digit. If it is, the program enters the number-reading state. Based on the type of number, the program decides how to handle subsequent characters.

The program does not support scientific notation (e.g., `1.23e+4` or `1.23E+4`) or negative numbers (e.g., `-1` or `-1.23`). Additionally, numbers starting with a `.` (e.g., `.23`) are not supported.
> For negative numbers, we prefer to handle them during the syntax analysis phase rather than the lexical analysis phase.
> For scientific notation, the program could implement support through simple string concatenation. However, to maintain simplicity and readability for this lab, it is not supported at this time.

The program only supports integers and floating-point numbers. Integers can be in decimal or hexadecimal format, and floating-point numbers can start with multiple leading zeros in decimal format.
> Modern programming languages generally support hexadecimal numbers, which start with `0x` or `0X` followed by digits and letters A-F (or a-f).
> The current implementation provides sufficient support for most programming language requirements. Readers can extend the functionality as needed.

The supported number formats include:
- Integers: `123`, `0x123`, `0X123`, `0xABCDEF`, `0XABCDEF`
- Floating-point numbers: `123.456`, `0.123456`, `000.123456`

> Note: Negative numbers such as `-1`, `-1.23`, `-0x123`, `-0xABCDEF`, `-0X123`, and `-0XABCDEF` do not need to be considered.
>  - In lexical analysis, considering numbers with a leading negative sign increases the complexity of the lexer, as it requires handling the special case of the negative sign.
>  - In syntax analysis, the negative sign can be treated as a separate operator, which can be handled during the syntax analysis phase rather than during the lexical analysis phase.
>    - For example, numbers can be defined as `a | b`, where `a` is a number and `b` is a negative sign followed by a number.
>    - By performing lookahead to identify symbols or similar markers, it can be determined whether `-` is a negative sign or a subtraction operator.
>  - This approach simplifies the design of the lexer, allowing it to focus on recognizing basic lexical tokens rather than handling complex syntactic rules.

```go
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
```

##### 2.3.7 Operators

Operators are symbols or keywords in programming languages used to perform specific operations. They are typically used to compute, compare, or perform logical operations on variables, constants, or expressions. Operators can be categorized into the following types:
- **Arithmetic Operators**: Used for mathematical operations, such as `+`, `-`, `*`, `/`, `%`, etc.
- **Relational Operators**: Used to compare two values, such as `<`, `>`, `<=`, `>=`, `==`, `!=`, etc.
- **Logical Operators**: Used for logical operations, such as `&&`, `||`, `!`, etc.
- **Bitwise Operators**: Used to manipulate binary bits, such as `&`, `|`, `^`, `<<`, `>>`, etc.
- **Assignment Operators**: Used to assign values to variables, such as `=`, `+=`, `-=`, etc.

Operators can vary in length, ranging from a single character (e.g., `+`, `-`, `*`, `/`, `=`) to multiple characters (e.g., `==`, `!=`, `<=`, `>=`, `&&`, `||`). Some operators are combinations of other operators (e.g., `+=`, `-=`, `&=`, `|=`).

An operator may serve as a prefix (e.g., `+`, `-`, `*`, `/`, `=`) or a suffix (e.g., `+=`, `-=`, `&=`, `|=`). To identify the type of an operator, it is necessary to search forward for the longest matching operator.

For example, consider the operators `<` and `<<`. When encountering `<`, the next character must be checked to determine if it forms `<<`. If it does, the operator is identified as `<<`; otherwise, it remains `<`. For `<<<`, it can be split into `<<` and `<`, where `<<` is the first operator and `<` is the next.

A highly efficient way to implement this is by using a prefix tree (Trie) to store operators. A Trie is a tree-like data structure where each node represents a character, and the combination of characters along a path represents an operator. Using a Trie allows for quick identification of the longest matching operator.
- Each Trie node represents a character, and the combination of characters along a path represents an operator.
- The root node represents an empty string, and child nodes represent operator prefixes.
- Leaf nodes represent complete operators.
- Non-leaf nodes represent operator prefixes and allow further searching.
- Each Trie node can include a flag indicating whether it represents a valid operator.

For example:
1. Define the operator set `{+, +=, ++, -, -=, (= }`. Note that `(` is not an operator in this set.
2. `+` and `+=` share a prefix relationship, as do `-` and `-=`, and `(` and `=`.
3. The program starts at the root node and searches layer by layer until the longest matching operator is found.
    - For the input `+`, the program starts at the root, finds the child node for `+`, and returns `+`.
    - For the input `+=`, the program starts at the root, finds the child node for `+`, then continues to the child node for `=`, and returns `+=` as the longest matching operator.
    - For the input `+++`, the program starts at the root, finds the child node for `+`, then continues to the child node for `+`, and returns `++` as the longest matching operator. However, since `+` is also an operator, `+++` can be split into `++` and `+`, where `++` is the first operator and `+` is the next operator, determined in the next search starting from the root.
    - For the input `(`, the program starts at the root, finds the child node for `(`, but since `(` is not marked as an operator, it is not returned as an operator.
    - For the input `(=`, the program starts at the root, finds the child node for `(`, then continues to the child node for `=`, and returns `(=`. Even though `(` is not marked as an operator, the combination `(=` is a valid operator (as indicated by the flag in the child node), so `(=` is returned.
   
![Trie Tree](/docs/img/lexer/2.png)

Since this program does not aim to introduce the complexity of a Trie data structure, a simple string set filter is used to match operators.

1. **Initialize Prefix and Candidate Set**  
    1. Use the current character as the initial prefix `prefix`.  
    2. Use the `Operators` set (containing all valid operators) as the initial candidate set `currentSet`.  
    3. Filter the `Operators` set to retain only those operators that start with `prefix`, narrowing down the candidate set.

2. **Extend Prefix Character by Character**  
    1. Read the next character and append it to `prefix`.  
    2. Update `currentSet` to retain only operators that start with the new `prefix`.  
    3. If `currentSet` becomes empty, it means no operator starts with the current `prefix`, and the extension stops.

3. **Record Longest Match**  
    1. During each prefix extension, check if `currentSet` contains the current `prefix`.  
    2. If it does, record it as the current best match `bestMatch`.

4. **Handle Termination Conditions**  
    1. If the end of the file (`EOF`) is reached, stop extending and backtrack to the position of the best match.  
    2. If `currentSet` becomes empty, backtrack to the previous character position.

5. **Backtrack to Best Match**  
    1. Calculate the number of characters to backtrack (`len(prefix) - len(bestMatch)`).  
    2. Use the backtrack function to return the extra characters to the buffer.

```go
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
	for range retractStep {
		l.retract()
	}

	return Token{Type: OPERATOR, Val: bestMatch, Line: l._line, Pos: l._pos}, errWhenPassed
}
```

#### 2.4 Token Structure
The `Token` structure represents the tokens generated by the lexical analyzer. It contains information such as the token's type, value, line number, and column position.

The definition of the `Token` structure is as follows:
- `Type`: The type of the token, represented using the `ItemType` enumeration.
- `Val`: The value of the token, represented as a string.
- `Line`: The line number where the token is located, represented as an integer.
- `Pos`: The position of the token within the line, represented as an integer.
- `_type`: The specific type of the token, represented using the `TokenSpecificType` enumeration.

```go
type Token struct {
	Type                ItemType
	Val                 string
	Line, Pos           int64

	_type TokenSpecificType
}
```

### 2.5 Tuple Output Format
The output format of the `Token` structure is a tuple, containing the type and value of the Token.

```go
func (t Token) String() string {
    return fmt.Sprintf("(%s, %s)", t.Type, t.Val)
}
```

The output format of the `Token` structure is as follows:
```
(保留字, package)
(标识符, log)
(保留字, import)
(分隔符, ()
(字符串, fmt)
(分隔符, ))
(保留字, type)
(标识符, Color)
(类型, string)
(保留字, const)
(分隔符, ()
(标识符, Black)
(标识符, Color)
```

### 2.6 How to use the Lexer
The `Lexer` requires an `io.Reader` as its input source, which can typically be implemented using `os.Stdin` or a file reader. The file reader can be either `os.File` or `mmap`.

Once the `Lexer` is initialized, the `NextToken()` method can be used to retrieve the next token. The `NextToken()` method returns a `Token` structure and an error. If no error occurs, the error will be `nil`.

``` go
type Lexer struct {
	// ...
}

func NewLexer(r io.Reader) *Lexer {
	// ...
}

func (l *Lexer) NextToken() (Token, error) {
	// ...
}
```

### 2.7 Testing the Lexer

The `Lexer` is tested using the `testing` package. The test cases cover various types of tokens, including reserved keywords, identifiers, strings, characters, numbers, and operators.

You can run the test cases using the `go test app/lexer` command. Adding the `-v` flag will display detailed test output.

```bash
go test app/lexer -v
```

### 2.8 Design and Implementation of the Symbol Table

Since the lexical analyzer cannot determine the context of a token, it cannot identify the type of a statement (e.g., declaration, assignment, or function call). This determination must be deferred to the syntax analysis phase.

However, we can briefly introduce the design and implementation of the symbol table here.

The symbol table is an important concept in programming languages, used to store information about identifiers such as variables, functions, and classes. It is typically utilized during the syntax and semantic analysis phases of a compiler. The design and implementation of a symbol table can be divided into the following components:
- **Structure of the Symbol Table**: The symbol table is often implemented using a hash table or a tree structure. Hash tables provide faster lookup speeds but slower insertion and deletion operations, while tree structures offer slower lookup speeds but faster insertion and deletion operations.
- **Scope Management**: The symbol table typically supports multiple levels of scope, allowing the same identifier to be defined in different scopes. The symbol table must support operations for entering and exiting scopes, so that a new symbol table is created when entering a new scope and the current symbol table is destroyed when exiting a scope.
- **Lookup Operations**: The symbol table must support looking up information based on an identifier's name. Lookup operations usually start from the current scope and proceed upward through enclosing scopes until the identifier is found.
- **Insertion Operations**: The symbol table must support inserting new identifier information. Insertions are typically performed in the current scope. If an identifier with the same name already exists in the current scope, an error should be reported.

From this, it is evident that the symbol table functions as a contextual environment, storing information about identifiers such as variables, functions, and classes, while supporting multi-level scopes and lookup operations.

A natural data structure to consider for this purpose is a singly linked list. Each node in the linked list can store information about an identifier and can point to the next node. However, the direction of linkage in a symbol table is opposite to that of a typical singly linked list. In a symbol table, the linkage direction is from the current scope upward, whereas in a singly linked list, the linkage direction is from the top downward.

Here is a simple example to illustrate this concept (using `<-` to indicate the currently processed statement block):

```
{
    int a;
    int b; <-
    {
        int c;
        a = 1;
        b = 2;
        int b;
        c = 3;
    }
    int c;
}

SymbolTable:
scope[0]: a:int:0x00000000; b:int:0x00000004;
```

```
{[0]
    int a;
    int b;
    {[1]
        int c;
        a = 1;
        b = 2; <-
        int b;
        c = 3;
    }
    int c;
}

SymbolTable:
scope[0]: a:int:0x00000000; b:int:0x00000004;
scope[1]: c:int:0x00000008; b:int:0x00000004; a:int:0x00000000;
```

```
{[0]
    int a;
    int b;
    {[1]
        int c;
        a = 1;
        b = 2;
        int b;
        c = 3; <-
    }
    int c;
}

SymbolTable:
scope[0]: a:int:0x00000000; b:int:0x00000004;
scope[1]: c:int:0x00000008; b:int:0x0000000c; a:int:0x00000000;
```

```
{[0]
    int a;
    int b;
    {[1]
        int c;
        a = 1;
        b = 2;
        int b;
        c = 3;
    } <- // leave scope[1]
    int c;
}

SymbolTable:
scope[0]: a:int:0x00000000; b:int:0x00000004;
scope[1]: droped;
```

```
{[0]
    int a;
    int b;
    {[1]
        int c;
        a = 1;
        b = 2;
        int b;
        c = 3;
    }
    int c; <-
}

SymbolTable:
scope[0]: a:int:0x00000000; b:int:0x00000004; c:int:0x00000008; 
// addr can be reused or allocated to other new like 0x00000010
```

We can design a simple symbol table structure using a linked list to store the symbol table entries for each scope. Each scope has a unique ID and a pointer to its parent scope, enabling support for multi-level scopes. Each symbol table entry contains information such as the variable name, type, address, line number, and column number.

```go

type SymbolTableItem struct {
	Variable string
	Type     SymbolTableItemType
	Address  int

	UnderlyingType string

	VariableSize int
	ArraySize    int

	Line, Pos int64
}

type Scope struct {
	ID     int
	Level  int
	Items  map[string]*SymbolTableItem
	Parent *Scope
}

type SymbolTable struct {
	LegacyScopes  []*Scope // for debugging purposes
	CurrentScope  *Scope
	EnterFunction func(*Scope) error
	ExitFunction  func(*Scope) error

	addrCounter  int
	constantAddr int
}
```

Note: The symbol table needs to be used in conjunction with the syntax analyzer. The specific algorithm implementation will not be introduced here and will be detailed when the syntax analyzer or intermediate code generator is implemented.

#### 2.8.x Symbol Table Generation Example

Input:
```
{
    int a;
    a = 0;
    if (a == 0) {
        a = 1;
        if (a == 1) {
            float a;
            a = 2;
        } else {
            int b;
            int c;
        }
        int c;
    } else {
        int c;
        c = 0;
        if (c >= 0) {
            float c;
        } else {
            do {} while (true);
        }
    }
}
```

Output:
```
Scope[1]: 
  Level: 1
  Symbols:
    0x10000004 -> a:int[alloc=4] << at line 1, pos 9


Scope[2]: 
  Level: 2
  Symbols:
    0x10000014 -> c:int[alloc=4] << at line 12, pos 13


Scope[3]: 
  Level: 3
  Symbols:
    0x10000008 -> a:float[alloc=4] << at line 6, pos 19


Scope[4]: 
  Level: 3
  Symbols:
    0x10000010 -> c:int[alloc=4] << at line 10, pos 17
    0x1000000c -> b:int[alloc=4] << at line 9, pos 17


Scope[5]: 
  Level: 2
  Symbols:
    0x10000018 -> c:int[alloc=4] << at line 14, pos 13


Scope[6]: 
  Level: 3
  Symbols:
    0x1000001c -> c:float[alloc=4] << at line 17, pos 19


Scope[7]: 
  Level: 3
  Symbols:


Scope[8]: 
  Level: 4
  Symbols:
```

## Testing

The [test file](/lexer/lexer_test.go) can be found in the `app/lexer` directory. The test cases cover various types of tokens, including reserved keywords, identifiers, strings, characters, numbers, and operators.

Alternatively, you can run the following command in the terminal to execute the tests (ensure the executable file is compiled first):

```bash
./bin/main -t lexer
```

Add the `-f <files>` parameter to the command to specify the files to test. By default, all files in the `app/lexer` directory will be tested.

```bash
./bin/main -t lexer -f 1.in|2.in|3.in|4.in
```