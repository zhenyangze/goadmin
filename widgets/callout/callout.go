// Package callout provides callout/alert widgets.
package callout

import "html/template"

// Callout is a highlighted message widget
type Callout struct {
	title   string
	content template.HTML
	style   string // info, success, warning, danger
	icon    string
}

// New creates a new callout widget
func New() *Callout {
	return &Callout{
		style: "info",
	}
}

// Title sets the callout title
func (c *Callout) Title(title string) *Callout {
	c.title = title
	return c
}

// Content sets the callout content
func (c *Callout) Content(content template.HTML) *Callout {
	c.content = content
	return c
}

// Style sets the callout style
func (c *Callout) Style(style string) *Callout {
	c.style = style
	return c
}

// Icon sets the callout icon
func (c *Callout) Icon(icon string) *Callout {
	c.icon = icon
	return c
}

// RenderContext provides data for rendering
type RenderContext struct {
	Title   string
	Content template.HTML
	Style   string
	Icon    string
}

// Render prepares the callout for rendering
func (c *Callout) Render() *RenderContext {
	return &RenderContext{
		Title:   c.title,
		Content: c.content,
		Style:   c.style,
		Icon:    c.icon,
	}
}

// Info creates an info callout
func Info(title string, content string) *Callout {
	return New().Title(title).Content(template.HTML(content)).Style("info")
}

// Success creates a success callout
func Success(title string, content string) *Callout {
	return New().Title(title).Content(template.HTML(content)).Style("success")
}

// Warning creates a warning callout
func Warning(title string, content string) *Callout {
	return New().Title(title).Content(template.HTML(content)).Style("warning")
}

// Danger creates a danger callout
func Danger(title string, content string) *Callout {
	return New().Title(title).Content(template.HTML(content)).Style("danger")
}
