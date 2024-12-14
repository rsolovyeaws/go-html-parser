package lexer

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
	l.skipWhitespace()

	switch l.ch {
	case '<':
		if l.peekChar() == '/' {
			l.readChar()
			return l.readEndTag()
		} else if l.peekChar() == '!' {
			if l.input[l.position+2:l.position+9] == doctypeDeclaration {
				return l.readDoctype()
			}
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

	l.skipWhitespace()
	attributes := l.readAttributes()

	if l.ch == '/' && l.peekChar() == '>' {
		l.readChar() // Consume '/'
		l.readChar() // Consume '>'
		return Token{Type: TokenSelfClosingTag, Value: trimSpaces(tagName + " " + attributes)}
	}

	l.readChar() // Consume '>'
	return Token{Type: TokenStartTag, Value: trimSpaces(tagName + " " + attributes)}
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
	text := l.readUntil("<")
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
