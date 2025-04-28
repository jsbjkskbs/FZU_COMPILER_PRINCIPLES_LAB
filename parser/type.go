package parser

import (
	"fmt"
	"slices"

	. "app/utils/collections"
)

type Parser struct {
	Grammar *Grammar
	Symbols Set[Symbol]

	FirstSet FirstSet

	States States

	Table *LRTable
}

func NewParser() *Parser {
	return &Parser{
		Grammar:  NewGrammar(),
		Symbols:  Set[Symbol]{},
		FirstSet: FirstSet{},
		States:   States{},
	}
}

type FirstSet map[Symbol]Set[Terminal]

type State struct {
	Index       int
	Items       LR1Items
	Transitions map[Symbol]*State
}

func (state *State) Equals(other *State) bool {
	if len(state.Items) != len(other.Items) {
		return false
	}

	for _, item := range state.Items {
		if !other.Items.Contains(item) {
			return false
		}
	}
	return true
}

type States []*State

func (states *States) Contains(state *State) bool {
	for _, s := range *states {
		if len(s.Items) != len(state.Items) {
			continue
		}
		if s.Equals(state) {
			return true
		}
	}
	return false
}

type LR1Item struct {
	Production Production
	Dot        int
	Lookahead  Terminal
}

func (i *LR1Item) AsKey() string {
	return fmt.Sprintf("%s\a%s\a%d\a%s", i.Production.Head, i.Production.Body, i.Dot, i.Lookahead)
}

func (i *LR1Item) String() string {
	s := fmt.Sprintf("%s -> ", i.Production.Head)
	for j, symbol := range i.Production.Body {
		if j == i.Dot {
			s += ". "
		}
		s += string(symbol) + " "
	}
	if i.Dot == len(i.Production.Body) {
		s += ". "
	}
	s += fmt.Sprintf("(%s)", i.Lookahead)
	return s
}

func (i *LR1Item) Equals(other LR1Item) bool {
	if !i.Production.Equals(other.Production) {
		return false
	}
	if i.Dot != other.Dot {
		return false
	}
	if i.Lookahead != other.Lookahead {
		return false
	}
	return true
}

type LR1Items []LR1Item

func (items *LR1Items) Contains(other LR1Item) bool {
	return slices.ContainsFunc(*items, func(item LR1Item) bool {
		return item.Equals(other)
	})
}
