package parser

import (
	"fmt"

	"app/lexer"
	. "app/utils/collections"
)

type Walker struct {
	Table   LRTable
	Grammar *Grammar

	States  Stack[int]
	Symbols Stack[Symbol] // Symbols is just for debugging purposes
	Tokens  Stack[*lexer.Token]

	SymbolTable *SymbolTable
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
		Grammar:     &g,
		States:      states,
		Symbols:     symbols,
		SymbolTable: NewSymbolTable(nil, nil),
	}
}

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

func (w *Walker) Reset() {
	w.States.Clear()
	w.Symbols.Clear()
	w.Tokens.Clear()
	w.States.Push(0)
}
