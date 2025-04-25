package parser_test

import (
	"fmt"
	"testing"

	. "app/parser"
	. "app/parser/grammar"
	. "app/parser/production"
	. "app/utils/collections"
)

func TestParser_BuildTable(t *testing.T) {
	tests := []struct {
		name    string
		grammar Grammar
	}{
		{
			name:    "Test1",
			grammar: grammars[0],
		},
		{
			name:    "Test2",
			grammar: grammars[1],
		},
		{
			name: "Test3",
			grammar: Grammar{
				AugmentedProduction: AugmentedProduction,
				Productions:         Productions,
				Terminals:           Terminals,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				Grammar:  &tt.grammar,
				Symbols:  Set[Symbol]{},
				FirstSet: FirstSet{},
				States:   States{},
			}
			p.EnsureStates()

			for _, state := range p.States {
				fmt.Printf("State %d: \n", state.Index)
				for _, item := range state.Items {
					fmt.Printf("%v\n", item)
				}
				for k, v := range state.Transitions {
					fmt.Printf("  %s -> State %d\n", k, v.Index)
				}
			}
			p.EnsureTable()

			fmt.Printf("Got %d states\n", len(p.States))

			for k, v := range p.Table.ActionTable {
				for k2, v2 := range v {
					fmt.Printf("Action[%v][%s] = %v\n", k, k2, v2)
				}
			}

			for k, v := range p.Table.GotoTable {
				for k2, v2 := range v {
					fmt.Printf("Goto[%v][%s] = %v\n", k, k2, v2)
				}
			}
		})
	}
}
