package parser

import (
	"errors"
	"fmt"
	"io"
	"slices"

	"app/lexer"
	"app/utils/log"
)

func (p *Parser) Parse(l *lexer.Lexer, logger func(string)) {
	walker := p.NewWalker()
	walker.SymbolTable.EnterScope()
	for {
		token, err := l.NextToken()
		if err != nil && !errors.Is(err, io.EOF) {
			logger(fmt.Sprintf("Error: %v", err))
			return
		}

		if errors.Is(err, io.EOF) {
			token.Type = lexer.EOF
		}
		symbol := p.Reflect(token)
		if token.SpecificType() == lexer.DelimiterLeftBrace {
			walker.SymbolTable.EnterScope()
		}

		for {
			logger(fmt.Sprintf("State: %v\nSymbols: %v\nSymbol: %s\n", walker.States, walker.Symbols, symbol))
			action, err := walker.Next(symbol)
			if err != nil {
				logger(fmt.Sprintf("Error: %v", err))
				return
			}
			logger(fmt.Sprintf("Token: (%s, %s), Action: %v\n\n", token.Type.ToString(), token.Val, action))
			if action.Type != REDUCE {
				break
			}
		}

		if token.SpecificType() == lexer.DelimiterRightBrace {
			walker.SymbolTable.ExitScope()
		}

		if symbol == TERMINATE {
			logger("Parsing completed successfully.")
			break
		}

		walker.Tokens.Push(&token)
	}

	logger("Symbol Table:")
	scopes := walker.SymbolTable.LegacyScopes[1:]
	for _, scope := range scopes {
		if scope == nil {
			continue
		}
		logger(fmt.Sprintf("\n\nScope[%d]: \n", scope.ID))
		logger(fmt.Sprintf("  Level: %d\n", scope.Level))
		logger(fmt.Sprintln("  Symbols:"))
		for _, symbol := range scope.Items {
			if symbol == nil {
				continue
			}
			logger(fmt.Sprintf("    0x%x -> %v:%v[alloc=%d] << at line %d, pos %d\n",
				symbol.Address, symbol.Variable, symbol.UnderlyingType, symbol.VariableSize, symbol.Line, symbol.Pos))
		}
	}
}

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
