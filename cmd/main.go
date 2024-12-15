package main

import (
	"fmt"

	"github.com/rsolovyeaws/go-html-parser/internal/parser"
)

func main() {
	html := `<div><p id="para1">Text 1</p><p class="text">Text 2</p></div>`
	p := parser.New(html)
	root := p.Parse()

	// Query nodes
	divs := root.FindByTag("div")
	fmt.Println("Divs found:", len(divs))

	p1 := root.FindByID("para1")
	fmt.Printf("Found element with ID 'para1': %+v\n", p1)

	textNodes := root.FindByClass("text")
	fmt.Println("Text nodes found:", len(textNodes))
	for _, node := range textNodes {
		fmt.Println("Content:", node.FirstChildNode().Content)
	}
}
