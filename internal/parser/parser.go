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
			for len(stack) > 1 && isImplicitClose(stack[len(stack)-1].TagName, node.TagName) {
				stack = stack[:len(stack)-1] // Implicitly close the tag
			}
			stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, node)
			if !isVoidElement(node.TagName) {
				stack = append(stack, node) // Push non-void elements onto the stack
			}
		case lexer.TokenSelfClosingTag:
			node := p.parseElement()
			stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, node)
		case lexer.TokenEndTag:
			// Match the current end tag with the stack
			for len(stack) > 1 {
				if stack[len(stack)-1].TagName == p.curr.Value {
					stack = stack[:len(stack)-1] // Pop the matching tag
					break
				}
				// If no match, implicitly close the unclosed tag
				stack = stack[:len(stack)-1]
			}
		case lexer.TokenText:
			content := p.curr.Value
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

	// Close any remaining unclosed tags
	for len(stack) > 1 {
		stack = stack[:len(stack)-1]
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

// isImplicitClose checks if a tag should be implicitly closed
func isImplicitClose(current, next string) bool {
	// Implicit closing rules for certain HTML elements
	implicitCloseRules := map[string][]string{
		"li":     {"li"},
		"p":      {"p", "div", "ul", "ol"},
		"dt":     {"dt", "dd"},
		"dd":     {"dt", "dd"},
		"thead":  {"tbody", "tfoot"},
		"tbody":  {"tbody", "tfoot"},
		"tfoot":  {"tbody"},
		"tr":     {"tr"},
		"td":     {"td", "th"},
		"th":     {"td", "th"},
		"option": {"option", "optgroup"},
	}

	if nextTags, ok := implicitCloseRules[current]; ok {
		for _, tag := range nextTags {
			if next == tag {
				return true
			}
		}
	}
	return false
}
