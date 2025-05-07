package parser

import (
	"fmt"
	"maps"
)

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

type LRTable struct {
	ActionTable ActionTable
	GotoTable   GotoTable
}

// Insert populates the LR table with actions and transitions based on the given state and grammar.
// It iterates through the items in the state and determines the appropriate action
// (SHIFT, REDUCE, ACCEPT) based on the grammar rules. It also updates the Goto table
// for non-terminal symbols. The function handles conflicts by checking for existing actions
// and reporting them if necessary.
// The function takes a state and a grammar as input and returns an error if any conflicts are found.
// It uses the state index and the lookahead symbol to determine the appropriate action.
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
			// fmt.Printf("when inserting : %v\n", err)
		}
	}
}

type Action struct {
	Type   ActionType
	Number int
}

type ActionTable map[int]map[Terminal]Action

// Copy creates a deep copy of the ActionTable.
func (t ActionTable) Copy() ActionTable {
	return maps.Clone(t)
}

// Register adds an action to the action table for a given state and terminal.
// It checks for conflicts and returns an error if a conflict is found.
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

type GotoTable map[int]map[Symbol]int

func (t GotoTable) Copy() GotoTable {
	return maps.Clone(t)
}

// Register adds a transition to the Goto table for a given state and non-terminal symbol.
// It checks for conflicts and returns an error if a conflict is found.
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

type ActionType string

const (
	SHIFT  ActionType = "shift"
	REDUCE ActionType = "reduce"
	ACCEPT ActionType = "accept"
	ERROR  ActionType = "error"
	GOTO   ActionType = "goto"
)

type SymbolTableItem struct {
	Variable string
	Type     SymbolTableItemType
	Address  int

	UnderlyingType string

	VariableSize int

	ArraySize        int
	ArrayElementSize int
	Dimension        []int
}

type SymbolTableItemType string

const (
	SymbolTableItemTypeVariable SymbolTableItemType = "variable"
	SymbolTableItemTypeArray    SymbolTableItemType = "array"
	SymbolTableItemTypeConstant SymbolTableItemType = "constant"
	SymbolTableItemTypeUnknown  SymbolTableItemType = "unknown"
)

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

	addrCounter int
}

const (
	initialAddr = 0x10000000
)

// NewSymbolTable creates a new symbol table with the specified enter and exit functions.
// The enter function is called when entering a new scope, and the exit function is called when exiting a scope.
func NewSymbolTable(enter, exit func(*Scope) error) *SymbolTable {
	return &SymbolTable{
		LegacyScopes:  make([]*Scope, 0),
		CurrentScope:  nil,
		EnterFunction: enter,
		ExitFunction:  exit,
		addrCounter:   initialAddr,
	}
}

// EnterScope creates a new scope and sets it as the current scope in the symbol table.
func (st *SymbolTable) EnterScope() error {
	if st.CurrentScope == nil {
		st.CurrentScope = &Scope{
			ID:     len(st.LegacyScopes),
			Level:  0,
			Items:  make(map[string]*SymbolTableItem),
			Parent: nil,
		}
	} else {
		st.CurrentScope = &Scope{
			ID:     len(st.LegacyScopes),
			Level:  st.CurrentScope.Level + 1,
			Items:  make(map[string]*SymbolTableItem),
			Parent: st.CurrentScope,
		}
	}
	st.LegacyScopes = append(st.LegacyScopes, st.CurrentScope)

	if st.EnterFunction != nil {
		if err := st.EnterFunction(st.CurrentScope); err != nil {
			return err
		}
	}
	return nil
}

// ExitScope exits the current scope and sets the parent scope as the current scope.
func (st *SymbolTable) ExitScope() error {
	if st.CurrentScope == nil {
		return fmt.Errorf("no scope to exit")
	}

	if st.ExitFunction != nil {
		if err := st.ExitFunction(st.CurrentScope); err != nil {
			return err
		}
	}

	st.CurrentScope = st.CurrentScope.Parent
	return nil
}

// Register adds a new item to the current scope in the symbol table.
// It checks for conflicts and ensures that the item is valid before adding it.
func (st *SymbolTable) Register(item *SymbolTableItem) (int, error) {
	if st.CurrentScope == nil {
		return -1, fmt.Errorf("no scope to register item")
	}

	if _, exists := st.CurrentScope.Items[item.Variable]; exists {
		return -1, fmt.Errorf("item %s already exists in scope", item.Variable)
	}

	if item.VariableSize <= 0 {
		return -1, fmt.Errorf("invalid variable size for item %s", item.Variable)
	}
	st.CurrentScope.Items[item.Variable] = item
	switch item.Type {
	case SymbolTableItemTypeVariable:
		item.Address = st.addrCounter
		st.addrCounter += item.VariableSize / 4
		if item.VariableSize/4*4 != item.VariableSize {
			st.addrCounter++
		}
	case SymbolTableItemTypeArray:
		if item.ArraySize <= 0 {
			return -1, fmt.Errorf("invalid array size for item %s", item.Variable)
		}
		if item.ArrayElementSize <= 0 {
			return -1, fmt.Errorf("invalid array element size for item %s", item.Variable)
		}
		item.Address = st.addrCounter
		st.addrCounter += item.ArrayElementSize * item.ArraySize / 4
		if item.ArrayElementSize*item.ArraySize/4*4 != item.ArrayElementSize*item.ArraySize {
			st.addrCounter++
		}
	}
	return item.Address, nil
}

func (st *SymbolTable) ArrayAddress(variable string, dimension []int) (int, int, error) {
	if st.CurrentScope == nil {
		return -1, -1, fmt.Errorf("no scope to lookup item")
	}
	item, _, err := st.Lookup(variable)
	if err != nil {
		return -1, -1, err
	}
	if item.Type != SymbolTableItemTypeArray {
		return -1, -1, fmt.Errorf("item %s is not an array", variable)
	}
	if len(dimension) != len(item.Dimension) {
		for i := len(dimension); i < len(item.Dimension); i++ {
			dimension = append(dimension, 0)
		}
	}
	offset := 0
	for i, dim := range dimension {
		if dim < 0 || dim >= item.Dimension[i] {
			return -1, -1, fmt.Errorf("index out of bounds for dimension %d of item %s", i, variable)
		}
		multiplier := 1
		for j := i + 1; j < len(item.Dimension); j++ {
			multiplier *= item.Dimension[j]
		}
		offset += dim * multiplier
	}

	if offset < 0 || offset >= item.ArraySize {
		return -1, -1, fmt.Errorf("index out of bounds for item %s", variable)
	}
	return item.Address + (item.ArrayElementSize * offset / 4), item.ArraySize, nil
}

// arrayAddress returns the address of the specified variable in the symbol table.
// It looks up the variable in the current scope and its parent scopes until it finds it or returns an error.
// It returns the address of the variable and an error if any.
func (st *SymbolTable) arrayAddress(variable string, offset int) (int, int, error) {
	if st.CurrentScope == nil {
		return -1, -1, fmt.Errorf("no scope to lookup item")
	}
	item, _, err := st.Lookup(variable)
	if err != nil {
		return -1, -1, err
	}
	if item.Type != SymbolTableItemTypeArray {
		return -1, -1, fmt.Errorf("item %s is not an array", variable)
	}
	return item.Address + (item.ArrayElementSize * offset / 4), item.ArraySize, nil
}

// Lookup searches for an item in the symbol table.
// It checks the current scope and its parent scopes until it finds the item or returns an error.
// It returns the item, a boolean indicating if it was found in the current scope, and an error if any.
// If the item is not found, it returns an error indicating that the item was not found in any scope.
func (st *SymbolTable) Lookup(variable string) (item *SymbolTableItem, findInCurrentScope bool, err error) {
	if st.CurrentScope == nil {
		return nil, false, fmt.Errorf("no scope to lookup item")
	}

	scope := st.CurrentScope
	for scope != nil {
		if item, exists := scope.Items[variable]; exists {
			return item, scope == st.CurrentScope, nil
		}
		scope = scope.Parent
	}

	return nil, false, fmt.Errorf("item %s not found in any scope", variable)
}

// TempAddr generates a temporary address for a variable in the symbol table.
// It uses the current address counter and increments it based on the size of the variable.
// It returns the address of the temporary variable.
func (st *SymbolTable) TempAddr(size int) int {
	addr := st.addrCounter
	st.addrCounter += size / 4
	if size/4*4 != size {
		st.addrCounter++
	}
	return addr
}
