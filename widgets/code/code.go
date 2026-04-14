// Package code provides code display widgets.
package code

import "html/template"

// Code displays code with syntax highlighting
type Code struct {
	content  string
	language string
	theme    string
	lineNums bool
	wrap     bool
	height   string
}

// New creates a new code widget
func New() *Code {
	return &Code{
		language: "go",
		theme:    "default",
		lineNums: true,
	}
}

// Content sets the code content
func (c *Code) Content(content string) *Code {
	c.content = content
	return c
}

// Language sets the programming language
func (c *Code) Language(lang string) *Code {
	c.language = lang
	return c
}

// Theme sets the color theme
func (c *Code) Theme(theme string) *Code {
	c.theme = theme
	return c
}

// LineNumbers enables/disables line numbers
func (c *Code) LineNumbers(show bool) *Code {
	c.lineNums = show
	return c
}

// Wrap enables/disables line wrapping
func (c *Code) Wrap(wrap bool) *Code {
	c.wrap = wrap
	return c
}

// Height sets the max height
func (c *Code) Height(height string) *Code {
	c.height = height
	return c
}

// RenderContext provides data for rendering
type RenderContext struct {
	Content     string
	Language    string
	Theme       string
	LineNumbers bool
	Wrap        bool
	Height      string
}

// Render prepares the code for rendering
func (c *Code) Render() *RenderContext {
	return &RenderContext{
		Content:     c.content,
		Language:    c.language,
		Theme:       c.theme,
		LineNumbers: c.lineNums,
		Wrap:        c.wrap,
		Height:      c.height,
	}
}

// EscapeHTML escapes HTML in code
func EscapeHTML(s string) template.HTML {
	return template.HTML(template.HTMLEscapeString(s))
}
