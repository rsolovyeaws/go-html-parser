package scraper

import (
	"fmt"

	"github.com/rsolovyeaws/go-html-parser/internal/httpclient"
	"github.com/rsolovyeaws/go-html-parser/internal/parser"
)

// Scraper represents a simple web scraper.
type Scraper struct {
	Headers map[string]string
}

// NewScraper creates a new instance of Scraper.
func NewScraper(headers map[string]string) *Scraper {
	return &Scraper{Headers: headers}
}

// Scrape fetches and parses the HTML from the given URL.
func (s *Scraper) Scrape(url string) (*parser.Node, error) {
	fmt.Printf("Fetching URL: %s\n", url)

	// Fetch the HTML content
	html, err := httpclient.FetchHTML(url, s.Headers)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %v", err)
	}

	fmt.Println("Parsing HTML...")
	// Parse the HTML content
	p := parser.New(html)
	root := p.Parse()

	fmt.Println("HTML parsing complete.")
	return root, nil
}

// PrintTree traverses and prints the parsed HTML tree.
func PrintTree(root *parser.Node) {
	printNode(root, "")
}

func printNode(node *parser.Node, indent string) {
	fmt.Printf("%sNode: Type=%s, TagName=%s, Content=%s\n", indent, node.Type, node.TagName, node.Content)
	for _, attr := range node.Attributes {
		fmt.Printf("%s  Attribute: %s\n", indent, attr)
	}
	for _, child := range node.Children {
		printNode(child, indent+"  ")
	}
}
