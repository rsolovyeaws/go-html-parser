package lexer

const (
	TokenStartTag       = "StartTag"
	TokenEndTag         = "EndTag"
	TokenSelfClosingTag = "SelfClosingTag"
	TokenText           = "Text"
	TokenAttribute      = "Attribute"
	TokenComment        = "Comment"
	TokenEOF            = "EOF"
)

type Token struct {
	Type     string
	Value    string
	Position int // Position in input for debugging
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

func (l *Lexer) readChar() {
	//fmt.Printf("readChar: current='%c', position=%d\n", l.ch, l.position) // Debug log
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

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	switch l.ch {
	case '<':
		if l.peekChar() == '/' {
			l.readChar()
			return l.readEndTag()
		} else if l.peekChar() == '!' {
			return l.readComment()
		} else {
			return l.readStartTag()
		}
	case 0:
		return Token{Type: TokenEOF, Value: ""}
	default:
		return l.readText()
	}
}

func (l *Lexer) readStartTag() Token {
	l.readChar() // Consume '<'
	tagName := l.readIdentifier()

	l.skipWhitespace() // Ensure no leading spaces before attributes
	attributes := l.readAttributes()

	if l.ch == '/' && l.peekChar() == '>' { // Self-closing tag
		l.readChar() // Consume '/'
		l.readChar() // Consume '>'
		if attributes != "" {
			return Token{Type: TokenSelfClosingTag, Value: tagName + " " + attributes}
		}
		return Token{Type: TokenSelfClosingTag, Value: tagName}
	}

	l.readChar() // Consume '>'
	if attributes != "" {
		return Token{Type: TokenStartTag, Value: tagName + " " + attributes}
	}
	return Token{Type: TokenStartTag, Value: tagName}
}

func (l *Lexer) readEndTag() Token {
	l.readChar()                  // Consume '<' and '/'
	tagName := l.readIdentifier() // Read the tag name starting at the correct position
	l.readChar()                  // Consume '>'

	return Token{Type: TokenEndTag, Value: tagName}
}

func (l *Lexer) readComment() Token {
	l.readChar() // Consume '<'
	l.readChar() // Consume '!'
	l.readChar() // Consume '-'
	l.readChar() // Consume '-'
	l.skipWhitespace()
	// Start reading after `<!--`
	start := l.position
	for {
		if l.position+3 <= len(l.input) && l.input[l.position:l.position+3] == "-->" {
			break
		}
		l.readChar()
	}

	// Extract the comment value, stopping before `-->`
	comment := trimTrailingSpaces(l.input[start:l.position])
	l.readChar() // Consume '-'
	l.readChar() // Consume '-'
	l.readChar() // Consume '>'

	return Token{Type: TokenComment, Value: comment}
}

func (l *Lexer) readText() Token {
	text := l.readUntil("<")
	return Token{Type: TokenText, Value: text}
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '-' || l.ch == '_' {
		l.readChar() // Move to the next character
	}
	identifier := l.input[start:l.position]
	// fmt.Printf("readIdentifier: '%s'\n", identifier) // Debugging final identifier
	return identifier
}

func (l *Lexer) readAttributes() string {
	start := l.position
	for l.ch != '>' && l.ch != '/' && l.ch != 0 {
		l.readChar()
	}
	return trimTrailingSpaces(l.input[start:l.position])
}

func trimTrailingSpaces(s string) string {
	end := len(s)
	for end > 0 && s[end-1] == ' ' {
		end--
	}
	return s[:end]
}

func (l *Lexer) readUntil(stop string) string {
	start := l.position
	for {
		// Prevent slicing out of bounds
		if l.position+len(stop) > len(l.input) {
			break
		}
		// Stop if the `stop` string is found
		if l.input[l.position:l.position+len(stop)] == stop {
			break
		}
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
