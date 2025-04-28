package parser_test

import (
	"fmt"
	"slices"
	"testing"

	. "app/parser"
	. "app/utils/collections"
	"app/utils/log"
)

func TestWalker_Next(t *testing.T) {
	walker := Walker{
		Grammar: &tableGrammar,
		Table:   table,
		States:  Stack[int]{},
		Symbols: Stack[Symbol]{},
	}
	walker.States.Push(0)
	seq := []Symbol{"id", "*", "id", "+", "id", TERMINATE}
	for i := 0; i < len(seq); i++ {
		fmt.Print(log.Sprintf(log.Argument{FrontColor: log.Blue, Highlight: true, Format: "State: %v", Args: []any{walker.States}}))
		fmt.Print(log.Sprintf(log.Argument{FrontColor: log.Yellow, Highlight: true, Format: ", Symbols: %v", Args: []any{walker.Symbols}}))
		fmt.Print(log.Sprintf(log.Argument{FrontColor: log.Magenta, Highlight: true, Format: ", Symbol: %s", Args: []any{seq[i]}}))
		action, err := walker.Next(seq[i])
		fmt.Print(log.Sprintf(log.Argument{FrontColor: log.Green, Highlight: true, Format: ", Action: %v", Args: []any{action}}))
		if err != nil {
			fmt.Print(log.Sprintf(log.Argument{FrontColor: log.Red, Highlight: true, Format: ", Error: %v", Args: []any{err}}))
			break
		}
		fmt.Println()
		if action.Type == REDUCE {
			i--
		}
	}
}

func TestWalker_Next2(t *testing.T) {
	p := NewParser()
	p.EnsureTable()

	//for first := range p.FirstSet {
	//	if p.Grammar.IsNonTerminal(first) {
	//		fmt.Printf("FirstSet[%s] = %v\n", first, p.FirstSet[first])
	//	}
	//}
	//
	//for _, state := range p.States {
	//	fmt.Printf("State %d: \n", state.Index)
	//	for _, item := range state.Items {
	//		fmt.Printf("%v\n", item)
	//	}
	//	for k, v := range state.Transitions {
	//		fmt.Printf("  %s -> State %d\n", k, v.Index)
	//	}
	//}
	//type Item struct {
	//	Action Action
	//	Symbol Symbol
	//	Index  int
	//}
	//var actions []Item
	//for k, v := range p.Table.ActionTable {
	//	for k2, v2 := range v {
	//		actions = append(actions, Item{Action: v2, Symbol: Symbol(k2), Index: k})
	//	}
	//}
	//
	//slices.SortFunc(actions, func(a, b Item) int {
	//	return a.Index - b.Index
	//})
	//
	//for _, action := range actions {
	//	fmt.Printf("Action[%d][%s] = %v\n", action.Index, action.Symbol, action.Action)
	//}

	walker := p.NewWalker()

	seqs := [][]Symbol{
		{"{", "basic", "id", ";", "}", TERMINATE},
		{"{", "basic", "id", ";", "basic", "id", ";", "}", TERMINATE},
		{"{", "basic", "id", ";", "id", "=", "num", ";", "}", TERMINATE},
		{"{", "basic", "id", ";", "id", "=", "(", "num", ">", "num", ")", ";", "}", TERMINATE},
		{"{", "basic", "id", ";", "{", "basic", "id", ";", "}", "}", TERMINATE},
		{"{", "if", "(", "bool", ")", "{", "basic", "id", ";", "}", "else", "{", "basic", "id", ";", "}", "}", TERMINATE},
		{"{", "if", "(", "bool", ")", "{", "basic", "id", ";", "}", "else", "{", "basic", "id", ";",
			"if", "(", "bool", ")", "{", "basic", "id", ";", "}", "else", "{", "basic", "id", ";", "}", "}", "}", TERMINATE},
	}
	fmt.Println("=======================")
	for _, seq := range seqs {
		for i := 0; i < len(seq); i++ {
			fmt.Print(log.Sprintf(log.Argument{FrontColor: log.Blue, Highlight: true, Format: "State: %v", Args: []any{walker.States}}))
			fmt.Print(log.Sprintf(log.Argument{FrontColor: log.Yellow, Highlight: true, Format: ", Symbols: %v", Args: []any{walker.Symbols}}))
			fmt.Print(log.Sprintf(log.Argument{FrontColor: log.Magenta, Highlight: true, Format: ", Symbol: %s", Args: []any{seq[i]}}))
			action, err := walker.Next(seq[i])
			fmt.Print(log.Sprintf(log.Argument{FrontColor: log.Green, Highlight: true, Format: ", Action: %v", Args: []any{action}}))
			if err != nil {
				fmt.Print(log.Sprintf(log.Argument{FrontColor: log.Red, Highlight: true, Format: ", Error: %v\n", Args: []any{err}}))
				break
			}
			fmt.Println()
			if action.Type == REDUCE {
				i--
			}
		}

		walker.Reset()
		fmt.Println("=======================")
	}
}

func TestWalker_Next3(t *testing.T) {
	p := Parser{
		Grammar: &Grammar{
			AugmentedProduction: Production{
				Head: "S'",
				Body: []Symbol{"S"},
			},
			Productions: []Production{
				{
					Head: "S",
					Body: []Symbol{"{", "C", "}"},
				},
				{
					Head: "C",
					Body: []Symbol{"A", "B"},
				},
				{
					Head: "C",
					Body: []Symbol{"A"},
				},
				{
					Head: "C",
					Body: []Symbol{"B"},
				},
				{
					Head: "A",
					Body: []Symbol{"a"},
				},
				{
					Head: "A",
					Body: []Symbol{"A", "a"},
				},
				{
					Head: "B",
					Body: []Symbol{"b"},
				},
				{
					Head: "B",
					Body: []Symbol{"B", "b"},
				},
			},
			Terminals: Set[Terminal]{}.AddAll("a", "b", "{", "}", EPSILON, TERMINATE),
		},
	}
	p.EnsureTable()

	for _, state := range p.States {
		fmt.Printf("State %d: \n", state.Index)
		for _, item := range state.Items {
			fmt.Printf("%v\n", item)
		}
		for k, v := range state.Transitions {
			fmt.Printf("  %s -> State %d\n", k, v.Index)
		}
	}

	type Item struct {
		Action Action
		Symbol Symbol
		Index  int
	}
	var actions []Item
	for k, v := range p.Table.ActionTable {
		for k2, v2 := range v {
			actions = append(actions, Item{Action: v2, Symbol: Symbol(k2), Index: k})
		}
	}

	slices.SortFunc(actions, func(a, b Item) int {
		return a.Index - b.Index
	})

	for _, action := range actions {
		fmt.Printf("Action[%d][%s] = %v\n", action.Index, action.Symbol, action.Action)
	}

	walker := p.NewWalker()
	seq := []Symbol{"{", "a", "a", "a", "b", "}", TERMINATE}
	for i := 0; i < len(seq); i++ {
		fmt.Print(log.Sprintf(log.Argument{FrontColor: log.Blue, Highlight: true, Format: "State: %v", Args: []any{walker.States}}))
		fmt.Print(log.Sprintf(log.Argument{FrontColor: log.Yellow, Highlight: true, Format: ", Symbols: %v", Args: []any{walker.Symbols}}))
		fmt.Print(log.Sprintf(log.Argument{FrontColor: log.Magenta, Highlight: true, Format: ", Symbol: %s", Args: []any{seq[i]}}))
		action, err := walker.Next(seq[i])
		fmt.Print(log.Sprintf(log.Argument{FrontColor: log.Green, Highlight: true, Format: ", Action: %v", Args: []any{action}}))
		if err != nil {
			fmt.Print(log.Sprintf(log.Argument{FrontColor: log.Red, Highlight: true, Format: ", Error: %v", Args: []any{err}}))
			break
		}
		fmt.Println()
		if action.Type == REDUCE {
			i--
		}
	}
}
