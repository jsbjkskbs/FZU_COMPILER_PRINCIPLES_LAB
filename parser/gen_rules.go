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
	Decl:      Declare,
	TypeArray: TypeArray,
	TypeBasic: TypeBasic,

	UnaryFactor: UnaryFactor,
	UnaryNeg:    UnaryNeg,
	UnaryNot:    UnaryNot,
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
	if v == nil {
		return fmt.Errorf("no tokens available for declaration")
	}
	return w.SymbolTable.Register(&SymbolTableItem{
		Variable: v.Val,
		Type:     w.Environment.CurrentType,
		Line:     v.Line,
		Pos:      v.Pos,

		UnderlyingType: w.Environment.CurrentDataType.ToString(),
		VariableSize:   w.Environment.CurrentDataSize,
		ArraySize:      w.Environment.CurrentArraySize,
	})
}

func TypeArray(w *Walker) error {
	num, _ := w.Tokens.PeekAtK(1)
	basicType, _ := w.Tokens.PeekAtK(3)
	if num == nil || basicType == nil {
		return fmt.Errorf("no tokens available for type array")
	}
	size, err := strconv.Atoi(num.Val)
	if err != nil {
		return fmt.Errorf("invalid array size: %s", num.Val)
	}
	if size <= 0 {
		return fmt.Errorf("invalid array size: %d", size)
	}
	w.Environment.CurrentArraySize = size
	w.Environment.CurrentType = SymbolTableItemTypeArray
	w.Environment.CurrentDataType = basicType.SpecificType()
	w.Environment.CurrentDataSize = basicType.AllocSize()
	return nil
}

func TypeBasic(w *Walker) error {
	basicType, _ := w.Tokens.PeekAtK(0)
	if basicType == nil {
		return fmt.Errorf("no tokens available for basic type")
	}
	w.Environment.CurrentType = SymbolTableItemTypeVariable
	w.Environment.CurrentDataType = basicType.SpecificType()
	w.Environment.CurrentDataSize = basicType.AllocSize()
	return nil
}

func UnaryNeg(w *Walker) error {
	neg, _ := w.Tokens.PeekAtK(0)
	if neg == nil {
		return fmt.Errorf("no tokens available for unary negation")
	}
	switch w.Environment.CurrentUnary.(type) {
	case int, int8, int16, int32, int64:
		w.Environment.CurrentUnary = -w.Environment.CurrentUnary.(int)
	case float32, float64:
		w.Environment.CurrentUnary = -w.Environment.CurrentUnary.(float64)
	case bool:
		w.Environment.CurrentUnary = !w.Environment.CurrentUnary.(bool)
	default:
		return fmt.Errorf("unsupported unary negation type: %T", w.Environment.CurrentUnary)
	}
	return nil
}

func UnaryNot(w *Walker) error {
	switch w.Environment.CurrentUnary.(type) {
	case bool:
		w.Environment.CurrentUnary = !w.Environment.CurrentUnary.(bool)
	default:
		return fmt.Errorf("unsupported unary not type: %T", w.Environment.CurrentUnary)
	}
	return nil
}

func UnaryFactor(w *Walker) error {
	factor, _ := w.Tokens.PeekAtK(0)
	if factor == nil {
		return fmt.Errorf("no tokens available for unary factor")
	}
	switch factor.SpecificType() {
	case lexer.ConstantInt:
		i, _ := strconv.Atoi(factor.Val)
		w.Environment.CurrentUnary = i
	case lexer.ConstantFloat:
		f, _ := strconv.ParseFloat(factor.Val, 64)
		w.Environment.CurrentUnary = f
	case lexer.ConstantBoolTrue, lexer.ConstantBoolFalse:
		b, _ := strconv.ParseBool(factor.Val)
		w.Environment.CurrentUnary = b
	default:
		return fmt.Errorf("unsupported unary factor type: %s", factor.SpecificType().ToString())
	}
	return nil
}
