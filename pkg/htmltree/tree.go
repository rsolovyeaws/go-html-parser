package htmltree

type Tree struct {
	Root *Node
}

type Node struct {
	Tag        string
	Attributes map[string]string
	Children   []*Node
	Text       string
	Parent     *Node
}

func NewTree() *Tree {
	return &Tree{Root: &Node{}}
}
