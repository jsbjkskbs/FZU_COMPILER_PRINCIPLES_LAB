# 中间代码生成

## 实验目的
通过语法制导技术生成中间代码，设计中间代码语义规则并嵌入语法规则中。在语法分析过程中，同时输出等价的中间代码。

## 试验任务
实现编译器的以下功能： 
1. 输出语法分析过程； 
2. 输出中间代码，要求三地址代码；

## 文法
```plaintext
program  → block 
block    → { decls stmts} 
decls    → decls decl | ε 
decl     → type id; 
type     → type[num] | basic 
stmts    → stmts stmt | ε 
stmt     → loc = bool; 
          | if ( bool ) stmt 
          | if ( bool ) stmt else stmt 
          | while ( bool ) stmt 
          | do stmt while ( bool ); 
          | break; 
          | block 
loc      → loc[num] | id 
bool     → bool || join | join 
join     → join && equality | equality 
equality → equality == rel | equality != rel | rel 
rel      → expr < expr | expr <= expr | expr >= expr | expr > expr | expr 
expr     → expr + term | expr - term | term 
term     → term * unary | term / unary | unary 
unary    → ! unary | - unary | factor 
factor   → ( bool ) | loc | num | real | true | false
```

在语法分析器中，提到文法存在问题，接下来我们将对问题进行描述。

### 文法问题

#### 左递归
在编译原理的语法制导翻译过程中，左递归文法的存在会对代码生成阶段产生特定的复杂性。以如下左递归文法片段为例：

```plaintext
bool → bool || join
join → join && equality
equality → equality == rel
equality → equality != rel
```

这类直接左递归产生式的结构特征是：

非终结符 `bool`、 `join`、 `equality`在产生式右部的最左端直接递归引用自身，形成形如`A → A α`的左递归模式。在基于 LR (1) 分析表的自底向上语法分析中，尽管 LR (1) 分析法能够处理包含直接左递归的文法（通过构造正确的分析表实现 **移进 - 规约** 决策），但在语法制导翻译时，这种递归结构会导致代码生成时机的不确定性问题。

以条件语句产生式`stmt → if ( bool ) stmt else stmt`为例：

当分析器处理 `bool` 表达式时，需要生成对应的条件判断代码（如跳转指令）。然而左递归结构使得 `bool` 的规约过程呈现链式特征：对于产生式`bool → bool || join`，每次对 `bool` 的规约实际上是对递归链条中一个中间节点的处理，而非最终完整表达式的终结。在自底向上的分析流程中，当遇到 `bool` 非终结符的规约时机时，分析器仅能确定当前处理的是递归结构中的某个子表达式（如左部 `bool` 对应已处理的子表达式，右部 `|| join` 对应待处理的后续操作），但无法判断该 `bool` 是否为最终需要生成完整条件结果的顶层表达式。

具体而言，左递归存在以下问题：

1. 按照语法制导翻译的常规逻辑，表达式的代码生成应遵循 **"先子表达式、后操作符"** 的顺序。但左递归结构导致非终结符 `bool` 在语法分析中被多次规约（每次规约对应递归链条的一个环节），而每次规约时均需决定是否生成完整表达式的结果代码。
    > 例如，处理 `bool → bool || join` 时，左部 `bool` 对应的子表达式代码已生成，但若该 `bool` 并非最终顶层表达式（而是作为另一个 `bool || join` 结构的左部），此时生成完整结果代码将导致后续 `||` 操作的语义无法正确关联。
2. 在 LR (1) 分析表的构造中，同心项（具有相同核心但不同展望符的项目）可能对应不同的语义处理逻辑。当左递归产生式的非终结符 `bool` 作为同心项的核心时，分析器在规约时仅能根据当前展望符（即后继输入符号）决定移进或规约，但无法预判该 `bool` 是否会在后续分析中作为更大左递归结构的左部继续参与规约。这种不确定性使得代码生成器无法在单次规约时确定是否需要生成最终结果的目标代码（如条件跳转的终节点），因为当前规约可能仅是递归链条中的中间步骤，而非表达式的终结点。
3. 对于逻辑表达式 `bool || join` 的短路求值语义（左操作数为真时无需计算右操作数），代码生成需要在左部 `bool` 求值后插入条件跳转以跳过右部 `join` 的计算。然而左递归结构导致 `bool` 可能作为右部 `join` 的祖先节点多次出现，每次规约 `bool` 时均需判断是否处于顶层表达式以决定跳转目标的绑定。若误将中间递归层的 `bool` 作为顶层表达式生成跳转，将导致后续操作数的代码被错误跳过，破坏语义正确性。

自非终结符的多次规约（作为递归链条的中间节点）与代码生成所需的 "终结点判断" 之间存在语义偏差。分析器仅能基于当前符号和展望符进行语法层面的规约决策，却无法感知该非终结符在语义层面是否已构成完整表达式的最终结果。这种语法结构与语义处理的差别，使得代码生成器难以在恰当的时机生成完整的结果代码。

#### else 悬挂问题
上述文法存在**else 悬挂问题**在编译原理的语法分析过程中，**else 悬挂问题**会对 `if - else` 语句的解析增加复杂度。以如下包含 `if - else` 语句的文法片段为例：

```plaintext
stmt → if ( bool ) stmt 
     | if ( bool ) stmt else stmt 
```

这类产生式的结构特征是：

- `if` 语句存在两种形式，一种是不含 `else` 子句的，另一种是包含 `else` 子句的。这种结构使得在嵌套的 `if - else` 语句中，`else` 子句的匹配对象不明确，从而引发了 **else 悬挂问题**。在基于自底向上的语法分析中，尽管语法分析器能够根据产生式进行**移进 - 规约**操作，但对于 `else` 子句与哪个 `if` 配对，该文法并未给出明确规则。

以嵌套的 `if - else` 语句 `if (condition1) if (condition2) statement1; else statement2;` 为例，当分析器处理该语句时，需要确定 `else` 子句与哪个 `if` 匹配。然而，参考文法结构特点， `else` 的匹配过程十分模糊：

- 对于这条语句，`else` 既可以与内层的 `if (condition2)` 匹配，也可以与外层的 `if (condition1)` 匹配。在自底向上的分析流程中，当遇到 `else` 符号时，分析器仅能确定当前处理的是一个 `if - else` 结构的一部分，但无法判断该 `else` 应与前面哪个未匹配的 `if` 结合。

具体而言，**else 悬挂问题**在语法分析过程中存在以下问题：
1. 按照语法分析的常规逻辑，每个 `else` 子句都应该明确地与一个 `if` 配对。但该文法缺乏明确的 `else` 匹配规则，导致在分析嵌套的 `if - else` 语句时，可能会出现多种不同的匹配方式。例如，对于上述嵌套语句，若没有明确规则，分析器无法确定 `else` 是与内层 `if` 还是外层 `if` 关联，这就使得代码的语义解释不唯一。
2. 在语法分析过程中，不同的 `else` 匹配方式会产生不同的语法树。这些不同的语法树代表了不同的代码语义，这就使得分析结果具有不确定性。例如，`else` 与内层 `if` 匹配和与外层 `if` 匹配所表达的条件判断逻辑和执行流程是不同的，这会影响后续的语义分析和代码生成等编译步骤。
3. 为了解决 **else 悬挂问题**，语法分析器需要额外的规则或策略来确定 `else` 的配对。这可能涉及到回溯、状态管理等复杂操作，从而增加了语法分析的复杂度和实现难度。例如，分析器可能需要记录已经遇到的 `if` 语句，以便在遇到 `else` 时能够准确地找到合适的匹配对象。

如果缺乏明确的 `else` 匹配规则，那么 `else` 与 `if` 的配对就会存在不确定性，偏离语法分析所需的明确性——分析器仅能基于当前符号进行语法层面的**移进 - 规约**决策，却无法确定 `else` 在语义层面应与哪个 `if` 正确匹配。

### 文法改进
#### 左递归
其实并不是所有左递归文法都需要改进，最主要的判断依据是该左递归文法是否会严重影响到语法分析器的实现。对于本实验中的文法，左递归的存在会导致语法分析器在处理 `bool` 表达式时，无法确定何时生成完整表达式的代码。

我们可以通过将左递归文法改写为非左递归的形式来解决这个问题。以下是对原文法的改进：

```plaintext
bool → bool'
bool' → bool' || join | join
...
```

在这个改进后的文法中，我们引入了一个新的非终结符 `bool'`，它表示 `bool` 的右部。这样，我们就可以保证 `bool` 被规约为一个完整的表达式时，`bool` 不属于任何递归链条的中间节点。

#### else 悬挂问题

在本实验中，我们可以通过引入一个新的非终结符 `matched_stmt` 和 `unmatched_stmt` 来解决 **else 悬挂问题**。以下是对原文法的改进：

```plaintext
stmt           → matched_stmt 
                 | unmatched_stmt
unmatched_stmt → if ( bool ) unmatched_stmt
                 | if ( bool ) matched_stmt else unmatched_stmt
matched_stmt   → loc = bool ;
                 | if ( bool ) matched_stmt else matched_stmt
                 | if ( bool ) matched_stmt
                 | while ( bool ) stmt
                 | do stmt while ( bool ) ;
                 | break ;
                 | block
```

在这个改进后的文法中，我们将 `stmt` 分为两类：`matched_stmt` 和 `unmatched_stmt`。`matched_stmt` 表示已经匹配了 `else` 的语句，而 `unmatched_stmt` 表示尚未匹配的语句。这样，我们就可以确保每个 `else` 子句都能正确地与一个 `if` 语句匹配，从而避免了 **else 悬挂问题**。

## 实现
### 分析过程中生成中间代码
我们可以给每个产生式添加一个语义动作，在语法分析的过程中生成中间代码，当分析器规约到某个产生式时，就执行对应的语义动作。

例如，对于产生式 `unary → - unary`：
```plaintext
-a

// 规约文法链
id    → loc      // 这时候要查询符号表，获取 id 的值
loc   → unary    // loc 作为 unary 的子节点
unary → - unary  // - unary 作为 unary 的子节点
```

在规约到 `unary → - unary` 时，我们可以生成中间代码：
```plaintext
mov $(t1), -$(t2) // 将 unary 的值取反，存入 t1
```

> 注：此后的三地址码（四元组）均使用 `op dist arg1 arg2` 的形式表示。
> 这个实验依赖于语法分析器的实现，首先要保证文法和语法分析器逻辑的正确，否则无法生成正确的中间代码。

### 符号表（符号作用域）
在中间代码生成过程中，我们需要维护一个符号表，用于存储变量的类型、作用域等信息。

符号表的设计可以使用一个栈来实现，每当进入一个新的作用域时，就将当前作用域的符号表压入栈中；当离开作用域时，就将栈顶的符号表弹出。这样可以实现嵌套作用域的管理。
> 实际上还可以用单向链表来实现符号表，链表的每个节点表示一个作用域，每个节点中存储该作用域的符号表。
> 当进入一个新的作用域时，就在链表头插入一个新的节点；当离开作用域时，就删除链表头的节点。
> 这样可以实现嵌套作用域的管理。
> 其实这样的单向链表和栈是等价的。

如图所示：
![符号表示意](/docs/img/intermediate-code-generation/1.png)

流程如下：
1. 外层声明 `int a`，符号表 **S1** 记录 `a = $(0)`；执行 `a = 1`，更新 `$(0) = 1`。
进入内部代码块：
2. 先执行 `a = 1`，仍关联 `$(0) = 1`（此时未重新声明，暂用外层关联）。
声明内部 `int a`，符号表切换到 **S2**，记录 `a = $(1)`；再执行 `a = 1`，更新 `$(1) = 1`。
3. 声明 `int c`，符号表 **S2** 记录 `c = $(2)`。
4. 内部块结束 `drop S2`，回到外层作用域。
5. 外层声明 `int c`，符号表 **S1** 记录 `c = $(3)`；执行 `c = 1`，更新 `$(3) = 1`，最终 **S1** 中保留 `a = $(0)` 和 `c = $(3)`。

> 提示：后续还需要维护一个循环体的栈，用于处理 `break` 语句。
> 此时单向链表没有太多的必要，因为我们只需要维护循环体代码结束位置的栈即可（不需要存储太复杂的结构）。

### 抽象语法树

抽象语法树（Abstract Syntax Tree，AST）是一种对源代码进行抽象表示的数据结构，它以树状形式精准地描述了源代码的语法结构。
- 抽象语法树由各种节点组成，每个节点代表了源代码中的一个特定语法元素。比如，变量声明、函数调用、表达式、语句块等都可以是一个节点。这些节点具有不同的类型，以区分它们所代表的语法元素的种类。例如，一个表示变量声明的节点具有 “Decl” 类型。
- 节点除了有类型信息外，还可能包含其他属性，用于存储与该语法元素相关的具体信息。例如，变量声明节点可能包含变量的名称、数据类型等属性；函数调用节点可能包含函数的名称以及传递给函数的参数列表等属性。这些属性能够更详细地描述语法元素的特征，为后续的分析和处理提供必要的数据。
- 节点之间通过父子关系和兄弟关系形成层次分明的树状结构。这种结构反映了源代码中语法元素的嵌套关系和逻辑顺序。例如，在一个函数定义中，函数体内部的语句块可以看作是函数定义节点的子节点，而语句块中的各个语句则是语句块节点的子节点。通过这种层次结构，可以清晰地展现出源代码的整体语法架构，使得对代码的分析和理解更加直观。

LR 分析器是一个自底向上的语法分析器，从抽象语法树的根节点开始，逐步向下分析，直到到达叶子节点——也就是说，LR 分析的本质是一次深度优先遍历（DFS）。
> 我们可以维护一个抽象语法树节点栈，在语法分析的过程中，将每个产生式的右部节点压入栈中，当规约到某个产生式时，就将对应的节点弹出栈，弹出节点合成一个新的节点，并将新节点压入栈中。
> 这样，我们就可以在语法分析的过程中构建出一棵完整的抽象语法树。

在 LR 的分析过程中：
- **SHIFT**：当执行移进操作时，把当前输入的词法单元转化为对应的 `AST` 叶子节点，如字面量、标识符等节点，然后将其压入节点栈。例如，输入词法单元为数字 `5`，就创建一个 `Num` 类型的 `AST` 节点并压栈。
- **REDUCE**：当进行归约时，按当前使用的产生式，从节点栈中弹出和产生式右部符号数量相同的节点。接着，创建一个新的 `AST` 节点代表产生式左部的非终结符，把弹出的节点按产生式右部的顺序作为子节点连接到新节点上，最后将新节点压入栈。例如，对于产生式 `E → E + T`，栈中弹出表示 `E`、`+`、`T` 的节点，创建一个新的 `E` 节点，将这三个节点作为其子节点，再把新的 `E` 节点压栈。

假设有简单算术表达式文法：
```plaintext
E → E + T
E → T
T → T * F
T → F
F → ( E )
F → num
```
在一次 LR 分析过程中，假设输入的词法单元序列为 `3 + 4 * 5`，分析器会按以下步骤构建抽象语法树：
1. 输入表达式 `3 + 4 * 5`，词法单元序列为 `[num(3), +, num(4), *, num(5)]`。
2. 遇到 `num(3)`，创建 `Num` 节点压入栈。
3. 依据产生式 `F → num` 归约，弹出 `Num` 节点，创建 `F` 节点并连接，再将 `F` 节点压栈。
4. 按 `T → F` 归约，弹出 `F` 节点，创建 `T` 节点连接，`T` 节点压栈。
5. 后续操作类似，直至完成整个表达式分析，栈顶节点就是表示整个表达式的 `AST` 根节点。

![3 + 4 * 5 的分析示例](/docs/img/intermediate-code-generation/2.png)

完成分析后，节点栈栈顶的节点就是构建好的抽象语法树的根节点，它能够完整地反映输入源代码的语法结构。

#### 抽象语法树与代码优化
抽象语法树（AST）在编译器的中间代码生成和优化阶段起着重要作用。它不仅提供了源代码的结构化表示，还为后续的代码优化和生成提供了基础。以下是 AST 在中间代码生成和优化中的一些关键作用：
- AST 的每个节点可以映射到一个或多个中间代码指令，从而实现源代码到中间代码的转换。通过遍历 AST，编译器可以生成对应的三地址码或其他形式的中间表示。
- 编译器可以分析 AST 的结构，识别出冗余的计算、常量表达式、死代码等，并进行相应的优化。例如，可以通过常量折叠（constant folding）将常量表达式直接计算为一个值，从而减少运行时的计算开销；可以通过死代码消除（dead code elimination）来删除不必要的变量定义和使用，从而减少内存开销。
- 编译器可以分析 AST 中的条件语句、循环等结构，识别出不必要的分支和循环，从而进行优化。例如，可以通过循环展开（loop unrolling）来减少循环的迭代次数，提高执行效率。

> 注：在本实验中，我打算直接在规约时生成中间代码，构建抽象语法树只是为了方便后续的代码生成。
> 
> 但实际编译过程中，抽象语法树的构建和中间代码生成是两个独立的阶段，通常会先构建抽象语法树，然后再进行中间代码生成。
> 
> 在中间代码生成阶段，编译器会遍历抽象语法树，根据节点的类型和属性生成对应的中间代码指令。
>
> 感兴趣的读者可以试着完善一下本仓库的`parser/ast.go`和`parser/gen_rules.go`，实现先构建抽象语法树，再通过语法树优化并生成中间代码。

### 关键规则的中间代码生成

我们不妨从简单的、抽象语法树叶子节点开始，逐步向上分析，直到根节点。

#### Token 抽象语法树叶子节点
从词法分析器中获取的 Token 直接作为抽象语法树的叶子节点，只需要进行简单的转换即可。

```go
func (p *Parser) Token2ASTNode(token *lexer.Token) *ASTNode {
	return &ASTNode{
		raw:      token.Val,
		Token:    token,
		Children: []*ASTNode{},
		Type:     p.Reflect(token),
		Payload:  nil,
	}
}
```
> - `p.Reflect(token)` 是一个反射函数，用于获取 Token 的类型。
> - `ASTNode` 是一个结构体，包含了 Token 的原始值、类型、子节点等信息。
>   - `raw` 用于存储节点关联表达式。
>   - `Token` 用于存储 Token 的信息（为了拓展，Token 还可以存储分配后的地址等信息）。
>   - `Children` 用于存储子节点（在语法分析过程中，子节点会被压入栈中）。
>   - `Type` 用于存储节点的类型（如 `Num`、`Id`、`Bool` 等）。
>   - `Payload` 用于存储额外的信息（虽然我没怎么用到）。

#### Factor 抽象语法树节点
`factor` 语法规则如下：
```plaintext
factor → ( bool ) | loc | num | real | true | false
```

在语法分析过程中，我们可以根据 `factor` 的类型来生成对应的中间代码。

以下是其中一个产生式的语义动作：
```go
// factor → false
func FactorFalse(w *Walker) error {
	children := w.Tokens.PopTopN(1)
	w.Tokens.Push(&ASTNode{
		raw:      "false",
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: "0"},
		Children: children,
		Type:     "factor-false",
		Payload:  "!const(size=1)",
	})
	return nil
}
```

> - 当语法分析器规约到 `factor → false` 时，执行 `FactorFalse` 函数。
> - 此时的栈顶必然是 `false` 的 抽象语法树节点，`PopTopN(1)` 将其弹出。
> - 然后创建一个新的 `ASTNode` 节点，类型为 `factor-false`，并将弹出的节点作为子节点。
> - 最后，将新的节点压入栈中（表示规约完成）。

#### Unary 抽象语法树节点
`unary` 语法规则如下：
```plaintext
unary → ! unary | - unary | factor
```

在语法分析过程中，我们可以根据 `unary` 的类型来生成对应的中间代码。

对于 `-1` 的中间代码生成：
```go
// unary -> !unary
func UnaryNot(w *Walker) error {
	result := w.SymbolTable.TempAddr(4)
	addr := fmt.Sprintf("$(%#x)", result)
	children := w.Tokens.PopTopN(2)
	w.Tokens.Push(&ASTNode{
		raw: joinChildren(children),
		Token: &lexer.Token{
			Type: lexer.EXTRA,
			Val:  fmt.Sprintf("$(%#x)", result),
		},
		Children: children,
		Type:     "not",
		Payload:  "!dist:!ptr(size=4)",
	})
	w.Emit(addr, "not", children[1].Token.Val)
	return nil
}
```
> - 当语法分析器规约到 `unary → ! unary` 时，执行 `UnaryNot` 函数。
> - `w.SymbolTable.TempAddr(4)` 分配一个临时变量，大小为 4 字节。
> - `addr` 是临时变量的地址，格式为 `$(0x00000000)`。
> - `children` 是弹出的节点，包含 `!` 和 `unary` 的 Token 节点。
> - `joinChildren(children)` 是一个函数，用于将子节点的原始值连接起来。
> - `w.Tokens.Push` 将新的节点压入栈中，表示规约完成。
> - 此时程序会生成一条中间代码：
> ```plaintext
> not $(dist) $(src)
> ```
> - `dist` 是临时变量的地址，`src` 是 `unary` 的值。

#### Term 抽象语法树节点
```plaintext
term → term * unary | term / unary | unary
```

在语法分析过程中，我们可以根据 `term` 的类型来生成对应的中间代码。
```go
// term -> term * unary
func TermMult(w *Walker) error {
	result := w.SymbolTable.TempAddr(4)
	children := w.Tokens.PopTopN(3)
	resultStr := fmt.Sprintf("$(%#x)", result)
	w.Tokens.Push(&ASTNode{
		raw: joinChildren(children),
		Token: &lexer.Token{
			Type: lexer.EXTRA,
			Val:  resultStr,
		},
		Children: children,
		Type:     "mult",
		Payload:  "!dist:!ptr(size=4)",
	})
	w.Emit("mul", resultStr, children[0].Token.Val, children[2].Token.Val)
	return nil
}
```

> - 当语法分析器规约到 `term → term * unary` 时，执行 `TermMult` 函数。
> - `w.SymbolTable.TempAddr(4)` 分配一个临时变量，大小为 4 字节。
> - `children` 是弹出的节点，包含 `term`、`*` 和 `unary` 的 Token 节点。
> - `w.Tokens.Push` 将新的节点压入栈中，表示规约完成。
> - 此时程序会生成一条中间代码：
> ```plaintext
> mul $(dist) $(src1) $(src2)
> ```
> - `dist` 是临时变量的地址，`src1` 是 `term` 的值，`src2` 是 `unary` 的值。
> - 这里的 `mul` 是乘法操作，`dist` 是结果存储地址，`src1` 和 `src2` 是操作数。

此后的 `expr`、`rel`、`equality`、`join` 和 `bool` 的中间代码生成与 `term` 类似，都是通过规约产生式来生成对应的中间代码。

#### 中间代码回填
在编译过程中，有些信息在中间代码生成的早期阶段是未知的，例如跳转指令的目标地址。回填的目的就是在这些信息变得可知时，将其补充到之前生成的中间代码中，使中间代码完整且能够正确地被目标机器理解和执行。

通常在完成对源程序的语法和语义分析之后，在生成目标代码之前进行。当编译器收集到足够的信息来确定那些先前未确定的部分时，就会执行回填操作。例如，在处理条件语句或循环语句时，一开始无法确定跳转指令的目标地址，只有在分析完整个语句块，确定了跳转的目标位置后，才能进行回填。

编译器通常会使用一些数据结构来记录需要回填的信息，比如使用符号表来记录变量的相关信息，包括其存储地址等。对于跳转指令的目标地址，可能会使用专门的列表或栈来记录那些尚未确定目标地址的跳转指令，当目标地址确定后，再从列表或栈中取出相应的指令进行回填。
> 以条件语句 `if (a > 10) goto L1;` 为例，在生成中间代码时，`goto L1` 中的 **L1** 地址是未知的，会将该指令以及相关信息记录下来。当后续分析到条件语句的结束位置，确定了 **L1** 对应的实际地址后，就可以将该地址回填到 `goto L1` 指令中。

##### 留空
一般的，我们需要确认产生留空跳转指令的时机——这可不是整个语句块被分析完后才进行的，而是当编译器分析到某个语句块的开始位置时，就可以预留跳转指令的空位。

我们需要通过文法来分析出一个语句块的开始位置，同时，这个位置必须能通过规约操作触发。

可以通过一个例子来简单说明：
```plaintext
matched_stmt → do stmt while ( bool ) ;
stmt → matched_stmt
matched_stmt → block
block → { decls stmts }
decls → decls decl | ε
```

在规约过程中，假设语法正确，存在尾部为栈顶的符号子序列：
```plaintext
do {
```

符号栈的下一个状态必然为：
```plaintext
do { decls
```
> 1. 为了规约出 `do stmt`，需要先规约出 `matched_stmt`。
> 2. 为了规约出 `matched_stmt`，需要先规约出 `block`。
> 3. 为了规约出 `block`，需要先规约出 `{ decls stmts }`。
> 4. 由此可知，`decls` 为规约的必需项。
> 5. 由 `decls` 产生式可知，`decls` 可能为空，则对于任意的内部表达式，必然存在尾部为栈顶的符号子序列：
>    ```plaintext
>    do { decls
>    ```

因此，`decls` 产生式的 `ε` 产生式可以作为留空跳转指令的时机。

##### 简单 If 语句的回填

在处理 if 语句的条件表达式并生成跳转相关中间代码时，编译器还未完全分析完整个 if 语句块及后续代码。例如图中第 2 行 `jmp`，在生成时，它跳转目标地址（对应 if 语句块结束位置等）还无法确定。

我们可以采用回填机制来处理这类情况。先预留跳转指令目标地址的空位，后续当分析完 if 语句块，确定了跳转目标的实际位置（如 if 语句块结束位置对应的代码行号），再将正确地址填充进去。像图中第 2 行 `jmp` 目标地址（虚线处） 以及第 6 行 `jmp` 目标地址（虚线处），都是先留空，待后续信息完备后通过回填来完善指令。

![简单 If 语句生成示例-第一部分](/docs/img/intermediate-code-generation/3.png)


随着编译进行，当编译器分析完 if 语句块及后续代码，确定了正确的跳转目标位置后，就进行回填。比如图中第 2 行 `jmp` 指令原本留空的目标地址被回填为 7，即直接跳转到程序结束的 `exit` 指令处，实现条件不满足时跳过 if 语句块。

第 6 行 `jmp` 指令的目标地址也被回填为 7，目的是在执行完 if 语句块内部语句后，跳转到程序结束位置 。通过回填，使跳转指令能够准确地控制程序执行流程，保证 if 语句逻辑在中间代码层面的正确实现。

![简单 If 语句生成示例-第二部分](/docs/img/intermediate-code-generation/4.png)

#### If-Else 链式语句的回填

在处理 if-else 语句时，编译器需要生成相应的中间代码来表示条件判断和执行流程。与简单的 if 语句不同，if-else 语句需要处理两个分支的跳转逻辑。

这里，我们不应当考虑单个 if-else 语句的回填，而是考虑 if-else 链式语句的回填。

假设我们有以下链式 if-else 语句：
```plaintext
if (cond1) {

} else if (cond2) {

} else if (cond3) {

...
} else if (condN) {

} [else {} | else if (condN+1) {} | ε] 
```

将其拆分为链式结构：
```plaintext
if cond1 stmt else -> if cond2 stmt else -> if cond3 stmt else -> ... -> if condN stmt [else stmt | else if condN+1 stmt]
```

我们可以明显地发现链式结构的特点：
- 每个 `if` 语句都对应一个 `else`，并且每个 `else` 都与前面的 `if` 语句相连。
- `else` 语句的存在使得每个 `if` 语句都可以有两个分支：一个是 `if` 的分支，另一个是 `else` 的分支。
- `else` 语句可以是一个空语句，也可以是一个新的 `if` 语句，形成嵌套的结构。

那么 `else`可以成为该链式结构的分隔符，既可以表示上一个 `if` 语句的结束，也可以表示下一个 `if` 语句的开始。

关于留空指令的插入时机：
- 对于任意的 `if` 语句，必然存在尾部为栈顶的符号子序列：
    ```plaintext
    if ( bool
    ```
  在 `bool` 语句分析完成后，可以生成留空的 `jmp`（称其为`jmp-next-if` 指令（`jnz` 指令的目标此时就可以填写为为 `jmp` 指令的下一个指令）。
- 对于任意的 `if` 语句，必然存在尾部为栈顶的符号子序列：
    ```plaintext
    if ( bool ) matched_stmt[matched_stmt -> block]
    ```
    在 `matched_stmt` 语句分析完成后，可以生成留空的 `jmp`（称其为`jmp-endif`） 指令（与简单 if 语句不同，此时还无法确定跳转的目标地址）。

关于回填指令的时机：
- 对于任意的 `if` 语句，必然存在尾部为栈顶的符号子序列：
    ```plaintext
    if ( bool ) matched_stmt else if ( bool ) matched_stmt else ... tail_stmt
    ```
    其中`tail_stmt` 是 `else matched_stmt | else if (bool) matched_stmt` 去掉 `}` 之后的部分。

    该尾部遇到 `}` 时，必然触发链式规约，此时就能够确定该链式结构的全部直接子项的结束行跳转位置，即这些子项的 `jmp-endif` 指令的目标地址。

    同样地，当触发链式规约时，我们肯定知道 `jmp-endif` 指令对应的生成代码行号（即 上一个`jmp-next-if` 指令的目标地址）。

    当触发链式规约时，我们可以得出以下公式：
    ```plaintext
    // 任意在链式结构中的 jmp-endif 的目标地址为当前规约状态下的下一个 codeline
    every(dist(jmp-endif)) = next-codeline(REDUCE)

    // 任意在链式结构中的 jmp-next-if 的目标地址为当前 jmp-endif 的所在代码的下一行，注意，留空 jmp-next-if[n] 时，jmp-endif[n] 并未生成
    dist(jmp-next-if[n]) = line_of(jmp-endif[n] + 1)
    ```

    ![If-Else 链式语句的回填](/docs/img/intermediate-code-generation/5.png)

    为了方便，我们可以将尾部 `if-else` 的 `jmp-endif` 指令的目标地址回填时机推迟到头部 `if` 语句被规约完成时。

    部分代码如下：
    ```go
    // matched_stmt → if ( bool ) matched_stmt else matched_stmt
    func MatchedStmtIfElse(w *Walker) error {
        prevEl, _ := w.Tokens.PeekAtK(7)
        children := w.Tokens.PopTopN(7)
        w.Tokens.Push(&ASTNode{
            raw:      joinChildren(children),
            Token:    &lexer.Token{Type: lexer.EXTRA, Val: "matched-stmt-if-else"},
            Children: children,
            Type:     "stmt-if-else",
            Payload:  "!<if-else>",
        })
        if prevEl.Token.SpecificType() != lexer.ReservedWordElse {
            n := w.Environment.LabelStack.PopTopN(2)
            m := w.Environment.EndIfStmtStack.PopTopN(2)
            w.EmitLabel(n[1], fmt.Sprintf("L%d", m[0]+1), "jmp")
            w.EmitGoto(m[0], w.GetCurrentLabelCount())
            w.EmitGoto(m[1], w.GetCurrentLabelCount())
        } else {
            n := w.Environment.LabelStack.PopTopN(2)
            m := w.Environment.EndIfStmtStack.PopTopN(min(2, w.Environment.EndIfStmtStack.Size()))
            w.EmitLabel(n[1], fmt.Sprintf("L%d", m[0]+1), "jmp")
            // only fill the first block
            // in case of `if (condition) { ... } else if (condition) { ... } else { ... }`
            // we should delegate the next endif to the next block
            w.EmitGoto(m[0], w.GetCurrentLabelCount())
            if len(m) > 1 {
                w.Environment.EndIfStmtStack.Push(m[1])
            }
        }
        return nil
    }
    ```

#### 循环语句的回填
在处理循环语句时，编译器需要生成相应的中间代码来表示循环的开始和结束位置。与 if 语句类似，循环语句也需要处理跳转指令的目标地址回填。

唯一与 if 语句不同的是，循环语句的中间代码的最后一个 `jmp` 指令的目标地址始终为循环体的开始位置，初次之外基本上与 if 语句相同，因此不再赘述。

#### Break 语句的回填
在 LR 分析中，`break` 语句的规约总是先于最近循环体的规约完成。这种特性使得我们在处理 `break` 语句时，需要留空跳转指令的目标地址。

我们可以借鉴符号表的作用域概念，维护一个循环体栈来处理 `break` 语句。这是因为 `break` 语句总是跳出离它最近的循环体，这和符号表中变量遵循最近作用域规则相类似。

在 LR 分析的过程中，有以下特性：
1. 对于任何一个循环体，在处理流程中，总是能在遇到 `break` 语句之前就确定该循环体的开始位置。
2. 对于任何一个 `break` 语句，在生成中间代码时，总是能保证其在循环体内部完成中间代码的生成。
3. 对于任何一个循环体，在进行规约操作时，总是能够获取跳出该循环体的目标地址，并且在此时触发对循环体内所有 `break` 语句的目标地址回填操作。

具体操作流程如下：
1. 循环体栈的元素是一个集合，集合中包含所有留空 `jmp` 指令所对应的中间代码位置。
2. 栈的下标对应着不同的循环体，栈顶下标所对应的是离当前 `break` 语句最近的循环体。
3. 当进入一个新的循环体时，创建一个空集合，并将其压入循环体栈中。
4. 当遇到 `break` 语句时，生成一条留空的 `jmp` 指令，同时把该指令在中间代码中的位置添加到栈顶集合里。
5. 当离开当前循环体时（即循环体触发规约操作时），将栈顶的集合弹出，然后把集合中所有留空 `jmp` 指令的目标地址回填为当前代码行号。

![break 语句的回填](/docs/img/intermediate-code-generation/6.png)

## 结果

这里给出一些简单的中间代码生成示例，供读者参考：

### 多重 if 语句的中间代码生成示例
输入：
```plaintext
{
    if (1 == 1) {
        int a;
        a = 1 + 1;
    } else if (2 == 2) {
        int a;
        a = 1 - 1;
        if (5 == 5) {} else if (6 == 6) {}
    } else if (3 == 3) {
        int a;
        a = 1 * 1;
    } else {
        int a;
        a = 1 / 1;
    }
}
```

输出：
```plaintext
Three Address Code:
L0             jmp               L1
L1              eq    $(0x10000000)                1                1
L2             cmp    $(0x10000001)    $(0x10000000)                0
L3             jnz               L5    $(0x10000001)
L4             jmp               L9
L5           alloc    $(0x10000002)                4                0
L6             add    $(0x10000003)                1                1
L7             mov    $(0x10000002)    $(0x10000003)
L8             jmp              L41
L9              eq    $(0x10000004)                2                2
L10            cmp    $(0x10000005)    $(0x10000004)                0
L11            jnz              L13    $(0x10000005)
L12            jmp              L29
L13          alloc    $(0x10000006)                4                0
L14            sub    $(0x10000007)                1                1
L15            mov    $(0x10000006)    $(0x10000007)
L16             eq    $(0x10000008)                5                5
L17            cmp    $(0x10000009)    $(0x10000008)                0
L18            jnz              L20    $(0x10000009)
L19            jmp              L22
L20            nop                 
L21            jmp              L28
L22             eq    $(0x1000000a)                6                6
L23            cmp    $(0x1000000b)    $(0x1000000a)                0
L24            jnz              L26    $(0x1000000b)
L25            jmp              L28
L26            nop                 
L27            jmp              L28
L28            jmp              L41
L29             eq    $(0x1000000c)                3                3
L30            cmp    $(0x1000000d)    $(0x1000000c)                0
L31            jnz              L33    $(0x1000000d)
L32            jmp              L37
L33          alloc    $(0x1000000e)                4                0
L34            mul    $(0x1000000f)                1                1
L35            mov    $(0x1000000e)    $(0x1000000f)
L36            jmp              L41
L37          alloc    $(0x10000010)                4                0
L38            div    $(0x10000011)                1                1
L39            mov    $(0x10000010)    $(0x10000011)
L40            jmp              L41
L41           exit                0
-------------------------------------
Scope: 1
-------------------------------------
Scope: 2
Variable: a, Type: int, Address: 0x10000002
-------------------------------------
Scope: 3
Variable: a, Type: int, Address: 0x10000006
-------------------------------------
Scope: 4
-------------------------------------
Scope: 5
-------------------------------------
Scope: 6
Variable: a, Type: int, Address: 0x1000000e
-------------------------------------
Scope: 7
Variable: a, Type: int, Address: 0x10000010
```

### 数组声明与变量赋值的中间代码生成示例
输入：
```plaintext
{
    int[6] a;
    int[8] b;
    a[2] = 1;
    if (a == 3) {
        int[2048] a;
        if (true) {
            a[1024] = 4;
        }
    }
    int c;
    c = a[2];
}
```

输出：
```plaintext
Three Address Code:
L0             jmp               L1
L1           alloc    $(0x10000000)               24                0
L2           alloc    $(0x10000006)               32                0
L3             mov    $(0x10000002)                1
L4              eq    $(0x1000000e)    $(0x10000000)                3
L5             cmp    $(0x1000000f)    $(0x1000000e)                0
L6             jnz               L8    $(0x1000000f)
L7             jmp              L15
L8           alloc    $(0x10000010)             8192                0
L9             cmp    $(0x10000810)                1                0
L10            jnz              L12    $(0x10000810)
L11            jmp              L14
L12            mov    $(0x10000410)                4
L13            jmp              L14
L14            jmp              L15
L15          alloc    $(0x10000811)                4                0
L16            mov    $(0x10000811)    $(0x10000002)
L17           exit                0
-------------------------------------
Scope: 1
Variable: a, Type: !ptr<int>, Address: 0x10000000
Array Size: 6, Element Size: 4
Variable: b, Type: !ptr<int>, Address: 0x10000006
Array Size: 8, Element Size: 4
Variable: c, Type: int, Address: 0x10000811
-------------------------------------
Scope: 2
Variable: a, Type: !ptr<int>, Address: 0x10000010
Array Size: 2048, Element Size: 4
-------------------------------------
Scope: 3
```

### 循环体与 Break 语句的中间代码生成示例
输入：
```plaintext
{
    do {
        if (1 == 1) {
            break;
        }
        while(4 == 4 == 2)  {
            if (5 == 5) {
                break;
                if (6 == 6) {
                    break;
                }
            }
        } 
        if (2 == 2) {
            break;
        }
    } while(3 == 3);
}
```

输出：
```plaintext
Three Address Code:
L0             jmp               L1
L1              eq    $(0x10000000)                1                1
L2             cmp    $(0x10000001)    $(0x10000000)                0
L3             jnz               L5    $(0x10000001)
L4             jmp               L7
L5             jmp              L35
L6             jmp               L7 <- if 语句的残留(jmp-endif)，后续可以优化(相邻两个 jmp 中，最后一个 jmp 失效)
L7              eq    $(0x10000002)                4                4
L8              eq    $(0x10000003)    $(0x10000002)                2
L9             cmp    $(0x10000004)    $(0x10000003)                0
L10            jnz              L12    $(0x10000004)
L11            jmp              L25
L12             eq    $(0x10000005)                5                5
L13            cmp    $(0x10000006)    $(0x10000005)                0
L14            jnz              L16    $(0x10000006)
L15            jmp              L24
L16            jmp              L25
L17             eq    $(0x10000007)                6                6
L18            cmp    $(0x10000008)    $(0x10000007)                0
L19            jnz              L21    $(0x10000008)
L20            jmp              L23
L21            jmp              L25
L22            jmp              L23
L23            jmp              L24
L24            jmp               L7
L25             eq    $(0x10000009)                2                2
L26            cmp    $(0x1000000a)    $(0x10000009)                0
L27            jnz              L29    $(0x1000000a)
L28            jmp              L31
L29            jmp              L35
L30            jmp              L31
L31             eq    $(0x1000000b)                3                3
L32            cmp    $(0x1000000c)    $(0x1000000b)                0
L33            jnz               L1
L34            jmp              L35
L35           exit                0
-------------------------------------
Scope: 1
-------------------------------------
Scope: 2
-------------------------------------
Scope: 3
-------------------------------------
Scope: 4
-------------------------------------
Scope: 5
-------------------------------------
Scope: 6
-------------------------------------
Scope: 7
```