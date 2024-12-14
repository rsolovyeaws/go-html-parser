package lexer

import (
	"testing"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedTokens []Token
	}{
		{
			name:  "Basic Tags",
			input: `<div class="test" id="1">Hello</div>`,
			expectedTokens: []Token{
				{Type: TokenStartTag, Value: "div class=\"test\" id=\"1\""},
				{Type: TokenText, Value: "Hello"},
				{Type: TokenEndTag, Value: "div"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "Self-Closing Tag",
			input: `<img src="image.jpg" />`,
			expectedTokens: []Token{
				{Type: TokenSelfClosingTag, Value: "img src=\"image.jpg\""},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "Nested Tags",
			input: `<div><p>Nested</p></div>`,
			expectedTokens: []Token{
				{Type: TokenStartTag, Value: "div"},
				{Type: TokenStartTag, Value: "p"},
				{Type: TokenText, Value: "Nested"},
				{Type: TokenEndTag, Value: "p"},
				{Type: TokenEndTag, Value: "div"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "Comment",
			input: `<!-- This is a comment -->`,
			expectedTokens: []Token{
				{Type: TokenComment, Value: "This is a comment"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "Malformed HTML",
			input: `<div><p>Unclosed Div`,
			expectedTokens: []Token{
				{Type: TokenStartTag, Value: "div"},
				{Type: TokenStartTag, Value: "p"},
				{Type: TokenText, Value: "Unclosed Div"},
				{Type: TokenEOF, Value: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)

			for i, expected := range tt.expectedTokens {
				tok := l.NextToken()

				if tok.Type != expected.Type {
					t.Fatalf("test '%s' [%d] - tokentype wrong. expected=%q, got=%q",
						tt.name, i, expected.Type, tok.Type)
				}

				if tok.Value != expected.Value {
					t.Fatalf("test '%s' [%d] - tokenvalue wrong. expected=%q, got=%q",
						tt.name, i, expected.Value, tok.Value)
				}
			}
		})
	}
}
