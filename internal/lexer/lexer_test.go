package lexer

import (
	"reflect"
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
				{Type: TokenStartTag, Value: "div", Attributes: map[string]string{"class": "test", "id": "1"}},
				{Type: TokenText, Value: "Hello"},
				{Type: TokenEndTag, Value: "div"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "Self-Closing Tag",
			input: `<img src="image.jpg" />`,
			expectedTokens: []Token{
				{Type: TokenSelfClosingTag, Value: "img", Attributes: map[string]string{"src": "image.jpg"}},
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
		{
			name:  "Attributes Without Quotes",
			input: `<div id=test class=test-class>Content</div>`,
			expectedTokens: []Token{
				{Type: TokenStartTag, Value: "div", Attributes: map[string]string{"id": "test", "class": "test-class"}},
				{Type: TokenText, Value: "Content"},
				{Type: TokenEndTag, Value: "div"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "Special Characters in Attributes",
			input: `<input value="Tom & Jerry" disabled>`,
			expectedTokens: []Token{
				{Type: TokenStartTag, Value: "input", Attributes: map[string]string{"value": "Tom & Jerry", "disabled": ""}},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "Nested Comments",
			input: `<!-- Outer <!-- Inner --> -->`,
			expectedTokens: []Token{
				{Type: TokenComment, Value: "Outer <!-- Inner -->"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "Missing End Tags",
			input: `<div><span>Text`,
			expectedTokens: []Token{
				{Type: TokenStartTag, Value: "div"},
				{Type: TokenStartTag, Value: "span"},
				{Type: TokenText, Value: "Text"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "Invalid Tags",
			input: `<123invalid>Text</123invalid>`,
			expectedTokens: []Token{
				{Type: TokenStartTag, Value: "123invalid"},
				{Type: TokenText, Value: "Text"},
				{Type: TokenEndTag, Value: "123invalid"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "HTML with Doctype Declaration",
			input: `<!DOCTYPE html><html></html>`,
			expectedTokens: []Token{
				{Type: TokenComment, Value: "DOCTYPE html"},
				{Type: TokenStartTag, Value: "html"},
				{Type: TokenEndTag, Value: "html"},
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

				if !reflect.DeepEqual(normalizeAttributes(tok.Attributes), normalizeAttributes(expected.Attributes)) {
					t.Fatalf("test '%s' [%d] - tokenattributes wrong. expected=%v, got=%v",
						tt.name, i, expected.Attributes, tok.Attributes)
				}
			}
		})
	}
}

func normalizeAttributes(attrs map[string]string) map[string]string {
	if len(attrs) == 0 {
		return nil
	}
	return attrs
}
