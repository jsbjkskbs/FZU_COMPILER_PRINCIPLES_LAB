package parser

import (
	"fmt"
	"slices"

	. "app/utils/collections"
	"app/utils/log"
)

type Production struct {
	Head Symbol
	Body []Symbol

	Rule Rule
}

type Rule func(*Walker) error

// Equals checks if two productions are equal by comparing their heads and bodies.
func (p *Production) Equals(other Production) bool {
	if p.Head != other.Head {
		return false
	}
	return slices.Equal(p.Body, other.Body)
}

// HandleRule executes the rule associated with the production if it is not nil.
func (p *Production) HandleRule(walker *Walker) error {
	if p.Rule == nil {
		fmt.Println(log.Sprintf(log.Argument{FrontColor: log.Yellow, Highlight: true, Format: "Warning: rule is nil for production %s -> %s", Args: []any{p.Head, p.Body}}))
		return nil
	}
	return p.Rule(walker)
}

type Symbol string

// IsEpsilon checks if the symbol is equal to EPSILON.
func (s *Symbol) IsEpsilon() bool {
	return *s == EPSILON
}

type Terminal string

// IsEpsilon checks if the terminal is equal to EPSILON.
func (t *Terminal) IsEpsilon() bool {
	return *t == EPSILON
}

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
		Rule: GenRules.Program,
	},
	// block → { decls stmts }
	// ** optimized to combined_decls_stmts **
	// block → { combined_decls_stmts }
	// combined_decls_stmts → decls stmts | stmts | decls | ε
	{
		Head: "block",
		Body: []Symbol{"{", "decls", "stmts", "}"},
		Rule: GenRules.BlockDeclsStmts,
	},
	{
		Head: "block",
		Body: []Symbol{"{", "decls", "}"},
		Rule: GenRules.BlockDecls,
	},
	{
		Head: "block",
		Body: []Symbol{"{", "stmts", "}"},
		Rule: GenRules.BlockStmts,
	},
	{
		Head: "block",
		Body: []Symbol{"{", "}"},
		Rule: GenRules.BlockEpsilon,
	},
	// decls → decls decl | ε
	{
		Head: "decls",
		Body: []Symbol{"decls", "decl"},
		Rule: GenRules.Decls,
	},
	{
		Head: "decls",
		Body: []Symbol{EPSILON}, // ε
		Rule: GenRules.DeclsEpsilon,
	},
	// decl → type id;
	{
		Head: "decl",
		Body: []Symbol{"type", "id", ";"},
		Rule: GenRules.Decl,
	},
	// type → type[num] | basic
	{
		Head: "type",
		Body: []Symbol{"type", "[", "num", "]"},
		Rule: GenRules.TypeArray,
	},
	{
		Head: "type",
		Body: []Symbol{"basic"},
		Rule: GenRules.TypeBasic,
	},
	// stmts → stmts stmt | ε
	{
		Head: "stmts",
		Body: []Symbol{"stmts", "stmt"},
		Rule: GenRules.Stmts,
	},
	{
		Head: "stmts",
		Body: []Symbol{EPSILON}, // ε
		Rule: GenRules.StmtsEpsilon,
	},
	// stmt → matched_stmt | unmatched_stmt
	// ** solved else-hanging problem **
	{
		Head: "stmt",
		Body: []Symbol{"matched_stmt"},
		Rule: GenRules.StmtMatchedStmt,
	},
	{
		Head: "stmt",
		Body: []Symbol{"unmatched_stmt"},
		Rule: GenRules.StmtUnmatchedStmt,
	},
	{
		Head: "stmt",
		Body: []Symbol{"decls"},
		Rule: GenRules.StmtDecls,
	},
	// unmatched_stmt → if ( bool ) unmatched_stmt
	{
		Head: "unmatched_stmt",
		Body: []Symbol{"if", "(", "bool", ")", "unmatched_stmt"},
		Rule: GenRules.UnmatchedStmtIf,
	},
	// unmatched_stmt → if ( bool ) matched_stmt else unmatched_stmt
	{
		Head: "unmatched_stmt",
		Body: []Symbol{"if", "(", "bool", ")", "matched_stmt", "else", "unmatched_stmt"},
		Rule: GenRules.UnmatchedStmtIfElse,
	},
	// matched_stmt → loc = bool ;
	{
		Head: "matched_stmt",
		Body: []Symbol{"loc", "=", "bool", ";"},
		Rule: GenRules.MatchedStmtAssign,
	},
	// matched_stmt → if ( bool ) matched_stmt else matched_stmt | if ( bool ) matched_stmt
	{
		Head: "matched_stmt",
		Body: []Symbol{"if", "(", "bool", ")", "matched_stmt", "else", "matched_stmt"},
		Rule: GenRules.MatchedStmtIfElse,
	},
	{
		Head: "matched_stmt",
		Body: []Symbol{"if", "(", "bool", ")", "matched_stmt"},
		Rule: GenRules.MatchedStmtIf,
	},
	// matched_stmt → while ( bool ) stmt
	{
		Head: "matched_stmt",
		Body: []Symbol{"while", "(", "bool", ")", "stmt"},
		Rule: GenRules.MatchedStmtWhile,
	},
	// matched_stmt → do stmt while ( bool ) ;
	{
		Head: "matched_stmt",
		Body: []Symbol{"do", "stmt", "while", "(", "bool", ")", ";"},
		Rule: GenRules.MatchedStmtDoWhile,
	},
	// matched_stmt → break ;
	{
		Head: "matched_stmt",
		Body: []Symbol{"break", ";"},
		Rule: GenRules.MatchedStmtBreak,
	},
	// matched_stmt → block
	{
		Head: "matched_stmt",
		Body: []Symbol{"block"},
		Rule: GenRules.MatchedStmtBlock,
	},
	// loc → loc[num] | id
	{
		Head: "loc",
		Body: []Symbol{"loc", "[", "num", "]"},
		Rule: GenRules.LocArray,
	},
	{
		Head: "loc",
		Body: []Symbol{"id"},
		Rule: GenRules.LocId,
	},
	// bool → bool || join | join
	{
		Head: "bool",
		Body: []Symbol{"bool'"},
		Rule: GenRules.Bool,
	},
	{
		Head: "bool'",
		Body: []Symbol{"bool'", "||", "join"},
		Rule: GenRules.BoolPrime,
	},
	{
		Head: "bool'",
		Body: []Symbol{"join"},
		Rule: GenRules.BoolPrimeJoin,
	},
	// join → join && equality | equality
	{
		Head: "join",
		Body: []Symbol{"join", "&&", "equality"},
		Rule: GenRules.Join,
	},
	{
		Head: "join",
		Body: []Symbol{"equality"},
		Rule: GenRules.JoinEquality,
	},
	// equality → equality == rel | equality != rel | rel
	{
		Head: "equality",
		Body: []Symbol{"equality", "==", "rel"},
		Rule: GenRules.Equality,
	},
	{
		Head: "equality",
		Body: []Symbol{"equality", "!=", "rel"},
		Rule: GenRules.NotEquality,
	},
	{
		Head: "equality",
		Body: []Symbol{"rel"},
		Rule: GenRules.EqualityRelational,
	},
	// rel → expr<expr | expr<=expr | expr>=expr | expr>expr | expr
	{
		Head: "rel",
		Body: []Symbol{"expr", "<", "expr"},
		Rule: GenRules.RelationalLess,
	},
	{
		Head: "rel",
		Body: []Symbol{"expr", "<=", "expr"},
		Rule: GenRules.RelationalLessEqual,
	},
	{
		Head: "rel",
		Body: []Symbol{"expr", ">=", "expr"},
		Rule: GenRules.RelationalGreaterEqual,
	},
	{
		Head: "rel",
		Body: []Symbol{"expr", ">", "expr"},
		Rule: GenRules.RelationalGreater,
	},
	{
		Head: "rel",
		Body: []Symbol{"expr"},
		Rule: GenRules.RelationalExpr,
	},
	// expr → expr+term | expr-term | term
	{
		Head: "expr",
		Body: []Symbol{"expr", "+", "term"},
		Rule: GenRules.ExprPlus,
	},
	{
		Head: "expr",
		Body: []Symbol{"expr", "-", "term"},
		Rule: GenRules.ExprMinus,
	},
	{
		Head: "expr",
		Body: []Symbol{"term"},
		Rule: GenRules.ExprTerm,
	},
	// term → term*unary | term/unary | unary
	{
		Head: "term",
		Body: []Symbol{"term", "*", "unary"},
		Rule: GenRules.TermMult,
	},
	{
		Head: "term",
		Body: []Symbol{"term", "/", "unary"},
		Rule: GenRules.TermDiv,
	},
	{
		Head: "term",
		Body: []Symbol{"unary"},
		Rule: GenRules.TermUnary,
	},
	// unary → !unary | -unary | factor
	{
		Head: "unary",
		Body: []Symbol{"!", "unary"},
		Rule: GenRules.UnaryNot,
	},
	{
		Head: "unary",
		Body: []Symbol{"-", "unary"},
		Rule: GenRules.UnaryNeg,
	},
	{
		Head: "unary",
		Body: []Symbol{"factor"},
		Rule: GenRules.UnaryFactor,
	},
	// factor → (bool) | loc | num | real | true | false
	{
		Head: "factor",
		Body: []Symbol{"(", "bool", ")"},
		Rule: GenRules.FactorBool,
	},
	{
		Head: "factor",
		Body: []Symbol{"loc"},
		Rule: GenRules.FactorLoc,
	},
	{
		Head: "factor",
		Body: []Symbol{"num"},
		Rule: GenRules.FactorNum,
	},
	{
		Head: "factor",
		Body: []Symbol{"real"},
		Rule: GenRules.FactorReal,
	},
	{
		Head: "factor",
		Body: []Symbol{"true"},
		Rule: GenRules.FactorTrue,
	},
	{
		Head: "factor",
		Body: []Symbol{"false"},
		Rule: GenRules.FactorFalse,
	},
}
