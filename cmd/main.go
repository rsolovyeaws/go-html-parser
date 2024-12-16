package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rsolovyeaws/go-html-parser/internal/httpclient"
	"github.com/rsolovyeaws/go-html-parser/internal/parser"
)

func main() {
	// Default URL if none is provided
	defaultURL := "https://elektrodistribucija.rs/planirana-iskljucenja-beograd/Dan_1_Iskljucenja.htm"

	// Check if a URL argument is provided
	var url string
	if len(os.Args) > 1 {
		url = os.Args[1]
	} else {
		url = defaultURL
	}

	// Headers to include in the request
	headers := map[string]string{
		"User-Agent": "Go-HTML-Parser",
	}

	// Fetch the HTML content from the provided URL
	fmt.Println("Fetching URL:", url)
	html, err := httpclient.FetchHTML(url, headers)
	if err != nil {
		log.Fatalf("Error fetching URL: %v", err)
	}

	// Parse the fetched HTML content
	p := parser.New(html)
	root := p.Parse()

	// Print the parsed tree
	fmt.Println("\nParsed Tree:")
	printTree(root, "")
}

// printTree recursively prints the parsed tree
func printTree(node *parser.Node, indent string) {
	fmt.Printf("%sNode: Type=%s, TagName=%s, Attributes=%v, Content=%s\n",
		indent, node.Type, node.TagName, node.Attributes, node.Content)

	for _, child := range node.Children {
		printTree(child, indent+"  ")
	}
}
