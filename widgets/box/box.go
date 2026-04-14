// Package box provides box/container widgets.
package box

import "html/template"

// Box is a container widget with optional header and footer
type Box struct {
	title       string
	subtitle    string
	content     template.HTML
	tools       []template.HTML
	footer      template.HTML
	solid       bool
	collapsible bool
	collapsed   bool
	style       string // default, primary, info, success, warning, danger
}

// New creates a new box widget
func New() *Box {
	return &Box{
		style: "default",
	}
}

// Title sets the box title
func (b *Box) Title(title string) *Box {
	b.title = title
	return b
}

// Subtitle sets the box subtitle
func (b *Box) Subtitle(subtitle string) *Box {
	b.subtitle = subtitle
	return b
}

// Content sets the box content
func (b *Box) Content(content template.HTML) *Box {
	b.content = content
	return b
}

// Tool adds a tool to the box header
func (b *Box) Tool(tool template.HTML) *Box {
	b.tools = append(b.tools, tool)
	return b
}

// Footer sets the box footer
func (b *Box) Footer(footer template.HTML) *Box {
	b.footer = footer
	return b
}

// Solid makes the box header solid colored
func (b *Box) Solid() *Box {
	b.solid = true
	return b
}

// Style sets the box style
func (b *Box) Style(style string) *Box {
	b.style = style
	return b
}

// Collapsible makes the box collapsible
func (b *Box) Collapsible(collapsed bool) *Box {
	b.collapsible = true
	b.collapsed = collapsed
	return b
}

// RenderContext provides data for rendering
type RenderContext struct {
	Title       string
	Subtitle    string
	Content     template.HTML
	Tools       []template.HTML
	Footer      template.HTML
	Style       string
	Solid       bool
	Collapsible bool
	Collapsed   bool
}

// Render prepares the box for rendering
func (b *Box) Render() *RenderContext {
	return &RenderContext{
		Title:       b.title,
		Subtitle:    b.subtitle,
		Content:     b.content,
		Tools:       b.tools,
		Footer:      b.footer,
		Style:       b.style,
		Solid:       b.solid,
		Collapsible: b.collapsible,
		Collapsed:   b.collapsed,
	}
}
