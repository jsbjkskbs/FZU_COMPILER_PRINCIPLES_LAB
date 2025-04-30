package parser

import (
	"fmt"
	"strconv"

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
	UnmatchedStmtIf, UnmatchedStmtIfElse                  Rule
	MatchedStmtAssign, MatchedStmtIf, MatchedStmtIfElse   Rule
	MatchedStmtWhile, MatchedStmtDoWhile                  Rule
	MatchedStmtBreak, MatchedStmtBlock                    Rule
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
	FactorBool, FactorLoc, FactorNum, FactorReal          Rule
	FactorTrue, FactorFalse                               Rule
}{
	// MatchedStmtIf: debugPrintWhenRuleTriggered,
	Equality:    Equality,
	NotEquality: NotEquality,
}

func debugPrintWhenRuleTriggered(w *Walker) error {
	fmt.Println("Rule triggered")
	fmt.Println("Current states:", w.States)
	fmt.Println("Current symbols: ", w.Symbols)
	return nil
}

type GenRuleTemplate struct{}

var GenRuleTemplates = GenRuleTemplate{}

func (g *GenRuleTemplate) NOP() Rule {
	return func(w *Walker) error { return nil }
}

func Equality(w *Walker) error {
	arg1, _ := w.Tokens.PeekAtK(2)
	arg2, _ := w.Tokens.PeekAtK(0)
	result := w.SymbolTable.TempAddr(4)
	w.Emit(strconv.Itoa(result), "eq", arg1, arg2)
	children := w.Tokens.PopTopN(3)
	w.Tokens.Push(&ASTNode{
		raw: fmt.Sprintf("%s == %s", arg1.Token.Val, arg2.Token.Val),
		Token: &lexer.Token{
			Type: lexer.EXTRA,
			Val:  fmt.Sprintf("$(0x%x)", result),
		},
		Children: children,

		Type:    "equality",
		Payload: "!dist:!ptr(size=4)",
	})
	n, _ := w.Tokens.Peek()
	fmt.Printf("emit($(%#x), eq, %s, %s)\n", result, arg1.Token.Val, arg2.Token.Val)
	fmt.Println("Node: ", n.raw)
	fmt.Println("Children:")
	for _, child := range n.Children {
		fmt.Println("  -", child.raw)
	}
	fmt.Println("Token:", n.Token.Val)
	fmt.Println("Type:", n.Type)
	fmt.Println("Payload:", n.Payload)
	fmt.Println()

	return nil
}

func NotEquality(w *Walker) error {
	arg1, _ := w.Tokens.PeekAtK(2)
	arg2, _ := w.Tokens.PeekAtK(0)
	result := w.SymbolTable.TempAddr(4)
	w.Emit(strconv.Itoa(result), "ne", arg1, arg2)
	children := w.Tokens.PopTopN(3)
	w.Tokens.Push(&ASTNode{
		raw: fmt.Sprintf("%s != %s", arg1.Token.Val, arg2.Token.Val),
		Token: &lexer.Token{
			Type: lexer.EXTRA,
			Val:  fmt.Sprintf("$(0x%x)", result),
		},
		Children: children,
		Type:     "not-equality",
		Payload:  "!dist:!ptr(size=4)",
	})
	n, _ := w.Tokens.Peek()
	fmt.Printf("emit($(%#x), eq, %s, %s)\n", result, arg1.Token.Val, arg2.Token.Val)
	fmt.Println("Node: ", n.raw)
	fmt.Println("Children:")
	for _, child := range n.Children {
		fmt.Println("  -", child.raw)
	}
	fmt.Println("Token:", n.Token.Val)
	fmt.Println("Type:", n.Type)
	fmt.Println("Payload:", n.Payload)
	fmt.Println()
	return nil
}
