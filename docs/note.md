# 备注

依我看，编译原理是大三阶段中最能够拉开差距的课程，如果你想要在大三下之后依靠推免（GPA）进入你想要的研究生院，**编译原理是你必须要学好的一门课**。
> 至少来说编译原理我看起来考得不怎么样，但居然能在前 5(/156) 名中提升 1 名，<s>你就能发现编译原理是多么强大</s>

这里，我将会告诉你在这个实验中你需要注意的地方。

## 偏差

### 词法解析

你可以发现，在词法实验中，我并没有按照老师的要求去写一个 DFA，而是直接使用了一个简单函数去完成字符流转 Token 的转换。原因是因为我觉得 DFA 过于复杂，且不利于调试。

问题在于，如果是读取一些简单的 Token，或许直接用函数就可以了，但如果是读取一些复杂的 Token（比如说比较灵活的浮点数 1.0e-2），你就会发现 DFA 的优势了。

如果你尝试用和我一样的方式，想要实现一个完备的浮点数 Tokenizer，你会发现你需要写出非常复杂的代码，这和我们的目标是相悖的。

所以，建议你还是按照老师的要求去实现 DFA。毕竟考试的时候，DFA是必考的（**将一个正则表达式转化为 NFA，再将 NFA 转化为 DFA，最后将 DFA 转化为最小化 DFA**）。

> 一个浮点数的 Tokenizer 性质如下：
>
> ``` plaintext
> digit → 0|1|2|…|9
> digits → digit digit*
> fraction → .digits | ε
> exponent → ( E(+|-|ε)digits ) | ε
> number → digits fraction exponent
> ```
>
>
> NFA 如下：
> ``` mermaid
> %%{ init: { 'flowchart': { 'curve': 'basis', 'nodeSpacing': 20, 'rankSpacing': 60, 'animate': true } } }%%
> flowchart LR
>    any@{ shape: text, label: " " }
>    start@{ shape: text, label: "start" }
>    
>    0((0))
>    1((1))
>    2((2))
>    3((3))
>    4((4))
>    5((5))
>    6(((6)))
>    start e@==> 0
>    0 e1@==>|d| 1
>    1 e2@==>|d| 1
>    1 e3@==>|.| 2
>    1 e4@==>|ε| 3
>    2 e5@==>|d| 3
>    3 e6@==>|d| 3
>    3 e7@==>|E| 4
>    3 e8@==>|ε| 6
>    4 e9@==>|\+| 5
>    4 e10@==>|\-| 5
>    4 e11@==>|ε| 5
>    5 e12@==>|d| 6
>    6 e13@==>|d| 6
>   
>    classDef animate stroke-dasharray: 9,5,stroke-dashoffset: 900,animation: dash 25s linear infinite;
>    class e,e1,e2,e3,e4,e5,e6,e7,e8,e9,e10,e11,e12,e13 animate;
> ```
>
> DFA 如下：
> ``` mermaid
> %%{ init: { 'flowchart': { 'curve': 'basis', 'nodeSpacing': 20, 'rankSpacing': 80, 'animate': true } } }%%
> flowchart LR
>    start@{ shape: text, label: "start" }
>    any@{ shape: text, label: " " }
>    
>    0((0))
>    1,3,6(((1,3,6)))
>    2((2))
>    3,6(((3,6)))
>    4,5((4,5)) ~~~ any
>    5((5))
>    6(((6)))  ~~~ any
>    start e@==> 0
>    0 e1@==>|d| 1,3,6
>    1,3,6 e2@==>|d| 1,3,6
>    1,3,6 e3@==>|.| 2
>    2 e4@==>|d| 3,6
>    3,6 e5@==>|d| 3,6
>    3,6 e6@==>|E| 4,5
>    1,3,6 e7@==>|E| 4,5
>    4,5 e8@==>|\+| 5
>    4,5 e9@==>|\-| 5
>    4,5 e10@==>|d| 6
>    5 e11@==>|d| 6
>    
>    classDef animate stroke-dasharray: 9,5,stroke-dashoffset: 900,animation: dash 25s linear infinite;
>    class e,e1,e2,e3,e4,e5,e6,e7,e8,e9,e10,e11 animate;
> ```

词法解析器的练手，可以尝试 LeetCode 上的 [正则表达式匹配](https://leetcode.cn/problems/regular-expression-matching)

### 语法解析

和课本一样，不作解释。

**最好拿些题目练练手**

### 语义解析

你会发现，语义解析的实验和课本的实验有很大的差异。

#### 结果

最明显的就是生成的中间代码，我生成的是类汇编语言，而课本生成的是伪代码形式。

建议**尽量**参考课本。

#### 设计

设计上，我做了一个 trick，通过语法栈特征来完成继承属性往综合属性的传递。

<div style="text-align: center;">不要这么做！！！</div>

<sub title="细节ppt小字😋">好吧，你也可以这么做</sub>

对于如下包含继承属性的语法规则：

``` plaintext

S → A B {B.Action(nullptr)} C {S.Action(A, B, C)}

```

你应该将其改为：

``` plaintext
S → A B Placeholder C {S.Action(A, B, C)}
Placeholder → ε {B.Action(nullptr)}
```

这样做的好处是：
1. 你可以完全默认在规约操作发生时执行继承属性和综合属性的动作（这是符合 LR 分析过程的）
2. ✔这是课本的操作（虽然老师估计也不会细看）

LR 分析过程中，我还用 AST 来记录上下文信息，我还是比较**推荐**这么做的。

布尔表达式中，我并没有针对每个布尔表达式进行 jmp，而是对布尔表达式链整体进行 jmp：
``` cpp
if (1 == 2 && 3 == 4) {}

// 我的：
// t0 := equal(1, 2)
// t1 := equal(3, 4)
// t2 := and(t0, t1)
// jmz(t2, label1)
// jmp(label)
// 课本的：
// t0 := equal(1, 2)
// jmz(t0, label1)
// t1 := equal(3, 4)
// jmz(t1, label1)
// jmp(label)
```

我推荐参考课本的做法。

## 推荐复习

1. Regex → NFA → DFA → 最小化 DFA
2. 怎么将一个正则表达式描述成中文（没错，考你的语文）
3. LL(1) 
   1. 必须掌握构造 LL(1) 分析表之后，对字符流的**分析过程**
4. LR(0) → SLR(1) → LALR(1) → LR(1) 
   1. 必须掌握构造 LR 分析表之后，对字符流的**分析过程**
   2. 注意存在左递归的情况下，怎么对状态进行去重
   3. LR(0) 分析表与 SLR(1) 分析表的相互转换（我这一年没考）
   4. LR(1) 分析表与 LALR(1) 分析表的相互转换
    > <span style="font-size: large; font-weight:bold;">LR(0)、 SLR(1)、LALR(1)、LR(1) 的区别？</span>
    > - **LR(0)**：不考虑符号的后文，体现在分析表中的**规约项**总是**占满所在行**
    > - **SLR(1)**：考虑符号的后文，体现在分析表中的**规约项**所在行只占满属于 **FOLLOW集** 的部分
    > - **LR(1)**：考虑符号的后文，体现在分析表中的**规约项**所在行只占满属于 **向前搜索符集合** 的部分
    > - **LALR(1)**：考虑符号的后文，体现在分析表中的**规约项**所在行只占满属于 **向前搜索符集合** 的部分，且合并了 **LR(1)** 分析表中的同心状态
5. 注释分析树
   1. 特别注意箭头的方向（指向被依赖节点指向依赖节点）
    > 比如：
    > ``` plaintext
    >     S → A B C { func(S.in, A.x) }
    >     A → a { A.x = a.integer }
    > ```
    >
    > 注意，我们将其改为
    > ``` plaintext
    >     S → A B C { · := func(S.in, A.x) }
    >     A → a { A.x = a.integer }
    > ```
    >
    > 同时将 `·` 画在 `S` 的附近。
    > 
    > 如果不这么做，你会发现 `func(S.in, A.x)` 中不存在依赖关系（你没有办法确认 `func` 是否会将 `S.in` 和 `A.x` **关联**起来！），没有办法在注释分析树中体现出来。
    2. 还需要注意一下这些：
        1. 数组的注释分析树及其表达式
           - `array(5, array(2, int))`
        2. 记录（结构体）的注释分析树及其表达式
           - `record((a × int) × (b × int))`
6. DAG 图的绘制
7. 中间代码生成（语句翻译）