package parser

import (
	"slices"

	. "app/parser/production"
	. "app/utils/collections"
)

func (p *Parser) BuildStates() {
	p.EnsureSymbols()

	initialItem := LR1Item{
		Production: p.Grammar.AugmentedProduction,
		Dot:        0,
		Lookahead:  TERMINATE,
	}

	initialState := &State{
		Index:       0,
		Items:       LR1Items{initialItem},
		Transitions: make(map[Symbol]*State),
	}

	initialState.Items = p.CLOSURE(initialState.Items)

	p.States = States{initialState}

	length := len(p.States)
	for i := 0; i < length; i++ {
		state := p.States[i]

		for symbol := range p.Symbols {
			gotoItems := p.GOTO(state.Items, symbol)
			if len(gotoItems) == 0 {
				continue
			}

			newState := &State{
				Index:       len(p.States),
				Items:       p.CLOSURE(gotoItems),
				Transitions: make(map[Symbol]*State),
			}
			index := slices.IndexFunc(p.States, func(s *State) bool {
				return s.Equals(newState)
			})
			if index == -1 {
				p.States = append(p.States, newState)
				state.Transitions[symbol] = newState
				length++
			} else {
				state.Transitions[symbol] = p.States[index]
			}
		}
	}
}

func (p *Parser) BuildSymbols() {
	p.Symbols = Set[Symbol]{}
	for _, production := range p.Grammar.Productions {
		p.Symbols.Add(production.Head)
		for _, symbol := range production.Body {
			p.Symbols.Add(symbol)
		}
	}

	p.Symbols.Remove(EPSILON)
}

func (p *Parser) BuildFirstSet() {
	p.FirstSet = make(FirstSet)

	for terminal := range p.Grammar.Terminals {
		p.FirstSet[Symbol(terminal)] = Set[Terminal]{}
		p.FirstSet[Symbol(terminal)].Add(terminal)
	}

	for _, production := range p.Grammar.Productions {
		if _, exists := p.FirstSet[production.Head]; !exists {
			p.FirstSet[production.Head] = Set[Terminal]{}
		}
	}

	loop := true
	for loop {
		loop = false
		for _, production := range p.Grammar.Productions {
			firstSet := p.FirstSet[production.Head]

			if len(production.Body) == 0 {
				if !firstSet.Contains(EPSILON) {
					firstSet.Add(EPSILON)
					loop = true
				}
			}

			for _, symbol := range production.Body {
				if symbol.IsEpsilon() {
					if !firstSet.Contains(EPSILON) {
						firstSet.Add(EPSILON)
						loop = true
					}
					break
				}

				if symbolFirstSet, isNonTerminal := p.FirstSet[symbol]; isNonTerminal {
					for terminal := range symbolFirstSet {
						if !terminal.IsEpsilon() && !firstSet.Contains(terminal) {
							firstSet.Add(terminal)
							loop = true
						}
					}
					if !symbolFirstSet.Contains(EPSILON) {
						break
					}

					if symbol == production.Body[len(production.Body)-1] && symbolFirstSet.Contains(EPSILON) {
						firstSet.Add(EPSILON)
						loop = true
					}
				} else {
					if !firstSet.Contains(Terminal(symbol)) {
						firstSet.Add(Terminal(symbol))
						loop = true
					}
					break
				}
			}
		}
	}
}

func (p *Parser) CLOSURE(items []LR1Item) []LR1Item {
	p.EnsureFirstSet()

	closure := make([]LR1Item, len(items))
	copy(closure, items)

	marks := Set[string]{}

	loop := true
	for loop {
		loop = false

		for _, item := range closure {
			if marks.Contains(item.AsKey()) {
				continue
			}

			marks.Add(item.AsKey())

			if item.Dot >= len(item.Production.Body) {
				continue
			}

			nextSymbol := item.Production.Body[item.Dot]
			if p.Grammar.IsTerminal(nextSymbol) {
				continue
			}

			for _, production := range p.Grammar.Productions {
				if production.Head == nextSymbol {
					if len(production.Body) == 0 {
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
					} else {
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
					}
				}
			}
		}
	}
	return closure
}

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
	return gotoItems
}

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
