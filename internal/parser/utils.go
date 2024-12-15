package parser

import (
	"fmt"
	"strings"
)

// Parses a tag string into a tag name and a map of attributes
func parseTag(tag string) (string, map[string]string) {
	tagNameEnd := strings.IndexAny(tag, " \t\n") // Find the end of the tag name
	if tagNameEnd == -1 {
		tagNameEnd = len(tag) // No attributes, the tag name is the entire string
	}
	tagName := tag[:tagNameEnd]
	attributes := map[string]string{}

	// Parse the remaining string for attributes
	rest := tag[tagNameEnd:]
	var key, value string
	var inKey, inValue, inQuote bool
	quoteChar := byte(0)

	for i := 0; i < len(rest); i++ {
		ch := rest[i]

		switch {
		case inKey && (ch == '=' || ch == ' ' || ch == '\t'):
			// Transition from key to value
			if ch == '=' {
				inKey = false
				inValue = true
			} else {
				attributes[key] = ""
				key = ""
				inKey = false
			}
		case inValue && inQuote && ch == quoteChar:
			// End of quoted value
			attributes[key] = value
			key, value = "", ""
			inValue, inQuote = false, false
		case inValue && !inQuote && (ch == '"' || ch == '\''):
			// Start of quoted value
			inQuote = true
			quoteChar = ch
		case inValue && !inQuote && ch == ' ':
			// End of unquoted value
			attributes[key] = value
			key, value = "", ""
			inValue = false
		case !inKey && !inValue && ch != ' ' && ch != '\t':
			// Start of a new key
			inKey = true
			key = string(ch)
		case inKey:
			// Accumulate key
			key += string(ch)
		case inValue:
			// Accumulate value
			value += string(ch)
		}
	}

	// Add any trailing key or value
	if inKey {
		attributes[key] = ""
	} else if inValue {
		attributes[key] = value
	}

	fmt.Printf("parseTag: tagName=%s, attributes=%v\n", tagName, attributes)
	return tagName, attributes
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
