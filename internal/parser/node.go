package parser

type NodeType string

const (
	NodeElement NodeType = "Element"
	NodeText    NodeType = "Text"
	NodeComment NodeType = "Comment"
)

type Node struct {
	Type       NodeType          // Element, Text, Comment
	TagName    string            // Only for Element nodes
	Attributes map[string]string // Only for Element nodes
	Content    string            // Only for Text and Comment nodes
	Children   []*Node
	Parent     *Node
}
