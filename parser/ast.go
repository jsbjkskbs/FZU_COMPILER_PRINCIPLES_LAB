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

	_genCodeStartLine int
	_genCodeEndLine   int
}

type ASTNodeType int

// TreeString generates a string representation of the AST node and its children.
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

// Token2ASTNode converts a lexer token to an AST node.
// It creates a new AST node with the token's value and type, and initializes
// its children to an empty slice. The node's payload is set to nil.
func (p *Parser) Token2ASTNode(token *lexer.Token) *ASTNode {
	return &ASTNode{
		raw:      token.Val,
		Token:    token,
		Children: []*ASTNode{},
		Type:     p.Reflect(token),
		Payload:  nil,
	}
}
