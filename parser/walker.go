package parser

import (
	"fmt"
	"strings"

	. "app/utils/collections"
)

// Walker is a structure that represents the current state of the parser
// during the parsing process. It contains information about the current
// state, the symbols being processed, and the grammar being used.
type Walker struct {
	Table   LRTable
	Grammar *Grammar

	States  Stack[int]
	Symbols Stack[Symbol] // Symbols is just for debugging purposes
	Tokens  Stack[*ASTNode]

	SymbolTable *SymbolTable

	Environment  *Environment
	ThreeAddress []string

	ast *AbstractSyntaxTree
}

type Environment struct {
	onBooleanExpr bool

	BreakLabelStack Stack[*[]int]
	LoopLabelStack  Stack[int]
	LabelStack      Stack[int]
	EndIfStmtStack  Stack[int]
}

// NewEnvironment creates a new Environment instance and initializes it.
// The Environment is used to store the current state of the parser, including
// the current type, data type, data size, array size, variable name, etc.
// It is used to keep track of the current context during parsing and code generation.
// The Environment is reset to its initial state when a new parsing context is created.
func NewEnvironment() *Environment {
	e := &Environment{}
	e.Reset()
	return e
}

// Reset resets the environment to its initial state.
// The Walker operates in a bottom-up manner, as dictated by the grammar, and
// does not access parts of the parse tree that are not directly involved in
// the current context. This ensures that the parsing process adheres strictly
// to the grammar's rules and structure. So, there's no need to reset the environment.
func (env *Environment) Reset() {
}

// NewWalker creates a new Walker instance and initializes it with the
// grammar and action tables. The Walker is used to traverse the parse tree
// and perform actions based on the grammar rules. It maintains a stack of
// states and symbols, as well as a symbol table for managing variables and types.
func (p *Parser) NewWalker() *Walker {
	p.EnsureTable()

	g := p.Grammar.Copy()
	states := Stack[int]{}
	states.Push(0)
	symbols := Stack[Symbol]{}
	w := &Walker{
		Table: LRTable{
			ActionTable: p.Table.ActionTable.Copy(),
			GotoTable:   p.Table.GotoTable.Copy(),
		},
		Grammar:     &g,
		States:      states,
		Symbols:     symbols,
		SymbolTable: NewSymbolTable(nil, nil),
		Environment: NewEnvironment(),
	}
	w.Emit("jmp", "L1")
	return w
}

// Next processes the next symbol in the parsing process. It takes a symbol as input
// and returns an action and an error. The action can be SHIFT, REDUCE, ACCEPT, or ERROR.
// The function uses the current state and the symbol to determine the appropriate action
// to take.
// If the action is ACCEPT, it indicates that the parsing is complete.
// If the action is SHIFT, it pushes the new state and symbol onto the stacks.
// If the action is REDUCE, it pops the appropriate number of symbols from the stacks
// and applies the corresponding production rule. If the action is ACCEPT, it indicates
// that the parsing is complete.
// If there is an error, it returns an error message.
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

// Reset resets the Walker's state, symbol, and token stacks to their initial state.
// It clears the stacks and pushes the initial state (0) onto the state stack.
func (w *Walker) Reset() {
	w.States.Clear()
	w.Symbols.Clear()
	w.Tokens.Clear()
	w.States.Push(0)
}

func (w *Walker) NewLabel() int {
	w.ThreeAddress = append(w.ThreeAddress, fmt.Sprintf("L%-8d %8s", len(w.ThreeAddress), "xxx"))
	w.Environment.LabelStack.Push(len(w.ThreeAddress) - 1)
	return len(w.ThreeAddress) - 1
}

func (w *Walker) Emit(op string, dist string, args ...any) {
	line := fmt.Sprintf("L%-8d %8s %16s", len(w.ThreeAddress), op, dist)
	for _, arg := range args {
		line += fmt.Sprintf(" %16s", arg)
	}
	w.ThreeAddress = append(w.ThreeAddress, line)
}

func (w *Walker) NewGotoLabel() int {
	w.ThreeAddress = append(w.ThreeAddress, fmt.Sprintf("L%-8d %8s", len(w.ThreeAddress), "yyy"))
	w.Environment.EndIfStmtStack.Push(len(w.ThreeAddress) - 1)
	return len(w.ThreeAddress) - 1
}

func (w *Walker) EmitGoto(label int, distLabel int) {
	line := fmt.Sprintf("L%-8d %8s %16s", label, "jmp", fmt.Sprintf("L%d", distLabel))
	w.ThreeAddress[label] = line
}

func (w *Walker) EmitLabel(label int, dist string, op string, args ...any) {
	line := fmt.Sprintf("L%-8d %8s %16s", label, op, dist)
	for _, arg := range args {
		line += fmt.Sprintf(" %16s", arg)
	}
	w.ThreeAddress[label] = line
}

func (w *Walker) AdjustJMP(label int, jmp int) error {
	line := w.ThreeAddress[label]
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return fmt.Errorf("invalid line format: %s", line)
	}
	if parts[1] != "jmp" && parts[1] != "jz" && parts[1] != "jnz" {
		return fmt.Errorf("invalid jump instruction: %s", parts[1])
	}
	w.ThreeAddress[label] = fmt.Sprintf("L%-8d %8s %16s", label, parts[1], fmt.Sprintf("L%d", jmp))
	return nil
}

func (w *Walker) GetCurrentLabelCount() int {
	return len(w.ThreeAddress)
}

func (w *Walker) EnterLoop() {
	a := make([]int, 0)
	w.Environment.BreakLabelStack.Push(&a)
}

func (w *Walker) AddBreakLabel() {
	if w.Environment.BreakLabelStack.IsEmpty() {
		println("AddBreakLabel: BreakLabelStack is empty")
		return
	}
	w.ThreeAddress = append(w.ThreeAddress, fmt.Sprintf("L%-8d %8s", len(w.ThreeAddress), "nop"))
	t, _ := w.Environment.BreakLabelStack.Pop()
	*t = append(*t, len(w.ThreeAddress)-1)
	w.Environment.BreakLabelStack.Push(t)
}

func (w *Walker) ExitLoop(exit int) {
	if w.Environment.BreakLabelStack.IsEmpty() {
		println("ExitLoop: BreakLabelStack is empty")
		return
	}
	l, _ := w.Environment.BreakLabelStack.Pop()
	for _, label := range *l {
		w.EmitLabel(label, fmt.Sprintf("L%d", exit), "jmp")
	}
}
