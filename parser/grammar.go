package parser

import (
	"slices"

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

// Copy creates a deep copy of the Grammar instance.
func (g *Grammar) Copy() Grammar {
	return Grammar{
		AugmentedProduction: g.AugmentedProduction,
		Productions:         slices.Clone(g.Productions),
		Terminals:           g.Terminals.Copy(),
	}
}

// IsTerminal checks if the given symbol is a terminal symbol in the grammar.
func (g *Grammar) IsTerminal(symbol Symbol) bool {
	return g.Terminals.Contains(Terminal(symbol))
}

// IsNonTerminal checks if the given symbol is a non-terminal symbol in the grammar.
func (g *Grammar) IsNonTerminal(symbol Symbol) bool {
	return !g.IsTerminal(symbol)
}

// GetIndex returns the index of the given production in the grammar's productions slice.
func (g *Grammar) GetIndex(production Production) int {
	if production.Equals(g.AugmentedProduction) {
		return 0
	}
	return slices.IndexFunc(g.Productions, func(p Production) bool {
		return p.Equals(production)
	})
}
