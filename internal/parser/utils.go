package parser

import "html"

// DecodeEntities decodes HTML entities using Go's built-in function.
func DecodeEntities(input string) string {
	return html.UnescapeString(input)
}

func isVoidElement(tagName string) bool {
	// List of void elements in HTML
	voidElements := map[string]bool{
		"area": true, "base": true, "br": true, "col": true,
		"embed": true, "hr": true, "img": true, "input": true,
		"link": true, "meta": true, "source": true, "track": true,
		"wbr": true,
	}
	return voidElements[tagName]
}
