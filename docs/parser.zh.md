# 语法分析器

## 实验目的

根据给出的文法编制LR（1）分析程序，以便对任意输入的符号串进行分析。本次实验的目的主要是加深对LR（1）分析法的理解。

## 实验任务
实现编译器的以下功能： 
1. 输出语法分析表； 
2. 输出分析表分析栈内容；

## 文法
```
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

注：这个文法存在问题，实验中可能需要对其进行修改。

## 程序整体结构

![program structure](/docs/img/parser/1.png)

## LR(1)

LR(1) 主要分为以下几个部分：
1. 项集族的构造
   1. 闭包运算
   2. GOTO运算
   3. 项集族的构造（状态转换图）
2. 分析表的构造
   1. ACTION表
   2. GOTO表
3. 分析过程
   1. 分析栈的构造
   2. 分析过程的实现

### 1. 项集族的构造

#### 1.1 闭包运算

闭包运算是将一个项目集扩展为一个闭包的过程。其基本思想是：对于项目集中每个项目，如果右部有非终结符且该非终结符的产生式存在于文法中，则将这些产生式的所有项目加入闭包。

闭包运算步骤：
1. 初始化一个空的项目集，将初始项目加入其中。
2. 对项目集中的每个项目，检查右部是否有非终结符。如果有，将该非终结符的所有产生式的项目加入项目集中。
3. 重复步骤 2，直到项目集不再变化。
4. 返回闭包后的项目集。

伪代码：
```plaintext
function CLOSURE(I):
    J = I.copy()
    repeat:
        for each item [A → α • Bβ, a] in J:
            for each production B → γ in G:
                for each b in FIRST(βa):
                    if [B → •γ, b] not in J:
                        add [B → •γ, b] to J
    until J does not change
    return J
```

具体实现：

1. 初始化：
   - 确保 FIRST 集已计算，因为闭包运算依赖它。
   - 初始化闭包为输入项目集的副本，并使用 `marks` 记录已处理的项目。
   ```go
   p.EnsureFirstSet()
   closure := make([]LR1Item, len(items))
   copy(closure, items)
   marks := Set[string]{}
   ```

2. 主循环：
   - 不断扩展闭包，直到闭包不再变化。
   ```go
   loop := true
   for loop {
       loop = false
       ...
   }
   ```

3. 遍历闭包中的每个项目：
   - 跳过已处理的项目，并标记当前项目为已处理。
   ```go
   for _, item := range closure {
       if marks.Contains(item.AsKey()) {
           continue
       }
       marks.Add(item.AsKey())
       ...
   }
   ```

4. 检查项目的点位置：
   - 如果点超出产生式右部长度，跳过该项目。
   ```go
   if item.Dot >= len(item.Production.Body) {
       continue
   }
   ```

5. 获取点后面的符号：
   - 如果点后面的符号是终结符，则无法扩展，跳过。
   ```go
   nextSymbol := item.Production.Body[item.Dot]
   if p.Grammar.IsTerminal(nextSymbol) {
       continue
   }
   ```

6. 遍历文法中的产生式：
   - 找到以 `nextSymbol` 为头部的产生式。
   ```go
   for _, production := range p.Grammar.Productions {
       if production.Head == nextSymbol {
           ...
       }
   }
   ```

7. 处理空产生式：
   - 如果产生式右部为空（ε），创建新项目并加入闭包。
   ```go
   if production.Body[0].IsEpsilon() {
       newItem := LR1Item{
           Production: production,
           Dot:        0,
           Lookahead:  item.Lookahead,
       }
       if !slices.ContainsFunc(closure, func(i LR1Item) bool {
           return i.Equals(newItem)
       }) {
           closure = append(closure, newItem)
           loop = true
       }
   }
   ```

8. 计算 Lookahead 集：
   - 调用 `p.findLookaheads` 计算 Lookahead 集，并将新项目加入闭包。
   ```go
   lookaheads := p.findLookaheads(item.Production.Body[item.Dot+1:], item.Lookahead)
   for lookahead := range lookaheads {
       newItem := LR1Item{
           Production: production,
           Dot:        0,
           Lookahead:  lookahead,
       }
       if !slices.ContainsFunc(closure, func(i LR1Item) bool {
           return i.Equals(newItem)
       }) {
           closure = append(closure, newItem)
           loop = true
       }
   }
   ```

9. 返回闭包：
   - 当闭包不再变化时，返回最终的闭包项目集。

##### findLookaheads 函数的实现

findLookaheads 函数用于计算给定符号串的 FIRST 集，并根据文法规则处理 ε（空串）。该函数的输入参数包括符号串和当前的 lookahead 符号。
伪代码如下：
```plaintext
function findLookaheads(symbols, lookahead):
    if symbols is empty:
        return {lookahead}
    result = {}
    for each symbol in symbols:
        if symbol is terminal:
            add symbol to result
        else:
            add FIRST(symbol) - {ε} to result
        if ε not in FIRST(symbol):
            break
    if ε in FIRST(all symbols):
        add lookahead to result
    return result
```

1. 空符号串处理：
   - 如果 `symbols` 为空，直接返回包含 `lookahead` 的集合。
2. 初始化：
   - 创建一个空的 `firstSet` 集合，用于存储符号串的 FIRST 集。
   - 使用 `flag` 标记是否需要将 `lookahead` 加入结果。
3. 遍历符号串：
   - 如果当前符号是终结符，直接加入 `firstSet`。
   - 如果是非终结符，将其 FIRST 集（去掉 ε）加入 `firstSet`。
   - 如果当前符号的 FIRST 集不包含 ε，则停止遍历。
4. 处理 ε：
   - 如果符号串的所有符号的 FIRST 集都包含 ε，则将 `lookahead` 加入结果。
5. 返回结果：
   - 返回最终计算的 `firstSet`。
  
以下是 `findLookaheads` 函数的实现代码：

```go
func (p *Parser) findLookaheads(symbols []Symbol, lookahead Terminal) Set[Terminal] {
    if len(symbols) == 0 {
        s := Set[Terminal]{}
        s.Add(lookahead)
        return s
    }

    flag := true
    firstSet := Set[Terminal]{}
    for _, symbol := range symbols {
        if p.Grammar.IsTerminal(symbol) {
            firstSet.Add(Terminal(symbol))
        }

        for terminal := range p.FirstSet[symbol] {
            if !terminal.IsEpsilon() {
                firstSet.Add(terminal)
            }
        }

        if !firstSet.Contains(EPSILON) {
            flag = false
            break
        }
    }

    if flag {
        firstSet.Add(lookahead)
    }

    return firstSet
}
```

##### FIRST 集的计算

FIRST 集计算规则：
1. 对于终结符，其 FIRST 集是其自身。
2. 对于非终结符，其 FIRST 集是所有产生式右部的 FIRST 集的并集。
3. 如果产生式右部包含 ε，则将 ε 加入 FIRST 集。
4. 如果产生式右部以非终结符开头，将该非终结符的 FIRST 集加入当前 FIRST 集。
5. 如果产生式右部以终结符开头，将该终结符加入当前 FIRST 集。
6. 如果产生式右部为空，将 ε 加入当前 FIRST 集。
7. 重复上述步骤，直到 FIRST 集不再变化。

伪代码：
```plaintext
function FIRST(X):
    if X is a terminal:
        return {X}
    if X is a non-terminal:
        result = {}
        for each production X → Y1 Y2 ... Yk:
            for i = 1 to k:
                temp = FIRST(Yi)
                result = result ∪ (temp - {ε})
                if ε not in temp:
                    break
            if ε in temp for all i:
                result = result ∪ {ε}
        return result
    if X is ε:
        return {ε}
```


> 在计算 FIRST 集时，很可能会出现相互依赖（递归调用）的情况：
> - 非终结符的 FIRST 集依赖于其他非终结符的 FIRST 集
>   - 产生式 A → B C，其中 B 的 FIRST 集依赖于 A 的 FIRST 集。
>   - 产生式 A → A b，其中 A 的 FIRST 集直接依赖于自身。
> - 如果直接递归调用，可能会导致无限递归或栈溢出。
> 
> 我们可以使用**迭代更新**的方式，逐步逼近最终的 FIRST 集，避免了直接递归调用。具体来说：
> - 每次迭代中，尝试更新所有非终结符的 FIRST 集。
>   - 尝试将右部符号的 FIRST 集加入左部非终结符的 FIRST 集。
>   - 如果右部符号是非终结符，则将其 FIRST 集（去掉 ε）加入左部非终结符的 FIRST 集。
>   - 如果右部符号的 FIRST 集包含 ε，则继续处理下一个符号，直到右部符号处理完毕或遇到终结符。
> - 如果某个非终结符的 FIRST 集发生了变化，则继续下一轮迭代。
> - 当所有 FIRST 集都不再变化时，算法终止。
>
> 示例：
> 
> 假设有以下文法：
> ```plaintext
> A → B C
> B → ε | a
> C → c | A
> ```
> 初始状态：
> ```plaintext
> FIRST(B) = {ε, a}
> FIRST(C) = {}
> FIRST(A) = {}
> ```
> 迭代过程：
> 1. 第一次迭代：
>   - 对于 A → B C：
>       - 将 FIRST(B) 加入 FIRST(A)，得到 FIRST(A) = {ε, a}。
>       - 因为 FIRST(B) 包含 ε，继续处理 C。
>       - 将 FIRST(C) 加入 FIRST(A)，但 FIRST(C) 为空，因此 FIRST(A) 不变。
>   - 对于 C → c | A：
>       - 将 c 加入 FIRST(C)，得到 FIRST(C) = {c}。
>       - 将 FIRST(A) 加入 FIRST(C)，得到 FIRST(C) = {c, ε, a}。
> 2. 第二次迭代
>   - 对于 A → B C：
>       - 将 FIRST(B) 加入 FIRST(A)，FIRST(A) 不变。
>       - 因为 FIRST(B) 包含 ε，继续处理 C。
>       - 将 FIRST(C) 加入 FIRST(A)，得到 FIRST(A) = {ε, a, c}。
>   - 对于 C → c | A：
>       - FIRST(C) 不变。
> 3. 所有 FIRST 集都不再变化，终止迭代。


具体实现：

1. 初始化：
   - 终结符的 FIRST 集为其自身。
   ```go
   for terminal := range p.Grammar.Terminals {
       p.FirstSet[Symbol(terminal)] = Set[Terminal]{}
       p.FirstSet[Symbol(terminal)].Add(terminal)
   }
   ```
   - 非终结符的 FIRST 集初始化为空。
   ```go
   for _, production := range p.Grammar.Productions {
       if _, exists := p.FirstSet[production.Head]; !exists {
           p.FirstSet[production.Head] = Set[Terminal]{}
       }
   }
   ```

2. 迭代计算：
   - 不断更新 FIRST 集，直到所有 FIRST 集稳定。
   ```go
   loop := true
   for loop {
       loop = false
       ...
   }
   ```

3. 处理每个产生式：
   - 空产生式： 如果右部为空或第一个符号是 ε，将 ε 加入 FIRST 集。
     ```go
     if len(production.Body) == 0 || production.Body[0].IsEpsilon() {
         if !firstSet.Contains(EPSILON) {
             firstSet.Add(EPSILON)
             loop = true
         }
     }
     ```
   - 遍历右部符号：
     - 如果符号是 ε，将 ε 加入 FIRST 集并停止。
     - 如果符号是非终结符，将其 FIRST 集（去掉 ε）加入当前 FIRST 集。
     - 如果符号是终结符，直接加入当前 FIRST 集并停止。
     ```go
     if symbol.IsEpsilon() {
         ...
     } else if symbolFirstSet, isNonTerminal := p.FirstSet[symbol]; isNonTerminal {
         ...
     } else {
         ...
     }
     ```

4. 终止条件：
   - 当所有 FIRST 集不再变化时，计算完成。

#### 1.2 GOTO 运算

GOTO 运算是将一个项目集转换为另一个项目集的过程。其基本思想是：对于给定的项目集和一个符号，计算出该符号对应的项目集。

GOTO 运算步骤：
1. 对于给定的项目集和符号，初始化一个空的项目集。
2. 遍历项目集中的每个项目，检查点后面的符号是否与给定符号相同。
3. 如果相同，将该项目的点向右移动一位，产生新项目集。
4. 对新项目集进行闭包运算，返回最终的项目集。

伪代码如下：
```plaintext
function GOTO(I, X):
    J = {}
    for each item [A → α • Xβ, a] in I:
        if X is the symbol after the dot in α • Xβ:
            add [A → αX • β, a] to J
    return CLOSURE(J)
```

具体实现如下：
```go
func (p *Parser) GOTO(items LR1Items, symbol Symbol) LR1Items {
	gotoItems := LR1Items{}
	for _, item := range items {
		if item.Dot < len(item.Production.Body) && item.Production.Body[item.Dot] == symbol {
			newItem := LR1Item{
				Production: item.Production,
				Dot:        item.Dot + 1,
				Lookahead:  item.Lookahead,
			}
			gotoItems = append(gotoItems, newItem)
		}
	}
	return p.CLOSURE(gotoItems)
}
```
   
#### 1.3 项集族的构造（状态转换图）

构造项集族的过程是将所有可能的项目集进行闭包运算和 GOTO 运算，直到没有新的项目集产生为止。

构造项集族的步骤：
1. 初始化一个空的项集族和一个初始项目集（一般是文法的开始符号的产生式）
2. 对初始项目集进行闭包运算，得到初始状态。
3. 将初始状态加入项集族。
4. 初始化一个空的待处理状态队列，将初始状态加入队列。
5. 循环处理队列中的状态：
   1. 从队列中取出一个状态。
   2. 对该状态中的每个符号进行 GOTO 运算，得到新的状态。
   3. 如果新状态不在项集族中，将其加入项集族，并将其加入队列。
   4. 如果新状态已经在项集族中，记录该状态的转换关系。
   5. 重复步骤 5，直到队列为空。
6. 记录项集族和状态转换关系。

伪代码如下：
```plaintext
function CONSTRUCT_ITEM_SET_FAMILY():
    I0 = CLOSURE([S' → •S, $])
    itemSetFamily = {I0}
    queue = [I0]
    while queue is not empty:
        I = queue.pop()
        for each symbol in G:
            J = GOTO(I, symbol)
            if J is not in itemSetFamily:
                itemSetFamily.add(J)
                queue.push(J)
            add (I, symbol) → J to transition table
    return itemSetFamily, transition table
```

具体实现：

1. 初始化符号集  
    - 调用 EnsureSymbols 方法，确保文法中的所有符号（终结符和非终结符）都已被添加到符号集中。  
    - 这些符号会在后续的 GOTO 运算中使用。  
    ```go
    p.EnsureSymbols()
    ```

2. 创建初始状态  
    - 创建一个初始项目 initialItem，表示文法的扩展产生式（S' → •S, $），其中：  
      - Dot 表示点的位置，初始为 0。  
      - Lookahead 是终止符 $，表示输入结束。  
    ```go
    initialItem := LR1Item{
         Production: p.Grammar.AugmentedProduction,
         Dot:        0,
         Lookahead:  TERMINATE,
    }
    ```
    - 创建初始状态 initialState，包含：  
      - Index：状态编号，初始为 0。  
      - Items：项目集，初始只包含 initialItem。  
      - Transitions：状态的转换关系，初始为空。  
    ```go
    initialState := &State{
         Index:       0,
         Items:       LR1Items{initialItem},
         Transitions: make(map[Symbol]*State),
    }
    ```
    - 对初始状态的项目集执行闭包运算，扩展项目集，生成完整的初始状态。  
    ```go
    initialState.Items = p.CLOSURE(initialState.Items)
    ```

3. 初始化状态集合  
    - 将初始状态加入状态集合 p.States，这是所有状态的起点。  
    ```go
    p.States = States{initialState}
    ```

4. 构造项集族  
    - 使用一个循环，逐一处理状态集合中的每个状态。  
    - length 记录当前状态集合的长度，随着新状态的加入会动态更新（模拟队列）。  
    ```go
    length := len(p.States)
    for i := 0; i < length; i++ {
         state := p.States[i]
         ...
    }
    ```
    1.  遍历符号  
        - 遍历所有符号（终结符和非终结符），对每个符号执行 GOTO 运算。  
        - 如果 GOTO 运算结果为空（即没有新的项目集），跳过该符号。  
        ```go
        for symbol := range p.Symbols {
             gotoItems := p.GOTO(state.Items, symbol)
             if len(gotoItems) == 0 {
                  continue
             }
             ...
        }
        ```

    2. 创建新状态  
       - 如果 GOTO 运算结果非空，创建一个新状态 newState，包含：  
         - Index：新状态的编号，等于当前状态集合的长度。  
         - Items：GOTO 运算生成的项目集。  
         - Transitions：初始化为空的转换关系。  
       ```go
       newState := &State{
            Index:       len(p.States),
            Items:       gotoItems,
            Transitions: make(map[Symbol]*State),
       }
       ```

    3.  检查状态是否已存在  
        - 检查新状态是否已经存在于状态集合中：  
          - 使用 slices.IndexFunc 遍历状态集合，比较每个状态的项目集。  
          - 如果找到相同的状态，返回其索引；否则返回 -1。 
    ```go
    index := slices.IndexFunc(p.States, func(s *State) bool {
         return s.Equals(newState)
    })
    ```

    4. 添加新状态或更新转换关系  
       - 如果新状态不存在：  
         - 将其加入状态集合 p.States。  
         - 更新当前状态的转换关系 state.Transitions，记录从当前状态通过 symbol 转换到新状态。  
         - 增加状态集合的长度 length，以便继续处理新状态。  
       - 如果新状态已存在：  
         - 直接更新当前状态的转换关系，指向已存在的状态。  
     
       ```go
       if index == -1 {
            p.States = append(p.States, newState)
            state.Transitions[symbol] = newState
            length++
       } else {
            state.Transitions[symbol] = p.States[index]
       }
       ```

5. 循环结束  
    当所有状态都被处理完毕（即队列为空），循环结束，构造完成。

### 2. LR(1)分析表
#### 2.1 ACTION 表
ACTION 表用于记录在每个状态下对于每个终结符的操作。其基本思想是：对于每个状态和终结符，确定是移进、规约还是接受。

- `Action` 结构体表示一个动作，包含动作类型和编号。
- `ActionTable` 是一个嵌套的 map，外层 map 的键是状态编号，内层 map 的键是终结符，值是 `Action` 结构体。
- `Copy` 方法用于复制 ACTION 表。
- `Register` 方法用于注册一个动作到 ACTION 表中，检查是否存在冲突。
  - 如果状态和终结符的组合已经存在，检查动作类型是否冲突。
  - 如果冲突，返回错误。
  - 如果没有冲突，将动作注册到表中。

```go
type Action struct {
	Type   ActionType
	Number int
}

type ActionTable map[int]map[Terminal]Action

func (t ActionTable) Copy() ActionTable {
	return maps.Clone(t)
}

func (t ActionTable) Register(stateIndex int, action Action, terminal Terminal) error {
	if t[stateIndex] == nil {
		t[stateIndex] = make(map[Terminal]Action)
	}

	if _, exists := t[stateIndex][terminal]; exists {
		if t[stateIndex][terminal].Type == SHIFT && action.Type == REDUCE {
			return fmt.Errorf("conflict in action table: state %d, terminal %s[shift] %d, [reduce] %d", stateIndex, terminal, t[stateIndex][terminal].Number, action.Number)
		} else if t[stateIndex][terminal].Type == REDUCE && action.Type == REDUCE {
			return fmt.Errorf("conflict in action table: state %d, terminal %s[reduce] %d, [reduce] %d", stateIndex, terminal, t[stateIndex][terminal].Number, action.Number)
		}
	}

	t[stateIndex][terminal] = action
	return nil
}
```

#### 2.2 GOTO 表
GOTO 表用于记录在每个状态下对于每个非终结符的转换关系。其基本思想是：对于每个状态和非终结符，确定下一个状态。

- `GotoTable` 是一个嵌套的 map，外层 map 的键是状态编号，内层 map 的键是非终结符，值是下一个状态的编号。
- `Copy` 方法用于复制 GOTO 表。
- `Register` 方法用于注册一个转换关系到 GOTO 表中，检查是否存在冲突。
  - 如果状态和非终结符的组合已经存在，返回错误。
  - 如果没有冲突，将转换关系注册到表中。

```go
type GotoTable map[int]map[Symbol]int

func (t GotoTable) Copy() GotoTable {
	return maps.Clone(t)
}

func (t GotoTable) Register(stateIndex, nextStateIndex int, symbol Symbol) error {
	if t[stateIndex] == nil {
		t[stateIndex] = make(map[Symbol]int)
	}

	// ignore conflict
	//if _, exists := t[stateIndex][symbol]; exists {
	//	return fmt.Errorf("conflict in goto table: state %d, symbol %s", stateIndex, symbol)
	//}

	t[stateIndex][symbol] = nextStateIndex
	return nil
}
```

### 2.3 LR(1)表

LR(1) 表是将 ACTION 表和 GOTO 表结合在一起的结果。其基本思想是：对于每个状态，记录对于每个终结符和非终结符的操作和转换关系。

- `LR1Table` 结构体包含 ACTION 表和 GOTO 表。
- `Insert` 方法用于将状态和文法插入到 LR(1) 表中。
  - 遍历状态的项目集，检查每个项目的点位置。
  - 如果点在产生式的末尾，检查是否接受或规约。
  - 如果点在产生式中间，检查下一个符号是终结符还是非终结符，并注册相应的动作或转换关系。
  - 如果出现冲突，返回错误。

```go
type LRTable struct {
	ActionTable ActionTable
	GotoTable   GotoTable
}

func (t LRTable) Insert(state *State, grammar *Grammar) {
	var err error
	for _, item := range state.Items {
		if item.Dot == len(item.Production.Body) || item.Production.Body[item.Dot].IsEpsilon() {
			if item.Lookahead == TERMINATE && item.Production.Equals(grammar.AugmentedProduction) {
				err = t.ActionTable.Register(state.Index, Action{Type: ACCEPT, Number: 0}, TERMINATE)
			} else {
				err = t.ActionTable.Register(state.Index, Action{Type: REDUCE, Number: grammar.GetIndex(item.Production)}, item.Lookahead)
			}
		} else {
			symbol := item.Production.Body[item.Dot]
			if symbol.IsEpsilon() {
				continue
			}
			if grammar.IsNonTerminal(symbol) {
				err = t.GotoTable.Register(state.Index, state.Transitions[symbol].Index, symbol)
			} else {
				err = t.ActionTable.Register(state.Index, Action{Type: SHIFT, Number: state.Transitions[symbol].Index}, Terminal(symbol))
			}
		}
		if err != nil {
			//fmt.Printf("when inserting : %v\n", err)
		}
	}
}
```

因此，LR(1) 表的构造过程只需要将第一步建立的状态集合和转换关系传入 LRTable 的 Insert 方法即可：
```go
func (p *Parser) BuildTable() {
	p.EnsureStates()

	p.Table = &LRTable{
		ActionTable: make(ActionTable),
		GotoTable:   make(GotoTable),
	}

	for _, state := range p.States {
		p.Table.Insert(state, p.Grammar)
	}
}
```

### 3. 分析过程

程序主要使用一个 `Walker` 结构体来实现分析过程。

```go
type Walker struct {
	Table   LRTable
	Grammar *Grammar

	States  Stack[int]
	Symbols Stack[Symbol]
	Tokens  Stack[*lexer.Token]

	SymbolTable *SymbolTable
}
```

该结构体包含了分析表、文法、状态栈、符号栈、token栈和符号表等信息。
其中，状态栈用于存储当前状态，符号栈用于存储当前符号串，token栈用于存储当前输入的 token，符号表用于存储变量和函数等信息。

注：token 栈实际上不参与分析过程，只是用于存储当前输入的 token，方便后续的语义分析和中间代码生成。

使用 `Parser` 结构体完成对文法的解析和分析表的构造，同时将分析结果保存或拷贝到新的 `Walker` 结构体中。

```go
type Parser struct {
	Grammar *Grammar
	Symbols Set[Symbol]

	FirstSet FirstSet

	States States

	Table *LRTable
}
```

该结构体包含了文法、符号集、FIRST 集、状态集合和分析表等信息。

提供了 `NewWalker` 方法用于创建一个新的 `Walker` 结构体，并将分析表和文法传入（以达成并发分析多个文件的目的）。

```go
func (p *Parser) NewWalker() *Walker {
	p.EnsureTable()

	g := p.Grammar.Copy()
	states := Stack[int]{}
	states.Push(0)
	symbols := Stack[Symbol]{}
	return &Walker{
		Table: LRTable{
			ActionTable: p.Table.ActionTable.Copy(),
			GotoTable:   p.Table.GotoTable.Copy(),
		},
		Grammar:     &g,
		States:      states,
		Symbols:     symbols,
		SymbolTable: NewSymbolTable(nil, nil),
	}
}
```

#### 3.1 分析栈的构造
分析栈用于存储当前状态和符号串。其基本思想是：对于每个输入的 token，检查当前状态和符号栈的状态，决定是移进、规约还是接受。

其体现为`Walker` 结构体中的 `States` 和 `Symbols` 字段。
- `States` 用于存储当前状态。
- `Symbols` 用于存储当前符号串。

#### 3.2 分析过程的实现
分析过程的实现主要分为以下几个步骤：
1. 初始化分析栈和输入 token。
2. 循环处理输入 token，直到接受或出错。
3. 根据当前状态和输入 token，决定是移进、规约还是接受。
4. 执行相应的操作，更新分析栈和输入 token。
5. 如果出现错误，输出错误信息并退出。
6. 如果接受，输出分析结果。
7. 如果规约，更新符号表和中间代码生成器。
8. 返回分析结果。
9. 结束分析过程。

程序将分析过程拆分为两个部分，一个部分用于token的传输（Parser），另一个部分用于语法分析和中间代码生成（Walker）。

##### Parser

`Parser` 结构体主要用于处理输入的 token，对外提供 `Parse` 方法。

该方法接收一个 `lexer.Lexer` 对象和一个日志函数作为参数，循环读取输入的 token，并将其传递给 `Walker` 进行分析。

主要逻辑如下：
1. 初始化
   - 创建一个新的 Walker 实例。
   - 调用 walker.SymbolTable.EnterScope() 进入一个新的作用域。
2. 循环读取 Token
   - 使用 l.NextToken() 从 lexer 中获取下一个 token。
   - 如果遇到错误且不是文件结束 (io.EOF)，则终止解析。
   - 如果是文件结束 (io.EOF)，将 token 的类型设置为 lexer.EOF。
3. Token 映射到 Symbol：
   - 调用 p.Reflect(token) 将 token 转换为对应的 symbol。
4. 作用域管理：
   - 如果 token 是左大括号，调用 walker.SymbolTable.EnterScope() 进入新的作用域。
    - 如果 token 是右大括号 ，调用 walker.SymbolTable.ExitScope() 退出当前作用域。
5. Walker 状态更新：
   - 调用 walker.Next(symbol) 根据当前 symbol 更新 Walker 的状态。
   - 如果 action.Type 是 REDUCE，继续处理当前 symbol，否则跳出内部循环（规约操作不可能一次完成，因为文法可能会不断嵌套或者递归，状态可能会不断变化，直到无法规约为止，因此需要不断尝试规约）。
6. 终止条件：
   - 如果 symbol 是 TERMINATE，表示解析完成，退出循环。

注：token 的压栈必须是在 Walker 完成一系列规约操作后进行的，为简化token 的压栈逻辑，Parser 在规约完成（分析表接受语法并完成状态转移）直接将 token 压栈到 Walker 中。

具体实现如下：
```go
func (p *Parser) Parse(l *lexer.Lexer, logger func(string)) {
	walker := p.NewWalker()
	walker.SymbolTable.EnterScope()
	for {
		token, err := l.NextToken()
		if err != nil && !errors.Is(err, io.EOF) {
			logger(fmt.Sprintf("Error: %v", err))
			return
		}

		if errors.Is(err, io.EOF) {
			token.Type = lexer.EOF
		}
		symbol := p.Reflect(token)
		if token.SpecificType() == lexer.DelimiterLeftBrace {
			walker.SymbolTable.EnterScope()
		}

		for {
			logger(fmt.Sprintf("State: %v\nSymbols: %v\nSymbol: %s\n", walker.States, walker.Symbols, symbol))
			action, err := walker.Next(symbol)
			if err != nil {
				logger(fmt.Sprintf("Error: %v", err))
				return
			}
			logger(fmt.Sprintf("Token: (%s, %s), Action: %v\n\n", token.Type.ToString(), token.Val, action))
			if action.Type != REDUCE {
				break
			}
		}

		if token.SpecificType() == lexer.DelimiterRightBrace {
			walker.SymbolTable.ExitScope()
		}

		if symbol == TERMINATE {
			logger("Parsing completed successfully.")
			break
		}

		walker.Tokens.Push(&token)
	}
}
```

##### Walker

`Walker` 结构体主要用于实现 LR(1) 分析过程，对外提供 `Next` 方法。

该方法接收一个 `Symbol` 对象作为参数，表示当前输入的符号。根据当前状态和输入符号，决定是移进、规约还是接受，并执行相应的操作。

仅接受 Symbol 作为参数，表示当前输入的符号。分析过程能够更加贴近语法分析的实际操作，而不会对输入的 token 产生过多的耦合。

主要逻辑如下：
1. 获取当前状态：
   - 使用 `w.States.Peek()` 获取当前状态的索引。
2. 判断输入符号类型：
   - 如果是终结符，查找 ACTION 表，获取对应的动作。
     - 如果是移进（SHIFT），将状态和符号压栈。
     - 如果是规约（REDUCE），调用对应的产生式处理规则。
       - 处理完成后，弹出状态和符号。
       - 查找 GOTO 表，获取下一个状态，将产生式头部压栈。
     - 如果是接受（ACCEPT），返回接受状态。
   - 如果是非终结符，查找 GOTO 表，获取下一个状态。
     - 将状态和符号压栈。 
   - 如果都没有，返回错误。

```go
func (w *Walker) Next(symbol Symbol) (action Action, err error) {
	topState, _ := w.States.Peek()
	if w.Grammar.IsTerminal(symbol) {
		action, ok := w.Table.ActionTable[topState][Terminal(symbol)]
		if !ok {
			return Action{Type: ERROR}, fmt.Errorf("no action found for state %d and symbol %s", topState, symbol)
		}
		switch action.Type {
		case SHIFT:
			w.States.Push(action.Number)
			w.Symbols.Push(symbol)
			return Action{Type: SHIFT, Number: action.Number}, nil
		case REDUCE:
			production := w.Grammar.Productions[action.Number]
			if err := production.HandleRule(w); err != nil {
				fmt.Println("Error handling rule:", err)
			}
			for i := range production.Body {
				if production.Body[i] == EPSILON {
					continue
				}
				w.States.Pop()
				w.Symbols.Pop()
			}
			topState, _ = w.States.Peek()
			gotoState, ok := w.Table.GotoTable[topState][production.Head]
			if !ok {
				return Action{Type: ERROR}, fmt.Errorf("no goto state found for state %d and symbol %s", topState, production.Head)
			}
			w.Symbols.Push(production.Head)
			w.States.Push(gotoState)
			return Action{Type: REDUCE, Number: action.Number}, nil
		case ACCEPT:
			return Action{Type: ACCEPT, Number: 0}, nil
		}
	} else {
		action, ok := w.Table.GotoTable[topState][symbol]
		if !ok {
			return Action{Type: ERROR}, fmt.Errorf("no goto state found for state %d and symbol %s", topState, symbol)
		}
		w.States.Push(action)
		w.Symbols.Push(symbol)
		return Action{Type: GOTO, Number: action}, nil
	}
	return Action{Type: ERROR}, fmt.Errorf("unexpected state %d and symbol %s", topState, symbol)
}
```

## 测试
### First 集测试
在 [algorithm_test.go](/parser/algorithm_test.go) 中的`TestParser_BuildFirstSet` 测试函数中，使用了一批简单的文法来测试 FIRST 集的构建。

#### 测试用例
<table><tb>
<tr><th style="text-align:center;">文法</th><th style="text-align:center;">FIRST 集</th><th style="text-align:center;">终结符</th></tr>
<tr><td valign="top">
    E -> T E' | + T E' | ε <br>
    T -> F T' | * F T' | ε <br>
    F -> ( E ) | id <br>
</td>
<td valign="top">
    FIRST(E) = { (, id, +, ε }<br>
    FIRST(E') = { +, ε }<br>
    FIRST(T) = { (, id }<br>
    FIRST(T') = { *, ε }<br>
    FIRST(F) = { (, id }<br>
</td>
<td valign="top">
id, +, *, (, ), ε
</td></tr>
<tr><td valign="top">
    S -> a A B b c d | ε<br>
    A -> A S d | ε<br>
    B -> S A h | e C | ε<br>
    C -> S f | C g | ε<br>
</td>
<td valign="top">
    FIRST(S) = { a, ε }<br>
    FIRST(A) = { a, d, ε }<br>
    FIRST(B) = { a, d, h, e, ε }<br>
    FIRST(C) = { a, f, g, ε }<br>
</td>
<td valign="top">
a, b, c, d, e, f, g, h, ε
</td></tr>
<tr><td valign="top">
    E -> T E' <br>
    E' -> + E | ε <br>
    T -> F T' <br>
    T' -> T | ε <br>
    F -> P F' <br>
    F' -> * F' | ε <br>
    P -> ( E ) | a | b | ^ <br>
</td>
<td valign="top">
    FIRST(E) = { (, a, b, ^ }<br>
    FIRST(E') = { +, ε }<br>
    FIRST(T) = { (, a, b, ^ }<br>
    FIRST(T') = { (, a, b, ^, ε }<br>
    FIRST(F) = { (, a, b, ^ }<br>
    FIRST(F') = { *, ε }<br>
    FIRST(P) = { (, a, b, ^ }<br>
</td>
<td valign="top">
+, (, ), a, b, ^, *, ε
</td></tr>
<tr><td valign="top">
    D -> B c | c <br>
    B -> b D | a <br>
</td>
<td valign="top">
    FIRST(B) = { a, b }<br>
    FIRST(D) = { a, b, c }<br>
</td>
<td valign="top">
a, b, c
</td></tr>
<tr><td valign="top">
    E -> T A <br>
    A -> + T A | ε <br>
    T -> F B <br>
    B -> * F B | ε <br>
    F -> i | ( E ) <br>
</td>
<td valign="top">
    FIRST(E) = { i, ( }<br>
    FIRST(A) = { +, ε }<br>
    FIRST(T) = { i, ( }<br>
    FIRST(B) = { *, ε }<br>
    FIRST(F) = { i, ( }<br>
</td>
<td valign="top">
+, *, i, (, ), ε
</td></tr>
</tb></table>
    

### 项集族测试（状态转换图）
在 [algorithm_test.go](/parser/algorithm_test.go) 中的`TestParser_BuildStates` 测试函数中，使用了一批简单的文法来测试项集族的构建。

#### 测试用例

<table>
<tr><th style="text-align:center;">增广文法</th><th style="text-align:center;">文法</th><th style="text-align:center;">终结符</th></tr>
<tr><td valign="top">
    S' -> S <br>
</td>
<td valign="top">
    S -> A a A b <br>
    S -> B b B a <br>
    A -> ε <br>
    B -> ε <br>
</td>
<td valign="top">
a, b, ε, $
</td></tr>
<tr><td valign="top">
    S' -> S <br>
</td>
<td valign="top">
    S -> L = R <br>
    S -> R <br>
    L -> * R <br>
    L -> id <br>
    R -> L <br>
</td>
<td valign="top">
*, =, id, ε, $
</td></tr>
<tr><td valign="top">
    S' -> S <br>
</td>
<td valign="top">
    S -> B B <br>
    B -> a B <br>
    B -> b <br>
</td>
<td valign="top">
a, b, ε, $
</td></tr>
</table>

### 分析表构建测试
在 [table_test.go](/parser/table_test.go) 中的`TestParser_BuildTable` 测试函数中，使用了一批简单的文法来测试分析表的构建。

<table>
<tr><th style="text-align:center;">增广文法</th><th style="text-align:center;">文法</th><th style="text-align:center;">终结符</th></tr>
<tr><td valign="top">
    S' -> S <br>
</td>
<td valign="top">
    S -> L = R <br>
    S -> R <br>
    L -> * R <br>
    L -> id <br>
    R -> L <br>
</td>
<td valign="top">
*, =, id, ε, $
</td></tr>
<tr><td valign="top">
    S' -> S <br>
</td>
<td valign="top">
    S -> V = E <br>
    S -> E <br>
    E -> V <br>
    V -> x <br>
    V -> * E <br>
</td>
<td valign="top">
=, x, *, ε, $
</td></tr>
<td valign="top">
    program' -> program <br>
</td>
<td valign="top">
    本实验文法
</td>
<td valign="top">
    本实验文法的终结符
</td></tr>
</table>

### 语法分析测试
在 [walker_test.go](/parser/walker_test.go) 中的`TestWalker_Next`、`TestWalker_Next2`、`TestWalker_Next3` 测试函数中，使用了一批简单的文法和 token 序列来测试语法分析的正确性。

#### 测试用例1

**文法：**
<table>
<tr><th style="text-align:center;">增广文法</th><th style="text-align:center;">文法</th><th style="text-align:center;">终结符</th></tr>
<tr><td valign="top">
    E' -> E <br>
</td>
<td valign="top">
    E -> E + T | T <br>
    T -> T * F | F <br>
    F -> ( E ) | id <br>
</td>
<td valign="top">
(, ), +, *, id, ε, $
</td></tr>
</table>

**LR(1) 分析表：**
<table><thead><tr><th rowspan="2" style="text-align:center;">状态</th><th colspan="6" style="text-align:center;">ACTION</th><th colspan="3" style="text-align:center;">GOTO</th></tr><tr><th style="text-align:center;">id</th><th style="text-align:center;">(</th><th style="text-align:center;">+</th><th style="text-align:center;">*</th><th style="text-align:center;">)</th><th style="text-align:center;">$</th><th style="text-align:center;">E</th><th style="text-align:center;">T</th><th style="text-align:center;">F</th></tr></thead><tbody><tr><td style="text-align:center;">0</td><td style="text-align:center;">S<sub>5</sub></td><td style="text-align:center;">S<sub>4</sub></td><td></td><td></td><td></td><td></td><td style="text-align:center;">1</td><td style="text-align:center;">2</td><td style="text-align:center;">3</td></tr><tr><td style="text-align:center;">1</td><td></td><td></td><td style="text-align:center;">S<sub>6</sub></td><td></td><td></td><td style="text-align:center;">Accept</td><td></td><td></td><td></td></tr><tr><td style="text-align:center;">2</td><td></td><td></td><td style="text-align:center;">R<sub>2</sub></td><td style="text-align:center;">S<sub>7</sub></td><td style="text-align:center;">R<sub>2</sub></td><td style="text-align:center;">R<sub>2</sub></td><td></td><td></td><td></td></tr><tr><td style="text-align:center;">3</td><td></td><td></td><td style="text-align:center;">R<sub>4</sub></td><td style="text-align:center;">R<sub>4</sub></td><td style="text-align:center;">R<sub>4</sub></td><td style="text-align:center;">R<sub>4</sub></td><td></td><td></td><td></td></tr><tr><td style="text-align:center;">4</td><td style="text-align:center;">S<sub>5</sub></td><td style="text-align:center;">S<sub>4</sub></td><td></td><td></td><td></td><td></td><td style="text-align:center;">8</td><td style="text-align:center;">2</td><td style="text-align:center;">3</td></tr><tr><td style="text-align:center;">5</td><td></td><td></td><td style="text-align:center;">R<sub>6</sub></td><td style="text-align:center;">R<sub>6</sub></td><td style="text-align:center;">R<sub>6</sub></td><td style="text-align:center;">R<sub>6</sub></td><td></td><td></td><td></td></tr><tr><td style="text-align:center;">6</td><td style="text-align:center;">S<sub>5</sub></td><td style="text-align:center;">S<sub>4</sub></td><td></td><td></td><td></td><td></td><td></td><td style="text-align:center;">9</td><td style="text-align:center;">3</td></tr><tr><td style="text-align:center;">7</td><td style="text-align:center;">S<sub>5</sub></td><td style="text-align:center;">S<sub>4</sub></td><td></td><td></td><td></td><td></td><td></td><td></td><td style="text-align:center;">10</td></tr><tr><td style="text-align:center;">8</td><td></td><td></td><td style="text-align:center;">S<sub>6</sub></td><td></td><td style="text-align:center;">S<sub>11</sub></td><td></td><td></td><td></td><td></td></tr><tr><td style="text-align:center;">9</td><td></td><td></td><td style="text-align:center;">R<sub>1</sub></td><td style="text-align:center;">S<sub>7</sub></td><td style="text-align:center;">R<sub>1</sub></td><td style="text-align:center;">R<sub>1</sub></td><td></td><td></td><td></td></tr><tr><td style="text-align:center;">10</td><td></td><td></td><td style="text-align:center;">R<sub>3</sub></td><td style="text-align:center;">R<sub>3</sub></td><td style="text-align:center;">R<sub>3</sub></td><td style="text-align:center;">R<sub>3</sub></td><td></td><td></td><td></td></tr><tr><td style="text-align:center;">11</td><td></td><td></td><td style="text-align:center;">R<sub>5</sub></td><td style="text-align:center;">R<sub>5</sub></td><td style="text-align:center;">R<sub>5</sub></td><td style="text-align:center;">R<sub>5</sub></td><td></td><td></td><td></td></tr></tbody></table>

**测试序列：**
<table>
    <thead>
        <tr>
            <th>输入</th>
            <th>预期输出</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>id * id + id $</td>
            <td>Accept</td>
        </tr>
    </tbody>
</table>


#### 测试用例2

**文法：** 本实验文法

**LR(1) 分析表：** 由算法自动生成

**测试序列：**
<table>
    <thead>
        <tr>
            <th>输入</th>
            <th>预期输出</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>{ basic id ; } $</td>
            <td>Accept</td>
        </tr>
        <tr>
            <td>{ basic id ; basic id ; } $</td>
            <td>Accept</td>
        </tr>
        <tr>
            <td>{ basic id ; id = num ; } $</td>
            <td>Accept</td>
        </tr>
        <tr>
            <td>{ basic id ; id = ( num > num ) ; } $</td>
            <td>Accept</td>
        </tr>
        <tr>
            <td>{ basic id ; { basic id ; } } $</td>
            <td>Accept</td>
        </tr>
        <tr>
            <td>{ if ( bool ) { basic id ; } else { basic id ; } } $</td>
            <td>Accept</td>
        </tr>
        <tr>
            <td>{ if ( bool ) { basic id ; } else { basic id ; if ( bool ) { basic id ; } else { basic id ; } } } $</td>
            <td>Accept</td>
        </tr>
    </tbody>
</table>

#### 测试用例3

**文法：**
<table>
<tr><th style="text-align:center;">增广文法</th><th style="text-align:center;">文法</th><th style="text-align:center;">终结符</th></tr>
<tr><td valign="top">
    S' -> S <br>
</td>
<td valign="top">
    S -> { C } <br>
    C -> A B | A | B <br>
    A -> a | A a <br>
    B -> b | B b <br>
</td>
<td valign="top">
a, b, {, }, ε, $
</td></tr>
</table>

**LR(1) 分析表：** 由算法自动生成

**测试序列：**
<table>
    <thead>
        <tr>
            <th>输入</th>
            <th>预期输出</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>{ a a a b } $</td>
            <td>Accept</td>
        </tr>
    </tbody>
</table>