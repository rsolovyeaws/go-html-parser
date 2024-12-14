package htmltree

import (
	"testing"
)

func TestNewTree(t *testing.T) {
	tree := NewTree()
	if tree.Root == nil {
		t.Error("Root node should not be nil")
	}
}
