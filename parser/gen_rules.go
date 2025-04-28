package parser

import (
	"fmt"

	"app/lexer"
)

var GenRules = struct {
	AugmentedProduction                                   Rule
	Program                                               Rule
	BlockDeclsStmts, BlockDecls, BlockStmts, BlockEpsilon Rule
	Decls, DeclsEpsilon                                   Rule
	Decl                                                  Rule
	TypeArray, TypeBasic                                  Rule
	Stmts, StmtsEpsilon                                   Rule
	StmtMatchedStmt, StmtUnmatchedStmt, StmtDecls         Rule
	UnmatchedStmtIf, UnmatchedStmtElse                    Rule
	MatchedStmtAssign, MatchedStmtIf, MatchedStmtWhile    Rule
	MatchedStmtDo, MatchedStmtBreak, MatchedStmtBlock     Rule
	LocArray, LocId                                       Rule
	Bool, BoolJoin                                        Rule
	Join, JoinEquality                                    Rule
	Equality, NotEquality, EqualityRelational             Rule
	RelationalLess, RelationalGreater                     Rule
	RelationalLessEqual, RelationalGreaterEqual           Rule
	RelationalExpr                                        Rule
	ExprPlus, ExprMinus, ExprTerm                         Rule
	TermMult, TermDiv, TermUnary                          Rule
	UnaryNot, UnaryNeg, UnaryFactor                       Rule
	FactorBool, FactorLoc, FactorNum                      Rule
	FactorTrue, FactorFalse                               Rule
}{
	Decl: func(w *Walker) error {
		debugPrintWhenRuleTriggered(w)
		return Declare(w)
	},
}

func debugPrintWhenRuleTriggered(w *Walker) error {
	fmt.Println("Rule triggered")
	fmt.Println("Current states:", w.States)
	fmt.Println("Current symbols: ", w.Symbols)
	fmt.Println("Current tokens: ")
	w.Tokens.Foreach(func(t *lexer.Token) {
		fmt.Printf("Token: %s, Val: %s, Pos: (%d line, %d col), Extra: %s, %v\n", t.String(), t.Val, t.Line, t.Pos, t.SpecificType().ToString(), t.AllocSize())
	})
	return nil
}

type GenRuleTemplate struct{}

var GenRuleTemplates = GenRuleTemplate{}

func (g *GenRuleTemplate) NOP() Rule {
	return func(w *Walker) error { return nil }
}

func Declare(w *Walker) error {
	v, _ := w.Tokens.PeekAtK(1)
	t, _ := w.Tokens.PeekAtK(2)
	if v == nil || t == nil {
		return fmt.Errorf("no tokens available for declaration")
	}
	return w.SymbolTable.Register(&SymbolTableItem{
		Variable: v.Val,
		Type:     SymbolTableItemTypeVariable,
		Line:     v.Line,
		Pos:      v.Pos,

		UnderlyingType: t.SpecificType().ToString(),
		VariableSize:   t.AllocSize(),
	})
}
