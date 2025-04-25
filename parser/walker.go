package parser

import (
	"fmt"

	. "app/parser/grammar"
	. "app/parser/production"
	. "app/utils/collections"
)

type Walker struct {
	Table   LRTable
	Grammar *Grammar

	States  Stack[int]
	Symbols Stack[Symbol]
}

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
		Grammar: &g,
		States:  states,
		Symbols: symbols,
	}
}

func (w *Walker) Next(symbol Symbol) (err error, action Action) {
	topState, _ := w.States.Peek()
	if w.Grammar.IsTerminal(symbol) {
		action, ok := w.Table.ActionTable[topState][Terminal(symbol)]
		if !ok {
			return fmt.Errorf("no action found for state %d and symbol %s", topState, symbol), Action{Type: ERROR}
		}
		switch action.Type {
		case SHIFT:
			w.States.Push(action.Number)
			w.Symbols.Push(symbol)
			return nil, Action{Type: SHIFT, Number: action.Number}
		case REDUCE:
			production := w.Grammar.Productions[action.Number]
			for i := 0; i < len(production.Body); i++ {
				if production.Body[i] == EPSILON {
					continue
				}
				w.States.Pop()
				w.Symbols.Pop()
			}
			topState, _ = w.States.Peek()
			gotoState, ok := w.Table.GotoTable[topState][production.Head]
			if !ok {
				return fmt.Errorf("no goto state found for state %d and symbol %s", topState, production.Head), Action{Type: ERROR}
			}
			w.Symbols.Push(production.Head)
			w.States.Push(gotoState)
			return nil, Action{Type: REDUCE, Number: action.Number}
		case ACCEPT:
			return nil, Action{Type: ACCEPT, Number: 0}
		}
	} else {
		action, ok := w.Table.GotoTable[topState][symbol]
		if !ok {
			return fmt.Errorf("no goto state found for state %d and symbol %s", topState, symbol), Action{Type: ERROR}
		}
		w.States.Push(action)
		w.Symbols.Push(symbol)
		fmt.Println("GOTO", action)
		return nil, Action{Type: GOTO, Number: action}
	}
	return fmt.Errorf("unexpected state %d and symbol %s", topState, symbol), Action{Type: ERROR}
}

func (w *Walker) Reset() {
	w.States.Clear()
	w.Symbols.Clear()
	w.States.Push(0)
}
