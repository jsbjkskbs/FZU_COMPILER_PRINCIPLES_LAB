package parser

import (
	"fmt"
	"strings"

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
	Payload any
}

type ASTNodeType int

func (a *ASTNode) TreeString(indent int) string {
	result := strings.Repeat("\t", indent)
	result += fmt.Sprintf("Raw: %s | ", a.raw)
	result += fmt.Sprintf("Token: %v | ", a.Token.Val)
	result += fmt.Sprintf("Type: %v | ", a.Type)
	result += fmt.Sprintf("Payload: %v\n", a.Payload)
	for _, child := range a.Children {
		result += child.TreeString(indent + 1)
	}
	return result
}

func (p *Parser) Token2ASTNode(token *lexer.Token) *ASTNode {
	return &ASTNode{
		raw:      token.Val,
		Token:    token,
		Children: []*ASTNode{},
		Type:     p.Reflect(token),
		Payload:  nil,
	}
}
