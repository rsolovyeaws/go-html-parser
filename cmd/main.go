package main

import (
	"fmt"
	"strings"

	"github.com/rsolovyeaws/go-html-parser/internal/parser"
)

func main() {
	// HTML as a single string
	html := `<HTML><BODY><TABLE bgcolor="#e5e5e5" border="0" bordercolor="#ffffff" cellspacing="0" cellpadding="3" width="100%"><TR bgcolor="#DDDDDD"><TD width="15%" align="left" bgcolor="#ffffff" valign="top" style="font-size: 12px; font-weight: normal; font-family: Verdana; height: 27px"><B>Table Header Placeholder</B></TD></TR></TABLE><TABLE bgcolor="#e5e5e5" border="0" bordercolor="#ffffff" cellspacing="0" cellpadding="3" width="100%"><TR height="20"><TD align="middle" bgcolor="#e21212" width="15%" style="height: 17px"><FONT class="tekst"><B><FONT color="#ffffff" size="1" style="font-size: 11px; font-weight: normal; font-family: Verdana">Column Header 1</FONT></B></FONT></TD><TD align="middle" bgcolor="#e21212" width="25%" style="height: 17px"><FONT class="tekst"><B><FONT color="#ffffff" size="1" style="font-size: 11px; font-weight: normal; font-family: Verdana">Column Header 2</FONT></B></FONT></TD><TD align="middle" bgcolor="#e21212" width="60%" style="height: 17px"><FONT class="tekst"><B><FONT color="#ffffff" size="1" style="font-size: 11px; font-weight: normal; font-family: Verdana">Column Header 3</FONT></B></FONT></TD></TR><TR bgcolor="#DDDDDD"><TD width="15%" align="middle" valign="top" bgcolor="#ffffff" style="font-size: 11px; font-weight: normal; font-family: Verdana; height: 27px">Row 1 Cell 1</TD><TD width="25%" align="middle" valign="top" bgcolor="#ffffff" style="font-size: 11px; font-weight: normal; font-family: Verdana; height: 27px">Row 1 Cell 2</TD><TD width="60%" align="left" valign="top" bgcolor="#ffffff" style="font-size: 11px; font-weight: normal; font-family: Verdana; height: 27px">Row 1 Cell 3</TD></TR><TR bgcolor="#DDDDDD"><TD width="15%" align="middle" valign="top" bgcolor="#ffffff" style="font-size: 11px; font-weight: normal; font-family: Verdana; height: 27px">Row 2 Cell 1</TD><TD width="25%" align="middle" valign="top" bgcolor="#ffffff" style="font-size: 11px; font-weight: normal; font-family: Verdana; height: 27px">Row 2 Cell 2</TD><TD width="60%" align="left" valign="top" bgcolor="#ffffff" style="font-size: 11px; font-weight: normal; font-family: Verdana; height: 27px">Row 2 Cell 3</TD></TR></TABLE></BODY></HTML>`

	// Parse the HTML
	fmt.Println("Parsing HTML...")
	p := parser.New(html)
	root := p.Parse()

	// Print the entire DOM tree
	fmt.Println("Parsed DOM Tree:")
	printNode(root, "")
}

// printNode recursively prints a node and its children with indentation
func printNode(node *parser.Node, indent string) {
	fmt.Printf("%sNode: Type=%s, TagName=%s, Attributes=%v, Content=%q\n",
		indent, node.Type, node.TagName, node.Attributes, strings.TrimSpace(node.Content))

	for _, child := range node.Children {
		printNode(child, indent+"  ")
	}
}
