package parser

import (
	"errors"
	"fmt"
	"io"
	"slices"

	"app/lexer"
	"app/utils/log"
)

// Parse is the main function that parses the input tokens using the LR(1) parser algorithm.
// It takes a lexer.Lexer instance and a logger function as arguments.
// The logger function is used to log messages during the parsing process.
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
		symbol := p.Reflect(&token)
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

		walker.Tokens.Push(p.Token2ASTNode(&token))
	}

	logger("\n\nThree Address Code:\n")
	for _, line := range walker.ThreeAddress {
		// fmt.Println(line)
		logger(fmt.Sprintln(line))
	}

	for _, scope := range walker.SymbolTable.LegacyScopes[1:] {
		logger("-------------------------------------\n")
		logger(fmt.Sprintf("Scope: %v\n", scope.ID))
		for _, item := range scope.Items {
			logger(fmt.Sprintf("Variable: %s, Type: %s, Address: %#x\n", item.Variable, item.UnderlyingType, item.Address))
			if item.Type == SymbolTableItemTypeArray {
				logger(fmt.Sprintf("Array Size: %d, Element Size: %d\n", item.ArraySize, item.ArrayElementSize))
			}
		}
	}
}

// Reflect converts a lexer.Token to a Symbol.
// It maps specific token types to corresponding symbols and returns the symbol representation.
func (p *Parser) Reflect(token *lexer.Token) Symbol {
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
	p._mu.Lock()
	defer p._mu.Unlock()
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
