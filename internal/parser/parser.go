package parser

import (
	"fmt"
	"strings"

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

	stack := []*Node{root} // Stack to track open elements

	for p.curr.Type != lexer.TokenEOF {
		switch p.curr.Type {
		case lexer.TokenStartTag:
			node := p.parseElement()
			// Implicitly close open tags based on HTML rules
			for len(stack) > 1 && isImplicitClose(stack[len(stack)-1].TagName, node.TagName) {
				stack = stack[:len(stack)-1] // Pop the stack
			}
			// Add the node to the current parent
			appendChild(stack[len(stack)-1], node)
			// Push non-void elements onto the stack
			if !isVoidElement(node.TagName) {
				stack = append(stack, node)
			}

		case lexer.TokenSelfClosingTag:
			node := p.parseElement()
			appendChild(stack[len(stack)-1], node)

		case lexer.TokenEndTag:
			// Pop the stack for matching end tags
			for len(stack) > 1 {
				if stack[len(stack)-1].TagName == p.curr.Value {
					stack = stack[:len(stack)-1]
					break
				}
				stack = stack[:len(stack)-1] // Implicitly close unclosed tags
			}

		case lexer.TokenText:
			content := DecodeEntities(p.curr.Value)
			if content != "" {
				textNode := &Node{
					Type:    NodeText,
					Content: content,
				}
				appendChild(stack[len(stack)-1], textNode)
			}

		case lexer.TokenComment:
			commentNode := &Node{
				Type:    NodeComment,
				Content: p.curr.Value,
			}
			appendChild(stack[len(stack)-1], commentNode)
		}

		// Move to the next token
		p.nextToken()
	}

	// Close any remaining unclosed tags
	for len(stack) > 1 {
		stack = stack[:len(stack)-1]
	}

	debugNode(root, "") // Debugging output for tree structure
	return root
}

// parseElement creates a Node from the current token
func (p *Parser) parseElement() *Node {
	decodedAttributes := make(map[string]string)
	for key, value := range p.curr.Attributes {
		decodedAttributes[key] = DecodeEntities(value) // Decode entities in attributes
	}
	return &Node{
		Type:       NodeElement,
		TagName:    p.curr.Value,
		Attributes: decodedAttributes,
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

func (n *Node) FindByTag(tag string) []*Node {
	var result []*Node
	if n.TagName == tag {
		result = append(result, n)
	}
	for _, child := range n.Children {
		result = append(result, child.FindByTag(tag)...)
	}
	return result
}

func (n *Node) FindByID(id string) *Node {
	if val, ok := n.Attributes["id"]; ok && val == id {
		return n
	}
	for _, child := range n.Children {
		if found := child.FindByID(id); found != nil {
			return found
		}
	}
	return nil
}

func (n *Node) FindByClass(class string) []*Node {
	var result []*Node
	if val, ok := n.Attributes["class"]; ok {
		classes := strings.Fields(val)
		for _, c := range classes {
			if c == class {
				result = append(result, n)
				break
			}
		}
	}
	for _, child := range n.Children {
		result = append(result, child.FindByClass(class)...)
	}
	return result
}

func appendChild(parent *Node, child *Node) {
	if len(parent.Children) > 0 {
		prev := parent.Children[len(parent.Children)-1]
		prev.NextSibling = child
		child.PrevSibling = prev
	}
	child.Parent = parent
	parent.Children = append(parent.Children, child)
}
