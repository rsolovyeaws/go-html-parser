package parser

import (
	"fmt"

	"github.com/rsolovyeaws/go-html-parser/internal/lexer"
)

type Parser struct {
	lexer *lexer.Lexer
	curr  lexer.Token
}

// New creates a new Parser instance
func New(input string) *Parser {
	l := lexer.New(input)
	return &Parser{lexer: l, curr: l.NextToken()}
}

func (p *Parser) nextToken() {
	p.curr = p.lexer.NextToken()
}

// Parse processes the input and returns the root Node of the parsed tree
func (p *Parser) Parse() *Node {
	root := &Node{
		Type:     NodeElement,
		TagName:  "root",
		Children: []*Node{},
	}

	stack := []*Node{root}

	for p.curr.Type != lexer.TokenEOF {
		switch p.curr.Type {
		case lexer.TokenStartTag:
			node := p.parseElement()
			if isVoidElement(node.TagName) {
				// Void elements are directly added to the parent
				stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, node)
			} else {
				// Add to parent and push onto stack
				stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, node)
				stack = append(stack, node)
			}
		case lexer.TokenSelfClosingTag:
			node := p.parseElement()
			fmt.Printf("SelfClosingTag - Node Created: TagName=%s, Attributes=%v\n", node.TagName, node.Attributes)
			stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, node)

		case lexer.TokenEndTag:
			// Pop the stack if there's a matching start tag
			if len(stack) > 1 && stack[len(stack)-1].TagName == p.curr.Value {
				stack = stack[:len(stack)-1]
			}
		case lexer.TokenText:
			content := p.curr.Value
			//content := strings.TrimSpace(p.curr.Value)
			if content != "" {
				node := &Node{
					Type:    NodeText,
					Content: content,
				}
				stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, node)
			}
		case lexer.TokenComment:
			node := &Node{
				Type:    NodeComment,
				Content: p.curr.Value,
			}
			stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, node)
		}
		p.nextToken()
	}

	debugNode(root, "")
	return root
}

// parseElement creates a Node from the current token
func (p *Parser) parseElement() *Node {
	return &Node{
		Type:       NodeElement,
		TagName:    p.curr.Value, // Only the tag name
		Attributes: p.curr.Attributes,
		Children:   []*Node{},
	}
}

// debugNode prints the Node structure for debugging purposes
func debugNode(node *Node, indent string) {
	fmt.Printf("%sNode: Type=%s, TagName=%s, Attributes=%v, Content=%s\n",
		indent, node.Type, node.TagName, node.Attributes, node.Content)
	for _, child := range node.Children {
		debugNode(child, indent+"  ")
	}
}
