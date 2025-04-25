package parser

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
		p.BuildTable()
	}
}
