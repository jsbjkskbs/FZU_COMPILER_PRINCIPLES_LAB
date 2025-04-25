package parser_test

import (
	. "app/parser"
	. "app/parser/grammar"
	. "app/parser/production"
	. "app/utils/collections"
)

var grammars = []Grammar{
	{
		AugmentedProduction: Production{Head: "S'", Body: []Symbol{"S"}},
		Productions: []Production{
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
		Terminals: Set[Terminal]{}.AddAll("*", "=", "id", EPSILON, TERMINATE),
	},
	{
		AugmentedProduction: Production{Head: "S'", Body: []Symbol{"S"}},
		Productions: []Production{
			{
				Head: "S",
				Body: []Symbol{"V", "=", "E"},
			},
			{
				Head: "S",
				Body: []Symbol{"E"},
			},
			{
				Head: "E",
				Body: []Symbol{"V"},
			},
			{
				Head: "V",
				Body: []Symbol{"x"},
			},
			{
				Head: "V",
				Body: []Symbol{"*", "E"},
			},
		},
		Terminals: Set[Terminal]{}.AddAll("=", "x", "*", EPSILON, TERMINATE),
	},
}

var tableGrammar = Grammar{
	AugmentedProduction: Production{Head: "E'", Body: []Symbol{"E"}},
	Productions: []Production{
		{},
		{
			Head: "E",
			Body: []Symbol{"E", "+", "T"},
		},
		{
			Head: "E",
			Body: []Symbol{"T"},
		},
		{
			Head: "T",
			Body: []Symbol{"T", "*", "F"},
		},
		{
			Head: "T",
			Body: []Symbol{"F"},
		},
		{
			Head: "F",
			Body: []Symbol{"(", "E", ")"},
		},
		{
			Head: "F",
			Body: []Symbol{"id"},
		},
	},
	Terminals: Set[Terminal]{}.AddAll("(", ")", "+", "*", "id", EPSILON, TERMINATE),
}

var table = LRTable{
	ActionTable: ActionTable{
		0:  {"id": Shift(5), "(": Shift(4)},
		1:  {"+": Shift(6), TERMINATE: Accept()},
		2:  {"+": Reduce(2), "*": Shift(7), ")": Reduce(2), TERMINATE: Reduce(2)},
		3:  {"+": Reduce(4), "*": Reduce(4), ")": Reduce(4), TERMINATE: Reduce(4)},
		4:  {"id": Shift(5), "(": Shift(4)},
		5:  {"+": Reduce(6), "*": Reduce(6), ")": Reduce(6), TERMINATE: Reduce(6)},
		6:  {"id": Shift(5), "(": Shift(4)},
		7:  {"id": Shift(5), "(": Shift(4)},
		8:  {"+": Shift(6), ")": Shift(11)},
		9:  {"+": Reduce(1), "": Shift(7), ")": Reduce(1), TERMINATE: Reduce(1)},
		10: {"+": Reduce(3), "*": Reduce(3), ")": Reduce(3), TERMINATE: Reduce(3)},
		11: {"+": Reduce(5), "": Reduce(5), ")": Reduce(5), TERMINATE: Reduce(5)},
	},
	GotoTable: GotoTable{
		0: {"E": 1, "T": 2, "F": 3},
		4: {"E": 8, "T": 2, "F": 3},
		6: {"T": 9, "F": 3},
		7: {"F": 10},
	},
}

func Shift(number int) Action {
	return Action{
		Type:   SHIFT,
		Number: number,
	}
}

func Reduce(number int) Action {
	return Action{
		Type:   REDUCE,
		Number: number,
	}
}

func Accept() Action {
	return Action{
		Type:   ACCEPT,
		Number: 0,
	}
}
