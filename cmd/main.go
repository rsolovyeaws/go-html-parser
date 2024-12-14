package main

import (
	"fmt"

	"github.com/rsolovyeaws/go-html-parser/pkg/htmltree"
)

func main() {
	fmt.Println("HTML Parser initialized")
	tree := htmltree.NewTree()
	fmt.Printf("Tree created: %v\n", tree)
}
