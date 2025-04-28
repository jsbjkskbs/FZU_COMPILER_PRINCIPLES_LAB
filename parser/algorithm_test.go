package parser_test

import (
	"fmt"
	"testing"

	. "app/parser"
	. "app/utils/collections"
	"app/utils/log"
)

func TestNewParser(t *testing.T) {
	p := NewParser()
	p.EnsureStates()
	l := len(p.States)
	for range 100 {
		pi := NewParser()
		pi.EnsureStates()
		if len(pi.States) != l {
			t.Errorf("Expected %d, got %d\n", l, len(pi.States))
		}
	}
	fmt.Println(log.Sprintf(log.Argument{FrontColor: log.Green, Highlight: true, Format: "Parser test passed[Stable States]", Args: []any{}}))
	// 注：只能测试LR(0)文法
	// tests := []struct {
	// 	name string
	// 	seq  []Symbol
	// }{
	// 	{
	// 		name: "Test1",
	// 		seq:  []Symbol{"{", "decls", "stmts", "}"},
	// 	},
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		p := NewParser()
	// 		p.EnsureStates()
	// 		state := p.States[0]
	// 		fmt.Printf("sequence: %s\n",
	// 			log.Sprintf(log.Argument{FrontColor: log.Green, Highlight: true, Format: "%s", Args: []any{tt.seq}}),
	// 		)
	// 		for _, symbol := range tt.seq {
	// 			if _, ok := state.Transitions[symbol]; !ok {
	// 				t.Errorf("Expected %s, got %v\n", symbol, state.Transitions[symbol])
	// 				break
	// 			}
	// 			fmt.Printf("%s%s%s\n",
	// 				log.Sprintf(log.Argument{FrontColor: log.Green, Highlight: true, Format: ">> from state: %d ", Args: []any{state.Index}}),
	// 				log.Sprintf(log.Argument{FrontColor: log.Red, Highlight: true, Format: ">> with symbol: %v ", Args: []any{symbol}}),
	// 				log.Sprintf(log.Argument{FrontColor: log.Green, Highlight: true, Format: ">> to state: %d", Args: []any{state.Transitions[symbol].Index}}),
	// 			)
	// 			state = state.Transitions[symbol]
	// 		}
	// 		fmt.Printf("final state: %s%s\n%s\n",
	// 			log.Sprintf(log.Argument{FrontColor: log.Green, Highlight: true, Format: "index: %d, ", Args: []any{state.Index}}),
	// 			log.Sprintf(log.Argument{FrontColor: log.Green, Highlight: true, Format: "items: %v", Args: []any{state.Items}}),
	// 			log.Sprintf(log.Argument{FrontColor: log.Cyan, Highlight: true, Format: "transitions: %v", Args: []any{state.Transitions}}),
	// 		)
	// 		fmt.Println()
	// 	})
	// }
}

func TestParser_BuildFirstSet(t *testing.T) {
	tests := []struct {
		name        string
		productions []Production
		terminals   []Terminal
		expected    FirstSet
	}{
		{
			name: "Test1",
			productions: []Production{
				{Head: "E", Body: []Symbol{"T", "E'"}},
				{Head: "E'", Body: []Symbol{"+", "T", "E'"}},
				{Head: "E'", Body: []Symbol{EPSILON}},
				{Head: "T", Body: []Symbol{"F", "T'"}},
				{Head: "T'", Body: []Symbol{"*", "F", "T'"}},
				{Head: "T'", Body: []Symbol{EPSILON}},
				{Head: "F", Body: []Symbol{"(", "E", ")"}},
				{Head: "F", Body: []Symbol{"id"}},
			},
			terminals: []Terminal{"id", "+", "*", "(", ")", EPSILON},
			expected: FirstSet{
				"E":  Set[Terminal]{}.AddAll("(", "id"),
				"E'": Set[Terminal]{}.AddAll("+", EPSILON),
				"T":  Set[Terminal]{}.AddAll("(", "id"),
				"T'": Set[Terminal]{}.AddAll("*", EPSILON),
				"F":  Set[Terminal]{}.AddAll("(", "id"),
			},
		},
		{
			name: "Test2",
			productions: []Production{
				{Head: "S", Body: []Symbol{"a", "A", "B", "b", "c", "d"}},
				{Head: "S", Body: []Symbol{"ε"}},
				{Head: "A", Body: []Symbol{"A", "S", "d"}},
				{Head: "A", Body: []Symbol{"ε"}},
				{Head: "B", Body: []Symbol{"S", "A", "h"}},
				{Head: "B", Body: []Symbol{"e", "C"}},
				{Head: "B", Body: []Symbol{"ε"}},
				{Head: "C", Body: []Symbol{"S", "f"}},
				{Head: "C", Body: []Symbol{"C", "g"}},
				{Head: "C", Body: []Symbol{"ε"}},
			},
			terminals: []Terminal{"a", "b", "c", "d", "e", "f", "g", "h", "ε"},
			expected: FirstSet{
				"S": Set[Terminal]{}.AddAll("a", "ε"),
				"A": Set[Terminal]{}.AddAll("a", "d", "ε"),
				"B": Set[Terminal]{}.AddAll("a", "d", "h", "e", "ε"),
				"C": Set[Terminal]{}.AddAll("a", "f", "g", "ε"),
			},
		},
		{
			name: "Test3",
			productions: []Production{
				{Head: "E", Body: []Symbol{"T", "E'"}},
				{Head: "E'", Body: []Symbol{"+", "E"}},
				{Head: "E'", Body: []Symbol{"ε"}},
				{Head: "T", Body: []Symbol{"F", "T'"}},
				{Head: "T'", Body: []Symbol{"T"}},
				{Head: "T'", Body: []Symbol{"ε"}},
				{Head: "F", Body: []Symbol{"P", "F'"}},
				{Head: "F'", Body: []Symbol{"*", "F'"}},
				{Head: "F'", Body: []Symbol{"ε"}},
				{Head: "P", Body: []Symbol{"(", "E", ")"}},
				{Head: "P", Body: []Symbol{"a"}},
				{Head: "P", Body: []Symbol{"b"}},
				{Head: "P", Body: []Symbol{"^"}},
			},
			terminals: []Terminal{"+", "(", ")", "a", "b", "^", "*", "ε"},
			expected: FirstSet{
				"E":  Set[Terminal]{}.AddAll("(", "a", "b", "^"),
				"E'": Set[Terminal]{}.AddAll("+", "ε"),
				"T":  Set[Terminal]{}.AddAll("(", "a", "b", "^"),
				"T'": Set[Terminal]{}.AddAll("(", "a", "b", "^", "ε"),
				"F":  Set[Terminal]{}.AddAll("(", "a", "b", "^"),
				"F'": Set[Terminal]{}.AddAll("*", "ε"),
				"P":  Set[Terminal]{}.AddAll("(", "a", "b", "^"),
			},
		},
		{
			name: "Test4",
			productions: []Production{
				{Head: "D", Body: []Symbol{"B", "c"}},
				{Head: "D", Body: []Symbol{"c"}},
				{Head: "B", Body: []Symbol{"b", "D"}},
				{Head: "B", Body: []Symbol{"a"}},
			},
			terminals: []Terminal{"a", "b", "c"},
			expected: FirstSet{
				"B": Set[Terminal]{}.AddAll("a", "b"),
				"D": Set[Terminal]{}.AddAll("a", "b", "c"),
			},
		},
		{
			name: "Test5",
			productions: []Production{
				{Head: "E", Body: []Symbol{"T", "A"}},
				{Head: "A", Body: []Symbol{"+", "T", "A"}},
				{Head: "A", Body: []Symbol{"ε"}},
				{Head: "T", Body: []Symbol{"F", "B"}},
				{Head: "B", Body: []Symbol{"*", "F", "B"}},
				{Head: "B", Body: []Symbol{"ε"}},
				{Head: "F", Body: []Symbol{"i"}},
				{Head: "F", Body: []Symbol{"(", "E", ")"}},
			},
			terminals: []Terminal{"+", "*", "i", "(", ")", "ε"},
			expected: FirstSet{
				"E": Set[Terminal]{}.AddAll("i", "("),
				"A": Set[Terminal]{}.AddAll("+", "ε"),
				"T": Set[Terminal]{}.AddAll("i", "("),
				"B": Set[Terminal]{}.AddAll("*", "ε"),
				"F": Set[Terminal]{}.AddAll("i", "("),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				Grammar: &Grammar{
					Productions: tt.productions,
					Terminals:   Set[Terminal]{}.AddAll(tt.terminals...),
				},
			}
			p.BuildFirstSet()
			for head, expected := range tt.expected {
				fmt.Printf("FIRST(%s) : %v", head, p.FirstSet[head])
				if !p.FirstSet[head].Equals(expected) {
					fmt.Println(log.Sprintf(log.Argument{FrontColor: log.Red, Highlight: true, Format: " !!! FAILED", Args: []any{}}))
					t.Errorf("Expected %v, got %v\n", expected, p.FirstSet[head])
				} else {
					fmt.Println(log.Sprintf(log.Argument{FrontColor: log.Green, Highlight: true, Format: " *** PASSED", Args: []any{}}))
				}
			}
		})
	}
}

func TestParser_BuildStates(t *testing.T) {
	tests := []struct {
		name                string
		augmentedProduction Production
		productions         []Production
		terminals           []Terminal
	}{
		{
			name:                "Test1",
			augmentedProduction: Production{Head: "S'", Body: []Symbol{"S"}},
			productions: []Production{
				{
					Head: "S",
					Body: []Symbol{"A", "a", "A", "b"},
				},
				{
					Head: "S",
					Body: []Symbol{"B", "b", "B", "a"},
				},
				{
					Head: "A",
					Body: []Symbol{EPSILON},
				},
				{
					Head: "B",
					Body: []Symbol{EPSILON},
				},
			},
			terminals: []Terminal{"a", "b", EPSILON, TERMINATE},
		},
		{
			name:                "Test2",
			augmentedProduction: Production{Head: "S'", Body: []Symbol{"S"}},
			productions: []Production{
				{
					Head: "S",
					Body: []Symbol{"L", "=", "R"},
				},
				{
					Head: "S",
					Body: []Symbol{"R"},
				},
				{
					Head: "L",
					Body: []Symbol{"*", "R"},
				},
				{
					Head: "L",
					Body: []Symbol{"id"},
				},
				{
					Head: "R",
					Body: []Symbol{"L"},
				},
			},
			terminals: []Terminal{"*", "=", "id", EPSILON, TERMINATE},
		},
		{
			name:                "Test3",
			augmentedProduction: Production{Head: "S'", Body: []Symbol{"S"}},
			productions: []Production{
				{
					Head: "S",
					Body: []Symbol{"B", "B"},
				},
				{
					Head: "B",
					Body: []Symbol{"a", "B"},
				},
				{
					Head: "B",
					Body: []Symbol{"b"},
				},
			},
			terminals: []Terminal{"a", "b", EPSILON, TERMINATE},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				Grammar: &Grammar{
					AugmentedProduction: tt.augmentedProduction,
					Productions:         tt.productions,
					Terminals:           Set[Terminal]{}.AddAll(tt.terminals...),
				},
			}
			p.BuildFirstSet()
			p.BuildStates()
			for _, state := range p.States {
				fmt.Printf("State %d:\n", state.Index)
				fmt.Println("Items:")
				for _, item := range state.Items {
					fmt.Println(item)
				}
				fmt.Println("Transitions:")
				for k, v := range state.Transitions {
					fmt.Printf("%s -> %d\n", k, v.Index)
				}
				fmt.Println()
			}
		})
	}
}
