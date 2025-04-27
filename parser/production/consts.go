package production

import (
	. "app/utils/collections"
)

const (
	EPSILON   = "ε"
	TERMINATE = "$"
)

var Terminals = Set[Terminal]{}.AddAll(
	// Brackets and punctuation
	"{", "}", ";", "[", "]", "(", ")",

	// Arithmetic operators
	"+", "-", "*", "/",

	// Logical and comparison operators
	"||", "&&", "==", "!=", "<", "<=", ">", ">=", "!", "=", "!=",

	// Keywords
	"if", "else", "while", "do", "break",

	// Literals
	"true", "false",

	// Types
	"basic", "id", "num", "real",

	// Special symbols
	EPSILON, TERMINATE,
)

var AugmentedProduction = Production{
	Head: "program'",
	Body: []Symbol{"program"},
}

var OptimizedSymbols = Set[Symbol]{}.AddAll()

var Productions = []Production{
	// program → block
	{
		Head: "program",
		Body: []Symbol{"block"},
	},
	// block → { decls stmts }
	// ** optimized to combined_decls_stmts **
	// block → { combined_decls_stmts }
	// combined_decls_stmts → decls stmts | stmts | decls | ε
	{
		Head: "block",
		Body: []Symbol{"{", "decls", "stmts", "}"},
	},
	{
		Head: "block",
		Body: []Symbol{"{", "decls", "}"},
	},
	{
		Head: "block",
		Body: []Symbol{"{", "stmts", "}"},
	},
	{
		Head: "block",
		Body: []Symbol{"{", "}"},
	},
	// decls → decls decl | ε
	{
		Head: "decls",
		Body: []Symbol{"decls", "decl"},
	},
	{
		Head: "decls",
		Body: []Symbol{EPSILON}, // ε
	},
	// decl → type id;
	{
		Head: "decl",
		Body: []Symbol{"type", "id", ";"},
	},
	// type → type[num] | basic
	{
		Head: "type",
		Body: []Symbol{"type", "[", "num", "]"},
	},
	{
		Head: "type",
		Body: []Symbol{"basic"},
	},
	// stmts → stmts stmt | ε
	{
		Head: "stmts",
		Body: []Symbol{"stmts", "stmt"},
	},
	{
		Head: "stmts",
		Body: []Symbol{EPSILON}, // ε
	},
	// stmt → matched_stmt | unmatched_stmt
	// ** solved else-hanging problem **
	{
		Head: "stmt",
		Body: []Symbol{"matched_stmt"},
	},
	{
		Head: "stmt",
		Body: []Symbol{"unmatched_stmt"},
	},
	{
		Head: "stmt",
		Body: []Symbol{"decls"},
	},
	// unmatched_stmt → if ( bool ) unmatched_stmt
	{
		Head: "unmatched_stmt",
		Body: []Symbol{"if", "(", "bool", ")", "unmatched_stmt"},
	},
	// unmatched_stmt → if ( bool ) matched_stmt else unmatched_stmt
	{
		Head: "unmatched_stmt",
		Body: []Symbol{"if", "(", "bool", ")", "matched_stmt", "else", "unmatched_stmt"},
	},
	// matched_stmt → loc = bool ;
	{
		Head: "matched_stmt",
		Body: []Symbol{"loc", "=", "bool", ";"},
	},
	// matched_stmt → if ( bool ) matched_stmt else matched_stmt
	{
		Head: "matched_stmt",
		Body: []Symbol{"if", "(", "bool", ")", "matched_stmt", "else", "matched_stmt"},
	},
	// matched_stmt → while ( bool ) stmt
	{
		Head: "matched_stmt",
		Body: []Symbol{"while", "(", "bool", ")", "stmt"},
	},
	// matched_stmt → do stmt while ( bool ) ;
	{
		Head: "matched_stmt",
		Body: []Symbol{"do", "stmt", "while", "(", "bool", ")", ";"},
	},
	// matched_stmt → break ;
	{
		Head: "matched_stmt",
		Body: []Symbol{"break", ";"},
	},
	// matched_stmt → block
	{
		Head: "matched_stmt",
		Body: []Symbol{"block"},
	},
	// loc → loc[num] | id
	{
		Head: "loc",
		Body: []Symbol{"loc", "[", "num", "]"},
	},
	{
		Head: "loc",
		Body: []Symbol{"id"},
	},
	// bool → bool || join | join
	{
		Head: "bool",
		Body: []Symbol{"bool", "||", "join"},
	},
	{
		Head: "bool",
		Body: []Symbol{"join"},
	},
	// join → join && equality | equality
	{
		Head: "join",
		Body: []Symbol{"join", "&&", "equality"},
	},
	{
		Head: "join",
		Body: []Symbol{"equality"},
	},
	// equality → equality == rel | equality != rel | rel
	{
		Head: "equality",
		Body: []Symbol{"equality", "==", "rel"},
	},
	{
		Head: "equality",
		Body: []Symbol{"equality", "!=", "rel"},
	},
	{
		Head: "equality",
		Body: []Symbol{"rel"},
	},
	// rel → expr<expr | expr<=expr | expr>=expr | expr>expr | expr
	{
		Head: "rel",
		Body: []Symbol{"expr", "<", "expr"},
	},
	{
		Head: "rel",
		Body: []Symbol{"expr", "<=", "expr"},
	},
	{
		Head: "rel",
		Body: []Symbol{"expr", ">=", "expr"},
	},
	{
		Head: "rel",
		Body: []Symbol{"expr", ">", "expr"},
	},
	{
		Head: "rel",
		Body: []Symbol{"expr"},
	},
	// expr → expr+term | expr-term | term
	{
		Head: "expr",
		Body: []Symbol{"expr", "+", "term"},
	},
	{
		Head: "expr",
		Body: []Symbol{"expr", "-", "term"},
	},
	{
		Head: "expr",
		Body: []Symbol{"term"},
	},
	// term → term*unary | term/unary | unary
	{
		Head: "term",
		Body: []Symbol{"term", "*", "unary"},
	},
	{
		Head: "term",
		Body: []Symbol{"term", "/", "unary"},
	},
	{
		Head: "term",
		Body: []Symbol{"unary"},
	},
	// unary → !unary | -unary | factor
	{
		Head: "unary",
		Body: []Symbol{"!", "unary"},
	},
	{
		Head: "unary",
		Body: []Symbol{"-", "unary"},
	},
	{
		Head: "unary",
		Body: []Symbol{"factor"},
	},
	// factor → (bool) | loc | num | real | true | false
	{
		Head: "factor",
		Body: []Symbol{"(", "bool", ")"},
	},
	{
		Head: "factor",
		Body: []Symbol{"loc"},
	},
	{
		Head: "factor",
		Body: []Symbol{"num"},
	},
	{
		Head: "factor",
		Body: []Symbol{"real"},
	},
	{
		Head: "factor",
		Body: []Symbol{"true"},
	},
	{
		Head: "factor",
		Body: []Symbol{"false"},
	},
}
