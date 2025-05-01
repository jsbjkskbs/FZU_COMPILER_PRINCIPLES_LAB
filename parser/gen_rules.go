package parser

import (
	"fmt"
	"strconv"
	"strings"

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
	Bool, BoolPrime, BoolPrimeJoin                        Rule
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
	Program:                Program,
	BlockDeclsStmts:        BlockDeclsStmts,
	BlockDecls:             BlockDecls,
	BlockStmts:             BlockStmts,
	BlockEpsilon:           BlockEpsilon,
	Decls:                  Decls,
	DeclsEpsilon:           DeclsEpsilon,
	Decl:                   Decl,
	TypeArray:              TypeArray,
	TypeBasic:              TypeBasic,
	Stmts:                  Stmts,
	StmtsEpsilon:           StmtsEpsilon,
	StmtMatchedStmt:        StmtMatchedStmt,
	StmtUnmatchedStmt:      StmtUnmatchedStmt,
	StmtDecls:              StmtDecls,
	UnmatchedStmtIf:        UnmatchedStmtIf,
	UnmatchedStmtIfElse:    UnmatchedStmtIfElse,
	MatchedStmtAssign:      MatchedStmtAssign,
	MatchedStmtIf:          MatchedStmtIf,
	MatchedStmtIfElse:      MatchedStmtIfElse,
	MatchedStmtWhile:       MatchedStmtWhile,
	MatchedStmtDoWhile:     MatchedStmtDoWhile,
	MatchedStmtBreak:       MatchedStmtBreak,
	MatchedStmtBlock:       MatchedStmtBlock,
	LocArray:               LocArray,
	LocId:                  LocId,
	Bool:                   Bool,
	BoolPrime:              BoolPrime,
	BoolPrimeJoin:          BoolPrimeJoin,
	Join:                   Join,
	JoinEquality:           JoinEquality,
	Equality:               Equality,
	NotEquality:            NotEquality,
	EqualityRelational:     EqualityRelational,
	RelationalLess:         RelationalLess,
	RelationalGreater:      RelationalGreater,
	RelationalLessEqual:    RelationalLessEqual,
	RelationalGreaterEqual: RelationalGreaterEqual,
	RelationalExpr:         RelationalExpr,
	ExprPlus:               ExprPlus,
	ExprMinus:              ExprMinus,
	ExprTerm:               ExprTerm,
	TermMult:               TermMult,
	TermDiv:                TermDiv,
	TermUnary:              TermUnary,
	UnaryNot:               UnaryNot,
	UnaryNeg:               UnaryNeg,
	UnaryFactor:            UnaryFactor,
	FactorBool:             FactorBool,
	FactorLoc:              FactorLoc,
	FactorNum:              FactorNum,
	FactorReal:             FactorReal,
	FactorTrue:             FactorTrue,
	FactorFalse:            FactorFalse,
}

func debugPrintWhenRuleTriggered(w *Walker) error {
	fmt.Println("Rule triggered")
	fmt.Println("Current states:", w.States)
	fmt.Println("Current symbols: ", w.Symbols)
	return nil
}

// program → block
func Program(w *Walker) error {
	children := w.Tokens.PopTopN(w.Tokens.Size())
	w.Tokens.Push(&ASTNode{
		raw:      joinChildren(children),
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: "program"},
		Children: children,
		Type:     "program",
		Payload:  nil,
	})
	w.Emit("exit", "0")
	// n, _ := w.Tokens.Pop()
	// fmt.Println(n.TreeString(0))
	return nil
}

// block → { decls stmts }
func BlockDeclsStmts(w *Walker) error {
	children := w.Tokens.PopTopN(4)
	w.Tokens.Push(&ASTNode{
		raw:      joinChildren(children),
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: fmt.Sprintf("{%s %s}", children[1].raw, children[2].raw)},
		Children: children,
		Type:     "block-decls-stmts",
		Payload:  "!<block>",
	})
	return nil
}

// block → { decls }
func BlockDecls(w *Walker) error {
	children := w.Tokens.PopTopN(3)
	w.Tokens.Push(&ASTNode{
		raw:      joinChildren(children),
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: fmt.Sprintf("{%s}", children[1].raw)},
		Children: children,
		Type:     "block-decls",
		Payload:  "!<block>",
	})
	return nil
}

// block → { stmts }
func BlockStmts(w *Walker) error {
	children := w.Tokens.PopTopN(3)
	w.Tokens.Push(&ASTNode{
		raw:      joinChildren(children),
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: fmt.Sprintf("{%s}", children[1].raw)},
		Children: children,
		Type:     "block-stmts",
		Payload:  "!<block>",
	})
	return nil
}

// block → { }
func BlockEpsilon(w *Walker) error {
	children := w.Tokens.PopTopN(2)
	w.Tokens.Push(&ASTNode{
		raw:      joinChildren(children),
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: "{}"},
		Children: children,
		Type:     "block-epsilon",
		Payload:  "!<block>",
	})
	w.Emit("nop", "")
	return nil
}

// decls → decls decl
func Decls(w *Walker) error {
	children := w.Tokens.PopTopN(2)
	w.Tokens.Push(&ASTNode{
		raw:      joinChildren(children),
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: "decls"},
		Children: children,
		Type:     "decls",
		Payload:  "!<decl>",
	})
	return nil
}

// decls → ε
func DeclsEpsilon(w *Walker) error {
	w.Environment.LoopLabelStack.Push(w.GetCurrentLabelCount())
	w.Tokens.Push(&ASTNode{
		raw:      "",
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: "decls-epsilon"},
		Children: nil,
		Type:     "decls-epsilon",
		Payload:  "!<decl>",
	})
	return nil
}

// decl → type id;
func Decl(w *Walker) error {
	children := w.Tokens.PopTopN(3)
	t, id := children[0].Token, children[1].Token
	item := &SymbolTableItem{
		Variable:       id.Val,
		VariableSize:   t.AllocSize(),
		Type:           SymbolTableItemTypeVariable,
		UnderlyingType: t.SpecificType().ToString(),
	}
	addr, err := w.SymbolTable.Register(item)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	w.Emit("alloc", fmt.Sprintf("$(%#x)", addr), strconv.Itoa(item.VariableSize), getInitialValue(t))
	w.Tokens.Push(&ASTNode{
		raw:      joinChildren(children),
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: id.Val},
		Children: children,
		Type:     "decl",
		Payload:  "!<decl>",
	})
	return nil
}

// type → type [ num ]
func TypeArray(w *Walker) error {
	children := w.Tokens.PopTopN(4)
	w.Tokens.Push(&ASTNode{
		raw:      joinChildren(children),
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: fmt.Sprintf("%s[%s]", children[0].raw, children[2].raw)},
		Children: children,
		Type:     "type-array",
		Payload:  "!<array>",
	})
	return nil
}

// type → basic
func TypeBasic(w *Walker) error {
	children := w.Tokens.PopTopN(1)
	w.Tokens.Push(&ASTNode{
		raw:      children[0].raw,
		Token:    children[0].Token,
		Children: children,
		Type:     "type-basic",
		Payload:  "!<basic>",
	})
	return nil
}

// stmts → stmts stmt
func Stmts(w *Walker) error {
	children := w.Tokens.PopTopN(2)
	w.Tokens.Push(&ASTNode{
		raw:      joinChildren(children),
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: "stmts"},
		Children: children,
		Type:     "stmts",
		Payload:  "!<stmt>",
	})
	return nil
}

// stmts → ε
func StmtsEpsilon(w *Walker) error {
	w.Tokens.Push(&ASTNode{
		raw:      "",
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: "stmts-epsilon"},
		Children: nil,
		Type:     "stmts-epsilon",
		Payload:  "!<stmt>",
	})
	return nil
}

// stmt → matched_stmt
func StmtMatchedStmt(w *Walker) error {
	children := w.Tokens.PopTopN(1)
	w.Tokens.Push(&ASTNode{
		raw:      children[0].raw,
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: "stmt-matched"},
		Children: children,
		Type:     "stmt-matched",
		Payload:  "!<matched-stmt>",
	})
	return nil
}

// stmt → unmatched_stmt
func StmtUnmatchedStmt(w *Walker) error {
	children := w.Tokens.PopTopN(1)
	w.Tokens.Push(&ASTNode{
		raw:      children[0].raw,
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: "stmt-unmatched"},
		Children: children,
		Type:     "stmt-unmatched",
		Payload:  "!<unmatched-stmt>",
	})
	return nil
}

// stmt → decls
func StmtDecls(w *Walker) error {
	children := w.Tokens.PopTopN(1)
	w.Tokens.Push(&ASTNode{
		raw:      children[0].raw,
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: "stmt-decls"},
		Children: children,
		Type:     "stmt-decls",
		Payload:  "!<decls>",
	})
	return nil
}

// unmatched_stmt → if ( bool ) unmatched_stmt
func UnmatchedStmtIf(w *Walker) error {
	children := w.Tokens.PopTopN(5)
	w.Tokens.Push(&ASTNode{
		raw:      joinChildren(children),
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: "unmatched-stmt-if"},
		Children: children,
		Type:     "stmt-if",
		Payload:  "!<if>",
	})
	return nil
}

// unmatched_stmt → if ( bool ) matched_stmt else unmatched_stmt
func UnmatchedStmtIfElse(w *Walker) error {
	children := w.Tokens.PopTopN(7)
	w.Tokens.Push(&ASTNode{
		raw:      joinChildren(children),
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: "unmatched-stmt-if-else"},
		Children: children,
		Type:     "stmt-if-else",
		Payload:  "!<if-else>",
	})
	return nil
}

// matched_stmt → loc = bool ;
func MatchedStmtAssign(w *Walker) error {
	children := w.Tokens.PopTopN(4)
	dist := children[0].Token.Val
	src := children[2].Token.Val
	i, _, err := w.SymbolTable.Lookup(dist)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	w.Emit("mov", fmt.Sprintf("$(%#x)", i.Address), src)

	w.Tokens.Push(&ASTNode{
		raw: joinChildren(children),
		Token: &lexer.Token{
			Type: lexer.EXTRA,
			Val:  "stmt",
		},
		Children: children,
		Type:     "stmt-assign",
		Payload:  "!copy(!dist:!src)",
	})
	return nil
}

// matched_stmt → if ( bool ) matched_stmt else matched_stmt
func MatchedStmtIfElse(w *Walker) error {
	prevEl, _ := w.Tokens.PeekAtK(7)
	children := w.Tokens.PopTopN(7)
	w.Tokens.Push(&ASTNode{
		raw:      joinChildren(children),
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: "matched-stmt-if-else"},
		Children: children,
		Type:     "stmt-if-else",
		Payload:  "!<if-else>",
	})
	if prevEl.Token.SpecificType() != lexer.ReservedWordElse {
		n := w.Environment.LabelStack.PopTopN(2)
		m := w.Environment.EndIfStmtStack.PopTopN(2)
		w.EmitLabel(n[1], fmt.Sprintf("L%d", m[0]+1), "jmp")
		w.EmitGoto(m[0], w.GetCurrentLabelCount())
		w.EmitGoto(m[1], w.GetCurrentLabelCount())
	} else {
		n := w.Environment.LabelStack.PopTopN(2)
		m := w.Environment.EndIfStmtStack.PopTopN(min(2, w.Environment.EndIfStmtStack.Size()))
		w.EmitLabel(n[1], fmt.Sprintf("L%d", m[0]+1), "jmp")
		// only fill the first block
		// in case of `if (condition) { ... } else if (condition) { ... } else { ... }`
		// we should delegate the next endif to the next block
		w.EmitGoto(m[0], w.GetCurrentLabelCount())
		if len(m) > 1 {
			w.Environment.EndIfStmtStack.Push(m[1])
		}
	}
	return nil
}

// matched_stmt → if ( bool ) matched_stmt
func MatchedStmtIf(w *Walker) error {
	prevEl, _ := w.Tokens.PeekAtK(5)
	children := w.Tokens.PopTopN(5)
	w.Tokens.Push(&ASTNode{
		raw:      joinChildren(children),
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: "matched-stmt-if"},
		Children: children,
		Type:     "stmt-if",
		Payload:  "!<if>",
	})
	n := w.Environment.LabelStack.PopTopN(2)
	m := w.Environment.EndIfStmtStack.PopTopN(1)
	w.EmitLabel(n[1], fmt.Sprintf("L%d", m[0]+1), "jmp")
	w.EmitGoto(m[0], w.GetCurrentLabelCount())
	if prevEl.Token.SpecificType() == lexer.ReservedWordElse {
		w.Environment.EndIfStmtStack.Push(m[0])
	}
	return nil
}

// matched_stmt → while ( bool ) stmt
func MatchedStmtWhile(w *Walker) error {
	children := w.Tokens.PopTopN(5)
	w.Tokens.Push(&ASTNode{
		raw:      joinChildren(children),
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: "matched-stmt-while"},
		Children: children,
		Type:     "stmt-while",
		Payload:  "!<while>",
	})
	n := w.Environment.LabelStack.PopTopN(2)
	m := w.Environment.EndIfStmtStack.PopTopN(1)
	w.EmitLabel(n[1], fmt.Sprintf("L%d", m[0]+1), "jmp")
	w.EmitGoto(m[0], n[0])
	return nil
}

// matched_stmt → do stmt while ( bool ) ;
func MatchedStmtDoWhile(w *Walker) error {
	children := w.Tokens.PopTopN(7)
	w.Tokens.Push(&ASTNode{
		raw:      joinChildren(children),
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: "matched-stmt-do-while"},
		Children: children,
		Type:     "stmt-do-while",
		Payload:  "!<do-while>",
	})
	n := w.Environment.LabelStack.PopTopN(2)
	m := w.Environment.LoopLabelStack.PopTopN(1)
	w.AdjustJMP(n[0], m[0])
	w.EmitLabel(n[1], fmt.Sprintf("L%d", n[1]+1), "jmp")
	return nil
}

// matched_stmt → break ;
func MatchedStmtBreak(w *Walker) error {
	children := w.Tokens.PopTopN(2)
	w.Tokens.Push(&ASTNode{
		raw:      children[0].raw,
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: "matched-stmt-break"},
		Children: children,
		Type:     "stmt-break",
		Payload:  "!<break>",
	})
	return nil
}

// matched_stmt → block
func MatchedStmtBlock(w *Walker) error {
	matchedStmtBlockIfWhileElse(w)
	children := w.Tokens.PopTopN(1)
	w.Tokens.Push(&ASTNode{
		raw:      children[0].raw,
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: "matched-stmt-block"},
		Children: children,
		Type:     "stmt-block",
		Payload:  "!<block>",
	})
	return nil
}

func matchedStmtBlockIfWhileElse(w *Walker) {
	doelse, _ := w.Tokens.PeekAtK(1)
	if doelse == nil {
		return
	}
	if doelse.Token.SpecificType() != lexer.ReservedWordDo {
		w.Environment.LoopLabelStack.PopTopN(1)
	}
	if doelse.Token.SpecificType() == lexer.ReservedWordElse {
		w.NewGotoLabel()
		return
	}
	ifwhile, _ := w.Tokens.PeekAtK(4)
	if ifwhile == nil {
		return
	}
	if ifwhile.Token.SpecificType() == lexer.ReservedWordIf || ifwhile.Token.SpecificType() == lexer.ReservedWordWhile {
		w.NewGotoLabel()
		return
	}
}

// loc → loc [ num ]
func LocArray(w *Walker) error {
	children := w.Tokens.PopTopN(4)
	w.Tokens.Push(&ASTNode{
		raw: joinChildren(children),
		Token: &lexer.Token{
			Type: lexer.EXTRA,
			Val:  fmt.Sprintf("%s[%s]", children[0].raw, children[2].raw),
		},
		Children: children,
		Type:     "loc-array",
		Payload:  "!dist:!ptr(size=4)",
	})
	return nil
}

// loc → id
func LocId(w *Walker) error {
	children := w.Tokens.PopTopN(1)
	w.Tokens.Push(&ASTNode{
		raw:      children[0].raw,
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: children[0].Token.Val},
		Children: children,
		Type:     "loc-id",
		Payload:  "!dist:!ptr(size=4)",
	})
	return nil
}

// bool → bool'
func Bool(w *Walker) error {
	children := w.Tokens.PopTopN(1)
	prev, _ := w.Tokens.Peek()
	if prev.Token.SpecificType() != lexer.OperatorAssignment {
		resultAddr := w.SymbolTable.TempAddr(4)
		resultStr := fmt.Sprintf("$(%#x)", resultAddr)
		w.Emit("cmp", resultStr, children[0].Token.Val, "0")
		w.Tokens.Push(&ASTNode{
			raw:      children[0].raw,
			Token:    &lexer.Token{Type: lexer.EXTRA, Val: resultStr},
			Children: children,
			Type:     "bool",
			Payload:  "!<bool'>",
		})
		boolLookbackIfWhile(w)
	} else {
		w.Tokens.Push(&ASTNode{
			raw:      children[0].raw,
			Token:    &lexer.Token{Type: lexer.EXTRA, Val: children[0].Token.Val},
			Children: children,
			Type:     "bool",
			Payload:  "!<bool'>",
		})
	}
	return nil
}

func boolLookbackIfWhile(w *Walker) {
	ifwhile, _ := w.Tokens.PeekAtK(2)
	if ifwhile == nil {
		return
	}
	if ifwhile.Token.SpecificType() == lexer.ReservedWordIf || ifwhile.Token.SpecificType() == lexer.ReservedWordWhile {
		jnz := w.NewLabel()
		w.NewLabel()
		top, _ := w.Tokens.Peek()
		w.EmitLabel(jnz, fmt.Sprintf("L%d", jnz+2), "jnz", top.Token.Val)
	}
}

// bool' → bool' || join
func BoolPrime(w *Walker) error {
	result := w.SymbolTable.TempAddr(4)
	children := w.Tokens.PopTopN(3)
	resultStr := fmt.Sprintf("$(%#x)", result)
	w.Tokens.Push(&ASTNode{
		raw: joinChildren(children),
		Token: &lexer.Token{
			Type: lexer.EXTRA,
			Val:  resultStr,
		},
		Children: children,
		Type:     "bool-prime",
		Payload:  "!dist:!ptr(size=4)",
	})
	w.Emit("or", resultStr, children[0].Token.Val, children[2].Token.Val)
	return nil
}

// bool' → join
func BoolPrimeJoin(w *Walker) error {
	children := w.Tokens.PopTopN(1)
	w.Tokens.Push(&ASTNode{
		raw:      children[0].raw,
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: children[0].Token.Val},
		Children: children,
		Type:     "bool-prime-join",
		Payload:  "!<join>",
	})
	return nil
}

// join → join && equality
func Join(w *Walker) error {
	result := w.SymbolTable.TempAddr(4)
	children := w.Tokens.PopTopN(3)
	resultStr := fmt.Sprintf("$(%#x)", result)
	w.Tokens.Push(&ASTNode{
		raw: joinChildren(children),
		Token: &lexer.Token{
			Type: lexer.EXTRA,
			Val:  resultStr,
		},
		Children: children,
		Type:     "join",
		Payload:  "!dist:!ptr(size=4)",
	})
	w.Emit("and", resultStr, children[0].Token.Val, children[2].Token.Val)
	return nil
}

// join → equality
func JoinEquality(w *Walker) error {
	children := w.Tokens.PopTopN(1)
	w.Tokens.Push(&ASTNode{
		raw:      children[0].raw,
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: children[0].Token.Val},
		Children: children,
		Type:     "join-equality",
		Payload:  "!<equality>",
	})
	return nil
}

// equality → equality == rel
func Equality(w *Walker) error {
	result := w.SymbolTable.TempAddr(4)
	children := w.Tokens.PopTopN(3)
	resultStr := fmt.Sprintf("$(%#x)", result)
	w.Tokens.Push(&ASTNode{
		raw: joinChildren(children),
		Token: &lexer.Token{
			Type: lexer.EXTRA,
			Val:  resultStr,
		},
		Children: children,

		Type:    "equality",
		Payload: "!dist:!ptr(size=4)",
	})
	w.Emit("eq", children[0].Token.Val, children[2].Token.Val)
	return nil
}

// equality → equality != rel
func NotEquality(w *Walker) error {
	result := w.SymbolTable.TempAddr(4)
	children := w.Tokens.PopTopN(3)
	resultStr := fmt.Sprintf("$(%#x)", result)
	w.Tokens.Push(&ASTNode{
		raw: joinChildren(children),
		Token: &lexer.Token{
			Type: lexer.EXTRA,
			Val:  resultStr,
		},
		Children: children,
		Type:     "not-equality",
		Payload:  "!dist:!ptr(size=4)",
	})
	w.Emit("ne", children[0].Token.Val, children[2].Token.Val)
	return nil
}

// equality → rel
func EqualityRelational(w *Walker) error {
	children := w.Tokens.PopTopN(1)
	w.Tokens.Push(&ASTNode{
		raw:      children[0].raw,
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: children[0].Token.Val},
		Children: children,
		Type:     "equality-relational",
		Payload:  "!<rel>",
	})
	return nil
}

// rel → expr < expr
func RelationalLess(w *Walker) error {
	result := w.SymbolTable.TempAddr(4)
	children := w.Tokens.PopTopN(3)
	resultStr := fmt.Sprintf("$(%#x)", result)
	w.Tokens.Push(&ASTNode{
		raw: joinChildren(children),
		Token: &lexer.Token{
			Type: lexer.EXTRA,
			Val:  resultStr,
		},
		Children: children,
		Type:     "less",
		Payload:  "!dist:!ptr(size=4)",
	})
	w.Emit("ls", children[0].Token.Val, children[2].Token.Val)
	return nil
}

// rel → expr > expr
func RelationalGreater(w *Walker) error {
	result := w.SymbolTable.TempAddr(4)
	children := w.Tokens.PopTopN(3)
	resultStr := fmt.Sprintf("$(%#x)", result)
	w.Tokens.Push(&ASTNode{
		raw: joinChildren(children),
		Token: &lexer.Token{
			Type: lexer.EXTRA,
			Val:  resultStr,
		},
		Children: children,
		Type:     "greater",
		Payload:  "!dist:!ptr(size=4)",
	})
	w.Emit("gt", children[0].Token.Val, children[2].Token.Val)
	return nil
}

// rel → expr <= expr
func RelationalLessEqual(w *Walker) error {
	result := w.SymbolTable.TempAddr(4)
	children := w.Tokens.PopTopN(3)
	resultStr := fmt.Sprintf("$(%#x)", result)
	w.Tokens.Push(&ASTNode{
		raw: joinChildren(children),
		Token: &lexer.Token{
			Type: lexer.EXTRA,
			Val:  resultStr,
		},
		Children: children,
		Type:     "less-equal",
		Payload:  "!dist:!ptr(size=4)",
	})
	w.Emit("le", children[0].Token.Val, children[2].Token.Val)
	return nil
}

// rel → expr >= expr
func RelationalGreaterEqual(w *Walker) error {
	result := w.SymbolTable.TempAddr(4)
	children := w.Tokens.PopTopN(3)
	resultStr := fmt.Sprintf("$(%#x)", result)
	w.Tokens.Push(&ASTNode{
		raw: joinChildren(children),
		Token: &lexer.Token{
			Type: lexer.EXTRA,
			Val:  resultStr,
		},
		Children: children,
		Type:     "greater-equal",
		Payload:  "!dist:!ptr(size=4)",
	})
	w.Emit("ge", children[0].Token.Val, children[2].Token.Val)
	return nil
}

// rel → expr
func RelationalExpr(w *Walker) error {
	children := w.Tokens.PopTopN(1)
	w.Tokens.Push(&ASTNode{
		raw:      children[0].raw,
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: children[0].Token.Val},
		Children: children,
		Type:     "rel-expr",
		Payload:  "!<expr>",
	})
	return nil
}

// expr -> expr + term
func ExprPlus(w *Walker) error {
	result := w.SymbolTable.TempAddr(4)
	children := w.Tokens.PopTopN(3)
	resultStr := fmt.Sprintf("$(%#x)", result)
	w.Tokens.Push(&ASTNode{
		raw: joinChildren(children),
		Token: &lexer.Token{
			Type: lexer.EXTRA,
			Val:  resultStr,
		},
		Children: children,
		Type:     "plus",
		Payload:  "!dist:!ptr(size=4)",
	})
	w.Emit("add", resultStr, children[0].Token.Val, children[2].Token.Val)
	return nil
}

// expr -> expr - term
func ExprMinus(w *Walker) error {
	result := w.SymbolTable.TempAddr(4)
	children := w.Tokens.PopTopN(3)
	resultStr := fmt.Sprintf("$(%#x)", result)
	w.Tokens.Push(&ASTNode{
		raw: joinChildren(children),
		Token: &lexer.Token{
			Type: lexer.EXTRA,
			Val:  resultStr,
		},
		Children: children,
		Type:     "minus",
		Payload:  "!dist:!ptr(size=4)",
	})
	w.Emit("sub", resultStr, children[0].Token.Val, children[2].Token.Val)
	return nil
}

// expr -> term
func ExprTerm(w *Walker) error {
	children := w.Tokens.PopTopN(1)
	w.Tokens.Push(&ASTNode{
		raw:      children[0].raw,
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: children[0].Token.Val},
		Children: children,
		Type:     "expr-term",
		Payload:  "!<term>",
	})
	return nil
}

// term -> term * unary
func TermMult(w *Walker) error {
	result := w.SymbolTable.TempAddr(4)
	children := w.Tokens.PopTopN(3)
	resultStr := fmt.Sprintf("$(%#x)", result)
	w.Tokens.Push(&ASTNode{
		raw: joinChildren(children),
		Token: &lexer.Token{
			Type: lexer.EXTRA,
			Val:  resultStr,
		},
		Children: children,
		Type:     "mult",
		Payload:  "!dist:!ptr(size=4)",
	})
	w.Emit("mul", resultStr, children[0].Token.Val, children[2].Token.Val)
	return nil
}

// term -> term / unary
func TermDiv(w *Walker) error {
	result := w.SymbolTable.TempAddr(4)
	children := w.Tokens.PopTopN(3)
	resultStr := fmt.Sprintf("$(%#x)", result)
	w.Tokens.Push(&ASTNode{
		raw: joinChildren(children),
		Token: &lexer.Token{
			Type: lexer.EXTRA,
			Val:  resultStr,
		},
		Children: children,
		Type:     "div",
		Payload:  "!dist:!ptr(size=4)",
	})
	w.Emit("div", resultStr, children[0].Token.Val, children[2].Token.Val)
	return nil
}

// term -> unary
func TermUnary(w *Walker) error {
	children := w.Tokens.PopTopN(1)
	w.Tokens.Push(&ASTNode{
		raw:      children[0].raw,
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: children[0].Token.Val},
		Children: children,
		Type:     "term-unary",
		Payload:  "!<unary>",
	})
	return nil
}

// unary -> -unary
func UnaryNeg(w *Walker) error {
	result := w.SymbolTable.TempAddr(4)
	addr := fmt.Sprintf("$(%#x)", result)
	children := w.Tokens.PopTopN(2)
	w.Tokens.Push(&ASTNode{
		raw: joinChildren(children),
		Token: &lexer.Token{
			Type: lexer.EXTRA,
			Val:  fmt.Sprintf("$(%#x)", result),
		},
		Children: children,
		Type:     "neg",
		Payload:  "!dist:!ptr(size=4)",
	})
	w.Emit(addr, "neg", children[1].Token.Val)
	return nil
}

// unary -> !unary
func UnaryNot(w *Walker) error {
	result := w.SymbolTable.TempAddr(4)
	addr := fmt.Sprintf("$(%#x)", result)
	children := w.Tokens.PopTopN(2)
	w.Tokens.Push(&ASTNode{
		raw: joinChildren(children),
		Token: &lexer.Token{
			Type: lexer.EXTRA,
			Val:  fmt.Sprintf("$(%#x)", result),
		},
		Children: children,
		Type:     "not",
		Payload:  "!dist:!ptr(size=4)",
	})
	w.Emit(addr, "not", children[1].Token.Val)
	return nil
}

// unary -> factor
func UnaryFactor(w *Walker) error {
	children := w.Tokens.PopTopN(1)
	w.Tokens.Push(&ASTNode{
		raw:      children[0].raw,
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: children[0].Token.Val},
		Children: children,
		Type:     "unary-factor",
		Payload:  "!<factor>",
	})
	return nil
}

// factor → ( bool )
func FactorBool(w *Walker) error {
	children := w.Tokens.PopTopN(3)
	w.Tokens.Push(&ASTNode{
		raw:      joinChildren(children),
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: children[1].Token.Val},
		Children: children,
		Type:     "factor-bool",
		Payload:  "!<bool>",
	})
	return nil
}

// factor → loc
func FactorLoc(w *Walker) error {
	children := w.Tokens.PopTopN(1)
	w.Tokens.Push(&ASTNode{
		raw:      children[0].raw,
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: children[0].Token.Val},
		Children: children,
		Type:     "factor-loc",
		Payload:  "!<loc>",
	})
	return nil
}

// factor → num
func FactorNum(w *Walker) error {
	children := w.Tokens.PopTopN(1)
	w.Tokens.Push(&ASTNode{
		raw:      children[0].raw,
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: children[0].raw},
		Children: children,
		Type:     "factor-num",
		Payload:  "!const(size=4)",
	})
	return nil
}

// factor → real
func FactorReal(w *Walker) error {
	children := w.Tokens.PopTopN(1)
	w.Tokens.Push(&ASTNode{
		raw:      children[0].raw,
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: children[0].raw},
		Children: children,
		Type:     "factor-real",
		Payload:  "!const(size=8)",
	})
	return nil
}

// factor → true
func FactorTrue(w *Walker) error {
	chidren := w.Tokens.PopTopN(1)
	w.Tokens.Push(&ASTNode{
		raw:      "true",
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: "1"},
		Children: chidren,
		Type:     "factor-true",
		Payload:  "!const(size=1)",
	})
	return nil
}

// factor → false
func FactorFalse(w *Walker) error {
	chidren := w.Tokens.PopTopN(1)
	w.Tokens.Push(&ASTNode{
		raw:      "false",
		Token:    &lexer.Token{Type: lexer.EXTRA, Val: "0"},
		Children: chidren,
		Type:     "factor-false",
		Payload:  "!const(size=1)",
	})
	return nil
}

func joinChildren(children []*ASTNode) string {
	res := []string{}
	for _, child := range children {
		res = append(res, child.raw)
	}
	return strings.Join(res, " ")
}

func getInitialValue(token *lexer.Token) string {
	switch token.SpecificType() {
	case lexer.TypeInt8, lexer.TypeUnsignedInt8, lexer.TypeByte, lexer.TypeBool:
		return "0"
	case lexer.TypeInt16, lexer.TypeUnsignedInt16:
		return "0"
	case lexer.TypeInt32, lexer.TypeUnsignedInt32, lexer.TypeInt, lexer.TypeUnsignedInt:
		return "0"
	case lexer.TypeInt64, lexer.TypeUnsignedInt64:
		return "0"
	case lexer.TypeFloat32, lexer.TypeFloat:
		return "0.0f"
	case lexer.TypeFloat64:
		return "0.0"
	}
	return "<nullptr>"
}
