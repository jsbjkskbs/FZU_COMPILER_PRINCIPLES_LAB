# Intermediate Code Generation

## Purpose of this Lab
Generate intermediate code using syntax-directed techniques. Design semantic rules for intermediate code and embed them into syntax rules. During the syntax analysis process, output equivalent intermediate code.

## Tasks
Implement the following compiler functionalities:
1. Output the syntax analysis process.
2. Output intermediate code in three-address code format.

## Grammar
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

### Grammar Issues

#### Left Recursion
In syntax-directed translation during compiler construction, the presence of left-recursive grammar introduces specific complexities in the code generation phase. Consider the following left-recursive grammar fragment:

```plaintext
bool → bool || join
join → join && equality
equality → equality == rel
equality → equality != rel
```

The defining characteristic of such directly left-recursive productions is that non-terminals like `bool`, `join`, and `equality` directly reference themselves at the leftmost position of their right-hand side, forming a pattern like `A → A α`. While LR (1) parsers can handle grammars with direct left recursion (by constructing appropriate parsing tables to make **shift-reduce** decisions), this recursive structure introduces uncertainty in the timing of code generation during syntax-directed translation.

For example, consider the production `stmt → if ( bool ) stmt else stmt`:

When the parser processes the `bool` expression, it needs to generate corresponding conditional evaluation code (e.g., jump instructions). However, the left-recursive structure causes the reduction of `bool` to exhibit a chain-like behavior: for the production `bool → bool || join`, each reduction of `bool` actually processes an intermediate node in the recursive chain rather than the final complete expression. During bottom-up parsing, when encountering a reduction involving the non-terminal `bool`, the parser can only determine that it is processing a subexpression within the recursive structure (e.g., the left-hand `bool` corresponds to an already processed subexpression, while the right-hand `|| join` corresponds to subsequent operations yet to be processed). However, it cannot ascertain whether this `bool` represents the top-level expression requiring the generation of a complete conditional result.

Specifically, left recursion introduces the following issues:

1. According to the standard logic of syntax-directed translation, code generation for expressions should follow the order **"subexpressions first, then operators"**. However, the left-recursive structure causes the non-terminal `bool` to be reduced multiple times during parsing (each reduction corresponding to a step in the recursive chain). At each reduction step, it must be decided whether to generate the code for the complete expression result.
    > For instance, when processing `bool → bool || join`, the code for the left-hand `bool` subexpression has already been generated. However, if this `bool` is not the final top-level expression (but rather part of another `bool || join` structure), generating the complete result code at this point would disrupt the semantic linkage of subsequent `||` operations.
2. In the construction of LR (1) parsing tables, items with the same core but different lookahead symbols (known as "kernel items") may correspond to different semantic actions. When the left-recursive production's non-terminal `bool` serves as the core of such items, the parser can only decide to shift or reduce based on the current lookahead symbol (i.e., the next input symbol). However, it cannot predict whether this `bool` will later serve as the left-hand side of a larger left-recursive structure. This uncertainty makes it difficult for the code generator to determine during a single reduction whether it needs to generate the final result's target code (e.g., the terminal node for conditional jumps), as the current reduction might only represent an intermediate step in the recursive chain rather than the expression's endpoint.
3. For the short-circuit evaluation semantics of logical expressions like `bool || join` (where the right-hand operand does not need to be evaluated if the left-hand operand is true), code generation requires inserting conditional jumps after evaluating the left-hand `bool` to skip the evaluation of the right-hand `join`. However, the left-recursive structure causes `bool` to appear multiple times as an ancestor node of the right-hand `join`. Each time `bool` is reduced, it must be determined whether it is part of the top-level expression to decide the binding of the jump target. If an intermediate recursive layer's `bool` is mistakenly treated as the top-level expression and generates a jump, subsequent operand code may be erroneously skipped, breaking semantic correctness.

The discrepancy between multiple reductions of non-terminals (as intermediate nodes in the recursive chain) and the need to identify "endpoint expressions" for code generation introduces semantic ambiguity. The parser can only make syntax-level reduction decisions based on the current symbol and lookahead, but it cannot determine whether the non-terminal represents a complete expression's final result at the semantic level. This mismatch between syntax structure and semantic processing makes it challenging for the code generator to produce complete result code at the appropriate time.

#### Dangling Else Problem

The above grammar exhibits the **dangling else problem**, which increases the complexity of parsing `if-else` statements during syntax analysis in compiler construction. Consider the following grammar fragment containing `if-else` statements:

```plaintext
stmt → if ( bool ) stmt 
    | if ( bool ) stmt else stmt 
```

The structural characteristics of such productions are:

- The `if` statement has two forms: one without an `else` clause and one with an `else` clause. This structure creates ambiguity in nested `if-else` statements, where the matching target of the `else` clause is unclear, leading to the **dangling else problem**. In bottom-up parsing, although the parser can perform **shift-reduce** operations based on the productions, the grammar does not provide explicit rules for determining which `if` an `else` clause should pair with.

For example, consider the nested `if-else` statement `if (condition1) if (condition2) statement1; else statement2;`. When the parser processes this statement, it must determine which `if` the `else` clause belongs to. However, based on the grammar's structural characteristics, the matching process for `else` is highly ambiguous:

- In this statement, the `else` can either match the inner `if (condition2)` or the outer `if (condition1)`. During the bottom-up parsing process, when encountering the `else` symbol, the parser can only determine that it is part of an `if-else` structure but cannot decide which unmatched `if` it should associate with.

Specifically, the **dangling else problem** introduces the following issues during syntax analysis:
1. According to the standard logic of syntax analysis, every `else` clause should explicitly pair with an `if`. However, this grammar lacks clear rules for `else` matching, making it difficult to analyze nested `if-else` statements. For instance, in the above nested statement, without explicit rules, the parser cannot determine whether the `else` belongs to the inner `if` or the outer `if`, leading to ambiguous semantic interpretations of the code.
2. Different ways of matching `else` during syntax analysis result in different parse trees. These different parse trees represent distinct code semantics, introducing uncertainty into the analysis results. For example, pairing `else` with the inner `if` versus the outer `if` leads to different conditional logic and execution flows, which affects subsequent stages like semantic analysis and code generation.
3. To resolve the **dangling else problem**, the parser requires additional rules or strategies to determine the pairing of `else`. This may involve backtracking, state management, or other complex operations, increasing the complexity and implementation difficulty of syntax analysis. For instance, the parser might need to track encountered `if` statements to accurately match `else` clauses when they are encountered.

Without explicit rules for `else` matching, the pairing of `else` with `if` becomes ambiguous, deviating from the determinism required for syntax analysis. The parser can only make **shift-reduce** decisions at the syntax level based on the current symbol but cannot determine at the semantic level which `if` the `else` should correctly match.

### Grammar Improvements

#### Left Recursion

Not all left-recursive grammars need to be rewritten. The primary criterion is whether the left-recursive grammar significantly impacts the implementation of the parser. In this lab, the presence of left recursion causes the parser to struggle with determining when to generate the complete code for `bool` expressions.

We can resolve this issue by rewriting the left-recursive grammar into a non-left-recursive form. Below is the improved grammar:

```plaintext
bool → bool'
bool' → bool' || join | join
...
```

In this improved grammar, we introduce a new non-terminal `bool'`, which represents the right-hand side of `bool`. This ensures that when `bool` is reduced to a complete expression, it is not part of any intermediate node in a recursive chain.

#### Dangling Else Problem

In this lab, we can address the **dangling else problem** by introducing two new non-terminals: `matched_stmt` and `unmatched_stmt`. Below is the improved grammar:

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

In this improved grammar, we divide `stmt` into two categories: `matched_stmt` and `unmatched_stmt`. `matched_stmt` represents statements where every `if` has a corresponding `else`, while `unmatched_stmt` represents statements where an `if` does not yet have a matching `else`. This ensures that every `else` clause is correctly paired with an `if` statement, thereby resolving the **dangling else problem**.

## Implementation
### Generating Intermediate Code During Parsing
We can add a semantic action to each production rule to generate intermediate code during parsing. When the parser reduces a production rule, the corresponding semantic action is executed.

For example, for the production `unary → - unary`:
```plaintext
-a

// Reduction chain
id    → loc      // Query the symbol table to get the value of id
loc   → unary    // loc becomes a child node of unary
unary → - unary  // - unary becomes a child node of unary
```

When reducing `unary → - unary`, we can generate the intermediate code:
```plaintext
mov $(t1), -$(t2) // Negate the value of unary and store it in t1
```

> Note: The three-address code (quadruples) in this lab is represented in the format `op dist arg1 arg2`.
> This lab depends on the implementation of the parser. The grammar and parser logic must be correct to generate accurate intermediate code.

### Symbol Table (Symbol Scope)
During intermediate code generation, a symbol table is maintained to store information about variables, such as their types and scopes.

The symbol table can be implemented using a stack. When entering a new scope, the current scope's symbol table is pushed onto the stack. When leaving a scope, the top symbol table is popped from the stack. This approach effectively manages nested scopes.
> Alternatively, a singly linked list can be used to implement the symbol table. Each node in the list represents a scope, and each node stores the symbol table for that scope.
> When entering a new scope, a new node is inserted at the head of the list. When leaving a scope, the head node is removed.
> This approach also effectively manages nested scopes.
> In essence, such a singly linked list is equivalent to a stack.

Illustration:
![Symbol Table Illustration](/docs/img/intermediate-code-generation/1.png)

Process:
1. The outer scope declares `int a`, and symbol table **S1** records `a = $(0)`. Executing `a = 1` updates `$(0) = 1`.
Entering the inner block:
2. Initially, `a = 1` still references `$(0) = 1` (since it hasn't been redeclared yet). Declaring `int a` in the inner block switches the symbol table to **S2**, recording `a = $(1)`. Executing `a = 1` updates `$(1) = 1`.
3. Declaring `int c` in the inner block adds `c = $(2)` to symbol table **S2**.
4. Exiting the inner block drops **S2**, returning to the outer scope.
5. Declaring `int c` in the outer scope adds `c = $(3)` to symbol table **S1**. Executing `c = 1` updates `$(3) = 1`. Finally, **S1** retains `a = $(0)` and `c = $(3)`.

> Note: A loop stack will also need to be maintained later to handle `break` statements.
> At this point, a singly linked list is not strictly necessary, as only a stack of loop end positions is required (no complex structures need to be stored).

### Abstract Syntax Tree

Abstract Syntax Tree (AST) is a data structure that provides an abstract representation of source code. It describes the syntax structure of the source code in a tree-like format.
- The AST consists of various nodes, each representing a specific syntactic element in the source code, such as variable declarations, function calls, expressions, and statement blocks. These nodes have different types to distinguish the kinds of syntactic elements they represent. For example, a node representing a variable declaration might have a type like "Decl."
- In addition to type information, nodes may also contain other attributes to store specific details about the syntactic element they represent. For instance, a variable declaration node might include attributes for the variable's name and data type, while a function call node might include the function's name and a list of arguments passed to it. These attributes provide detailed descriptions of the syntactic elements, enabling further analysis and processing.
- Nodes are organized in a hierarchical tree structure through parent-child and sibling relationships. This structure reflects the nesting and logical order of syntactic elements in the source code. For example, in a function definition, the statement block inside the function can be seen as a child node of the function definition node, while the individual statements within the block are child nodes of the block node. This hierarchical structure clearly represents the overall syntax architecture of the source code, making it easier to analyze and understand.

An LR parser is a bottom-up parser that starts from the root node of the AST and progressively analyzes down to the leaf nodes. Essentially, LR parsing is a depth-first traversal (DFS).
> A stack of AST nodes can be maintained during parsing. As the parser processes the grammar, it pushes the right-hand side nodes of each production onto the stack. When a production is reduced, the corresponding nodes are popped from the stack, combined into a new node, and the new node is pushed back onto the stack.
> This approach allows the construction of a complete AST during the parsing process.

During LR parsing:
- **SHIFT**: When a shift operation is performed, the current input token is converted into a corresponding AST leaf node, such as a literal or identifier node, and pushed onto the node stack. For example, if the input token is the number `5`, an `AST` node of type `Num` is created and pushed onto the stack.
- **REDUCE**: During a reduction, the parser pops nodes from the stack corresponding to the right-hand side symbols of the current production. A new AST node is then created to represent the left-hand side non-terminal of the production. The popped nodes are connected as children of the new node in the order of the production's right-hand side. Finally, the new node is pushed onto the stack. For example, for the production `E → E + T`, the nodes representing `E`, `+`, and `T` are popped from the stack, a new `E` node is created with these nodes as its children, and the new `E` node is pushed onto the stack.

Consider a simple arithmetic expression grammar:
```plaintext
E → E + T
E → T
T → T * F
T → F
F → ( E )
F → num
```
During an LR parsing process, given the input token sequence `3 + 4 * 5`, the parser constructs the AST as follows:
1. The input expression `3 + 4 * 5` is tokenized into `[num(3), +, num(4), *, num(5)]`.
2. Upon encountering `num(3)`, a `Num` node is created and pushed onto the stack.
3. Using the production `F → num`, the `Num` node is popped, an `F` node is created and connected, and the `F` node is pushed onto the stack.
4. The production `T → F` is applied, the `F` node is popped, a `T` node is created and connected, and the `T` node is pushed onto the stack.
5. Similar steps are followed for subsequent tokens until the entire expression is analyzed, and the root node of the AST representing the entire expression is at the top of the stack.

![Example of Parsing 3 + 4 * 5](/docs/img/intermediate-code-generation/2.png)

After parsing is complete, the node at the top of the stack is the root of the constructed AST, which fully represents the syntax structure of the input source code.

#### AST and Code Optimization
The Abstract Syntax Tree (AST) plays a crucial role in the intermediate code generation and optimization phases of a compiler. It not only provides a structured representation of the source code but also serves as the foundation for subsequent code optimization and generation. Below are some key roles of the AST in intermediate code generation and optimization:
- Each node in the AST can be mapped to one or more intermediate code instructions, enabling the transformation of source code into intermediate code. By traversing the AST, the compiler can generate corresponding three-address code or other forms of intermediate representation.
- The compiler can analyze the structure of the AST to identify redundant computations, constant expressions, dead code, and more, and perform corresponding optimizations. For example, constant folding can evaluate constant expressions at compile time, reducing runtime computation overhead. Dead code elimination can remove unnecessary variable definitions and usages, reducing memory usage.
- The compiler can analyze conditional statements, loops, and other structures in the AST to identify unnecessary branches and loops for optimization. For instance, loop unrolling can reduce the number of iterations in a loop, improving execution efficiency.

> Note: In this lab, intermediate code is generated directly during reduction, and the AST is constructed primarily to facilitate subsequent code generation.
>
> However, in practical compilation, constructing the AST and generating intermediate code are two separate stages. Typically, the AST is constructed first, followed by intermediate code generation.
>
> During the intermediate code generation phase, the compiler traverses the AST and generates corresponding intermediate code instructions based on the types and attributes of the nodes.
>
> Readers interested in this topic can try enhancing the repository's `parser/ast.go` and `parser/gen_rules.go` to implement AST construction first, followed by optimization and intermediate code generation.

### Intermediate Code Generation for Key Rules

Let's start from the simple leaf nodes of the Abstract Syntax Tree (AST) and gradually analyze upwards to the root node.

#### Token as AST Leaf Nodes
Tokens obtained from the lexical analyzer are directly treated as leaf nodes of the AST and only require simple conversion.

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
> - `p.Reflect(token)` is a reflection function used to retrieve the type of the token.
> - `ASTNode` is a structure containing information about the token's raw value, type, child nodes, etc.
>   - `raw` stores the expression associated with the node.
>   - `Token` stores information about the token (it can also store allocated addresses for extension).
>   - `Children` stores child nodes (pushed onto the stack during syntax analysis).
>   - `Type` stores the type of the node (e.g., `Num`, `Id`, `Bool`).
>   - `Payload` stores additional information (though it is not heavily used here).

#### Factor AST Nodes
The `factor` grammar rule is as follows:
```plaintext
factor → ( bool ) | loc | num | real | true | false
```

During syntax analysis, intermediate code can be generated based on the type of `factor`.

Here is the semantic action for one of the productions:
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

> - When the syntax analyzer reduces `factor → false`, the `FactorFalse` function is executed.
> - The top of the stack must be the AST node for `false`, which is popped using `PopTopN(1)`.
> - A new `ASTNode` is created with the type `factor-false`, and the popped node is added as a child.
> - Finally, the new node is pushed onto the stack (indicating the reduction is complete).

#### Unary AST Nodes
The `unary` grammar rule is as follows:
```plaintext
unary → ! unary | - unary | factor
```

During syntax analysis, intermediate code can be generated based on the type of `unary`.

For generating intermediate code for `-1`:
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
> - When the syntax analyzer reduces `unary → ! unary`, the `UnaryNot` function is executed.
> - `w.SymbolTable.TempAddr(4)` allocates a temporary variable of size 4 bytes.
> - `addr` is the address of the temporary variable, formatted as `$(0x00000000)`.
> - `children` are the popped nodes, including the `!` and `unary` token nodes.
> - `joinChildren(children)` is a function that concatenates the raw values of the child nodes.
> - `w.Tokens.Push` pushes the new node onto the stack, indicating the reduction is complete.
> - The program generates the following intermediate code:
> ```plaintext
> not $(dist) $(src)
> ```
> - `dist` is the address of the temporary variable, and `src` is the value of `unary`.

#### Term AST Nodes
```plaintext
term → term * unary | term / unary | unary
```

During syntax analysis, intermediate code can be generated based on the type of `term`.
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

> - When the syntax analyzer reduces `term → term * unary`, the `TermMult` function is executed.
> - `w.SymbolTable.TempAddr(4)` allocates a temporary variable of size 4 bytes.
> - `children` are the popped nodes, including the `term`, `*`, and `unary` token nodes.
> - `w.Tokens.Push` pushes the new node onto the stack, indicating the reduction is complete.
> - The program generates the following intermediate code:
> ```plaintext
> mul $(dist) $(src1) $(src2)
> ```
> - `dist` is the address of the temporary variable, `src1` is the value of `term`, and `src2` is the value of `unary`.
> - Here, `mul` represents the multiplication operation, `dist` is the result storage address, and `src1` and `src2` are the operands.

The intermediate code generation for `expr`, `rel`, `equality`, `join`, and `bool` follows a similar pattern, where corresponding intermediate code is generated during the reduction of production rules.

#### Intermediate Code Backfilling

During the compilation process, some information is unknown in the early stages of intermediate code generation, such as the target address of jump instructions. Backfilling aims to supplement this information once it becomes available, ensuring that the intermediate code is complete and can be correctly understood and executed by the target machine.

Backfilling is typically performed after completing the syntax and semantic analysis of the source program and before generating the target code. When the compiler gathers enough information to determine previously unknown parts, it executes the backfilling operation. For example, when handling conditional or loop statements, the target address of jump instructions may initially be unknown. Only after analyzing the entire statement block and determining the jump target can backfilling occur.

Compilers often use data structures like symbol tables to record information that needs backfilling, such as variable storage addresses. For jump instruction target addresses, specialized lists or stacks may be used to record jump instructions with unresolved targets. Once the target address is determined, the corresponding instructions are retrieved from the list or stack for backfilling.

> For example, in the conditional statement `if (a > 10) goto L1;`, the address of **L1** in `goto L1` is unknown during intermediate code generation. The instruction and related information are recorded. When the end of the conditional statement is analyzed, and the actual address of **L1** is determined, the address can be backfilled into the `goto L1` instruction.

##### Placeholder Instructions

In general, it is essential to determine when to generate placeholder jump instructions. This is not done after analyzing the entire statement block but rather when the compiler identifies the start of a statement block.

The grammar can help identify the start of a statement block, which must trigger a reduction operation.

Consider the following example:
```plaintext
matched_stmt → do stmt while ( bool ) ;
stmt → matched_stmt
matched_stmt → block
block → { decls stmts }
decls → decls decl | ε
```

During the reduction process, assuming the syntax is correct, there exists a sequence of symbols at the top of the stack:
```plaintext
do {
```

The next state of the symbol stack must be:
```plaintext
do { decls
```
> 1. To reduce `do stmt`, `matched_stmt` must first be reduced.
> 2. To reduce `matched_stmt`, `block` must first be reduced.
> 3. To reduce `block`, `{ decls stmts }` must first be reduced.
> 4. Therefore, `decls` is a necessary component for reduction.
> 5. From the `decls` production, `decls` may be empty. Thus, for any internal expression, there must exist a sequence of symbols at the top of the stack:
>    ```plaintext
>    do { decls
>    ```

Hence, the `ε` production of `decls` can serve as the timing for generating placeholder jump instructions.

##### Backfilling for Simple If Statements

When processing the conditional expression of an `if` statement and generating jump-related intermediate code, the compiler has not yet fully analyzed the entire `if` statement block and subsequent code. For example, in line 2 of the diagram, the `jmp` instruction's target address (corresponding to the end of the `if` statement block) cannot be determined at the time of generation.

We can use a backfilling mechanism to handle such cases. Initially, the target address of the jump instruction is left blank. Later, after analyzing the `if` statement block and determining the actual location of the jump target (e.g., the code line number corresponding to the end of the `if` statement block), the correct address is filled in. For instance, in line 2 of the diagram, the `jmp` instruction's target address (dashed line) and the target address of the `jmp` instruction in line 6 (dashed line) are initially left blank and later completed through backfilling once sufficient information is available.

![Simple If Statement Code Generation Example - Part 1](/docs/img/intermediate-code-generation/3.png)

As the compilation progresses, when the compiler completes the analysis of the `if` statement block and subsequent code, the correct jump target locations are determined, and backfilling is performed. For example, in the second line of the diagram, the `jmp` instruction's originally blank target address is backfilled to `7`, which directly jumps to the program's end `exit` instruction, ensuring that the `if` statement block is skipped when the condition is not met.

The target address of the sixth line `jmp` instruction is also backfilled to `7`, ensuring that after executing the internal statements of the `if` block, the program jumps to the end. Through backfilling, jump instructions can accurately control program execution flow, ensuring the logical correctness of the `if` statement at the intermediate code level.

![Simple If Statement Code Generation Example - Part 2](/docs/img/intermediate-code-generation/4.png)

#### Backfilling for If-Else Chain Statements

When handling if-else statements, the compiler needs to generate intermediate code to represent conditional evaluations and execution flows. Unlike simple if statements, if-else statements require handling the jump logic for two branches.

Here, we focus on backfilling for if-else chain statements rather than individual if-else statements.

Suppose we have the following chain of if-else statements:
```plaintext
if (cond1) {

} else if (cond2) {

} else if (cond3) {

...
} else if (condN) {

} [else {} | else if (condN+1) {} | ε] 
```

This can be broken down into a chain structure:
```plaintext
if cond1 stmt else -> if cond2 stmt else -> if cond3 stmt else -> ... -> if condN stmt [else stmt | else if condN+1 stmt]
```

The chain structure has the following characteristics:
- Each `if` statement corresponds to an `else`, and each `else` is connected to the preceding `if` statement.
- The presence of an `else` statement allows each `if` statement to have two branches: one for the `if` branch and one for the `else` branch.
- The `else` statement can be empty or a new `if` statement, forming a nested structure.

The `else` statement acts as a delimiter for the chain structure, marking the end of the preceding `if` statement and the start of the next `if` statement.

**Timing for Placeholder Instructions:**
- For any `if` statement, there must be a sequence of symbols at the top of the stack:
    ```plaintext
    if ( bool
    ```
  After analyzing the `bool` expression, a placeholder `jmp` instruction (referred to as `jmp-next-if`) can be generated. The target of the `jnz` instruction can now be set to the next instruction after the `jmp`.
- For any `if` statement, there must also be a sequence of symbols at the top of the stack:
    ```plaintext
    if ( bool ) matched_stmt[matched_stmt -> block]
    ```
    After analyzing the `matched_stmt`, a placeholder `jmp` instruction (referred to as `jmp-endif`) can be generated. Unlike simple if statements, the jump target address cannot yet be determined.

**Timing for Backfilling Instructions:**
- For any `if` statement, there must be a sequence of symbols at the top of the stack:
    ```plaintext
    if ( bool ) matched_stmt else if ( bool ) matched_stmt else ... tail_stmt
    ```
    Here, `tail_stmt` represents `else matched_stmt | else if (bool) matched_stmt` without the closing `}`.

    When the closing `}` is encountered, it triggers a chain reduction. At this point, the end jump locations for all direct children of the chain structure can be determined, i.e., the target addresses for their `jmp-endif` instructions.

    Similarly, during chain reduction, the code line number for the `jmp-next-if` instruction's target address can also be determined.

    During chain reduction, the following formulas apply:
    ```plaintext
    // The target address for any jmp-endif in the chain structure is the next code line during the REDUCE operation.
    every(dist(jmp-endif)) = next-codeline(REDUCE)

    // The target address for any jmp-next-if[n] in the chain structure is the code line of jmp-endif[n] + 1.
    dist(jmp-next-if[n]) = line_of(jmp-endif[n] + 1)
    ```

    ![Backfilling for If-Else Chain Statements](/docs/img/intermediate-code-generation/5.png)

    To simplify, the backfilling of the tail `if-else` statement's `jmp-endif` target address can be deferred until the head `if` statement is fully reduced.

    Partial code example:
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

#### Backfilling for Loop Statements
When handling loop statements, the compiler needs to generate intermediate code to represent the start and end of the loop. Similar to if statements, loop statements also require backfilling the target addresses of jump instructions.

The only difference from if statements is that the last `jmp` instruction in the intermediate code for a loop always targets the start of the loop body. Other aspects are essentially the same as for if statements and will not be repeated here.

#### Backfilling for Break Statements
In LR parsing, the reduction of a `break` statement always occurs before the reduction of the nearest enclosing loop. This characteristic means that when handling `break` statements, the target address of the jump instruction must be left blank initially.

We can use a loop stack, similar to the concept of a symbol table's scope, to handle `break` statements. This is because `break` statements always exit the nearest enclosing loop, analogous to how variables follow the nearest scope rule in a symbol table.

During LR parsing, the following characteristics hold:
1. For any loop, its start position is always determined before encountering a `break` statement.
2. For any `break` statement, its intermediate code generation is always completed within the loop body.
3. For any loop, during its reduction, the target address for exiting the loop is determined, triggering the backfilling of all `break` statements within the loop.

**Steps:**
1. The loop stack contains sets, where each set includes all placeholder `jmp` instructions corresponding to the intermediate code positions.
2. The stack index corresponds to different loops, with the top index representing the nearest enclosing loop for the current `break` statement.
3. When entering a new loop, create an empty set and push it onto the loop stack.
4. When encountering a `break` statement, generate a placeholder `jmp` instruction and add its position in the intermediate code to the top set.
5. When exiting the current loop (i.e., during the loop's reduction), pop the top set from the stack and backfill the target addresses of all placeholder `jmp` instructions in the set with the current code line number.

![Backfilling for Break Statements](/docs/img/intermediate-code-generation/6.png)

## Results

Here are some simple examples of intermediate code generation for reference:

### Intermediate Code Generation for Nested If Statements
Input:
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

Output:
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

### Array Declaration and Assignment Intermediate Code Generation Example
Input:
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

Output:
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

### Loop Statement Intermediate Code Generation Example
Input:
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

Output:
```plaintext
Three Address Code:
L0             jmp               L1
L1              eq    $(0x10000000)                1                1
L2             cmp    $(0x10000001)    $(0x10000000)                0
L3             jnz               L5    $(0x10000001)
L4             jmp               L7
L5             jmp              L35
L6             jmp               L7 <- Residual from the if statement (jmp-endif), can be optimized later (the last jmp in two consecutive jmps becomes redundant)
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