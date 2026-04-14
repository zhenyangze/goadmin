// Package async provides async content loading components.
package async

import (
	"html/template"
)

// Async is an async loading component.
type Async struct {
	id       string
	url      string
	interval int // Auto-refresh interval in seconds (0 = no auto-refresh)
	loading  template.HTML
	error    template.HTML
}

// New creates a new async component.
func New(url string) *Async {
	return &Async{
		url:     url,
		loading: template.HTML(`<div class="spinner"></div>`),
	}
}

// ID sets the component ID.
func (a *Async) ID(id string) *Async {
	a.id = id
	return a
}

// Interval sets auto-refresh interval in seconds.
func (a *Async) Interval(seconds int) *Async {
	a.interval = seconds * 1000 // Convert to milliseconds
	return a
}

// Loading sets custom loading content.
func (a *Async) Loading(content template.HTML) *Async {
	a.loading = content
	return a
}

// Error sets custom error content.
func (a *Async) Error(content template.HTML) *Async {
	a.error = content
	return a
}

// RenderContext provides data for rendering.
type RenderContext struct {
	ID       string
	URL      string
	Interval int
	Loading  template.HTML
	Error    template.HTML
}

// Render prepares the async component for rendering.
func (a *Async) Render() *RenderContext {
	return &RenderContext{
		ID:       a.id,
		URL:      a.url,
		Interval: a.interval,
		Loading:  a.loading,
		Error:    a.error,
	}
}

// Card extends card.Card with async loading support.
// This is a helper that combines Card with Async loading.
type Card struct {
	*Async
	title       string
	subtitle    string
	tools       []template.HTML
	class       string
	collapsible bool
	removable   bool
}

// NewCard creates a new async card.
func NewCard(url string) *Card {
	return &Card{
		Async: New(url),
	}
}

// Title sets the card title.
func (c *Card) Title(title string) *Card {
	c.title = title
	return c
}

// Subtitle sets the card subtitle.
func (c *Card) Subtitle(subtitle string) *Card {
	c.subtitle = subtitle
	return c
}

// Tool adds a tool button.
func (c *Card) Tool(tool template.HTML) *Card {
	c.tools = append(c.tools, tool)
	return c
}

// Class adds CSS classes.
func (c *Card) Class(class string) *Card {
	c.class = class
	return c
}

// Collapsible makes the card collapsible.
func (c *Card) Collapsible(collapsible bool) *Card {
	c.collapsible = collapsible
	return c
}

// Removable makes the card removable.
func (c *Card) Removable(removable bool) *Card {
	c.removable = removable
	return c
}

// CardRenderContext provides data for rendering the async card.
type CardRenderContext struct {
	ID          string
	Title       string
	Subtitle    string
	Tools       []template.HTML
	Class       string
	Collapsible bool
	Removable   bool
	URL         string
	Interval    int
	Loading     template.HTML
}

// RenderCard prepares the async card for rendering.
func (c *Card) RenderCard() *CardRenderContext {
	return &CardRenderContext{
		ID:          c.id,
		Title:       c.title,
		Subtitle:    c.subtitle,
		Tools:       c.tools,
		Class:       c.class,
		Collapsible: c.collapsible,
		Removable:   c.removable,
		URL:         c.url,
		Interval:    c.interval,
		Loading:     c.loading,
	}
}
