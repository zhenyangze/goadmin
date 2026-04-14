// Package markdown provides markdown rendering widgets.
package markdown

import (
	"html/template"
)

// Markdown renders markdown content
type Markdown struct {
	content string
	theme   string
}

// New creates a new markdown widget
func New() *Markdown {
	return &Markdown{
		theme: "default",
	}
}

// Content sets the markdown content
func (m *Markdown) Content(content string) *Markdown {
	m.content = content
	return m
}

// Theme sets the render theme
func (m *Markdown) Theme(theme string) *Markdown {
	m.theme = theme
	return m
}

// RenderContext provides data for rendering
type RenderContext struct {
	Content    string
	Theme      string
	HTML       template.HTML
}

// Render prepares the markdown for rendering
func (m *Markdown) Render() *RenderContext {
	// Basic markdown to HTML conversion (simplified)
	html := renderMarkdown(m.content)

	return &RenderContext{
		Content:    m.content,
		Theme:      m.theme,
		HTML:       html,
	}
}

// Simple markdown renderer
func renderMarkdown(content string) template.HTML {
	// This is a simplified version - in production use a proper markdown library
	return template.HTML(content)
}
