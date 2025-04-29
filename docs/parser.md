# Parser

## Purpose of the Lab

Develop an LR(1) parsing program based on the given grammar to analyze any input symbol string. The primary goal of this experiment is to deepen the understanding of the LR(1) parsing method.

## Tasks
Implement the following functionalities of the compiler:
1. Output the syntax analysis table;
2. Output the contents of the analysis stack.

## Grammar
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

Note: This grammar has some issues and may require modifications during the experiment.

## Overall Program Structure

![program structure](/docs/img/parser/1.png)

## LR(1)

LR(1) consists of the following parts:
1. Construction of the Item Sets
    1. Closure Operation
    2. GOTO Operation
    3. Construction of the Item Sets (State Transition Diagram)
2. Construction of the Parsing Table
    1. ACTION Table
    2. GOTO Table
3. Parsing Process
    1. Construction of the Parsing Stack
    2. Implementation of the Parsing Process

### 1. Construction of the Item Sets

#### 1.1 Closure Operation

The closure operation expands an item set into its closure. The basic idea is: for each item in the item set, if there is a non-terminal on the body and that non-terminal has productions in the grammar, then add all items of those productions to the closure.

Steps for the closure operation:
1. Initialize an empty item set and add the initial item to it.
2. For each item in the item set, check if there is a non-terminal on the body. If so, add all items of the productions of that non-terminal to the item set.
3. Repeat step 2 until the item set no longer changes.
4. Return the item set after closure.

Pseudocode:
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
Implementation:

1. Initialization:
    - Ensure the FIRST set is computed, as the closure operation depends on it.
    - Initialize the closure as a copy of the input item set and use `marks` to track processed items.
    ```go
    p.EnsureFirstSet()
    closure := make([]LR1Item, len(items))
    copy(closure, items)
    marks := Set[string]{}
    ```

2. Main Loop:
    - Continuously expand the closure until no further changes occur.
    ```go
    loop := true
    for loop {
         loop = false
         ...
    }
    ```

3. Iterate Over Items in the Closure:
    - Skip already processed items and mark the current item as processed.
    ```go
    for _, item := range closure {
         if marks.Contains(item.AsKey()) {
              continue
         }
         marks.Add(item.AsKey())
         ...
    }
    ```

4. Check the Dot Position:
    - If the dot exceeds the length of the production's body, skip the item.
    ```go
    if item.Dot >= len(item.Production.Body) {
         continue
    }
    ```

5. Get the Symbol After the Dot:
    - If the symbol after the dot is a terminal, it cannot be expanded, so skip it.
    ```go
    nextSymbol := item.Production.Body[item.Dot]
    if p.Grammar.IsTerminal(nextSymbol) {
         continue
    }
    ```

6. Iterate Over Productions in the Grammar:
    - Find productions where the head matches `nextSymbol`.
    ```go
    for _, production := range p.Grammar.Productions {
         if production.Head == nextSymbol {
              ...
         }
    }
    ```

7. Handle Empty Productions:
    - If the production's body is empty (ε), create a new item and add it to the closure.
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

8. Compute Lookahead Set:
    - Call `p.findLookaheads` to compute the lookahead set and add new items to the closure.
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

9. Return the Closure:
    - Once the closure no longer changes, return the final set of closure items.

##### Implementation of the `findLookaheads` Function

The `findLookaheads` function is used to compute the FIRST set of a given sequence of symbols and handle ε (empty string) according to grammar rules. The function takes as input a sequence of symbols and the current lookahead symbol.

Pseudocode:
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

1. **Handling Empty Symbol Sequences**:
   - If `symbols` is empty, directly return a set containing the `lookahead`.

2. **Initialization**:
   - Create an empty set `firstSet` to store the FIRST set of the symbol sequence.
   - Use a `flag` to indicate whether the `lookahead` should be added to the result.

3. **Iterating Over the Symbol Sequence**:
   - If the current symbol is a terminal, add it directly to `firstSet`.
   - If the current symbol is a non-terminal, add its FIRST set (excluding ε) to `firstSet`.
   - If the FIRST set of the current symbol does not contain ε, stop the iteration.

4. **Handling ε**:
   - If the FIRST sets of all symbols in the sequence contain ε, add the `lookahead` to the result.

5. **Returning the Result**:
   - Return the final computed `firstSet`.

Below is the implementation of the `findLookaheads` function:

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
##### Calculation of FIRST Set

Rules for calculating the FIRST set:
1. For a terminal, its FIRST set is itself.
2. For a non-terminal, its FIRST set is the union of the FIRST sets of all the bodies of its productions.
3. If the body of a production contains ε, add ε to the FIRST set.
4. If the body starts with a non-terminal, add the FIRST set of that non-terminal to the current FIRST set.
5. If the body starts with a terminal, add that terminal to the current FIRST set.
6. If the body is empty, add ε to the current FIRST set.
7. Repeat the above steps until the FIRST set no longer changes.

Pseudocode:
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

> When calculating the FIRST set, mutual dependencies (recursive calls) may occur:
> - The FIRST set of a non-terminal depends on the FIRST sets of other non-terminals.
>   - For example, in the production A → B C, the FIRST set of B depends on the FIRST set of A.
>   - In the production A → A b, the FIRST set of A directly depends on itself.
> - Direct recursive calls may lead to infinite recursion or stack overflow.
>
> To avoid this, we can use an **iterative update** approach to gradually approximate the final FIRST set, avoiding direct recursive calls. Specifically:
> - In each iteration, attempt to update the FIRST sets of all non-terminals.
>   - Try to add the FIRST set of the body symbols to the FIRST set of the head non-terminal.
>   - If the body symbol is a non-terminal, add its FIRST set (excluding ε) to the FIRST set of the head non-terminal.
>   - If the FIRST set of the body symbol contains ε, continue processing the next symbol until all symbols are processed or a terminal is encountered.
> - If the FIRST set of any non-terminal changes, proceed to the next iteration.
> - The algorithm terminates when all FIRST sets no longer change.
>
> Example:
>
> Given the following grammar:
> ```plaintext
> A → B C
> B → ε | a
> C → c | A
> ```
> Initial state:
> ```plaintext
> FIRST(B) = {ε, a}
> FIRST(C) = {}
> FIRST(A) = {}
> ```
> Iterative process:
> 1. First iteration:
>   - For A → B C:
>       - Add FIRST(B) to FIRST(A), resulting in FIRST(A) = {ε, a}.
>       - Since FIRST(B) contains ε, continue processing C.
>       - Add FIRST(C) to FIRST(A), but FIRST(C) is empty, so FIRST(A) remains unchanged.
>   - For C → c | A:
>       - Add c to FIRST(C), resulting in FIRST(C) = {c}.
>       - Add FIRST(A) to FIRST(C), resulting in FIRST(C) = {c, ε, a}.
> 2. Second iteration:
>   - For A → B C:
>       - Add FIRST(B) to FIRST(A), FIRST(A) remains unchanged.
>       - Since FIRST(B) contains ε, continue processing C.
>       - Add FIRST(C) to FIRST(A), resulting in FIRST(A) = {ε, a, c}.
>   - For C → c | A:
>       - FIRST(C) remains unchanged.
> 3. All FIRST sets no longer change, terminate the iteration.

Implementation:

1. **Initialization**:
    - The FIRST set of terminals is initialized to themselves.
      ```go
      for terminal := range p.Grammar.Terminals {
            p.FirstSet[Symbol(terminal)] = Set[Terminal]{}
            p.FirstSet[Symbol(terminal)].Add(terminal)
      }
      ```
    - The FIRST set of non-terminals is initialized as empty.
      ```go
      for _, production := range p.Grammar.Productions {
            if _, exists := p.FirstSet[production.Head]; !exists {
                 p.FirstSet[production.Head] = Set[Terminal]{}
            }
      }
      ```

2. **Iterative Computation**:
    - Continuously update the FIRST sets until they stabilize.
      ```go
      loop := true
      for loop {
            loop = false
            ...
      }
      ```

3. **Processing Each Production**:
    - **Empty Productions**: If the production body is empty or the first symbol is ε, add ε to the FIRST set.
      ```go
      if len(production.Body) == 0 || production.Body[0].IsEpsilon() {
            if !firstSet.Contains(EPSILON) {
                 firstSet.Add(EPSILON)
                 loop = true
            }
      }
      ```
    - **Iterate Over Symbols in the Body**:
      - If the symbol is ε, add ε to the FIRST set and stop.
      - If the symbol is a non-terminal, add its FIRST set (excluding ε) to the current FIRST set.
      - If the symbol is a terminal, directly add it to the current FIRST set and stop.
      ```go
      if symbol.IsEpsilon() {
            ...
      } else if symbolFirstSet, isNonTerminal := p.FirstSet[symbol]; isNonTerminal {
            ...
      } else {
            ...
      }
      ```

4. **Termination Condition**:
    - The computation completes when all FIRST sets no longer change.

#### 1.2 GOTO Operation

The GOTO operation transitions one item set to another. The basic idea is: for a given item set and a symbol, compute the item set corresponding to that symbol.

Steps for the GOTO operation:
1. For the given item set and symbol, initialize an empty item set.
2. Iterate through each item in the item set and check if the symbol after the dot matches the given symbol.
3. If they match, move the dot one position to the right, generating a new item set.
4. Perform the closure operation on the new item set and return the final item set.

Pseudocode:
```plaintext
function GOTO(I, X):
    J = {}
    for each item [A → α • Xβ, a] in I:
        if X is the symbol after the dot in α • Xβ:
            add [A → αX • β, a] to J
    return CLOSURE(J)
```

Implementation:
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
   
#### 1.3 Construction of the Item Sets (State Transition Diagram)

The process of constructing the item sets involves performing closure and GOTO operations on all possible item sets until no new item sets are generated.

Steps to construct the item sets:
1. Initialize an empty set of item sets and an initial item set (usually the start symbol's production of the grammar).
2. Perform the closure operation on the initial item set to obtain the initial state.
3. Add the initial state to the set of item sets.
4. Initialize an empty queue of states to process and add the initial state to the queue.
5. Process the states in the queue:
    1. Dequeue a state from the queue.
    2. Perform the GOTO operation for each symbol in the grammar to generate new states.
    3. If a new state is not already in the set of item sets, add it to the set and enqueue it.
    4. If the new state already exists in the set, record the transition relationship.
    5. Repeat step 5 until the queue is empty.
6. Record the set of item sets and the state transition relationships.

Pseudocode:
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

Implementation:

1. **Initialize the Symbol Set**  
    - Call the `EnsureSymbols` method to ensure that all symbols (terminals and non-terminals) in the grammar are added to the symbol set.  
    - These symbols will be used in subsequent GOTO operations.  
    ```go
    p.EnsureSymbols()
    ```

2. **Create the Initial State**  
    - Create an initial item `initialItem` representing the augmented production of the grammar (`S' → •S, $`), where:  
      - `Dot` indicates the position of the dot, initially set to 0.  
      - `Lookahead` is the terminal `$`, representing the end of input.  
    ```go
    initialItem := LR1Item{
          Production: p.Grammar.AugmentedProduction,
          Dot:        0,
          Lookahead:  TERMINATE,
    }
    ```
    - Create an initial state `initialState` containing:  
      - `Index`: The state index, initially set to 0.  
      - `Items`: The item set, initially containing only `initialItem`.  
      - `Transitions`: The state transition relationships, initially empty.  
    ```go
    initialState := &State{
          Index:       0,
          Items:       LR1Items{initialItem},
          Transitions: make(map[Symbol]*State),
    }
    ```
    - Perform the closure operation on the item set of the initial state to expand it and generate the complete initial state.  
    ```go
    initialState.Items = p.CLOSURE(initialState.Items)
    ```

3. **Initialize the Set of States**  
    - Add the initial state to the set of states `p.States`, which serves as the starting point for all states.  
    ```go
    p.States = States{initialState}
    ```

4. **Construct the Item Sets**  
    - Use a loop to process each state in the set of states.  
    - `length` keeps track of the current length of the state set, dynamically updated as new states are added (simulating a queue).  
    ```go
    length := len(p.States)
    for i := 0; i < length; i++ {
          state := p.States[i]
          ...
    }
    ```
    1. **Iterate Over Symbols**  
        - Iterate over all symbols (terminals and non-terminals) and perform the GOTO operation for each symbol.  
        - If the result of the GOTO operation is empty (i.e., no new item set), skip the symbol.  
        ```go
        for symbol := range p.Symbols {
              gotoItems := p.GOTO(state.Items, symbol)
              if len(gotoItems) == 0 {
                     continue
              }
              ...
        }
        ```

    2. **Create a New State**  
        - If the result of the GOTO operation is non-empty, create a new state `newState` containing:  
          - `Index`: The index of the new state, equal to the current length of the state set.  
          - `Items`: The item set generated by the GOTO operation.  
          - `Transitions`: An initially empty map of transition relationships.  
        ```go
        newState := &State{
              Index:       len(p.States),
              Items:       gotoItems,
              Transitions: make(map[Symbol]*State),
        }
        ```

    3. **Check if the State Already Exists**  
        - Check if the new state already exists in the set of states:  
          - Use `slices.IndexFunc` to iterate through the state set and compare the item sets of each state.  
          - If a matching state is found, return its index; otherwise, return -1.  
        ```go
        index := slices.IndexFunc(p.States, func(s *State) bool {
              return s.Equals(newState)
        })
        ```

    4. **Add the New State or Update the Transition Relationship**  
        - If the new state does not exist:  
          - Add it to the set of states `p.States`.  
          - Update the transition relationship of the current state `state.Transitions`, recording the transition from the current state to the new state via the symbol.  
          - Increment the length of the state set `length` to continue processing the new state.  
        - If the new state already exists:  
          - Directly update the transition relationship of the current state to point to the existing state.  
        ```go
        if index == -1 {
              p.States = append(p.States, newState)
              state.Transitions[symbol] = newState
              length++
        } else {
              state.Transitions[symbol] = p.States[index]
        }
        ```

5. **End of Loop**  
    - When all states have been processed (i.e., the queue is empty), the loop ends, and the construction is complete.

### 2. LR(1) Parsing Table
#### 2.1 ACTION Table
The ACTION table is used to record the operation for each terminal in every state. The basic idea is: for each state and terminal, determine whether to shift, reduce, or accept.

- The `Action` structure represents an operation, containing the action type and number.
- `ActionTable` is a nested map where the outer map's key is the state index, the inner map's key is the terminal, and the value is an `Action` structure.
- The `Copy` method is used to duplicate the ACTION table.
- The `Register` method is used to register an action into the ACTION table and check for conflicts.
    - If the combination of state and terminal already exists, check whether the action types conflict.
    - If there is a conflict, return an error.
    - If there is no conflict, register the action into the table.

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

#### 2.2 GOTO Table
The GOTO table is used to record the transition relationships for each non-terminal in every state. The basic idea is: for each state and non-terminal, determine the next state.

- `GotoTable` is a nested map where the outer map's key is the state index, the inner map's key is the non-terminal, and the value is the next state's index.
- The `Copy` method is used to duplicate the GOTO table.
- The `Register` method is used to register a transition relationship into the GOTO table and check for conflicts.
    - If the combination of state and non-terminal already exists, return an error.
    - If there is no conflict, register the transition relationship into the table.

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

### 2.3 LR(1) Table

The LR(1) table combines the ACTION and GOTO tables. The basic idea is to record the operations and transitions for each terminal and non-terminal in every state.

- The `LR1Table` structure contains the ACTION and GOTO tables.
- The `Insert` method is used to insert states and grammar into the LR(1) table.
    - Iterate through the item set of the state and check the position of the dot in each item.
    - If the dot is at the end of a production, check whether to accept or reduce.
    - If the dot is in the middle of a production, check whether the next symbol is a terminal or non-terminal, and register the corresponding action or transition.
    - If a conflict occurs, return an error.

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

Therefore, the construction of the LR(1) table only requires passing the state set and transition relationships established in the first step into the `Insert` method of the `LRTable`:
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

### 3. Parsing Process

The program primarily uses a `Walker` structure to implement the parsing process.

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

This structure includes information such as the parsing table, grammar, state stack, symbol stack, token stack, and symbol table. The state stack is used to store the current state, the symbol stack is used to store the current sequence of symbols, the token stack is used to store the current input tokens, and the symbol table is used to store variables, functions, and other information.

Note: The token stack does not participate in the parsing process itself. It is only used to store the current input tokens, facilitating subsequent semantic analysis and intermediate code generation.

The `Parser` structure is used to handle grammar parsing and construct the parsing table. It also saves or copies the parsing results into a new `Walker` structure.

```go
type Parser struct {
    Grammar *Grammar
    Symbols Set[Symbol]

    FirstSet FirstSet

    States States

    Table *LRTable
}
```

This structure includes information such as the grammar, symbol set, FIRST set, state collection, and parsing table.

The `NewWalker` method is provided to create a new `Walker` structure and pass in the parsing table and grammar (to enable concurrent parsing of multiple files).

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

#### 3.1 Construction of the Parsing Stack
The parsing stack is used to store the current state and symbol string. Its basic idea is: for each input token, check the current state and the state of the symbol stack to decide whether to shift, reduce, or accept.

This is reflected in the `Walker` structure's `States` and `Symbols` fields:
- `States` is used to store the current state.
- `Symbols` is used to store the current symbol string.

#### 3.2 Implementation of the Parsing Process
The implementation of the parsing process is mainly divided into the following steps:
1. Initialize the parsing stack and input tokens.
2. Loop through the input tokens until acceptance or an error occurs.
3. Based on the current state and input token, decide whether to shift, reduce, or accept.
4. Perform the corresponding operation, updating the parsing stack and input tokens.
5. If an error occurs, output the error message and exit.
6. If accepted, output the parsing result.
7. If reduced, update the symbol table and intermediate code generator.
8. Return the parsing result.
9. End the parsing process.

The program separates the parsing process into two parts: one for token transmission (`Parser`) and the other for syntax analysis and intermediate code generation (`Walker`).

##### Parser

The `Parser` structure is primarily used to handle input tokens and provides the `Parse` method.

This method takes a `lexer.Lexer` object and a logging function as parameters, loops through the input tokens, and passes them to the `Walker` for analysis.

Main logic:
1. **Initialization**:
    - Create a new `Walker` instance.
    - Call `walker.SymbolTable.EnterScope()` to enter a new scope.
2. **Token Loop**:
    - Use `l.NextToken()` to fetch the next token from the lexer.
    - If an error occurs and it is not the end of the file (`io.EOF`), terminate parsing.
    - If it is the end of the file (`io.EOF`), set the token type to `lexer.EOF`.
3. **Token Mapping to Symbol**:
    - Call `p.Reflect(token)` to convert the token into the corresponding symbol.
4. **Scope Management**:
    - If the token is a left brace, call `walker.SymbolTable.EnterScope()` to enter a new scope.
    - If the token is a right brace, call `walker.SymbolTable.ExitScope()` to exit the current scope.
5. **Walker State Update**:
    - Call `walker.Next(symbol)` to update the `Walker`'s state based on the current symbol.
    - If `action.Type` is `REDUCE`, continue processing the current symbol; otherwise, exit the inner loop (reduction operations may not complete in one step due to nested or recursive grammar, so repeated attempts are needed until no further reductions are possible).
6. **Termination Condition**:
    - If the symbol is `TERMINATE`, parsing is complete, and the loop exits.

Note: Token stacking must occur after the `Walker` completes a series of reduction operations. To simplify token stacking logic, the `Parser` directly pushes the token onto the `Walker` after reductions are complete (when the parsing table accepts the grammar and completes state transitions).

The implementation is as follows:
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

The `Walker` structure is primarily used to implement the LR(1) parsing process and provides the `Next` method.

This method takes a `Symbol` object as a parameter, representing the current input symbol. Based on the current state and input symbol, it determines whether to shift, reduce, or accept, and performs the corresponding operation.

By only accepting a `Symbol` as a parameter, the parsing process closely aligns with the actual syntax analysis operation, avoiding excessive coupling with the input tokens.

The main logic is as follows:
1. **Retrieve the Current State**:
    - Use `w.States.Peek()` to get the index of the current state.
2. **Determine the Type of Input Symbol**:
    - If it is a terminal, look up the ACTION table to get the corresponding action.
      - If it is a shift (SHIFT), push the state and symbol onto the stack.
      - If it is a reduction (REDUCE), invoke the corresponding production rule handler.
         - After handling, pop the states and symbols from the stack.
         - Look up the GOTO table to get the next state and push the production head onto the stack.
      - If it is an accept (ACCEPT), return the accept state.
    - If it is a non-terminal, look up the GOTO table to get the next state.
      - Push the state and symbol onto the stack.
    - If neither is found, return an error.

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

## Testing
### FIRST Set Testing
In the `TestParser_BuildFirstSet` test function located in [algorithm_test.go](/parser/algorithm_test.go), a set of simple grammars is used to test the construction of the FIRST set.

#### Test Cases
<table><tb>
<tr><th style="text-align:center;">Grammar</th><th style="text-align:center;">FIRST Set</th><th style="text-align:center;">Terminals</th></tr>
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
    

### Item Set Family Testing (State Transition Diagram)
In the `TestParser_BuildStates` test function located in [algorithm_test.go](/parser/algorithm_test.go), a set of simple grammars is used to test the construction of the item set family.

#### Test Cases

<table>
<tr><th style="text-align:center;">Augmented Grammar</th><th style="text-align:center;">Grammar</th><th style="text-align:center;">Terminals</th></tr>
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

### Parsing Table Construction Testing
In the `TestParser_BuildTable` test function located in [table_test.go](/parser/table_test.go), a set of simple grammars is used to test the construction of the parsing table.

<table>
<tr><th style="text-align:center;">Augmented Grammar</th><th style="text-align:center;">Grammar</th><th style="text-align:center;">Terminals</th></tr>
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
    The grammar used in this experiment
</td>
<td valign="top">
    The terminals of the grammar used in this experiment
</td></tr>
</table>

### Syntax Analysis Testing
In the `TestWalker_Next`, `TestWalker_Next2`, and `TestWalker_Next3` test functions located in [walker_test.go](/parser/walker_test.go), a set of simple grammars and token sequences are used to test the correctness of syntax analysis.

#### Test Case 1

**Grammar:**
<table>
<tr><th style="text-align:center;">Augmented Grammar</th><th style="text-align:center;">Grammar</th><th style="text-align:center;">Terminals</th></tr>
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

**LR(1) Parsing Table:**
<table><thead><tr><th rowspan="2" style="text-align:center;">State</th><th colspan="6" style="text-align:center;">ACTION</th><th colspan="3" style="text-align:center;">GOTO</th></tr><tr><th style="text-align:center;">id</th><th style="text-align:center;">(</th><th style="text-align:center;">+</th><th style="text-align:center;">*</th><th style="text-align:center;">)</th><th style="text-align:center;">$</th><th style="text-align:center;">E</th><th style="text-align:center;">T</th><th style="text-align:center;">F</th></tr></thead><tbody><tr><td style="text-align:center;">0</td><td style="text-align:center;">S<sub>5</sub></td><td style="text-align:center;">S<sub>4</sub></td><td></td><td></td><td></td><td></td><td style="text-align:center;">1</td><td style="text-align:center;">2</td><td style="text-align:center;">3</td></tr><tr><td style="text-align:center;">1</td><td></td><td></td><td style="text-align:center;">S<sub>6</sub></td><td></td><td></td><td style="text-align:center;">Accept</td><td></td><td></td><td></td></tr><tr><td style="text-align:center;">2</td><td></td><td></td><td style="text-align:center;">R<sub>2</sub></td><td style="text-align:center;">S<sub>7</sub></td><td style="text-align:center;">R<sub>2</sub></td><td style="text-align:center;">R<sub>2</sub></td><td></td><td></td><td></td></tr><tr><td style="text-align:center;">3</td><td></td><td></td><td style="text-align:center;">R<sub>4</sub></td><td style="text-align:center;">R<sub>4</sub></td><td style="text-align:center;">R<sub>4</sub></td><td style="text-align:center;">R<sub>4</sub></td><td></td><td></td><td></td></tr><tr><td style="text-align:center;">4</td><td style="text-align:center;">S<sub>5</sub></td><td style="text-align:center;">S<sub>4</sub></td><td></td><td></td><td></td><td></td><td style="text-align:center;">8</td><td style="text-align:center;">2</td><td style="text-align:center;">3</td></tr><tr><td style="text-align:center;">5</td><td></td><td></td><td style="text-align:center;">R<sub>6</sub></td><td style="text-align:center;">R<sub>6</sub></td><td style="text-align:center;">R<sub>6</sub></td><td style="text-align:center;">R<sub>6</sub></td><td></td><td></td><td></td></tr><tr><td style="text-align:center;">6</td><td style="text-align:center;">S<sub>5</sub></td><td style="text-align:center;">S<sub>4</sub></td><td></td><td></td><td></td><td></td><td></td><td style="text-align:center;">9</td><td style="text-align:center;">3</td></tr><tr><td style="text-align:center;">7</td><td style="text-align:center;">S<sub>5</sub></td><td style="text-align:center;">S<sub>4</sub></td><td></td><td></td><td></td><td></td><td></td><td></td><td style="text-align:center;">10</td></tr><tr><td style="text-align:center;">8</td><td></td><td></td><td style="text-align:center;">S<sub>6</sub></td><td></td><td style="text-align:center;">S<sub>11</sub></td><td></td><td></td><td></td><td></td></tr><tr><td style="text-align:center;">9</td><td></td><td></td><td style="text-align:center;">R<sub>1</sub></td><td style="text-align:center;">S<sub>7</sub></td><td style="text-align:center;">R<sub>1</sub></td><td style="text-align:center;">R<sub>1</sub></td><td></td><td></td><td></td></tr><tr><td style="text-align:center;">10</td><td></td><td></td><td style="text-align:center;">R<sub>3</sub></td><td style="text-align:center;">R<sub>3</sub></td><td style="text-align:center;">R<sub>3</sub></td><td style="text-align:center;">R<sub>3</sub></td><td></td><td></td><td></td></tr><tr><td style="text-align:center;">11</td><td></td><td></td><td style="text-align:center;">R<sub>5</sub></td><td style="text-align:center;">R<sub>5</sub></td><td style="text-align:center;">R<sub>5</sub></td><td style="text-align:center;">R<sub>5</sub></td><td></td><td></td><td></td></tr></tbody></table>

**Test Sequences:**
<table>
    <thead>
        <tr>
            <th>Input</th>
            <th>Expected Output</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>id * id + id $</td>
            <td>Accept</td>
        </tr>
    </tbody>
</table>


#### Test Case 2

**Grammar:** The grammar used in this experiment

**LR(1) Parsing Table:** Automatically generated by the algorithm

**Test Sequences:**
<table>
    <thead>
        <tr>
            <th>Input</th>
            <th>Expected Output</th>
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

#### Test Case 3

**Grammar:**
<table>
<tr><th style="text-align:center;">Augmented Grammar</th><th style="text-align:center;">Grammar</th><th style="text-align:center;">Terminals</th></tr>
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

**LR(1) Parsing Table:** Automatically generated by the algorithm

**Test Sequences:**
<table>
    <thead>
        <tr>
            <th>Input</th>
            <th>Expected Output</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>{ a a a b } $</td>
            <td>Accept</td>
        </tr>
    </tbody>
</table>