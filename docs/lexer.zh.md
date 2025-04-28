# 词法分析器

## 实验目的
本实验要求学生编写词法分析函数，从输入的源程序中，识别出各个具有独立意义的单词，即基本保留字、标识符、常数、运算符、分隔符五大类。并依次输出各个单词的内部编码及单词符号自身值。词法分析的目的是将输入的源程序进行划分，给出基本符号（token）的序列，并过滤注解和空格等分隔符，词法分析为下一个阶段的语法分析提供单词序列。

## 实验任务
实现编译器的以下功能：
1. 分隔单词，二元组形式输出(单词名称，单词类别)；
2. 删除注释；
3. 删除空白符（空格、回车符、制表符）；
4. 发现并定位错误；
5. 建立符号表，并按照顺序打印输出；

## 具体实现

### 1. 文件读取

使用 `mmap` 和 `bufio.Reader` 读取文件内容。

- `mmap` 通过将文件映射到内存，避免了传统文件 I/O 的多次系统调用，从而显著提高读取效率，同时允许直接在内存中访问文件内容，适合处理大文件。
- `bufio.Reader` 提供了缓冲机制，通过减少对底层文件的直接读取次数，优化了小块数据的读取性能。

结合两者的优势，`mmap` 提供了类似内存的访问方式，简化了文件内容的处理逻辑，而 `bufio.Reader` 提供了灵活的读取接口（如 `ReadRune` 和 `UnreadRune`），便于读取`UTF-8`字符。此外，`mmap` 支持随机访问文件内容，适合需要频繁跳转读取的场景，而 `bufio.Reader` 的缓冲区大小可控，可以根据需求在内存使用和性能之间取得平衡。

> 但实际上，`mmap` 和 传统 I/O 方式在此次实验中差距不大
> 对于一个由 1k+ 代码行重复拷贝而来的 10w+ 的代码文件，`mmap` 和传统 I/O 方式不管是在使用 4KB 缓冲区，还是 16B 缓冲区，二者的耗时差距不大。
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
> - 实验中，文件的读取是顺序进行的。这种情况下，操作系统的文件缓存机制（如页缓存）会显著优化传统 I/O 的性能。`mmap` 的优势主要体现在随机访问场景中，因为它允许直接通过内存地址访问文件的任意部分，而无需频繁的系统调用。但在顺序读取中，这种优势并不明显。
> - `bufio.Reader` 的 `UnreadRune` 操作允许将读取的字符退**回到**缓冲区中，但这种操作通常只涉及缓冲区内的数据，不会触发底层文件的随机读取。因此， `UnreadRune` 的使用不会显著增加随机访问的频率，文件读取仍然是以顺序读取为主。
> - 在较大的缓冲区（如 4KB）下，传统 I/O 的性能已经非常接近 `mmap`，因为较少的系统调用和较大的数据块读取可以减少性能开销。即使在较小的缓冲区（如 16B）下，顺序读取的特性仍然使得传统 I/O 的性能不会显著下降。

### 2. 代码文件 Token 化

#### 2.1 参考规范
在本实验中，词法分析器的设计参考了 Go 语言的语法规则。具体而言，关键字、字符串等基本符号（token）的定义和处理方式均遵循 Go 语言规范。这种设计选择的目的是利用 Go 语言清晰的语法结构和广泛的应用场景，为实验提供一个可靠的参考框架（用 Go 写的当然用 Go 参考更直接）。


#### 2.2 实现方案
本程序并未直接实现有限状态机，而是通过简单的分支判断来模拟类似的功能。这种实现方式主要是因为高级语言的词法规则通常是静态的，在编译时已明确定义，分支判断可以更直观地实现这些规则，而无需引入有限状态机的复杂性。同时，有限状态机更适合处理动态、运行时的需求，尤其是在需要频繁扩展或修改状态和转换规则的场景中。然而，本实验的词法分析任务相对固定，分支判断的实现方式足以满足需求。此外，分支判断的逻辑清晰，易于理解和调试，能够快速实现词法分析功能，同时避免了有限状态机的状态设计和转换表维护的额外开销。

简单来说，程序逻辑如下：

1. 分支判断代替状态转移
有限状态机通常通过状态和状态转移表来处理输入字符，而本程序直接使用条件分支（如 if、switch）来判断当前字符的类型和处理逻辑。每种字符类型（如字母、数字、运算符等）对应一个分支，直接调用相应的处理函数（如 ReadWord、ReadNumber 等）。
2. 静态规则代替动态状态
高级语言的词法规则通常是静态的，已在编译时明确定义。本程序通过硬编码的规则（如 Operators、Delimiters 集合）直接匹配输入字符，而无需动态维护状态和转移表。例如，运算符的匹配通过 Operators 集合的前缀匹配实现，而不是通过状态转移。
3. 逐字符读取代替状态流转
有限状态机依赖状态流转，而本程序通过逐字符读取输入流，并在每次读取后立即判断字符的类型和处理逻辑。例如，遇到字母或下划线 _ 开头时，调用 ReadWord 读取完整单词；遇到数字时，调用 ReadNumber 解析完整数字。
4. 错误处理内嵌于分支逻辑
有限状态机通常通过进入特定的错误状态来处理非法输入，而本程序直接在分支逻辑中检测非法字符或格式，并返回错误。

![lexer logic](/docs/img/lexer/1.png)

#### 2.3 具体实现

`Lexer` 在进入主要判断逻辑之前，首先会跳过注释和空格，需要注意的是，某些逻辑序列的进入条件可能需要多读取几个字符，如果不满足条件，则需要将读取的字符退回到缓冲区中。

`Lexer` 可以维持 Line 和 Column 的状态，方便后续的错误定位。

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

##### 2.3.1 空白符

程序的空白符包括空格、制表符和换行符。我们需要跳过这些字符，直到遇到下一个有效字符为止。

```go
// 在进入判断分支之前，跳过空白符
// 这里的空白符包括空格、制表符和换行符
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

##### 2.3.2 注释

根据 Go 的语法规则，注释被分为两类：
1. 以 `\\` 开头的单行注释
2. 以 `\*` 开头，`*\` 结尾的跨行注释

对于情况`1. `只需要将整行字符读取完毕即可。
> 目前的主流操作系统，行尾序列大都是 CRLF(`\r\n`)或 LF (`\n`)，少数操作系统使用 CR (`\r`)。程序中使用 `\n` 作为行结束符，能够兼容大多数操作系统。

```go
// 程序判断分支已经读取到 `//`，接下来需要跳过注释
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

对于情况`2. `，我们需要判断是否存在跨行注释的结束符 `*\`。如果没有，则继续读取下一行，直到找到结束符为止。

```go
// 程序判断分支已经读取到 `/*`，接下来需要跳过注释
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

##### 2.3.3 字符串

程序支持在字符串中使用转义字符，包括 `\n`、`\t`、`\r`、`\b`、`\f`、`\a`、`\v` 和 `\"`；同时支持 Unicode 字符的转义（`\u` 和 `\U`），以及八进制字符的转义（`\0`），这些字符的转义会被转换为对应的字符，体现在 Token 的 Val 中。

Go 语言的字符串由双引号 `"` 或反引号 `` ` `` 包围。

双引号字符串中可以包含转义字符（如 `\n`、`\t` 等）。在读取双引号字符串时，我们需要处理转义字符，并将其转换为对应的字符。

由于字符串中存在转义字符，用有限状态机来处理字符串的读取会比较复杂，这时，将转义状态单独提取出来会更清晰。在读取字符串时，我们需要判断当前字符是否为转义字符 `\`，如果是，则进入转义状态。根据转义字符的类型，我们可以决定如何处理后续的字符。

程序支持的字符串格式包括：
- 普通字符串：`"hello"`、`"world"`、`"hello world"`
- 转义字符串：`"hello\nworld"`、`"hello\tworld"`、`"hello\rworld"`、`"hello\bworld"`、`"hello\fworld"`、`"hello\aworld"`、`"hello\vworld"`、`"hello\"world"`
- Unicode 字符串：`"好"`
- Unicode 转义字符串：`"\u4e2d"`、`"\U00004e2d"`
- 八进制转义字符串：`"\040"`
- 反引号字符串：`` `hello` ``、`` `hello\nworld` ``、`` `hello\tworld` ``、`` `hello\rworld` ``、`` `hello\bworld` ``、`` `hello\fworld` ``、`` `hello\aworld` ``、`` `hello\vworld` ``、`` `hello\"world` ``

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

对于反引号字符串，程序会直接读取到下一个反引号为止。

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

##### 2.3.4 字符

字符常量由单引号 `''` 包围，支持转义字符（如 `\'`、`\n` 等）。在读取字符时，我们需要判断当前字符是否为转义字符 `\`，如果是，则进入转义状态。根据转义字符的类型，我们可以决定如何处理后续的字符。

程序可以判断字符常量是否合法，即是否只包含一个字符，如果字符常量中包含多个字符，则返回错误。

需要注意的是，现代编程语言均提供对 Unicode 字符的支持，字符常量也可以包含 Unicode 字符，而 Unicode 字符的字长往往大于一个字节，因此在读取字符常量时，我们需要考虑Unicode字符的长度。

程序对 Unicode 字符的支持是通过转义字符 `\u` 和 `\U` 实现的。对于 `\u` 转义字符，程序会读取 4 个十六进制数字，并将其转换为对应的 Unicode 字符；对于 `\U` 转义字符，程序会读取 8 个十六进制数字，并将其转换为对应的 Unicode 字符（在 Token 输出的详情中，可以看到对应的 Unicode 字符）。

程序支持的字符格式包括：
- 单个字符：`'a'`、`'_'`、`'1'`
- 转义字符：`'\n'`、`'\t'`、`'\r'`、`'\b'`、`'\f'`、`'\a'`、`'\v'`、`'\\'`
- Unicode 字符：`'好'`
- Unicode 转义字符：`'\u4e2d'`、`'\U00004e2d'`

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

##### 2.3.5 单词

单词由字母、数字和下划线 `_` 组成，且必须以字母或下划线 `_` 开头。我们需要判断当前字符是否为字母或下划线 `_`，如果是，则进入单词读取状态。根据单词的类型，我们可以决定如何处理后续的字符。

单词可能是保留字。保留字是编程语言中预定义的关键字，具有特殊的语法意义。我们需要判断当前单词是否为保留字，如果是，则将其标记为保留字；如果不是，则将其标记为标识符。

由此可见，Identifier 和 Keyword 的判断在逻辑上相互独立，而在实现上却是相互依赖的。

遗憾的是，程序并不支持除了字母、数字和下划线 `_` 之外的其他字符作为单词的组成部分（如 Unicode 中的汉字）。读者可以根据需要自行扩展。

程序支持的单词格式包括：
- 标识符：`a`、`_a`、`a1`、`_a1`、`a_1`、`_a_1`
- 保留字：`int`、`float`、`string`、`if`、`else`、`for`、`while`、`return`、`break`、`continue`
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

##### 2.3.6 数字

数字常量由数字组成，支持整数和浮点数。我们需要判断当前字符是否为数字，如果是，则进入数字读取状态。根据数字的类型，我们可以决定如何处理后续的字符。

程序不支持科学计数法（如 `1.23e+4` 或 `1.23E+4`），也不支持负数（如 `-1` 或 `-1.23`），当然也不支持 `.`开头的数字（如 `.23`）。
> 对于负数，我们希望在语法分析阶段进行处理，而不是在词法分析阶段进行处理。
> 对于科学计数法，程序可以通过简单的字符串拼接来实现，但考虑到实验的复杂性和可读性，暂时不支持。

程序只支持整数和浮点数的表示，整数可以是十进制或十六进制，浮点数可以是以多个前导零开头的十进制数。
> 目前的主流编程语言均支持十六进制数的表示，十六进制数以 `0x` 或 `0X` 开头，后面跟着数字和字母 A-F（或 a-f）。
> 对于数字的支持已经足够涵盖大多数编程语言的需求，读者可以根据需要自行扩展。

程序支持的数字格式包括：
- 整数：`123`、`0x123`、`0X123`、`0xABCDEF`、`0XABCDEF`
- 浮点数：`123.456`、`0.123456`、`000.123456`

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

##### 2.3.7 操作符

操作符是编程语言中用于执行特定操作的符号或关键字。它们通常用于对变量、常量或表达式进行计算、比较或逻辑运算。操作符可以分为以下几类：
- 算术操作符：用于执行数学运算，如`+`、`-`、`*`、`/`、`%` 等。
- 关系操作符：用于比较两个值的大小关系，如`<`、`>`、`<=`、`>=`、`==`、`!=` 等。
- 逻辑操作符：用于执行逻辑运算，如`&&`、`||`、`!` 等。
- 位操作符：用于对二进制位进行操作，如`&`、`|`、`^`、`<<`、`>>` 等。
- 赋值操作符：用于将值赋给变量，如 `=`、`+=`、`-=` 等。

操作符的长度不固定，可能是一个字符（如 `+`、`-`、`*`、`/`、`=` 等），也可能是两个字符（如 `==`、`!=`、`<=`、`>=`、`&&`、`||` 等）；同时可以发现，一些操作符是由其他操作符组合而成的（如 `+=`、`-=`、`&=`、`|=` 等）。

由此可见，操作符可能是某个操作符的前缀（如 `+`、`-`、`*`、`/`、`=` 等），也可能是某个操作符的后缀（如 `+=`、`-=`、`&=`、`|=` 等），我们不能简单地用单个字符来判断操作符的类型，需要向前搜索操作符，找到的`最长匹配`的操作符即为当前操作符。

简单来说，对于操作符 `<` 和 `<<`，当我们读取到 `<` 时，我们需要判断下一个字符是否为 `<`，如果是，则将其标记为 `<<`；如果不是，则将其标记为 `<`。对于 `<<<`，可以将其分割为 `<<` 和 `<`，即 `<<` 是第一个操作符，`<` 是下一个操作符。

一种高效的实现是是用前缀搜索树（Trie）来存储操作符。前缀搜索树是一种树形数据结构，其中每个节点表示一个字符，路径上的字符组合表示一个操作符。通过前缀搜索树，我们可以快速地找到最长匹配的操作符。
- Trie 的每个节点表示一个字符，路径上的字符组合表示一个操作符。
- Trie 的根节点表示空字符串，子节点表示操作符的前缀。
- Trie 的叶子节点表示完整的操作符。
- Trie 的非叶子节点表示操作符的前缀，可以继续向下搜索。
- Trie 的每个节点可以通过一个标志位表示是否存在这样的操作符。

如下图：
1. 定义了操作符集合 `{+, +=, ++, -, -=, (= }`，注意集合中不存在 `(`。
2. `+` 和 `+=` 是前缀关系，`-` 和 `-=` 是前缀关系，`(` 和 `=` 是前缀关系。
3. 程序会从根节点开始，逐层向下搜索操作符，直到找到最长匹配的操作符为止。
   1. 对于输入串 `+`，程序会从根节点开始，找到 `+` 的子节点，返回 `+`。
   2. 对于输入串 `+=`，程序会从根节点开始，找到 `+` 的子节点，然后继续向下搜索 `=` 的子节点，返回 `+=`，因为 `+=` 是最长匹配的操作符。
   3. 对于输入串 `+++`，程序会从根节点开始，找到 `+` 的子节点，然后继续向下搜索 `+` 的子节点，返回 `++`，因为 `++` 是最长匹配的操作符；但注意，`+` 也是一个操作符，所以 `+++` 可以分割为 `++` 和 `+`，即 `++` 是第一个操作符，`+` 是下一个操作符，而 `+` 会在下一次从根节点开始的搜索得出。
   4. 对于输入串 `(`，程序会从根节点开始，找到 `(` 的子节点，由于 `(` 的标志位表示其并不是一个操作符，所以 `(` 并不会作为操作符返回。
   5. 对于输入串 `(=`，程序会从根节点开始，找到 `(` 的子节点，然后继续向下搜索 `=` 的子节点，返回 `(=` ，即使 `(` 的标志位表示其并不是一个操作符，但 `(` 和 `=` 的组合是一个操作符（搜索到的子节点的标志位表示其是一个操作符），所以返回 `(=`。
   
![Trie Tree](/docs/img/lexer/2.png)

由于本程序并不想新增复杂的 Trie 数据结构，而是直接使用简单的字符串集合过滤器来实现操作符的匹配。
1. 初始化前缀和候选集合
   1. 将当前字符作为初始前缀 `prefix`
   2. 使用 `Operators` 集合（包含所有合法操作符）作为初始候选集合 `currentSet`
   3. 通过过滤器筛选出以 `prefix` 为前缀的操作符，缩小候选范围
2. 逐字符扩展前缀
   1. 读取下一个字符，将其追加到 `prefix` 中
   2. 更新 `currentSet`，保留所有以新 `prefix` 为前缀的操作符
   3. 如果 `currentSet` 为空，说明没有操作符以当前 `prefix` 为前缀，停止扩展。
3. 记录最长匹配
   1. 在每次扩展前缀时，检查 `currentSet` 是否包含当前 `prefix`
   2. 如果包含，将其记录为当前的最佳匹配 `bestMatch`
4. 处理结束条件
   1. 如果读取到文件末尾（`EOF`），停止扩展并回退到最佳匹配位置
   2. 如果 `currentSet` 为空，回退到上一个字符位置
5. 回退到最佳匹配
   1. 计算需要回退的字符数`（len(prefix) - len(bestMatch)）`
   2. 调用回退函数，将多余的字符退回缓冲区

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

#### 2.4 Token 结构体
`Token` 结构体用于表示词法分析器生成的 Token。它包含了 Token 的类型、值、行号、列号等信息。

`Token` 结构体的定义如下：
- `Type`：Token 的类型，使用 `ItemType` 枚举类型表示。
- `Val`：Token 的值，使用字符串表示。
- `Line`：Token 所在的行号，使用整数表示。
- `Pos`：Token 在行中的位置，使用整数表示。
- `_type`：Token 的具体类型，使用 `TokenSpecificType` 枚举类型表示。

```go
type Token struct {
	Type                ItemType
	Val                 string
	Line, Pos           int64

	_type TokenSpecificType
}
```

### 2.5 二元组输出格式
`Token` 结构体的输出格式为二元组，包含 Token 的类型和值。

```go
func (t Token) String() string {
    return fmt.Sprintf("(%s, %s)", t.Type, t.Val)
}
```

具体体现为
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

### 2.6 词法分析器的使用
`Lexer` 需要一个 `io.Reader` 作为输入源，通常可以通过 `os.Stdin` 或文件读取器来实现，文件读取的具体实现可以是`os.File` 或 `mmap`。

当 `Lexer` 初始化完成后，可以通过 `NextToken()` 方法来获取下一个 Token。`NextToken()` 方法会返回一个 `Token` 结构体和一个错误信息，如果没有错误，则错误信息为 `nil`。

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

### 2.7 词法分析器的测试

`Lexer` 的测试使用了 `testing` 包，测试用例包括了对不同类型 Token 的测试，包括保留字、标识符、字符串、字符、数字、操作符等。

可以使用`go test app/lexer` 命令来运行测试用例，添加`-v` 参数可以查看详细的测试输出。

```bash
go test app/lexer -v
```

### 2.8 符号表的设计与实现

由于词法分析器无法判断 Token 所在的上下文环境，因此无法判断语句的类型（声明、赋值、调用等）。为了实现这一点，需要等到语法分析器阶段进行处理。

不过这里可以简单介绍一下符号表的设计与实现。

符号表是编程语言中的一个重要概念，用于存储变量、函数、类等标识符的信息。它通常用于编译器的语法分析和语义分析阶段。符号表的设计与实现可以分为以下几个部分：
- 符号表的结构：符号表通常使用哈希表或树形结构来存储标识符的信息。哈希表的查找速度较快，但插入和删除操作较慢；树形结构的查找速度较慢，但插入和删除操作较快。
- 符号表的作用域：符号表通常支持多级作用域，即在不同的作用域中可以定义相同的标识符。符号表需要支持作用域的嵌套和退出操作，以便在进入新的作用域时创建新的符号表，在退出作用域时销毁当前的符号表。
- 符号表的查找：符号表需要支持根据标识符的名称查找对应的信息。查找操作通常是从当前作用域开始，逐级向上查找，直到找到对应的标识符为止。
- 符号表的插入：符号表需要支持插入新的标识符信息。插入操作通常是在当前作用域中进行，如果当前作用域中已经存在相同名称的标识符，则需要报错。

由此可见，符号表是类似于上下文环境的概念，它可以存储变量、函数、类等标识符的信息，并支持多级作用域和查找操作。

这里比较容易联想到的数据结构就是单向链表，单向链表的每个节点可以存储一个标识符的信息，并且可以通过指针连接到下一个节点。只不过，对于符号表的连接方向与一般的单向链表相反，符号表的连接方向是从下往上，而单向链表的连接方向是从上往下。

这里可以给出一个简单的举例（以 `<-` 表示当前读取到的语句块）：

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

我们可以设计一个简单的符号表结构，使用链表来存储每个作用域的符号表项。每个作用域都有一个唯一的 ID 和一个父作用域指针，用于支持多级作用域。每个符号表项包含变量名、类型、地址、行号和列号等信息。

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

注意，符号表需要配合语法分析器使用，这里不介绍具体算法的实现，等到语法分析器或中间代码生成器实现时再进行详细介绍。

#### 2.8.x 符号表生成样例

输入
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

输出
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

## 测试

[测试文件](/lexer/lexer_test.go)可以在 `app/lexer` 目录下找到，测试用例包括了对不同类型 Token 的测试，包括保留字、标识符、字符串、字符、数字、操作符等。

或者可以在命令行中运行以下命令来执行测试（首先要编译成可执行文件）：

```bash
./bin/main -t lexer
```

在命令中添加 `-f <files>` 参数可以指定要测试的文件，默认会测试 `app/lexer` 目录下的所有文件。

```bash
./bin/main -t lexer -f 1.in|2.in|3.in|4.in
```