package parser

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedRoot *Node
	}{
		{
			name:  "Simple HTML",
			input: `<div>Hello</div>`,
			expectedRoot: &Node{
				Type:    NodeElement,
				TagName: "root",
				Children: []*Node{
					{
						Type:    NodeElement,
						TagName: "div",
						Children: []*Node{
							{
								Type:    NodeText,
								Content: "Hello",
							},
						},
					},
				},
			},
		},
		{
			name:  "Nested Tags",
			input: `<div><p>Nested</p></div>`,
			expectedRoot: &Node{
				Type:    NodeElement,
				TagName: "root",
				Children: []*Node{
					{
						Type:    NodeElement,
						TagName: "div",
						Children: []*Node{
							{
								Type:    NodeElement,
								TagName: "p",
								Children: []*Node{
									{
										Type:    NodeText,
										Content: "Nested",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "Comments",
			input: `<div><!-- A comment --></div>`,
			expectedRoot: &Node{
				Type:    NodeElement,
				TagName: "root",
				Children: []*Node{
					{
						Type:    NodeElement,
						TagName: "div",
						Children: []*Node{
							{
								Type:    NodeComment,
								Content: "A comment",
							},
						},
					},
				},
			},
		},
		{
			name:  "Attributes and Self-Closing Tags",
			input: `<img src="image.jpg" alt="An image" />`,
			expectedRoot: &Node{
				Type:    NodeElement,
				TagName: "root",
				Children: []*Node{
					{
						Type:       NodeElement,
						TagName:    "img",
						Attributes: map[string]string{"src": "image.jpg", "alt": "An image"},
					},
				},
			},
		},
		{
			name:  "Malformed HTML",
			input: `<div><p>Unclosed Div`,
			expectedRoot: &Node{
				Type:    NodeElement,
				TagName: "root",
				Children: []*Node{
					{
						Type:    NodeElement,
						TagName: "div",
						Children: []*Node{
							{
								Type:    NodeElement,
								TagName: "p",
								Children: []*Node{
									{
										Type:    NodeText,
										Content: "Unclosed Div",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "Decode Entities in Text",
			input: `<p>Tom &amp; Jerry</p>`,
			expectedRoot: &Node{
				Type:    NodeElement,
				TagName: "root",
				Children: []*Node{
					{
						Type:    NodeElement,
						TagName: "p",
						Children: []*Node{
							{
								Type:    NodeText,
								Content: "Tom & Jerry",
							},
						},
					},
				},
			},
		},
		{
			name:  "Decode Entities in Attributes",
			input: `<img src="image.jpg" alt="Tom &amp; Jerry" />`,
			expectedRoot: &Node{
				Type:    NodeElement,
				TagName: "root",
				Children: []*Node{
					{
						Type:       NodeElement,
						TagName:    "img",
						Attributes: map[string]string{"src": "image.jpg", "alt": "Tom & Jerry"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.input)
			root := p.Parse()

			if !compareNodes(root, tt.expectedRoot) {
				t.Fatalf("Test '%s' failed: Expected %v, got %v", tt.name, tt.expectedRoot, root)
			}
		})
	}
}

func TestParserExpanded(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedRoot *Node
	}{
		{
			name:  "Attributes Without Values",
			input: `<input type="checkbox" checked>`,
			expectedRoot: &Node{
				Type:    NodeElement,
				TagName: "root",
				Children: []*Node{
					{
						Type:       NodeElement,
						TagName:    "input",
						Attributes: map[string]string{"type": "checkbox", "checked": ""},
					},
				},
			},
		},
		{
			name:  "Mixed Attribute Quoting",
			input: `<tag key1="value1" key2='value2'>`,
			expectedRoot: &Node{
				Type:    NodeElement,
				TagName: "root",
				Children: []*Node{
					{
						Type:       NodeElement,
						TagName:    "tag",
						Attributes: map[string]string{"key1": "value1", "key2": "value2"},
					},
				},
			},
		},
		{
			name:  "Nested Self-Closing Tags",
			input: `<div><img src="logo.png" /><br /></div>`,
			expectedRoot: &Node{
				Type:    NodeElement,
				TagName: "root",
				Children: []*Node{
					{
						Type:    NodeElement,
						TagName: "div",
						Children: []*Node{
							{
								Type:       NodeElement,
								TagName:    "img",
								Attributes: map[string]string{"src": "logo.png"},
							},
							{
								Type:    NodeElement,
								TagName: "br",
							},
						},
					},
				},
			},
		},
		{
			name:  "Mixed Content",
			input: `<div>Hello <span>world</span></div>`,
			expectedRoot: &Node{
				Type:    NodeElement,
				TagName: "root",
				Children: []*Node{
					{
						Type:    NodeElement,
						TagName: "div",
						Children: []*Node{
							{
								Type:    NodeText,
								Content: "Hello ",
							},
							{
								Type:    NodeElement,
								TagName: "span",
								Children: []*Node{
									{
										Type:    NodeText,
										Content: "world",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "Malformed HTML",
			input: `<div><span>Missing End Tags`,
			expectedRoot: &Node{
				Type:    NodeElement,
				TagName: "root",
				Children: []*Node{
					{
						Type:    NodeElement,
						TagName: "div",
						Children: []*Node{
							{
								Type:    NodeElement,
								TagName: "span",
								Children: []*Node{
									{
										Type:    NodeText,
										Content: "Missing End Tags",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "Void Elements",
			input: `<div><input type="text"><br></div>`,
			expectedRoot: &Node{
				Type:    NodeElement,
				TagName: "root",
				Children: []*Node{
					{
						Type:    NodeElement,
						TagName: "div",
						Children: []*Node{
							{
								Type:       NodeElement,
								TagName:    "input",
								Attributes: map[string]string{"type": "text"},
							},
							{
								Type:    NodeElement,
								TagName: "br",
							},
						},
					},
				},
			},
		},
		{
			name:  "Deeply Nested Structure",
			input: `<div><ul><li><a href="link">Item</a></li></ul></div>`,
			expectedRoot: &Node{
				Type:    NodeElement,
				TagName: "root",
				Children: []*Node{
					{
						Type:    NodeElement,
						TagName: "div",
						Children: []*Node{
							{
								Type:    NodeElement,
								TagName: "ul",
								Children: []*Node{
									{
										Type:    NodeElement,
										TagName: "li",
										Children: []*Node{
											{
												Type:       NodeElement,
												TagName:    "a",
												Attributes: map[string]string{"href": "link"},
												Children: []*Node{
													{
														Type:    NodeText,
														Content: "Item",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "Complex Attribute Combinations",
			input: `<input id="input1" class="form-input" type='text' data-value="123" />`,
			expectedRoot: &Node{
				Type:    NodeElement,
				TagName: "root",
				Children: []*Node{
					{
						Type:       NodeElement,
						TagName:    "input",
						Attributes: map[string]string{"id": "input1", "class": "form-input", "type": "text", "data-value": "123"},
					},
				},
			},
		},
		{
			name:  "Handling Doctype",
			input: `<!DOCTYPE html><html><body>Content</body></html>`,
			expectedRoot: &Node{
				Type:    NodeElement,
				TagName: "root",
				Children: []*Node{
					{
						Type:    NodeComment,
						Content: "DOCTYPE html",
					},
					{
						Type:    NodeElement,
						TagName: "html",
						Children: []*Node{
							{
								Type:    NodeElement,
								TagName: "body",
								Children: []*Node{
									{
										Type:    NodeText,
										Content: "Content",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "Real-World Sample",
			input: `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><title>Test</title></head><body><div class="container"><h1>Main Heading</h1><p>A <strong>bold</strong> statement.</p></div></body></html>`,
			expectedRoot: &Node{
				Type:    NodeElement,
				TagName: "root",
				Children: []*Node{
					{
						Type:    NodeComment,
						Content: "DOCTYPE html",
					},
					{
						Type:       NodeElement,
						TagName:    "html",
						Attributes: map[string]string{"lang": "en"},
						Children: []*Node{
							{
								Type:    NodeElement,
								TagName: "head",
								Children: []*Node{
									{
										Type:       NodeElement,
										TagName:    "meta",
										Attributes: map[string]string{"charset": "UTF-8"},
									},
									{
										Type:    NodeElement,
										TagName: "title",
										Children: []*Node{
											{
												Type:    NodeText,
												Content: "Test",
											},
										},
									},
								},
							},
							{
								Type:    NodeElement,
								TagName: "body",
								Children: []*Node{
									{
										Type:       NodeElement,
										TagName:    "div",
										Attributes: map[string]string{"class": "container"},
										Children: []*Node{
											{
												Type:    NodeElement,
												TagName: "h1",
												Children: []*Node{
													{
														Type:    NodeText,
														Content: "Main Heading",
													},
												},
											},
											{
												Type:    NodeElement,
												TagName: "p",
												Children: []*Node{
													{
														Type:    NodeText,
														Content: "A ",
													},
													{
														Type:    NodeElement,
														TagName: "strong",
														Children: []*Node{
															{
																Type:    NodeText,
																Content: "bold",
															},
														},
													},
													{
														Type:    NodeText,
														Content: " statement.",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "Malformed Large Document",
			input: `<div><p>Paragraph 1<p>Paragraph 2</div>`,
			expectedRoot: &Node{
				Type:    NodeElement,
				TagName: "root",
				Children: []*Node{
					{
						Type:    NodeElement,
						TagName: "div",
						Children: []*Node{
							{
								Type:    NodeElement,
								TagName: "p",
								Children: []*Node{
									{
										Type:    NodeText,
										Content: "Paragraph 1",
									},
								},
							},
							{
								Type:    NodeElement,
								TagName: "p",
								Children: []*Node{
									{
										Type:    NodeText,
										Content: "Paragraph 2",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.input)
			root := p.Parse()

			if !compareNodes(root, tt.expectedRoot) {
				t.Fatalf("Test '%s' failed: Trees do not match.\nExpected:\n%+v\nGot:\n%+v", tt.name, tt.expectedRoot, root)
			}
		})
	}
}

func compareNodes(a, b *Node) bool {
	if a.Type != b.Type || a.TagName != b.TagName || a.Content != b.Content {
		fmt.Printf("Node mismatch:\nExpected: %+v\nGot: %+v\n", b, a)
		return false
	}

	// Compare attributes
	if len(a.Attributes) != len(b.Attributes) {
		fmt.Printf("Attribute mismatch:\nExpected: %+v\nGot: %+v\n", b.Attributes, a.Attributes)
		return false
	}
	for key, val := range a.Attributes {
		if b.Attributes[key] != val {
			fmt.Printf("Attribute mismatch for key '%s':\nExpected: %s\nGot: %s\n", key, val, b.Attributes[key])
			return false
		}
	}

	// Compare children
	if len(a.Children) != len(b.Children) {
		fmt.Printf("Children count mismatch:\nExpected: %d\nGot: %d\n", len(b.Children), len(a.Children))
		return false
	}
	for i := range a.Children {
		if !compareNodes(a.Children[i], b.Children[i]) {
			fmt.Printf("Mismatch in child %d:\nExpected: %+v\nGot: %+v\n", i, b.Children[i], a.Children[i])
			return false
		}
	}

	return true
}
