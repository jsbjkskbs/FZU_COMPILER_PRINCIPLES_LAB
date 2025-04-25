package parser

import (
	"fmt"
	"maps"

	. "app/parser/grammar"
	. "app/parser/production"
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

type ActionType string

const (
	SHIFT  ActionType = "shift"
	REDUCE ActionType = "reduce"
	ACCEPT ActionType = "accept"
	ERROR  ActionType = "error"
	GOTO   ActionType = "goto"
)
