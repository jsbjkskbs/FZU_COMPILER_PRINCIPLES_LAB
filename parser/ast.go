package parser

import (
	"app/lexer"
)

type AbstractSyntaxTree struct {
	Root *ASTNode
}

type ASTNode struct {
	raw   string
	Token *lexer.Token

	Children []*ASTNode

	Type    Symbol // Type of the node (e.g., statement, expression, declaration, etc.)
	Payload any    // Additional data associated with the node (e.g., variable name, value, etc.)
}

type ASTNodeType int

func (p *Parser) Token2ASTNode(token *lexer.Token) *ASTNode {
	return &ASTNode{
		raw:      token.Val,
		Token:    token,
		Children: []*ASTNode{},
		Type:     p.Reflect(token),
		Payload:  nil,
	}
}
