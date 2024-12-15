package parser

type NodeType string

const (
	NodeElement NodeType = "Element"
	NodeText    NodeType = "Text"
	NodeComment NodeType = "Comment"
)

type Node struct {
	Type        NodeType          // Element, Text, Comment
	TagName     string            // Only for Element nodes
	Attributes  map[string]string // Only for Element nodes
	Content     string            // Only for Text and Comment nodes
	Children    []*Node           // Child nodes
	Parent      *Node             // Pointer to parent node
	PrevSibling *Node             // Previous sibling
	NextSibling *Node             // Next sibling
}

func (n *Node) ParentNode() *Node {
	return n.Parent
}

func (n *Node) NextSiblingNode() *Node {
	return n.NextSibling
}

func (n *Node) PreviousSiblingNode() *Node {
	return n.PrevSibling
}

func (n *Node) FirstChildNode() *Node {
	if len(n.Children) > 0 {
		return n.Children[0]
	}
	return nil
}

func (n *Node) LastChildNode() *Node {
	if len(n.Children) > 0 {
		return n.Children[len(n.Children)-1]
	}
	return nil
}
