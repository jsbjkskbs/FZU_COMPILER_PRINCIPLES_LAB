package grammar

import (
	. "app/parser/production"
	. "app/utils/collections"
)

type Grammar struct {
	AugmentedProduction Production
	Productions         []Production
	Terminals           Set[Terminal]
}

func NewGrammar() *Grammar {
	return &Grammar{
		AugmentedProduction: AugmentedProduction,
		Productions:         Productions,
		Terminals:           Terminals,
	}
}

func (g Grammar) IsTerminal(symbol Symbol) bool {
	return g.Terminals.Contains(Terminal(symbol))
}

func (g Grammar) IsNonTerminal(symbol Symbol) bool {
	return !g.IsTerminal(symbol)
}
