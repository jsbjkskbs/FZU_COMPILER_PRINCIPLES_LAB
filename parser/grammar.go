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

func (g *Grammar) Copy() Grammar {
	return Grammar{
		AugmentedProduction: g.AugmentedProduction,
		Productions:         slices.Clone(g.Productions),
		Terminals:           g.Terminals.Copy(),
	}
}

func (g *Grammar) IsTerminal(symbol Symbol) bool {
	return g.Terminals.Contains(Terminal(symbol))
}

func (g *Grammar) IsNonTerminal(symbol Symbol) bool {
	return !g.IsTerminal(symbol)
}

func (g *Grammar) GetIndex(production Production) int {
	if production.Equals(g.AugmentedProduction) {
		return 0
	}
	return slices.IndexFunc(g.Productions, func(p Production) bool {
		return p.Equals(production)
	})
}
