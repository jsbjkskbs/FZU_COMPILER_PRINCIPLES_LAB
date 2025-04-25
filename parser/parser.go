package parser

import (
	"fmt"
	"slices"

	"app/lexer"
	. "app/parser/production"
	"app/utils/log"
)

func (p *Parser) Reflect(token lexer.Token) Symbol {
	switch token.Type {
	case lexer.INTEGER:
		return "num"
	case lexer.FLOAT:
		return "real"
	case lexer.IDENTIFIER:
		return "id"
	case lexer.TYPE:
		return "basic"
	case lexer.EOF:
		return TERMINATE
	default:
		return Symbol(token.Val)
	}
}

func (p *Parser) EnsureFirstSet() {
	if len(p.FirstSet) == 0 {
		p.BuildFirstSet()
	}
}

func (p *Parser) EnsureSymbols() {
	if len(p.Symbols) == 0 {
		p.BuildSymbols()
	}
}

func (p *Parser) EnsureStates() {
	if len(p.States) == 0 {
		p.BuildStates()
	}
}

func (p *Parser) EnsureTable() {
	if p.Table == nil {
		p.OptimizedHeadsCheck()
		p.BuildTable()
	}
}

func (p *Parser) OptimizedHeadsCheck() {
	for _, production := range p.Grammar.Productions {
		if OptimizedSymbols.Contains(production.Head) || slices.ContainsFunc(production.Body, func(symbol Symbol) bool {
			return OptimizedSymbols.Contains(symbol)
		}) {
			fmt.Println(log.Sprintf(log.Argument{FrontColor: log.Yellow, Highlight: true, Format: "Warning: Optimized symbols may cause reduce-reduce conflict", Args: []any{}}))
			if OptimizedSymbols.Contains(production.Head) {
				fmt.Print(log.Sprintf(log.Argument{FrontColor: log.Red, Underline: true, Format: "%v", Args: []any{production.Head}}))
			} else {
				fmt.Print(log.Sprintf(log.Argument{FrontColor: log.Red, Format: "%v", Args: []any{production.Head}}))
			}
			fmt.Print(" -> ")
			for _, symbol := range production.Body {
				if OptimizedSymbols.Contains(symbol) {
					fmt.Print(" ", log.Sprintf(log.Argument{FrontColor: log.Red, Underline: true, Format: "%v", Args: []any{symbol}}))
				} else {
					fmt.Print(" ", log.Sprintf(log.Argument{FrontColor: log.Red, Format: "%v", Args: []any{symbol}}))
				}
			}
			fmt.Println()
		}
	}
}
