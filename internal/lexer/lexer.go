package lexer

import (
	"fmt"
	"strings"
)

const (
	TokenStartTag       = "StartTag"
	TokenEndTag         = "EndTag"
	TokenSelfClosingTag = "SelfClosingTag"
	TokenText           = "Text"
	TokenComment        = "Comment"
	TokenEOF            = "EOF"
	doctypeDeclaration  = "DOCTYPE"
	commentOpen         = "<!--"
	commentClose        = "-->"
)

type Token struct {
	Type       string
	Value      string
	Position   int               // Position in input for debugging
	Attributes map[string]string // Add this field for tag attributes
}

type Lexer struct {
	input        string // The HTML input
	position     int    // Current position in the input
	readPosition int    // Position after the current character
	ch           byte   // Current character
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

//////////////////////////
// Character Processing //
//////////////////////////

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII code for "NUL"
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

//////////////////////
// Token Processing //
//////////////////////

func (l *Lexer) NextToken() Token {
	// Skip leading whitespace only when outside of text
	if l.ch == '<' {
		switch {
		case l.peekChar() == '/':
			l.readChar()
			return l.readEndTag()
		case l.peekChar() == '!':
			if l.input[l.position+2:l.position+9] == doctypeDeclaration {
				return l.readDoctype()
			}
			return l.readComment()
		default:
			return l.readStartTag()
		}
	} else if l.ch == 0 {
		return Token{Type: TokenEOF, Value: ""}
	} else {
		return l.readText()
	}
}

func (l *Lexer) readDoctype() Token {
	l.readChar() // Consume '<'
	l.readChar() // Consume '!'

	start := l.position
	for l.ch != '>' && l.ch != 0 {
		l.readChar()
	}

	value := l.input[start:l.position]
	l.readChar() // Consume '>'
	return Token{Type: TokenComment, Value: value}
}

func (l *Lexer) readStartTag() Token {
	l.readChar() // Consume '<'
	tagName := l.readIdentifier()

	attributes := make(map[string]string)
	var fullTag strings.Builder
	fullTag.WriteString(tagName) // Start with the tag name

	for {
		l.skipWhitespace()
		if l.ch == '/' && l.peekChar() == '>' {
			return l.readSelfClosingTag(tagName, attributes)
		}
		if l.ch == '>' {
			break
		}
		key := l.readIdentifier()
		var value string
		isQuoted := false
		if l.ch == '=' {
			l.readChar() // Consume '='
			if l.ch == '"' || l.ch == '\'' {
				isQuoted = true
				quote := l.ch
				l.readChar() // Consume opening quote
				value = l.readUntil(string(quote))
				l.readChar() // Consume closing quote
			} else {
				value = l.readIdentifier()
			}
		}
		attributes[key] = value

		// Append the key-value pair to the tag string
		if value != "" {
			if isQuoted {
				fullTag.WriteString(fmt.Sprintf(` %s="%s"`, key, value))
			} else {
				fullTag.WriteString(fmt.Sprintf(` %s=%s`, key, value))
			}
		} else {
			fullTag.WriteString(fmt.Sprintf(" %s", key))
		}
	}

	l.readChar() // Consume '>'

	return Token{
		Type:       TokenStartTag,
		Value:      tagName, // fullTag.String(),
		Attributes: attributes,
	}
}

func (l *Lexer) readSelfClosingTag(tagName string, attributes map[string]string) Token {
	l.readChar() // Consume '/'
	l.readChar() // Consume '>'

	return Token{
		Type:       TokenSelfClosingTag,
		Value:      tagName,
		Attributes: attributes,
	}
}

func (l *Lexer) readEndTag() Token {
	l.readChar()                  // Consume '<' and '/'
	tagName := l.readIdentifier() // Read the tag name
	l.readChar()                  // Consume '>'
	return Token{Type: TokenEndTag, Value: tagName}
}

func (l *Lexer) readComment() Token {
	l.readChar() // Consume '<'
	l.readChar() // Consume '!'
	l.readChar() // Consume '-'
	l.readChar() // Consume '-'

	start := l.position
	depth := 1

	for {
		// Check if we've reached the end of input
		if l.position >= len(l.input) {
			break
		}

		// Handle nested `<!--` safely
		if l.position+4 <= len(l.input) && l.input[l.position:l.position+4] == "<!--" {
			depth++
			l.position += 4
			continue
		}

		// Handle closing `-->` safely
		if l.position+3 <= len(l.input) && l.input[l.position:l.position+3] == "-->" {
			depth--
			if depth == 0 {
				break
			}
			l.position += 3
			continue
		}

		l.readChar()
	}

	// Extract the comment value and trim spaces
	comment := trimSpaces(l.input[start:l.position])
	l.readChar() // Consume '-'
	l.readChar() // Consume '-'
	l.readChar() // Consume '>'

	return Token{Type: TokenComment, Value: comment}
}

func (l *Lexer) readText() Token {
	start := l.position
	for l.ch != '<' && l.ch != 0 { // Read until a '<' or EOF
		l.readChar()
	}
	text := l.input[start:l.position] // Extract raw text
	return Token{Type: TokenText, Value: text}
}

////////////////////////
// Helper Methods     //
////////////////////////

func (l *Lexer) readUntil(stop string) string {
	start := l.position
	for {
		if l.position >= len(l.input) {
			break
		}
		if l.input[l.position:l.position+len(stop)] == stop {
			break
		}
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '-' || l.ch == '_' {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readAttributes() string {
	start := l.position
	for l.ch != '>' && l.ch != '/' && l.ch != 0 {
		l.readChar()
	}
	return trimSpaces(l.input[start:l.position])
}

func trimSpaces(s string) string {
	start, end := 0, len(s)
	for start < end && s[start] == ' ' {
		start++
	}
	for end > start && s[end-1] == ' ' {
		end--
	}
	return s[start:end]
}

////////////////////////
// Character Checks   //
////////////////////////

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
