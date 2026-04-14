// Package card provides card widgets for containing other widgets like forms.
package card

import (
	"html/template"
)

// Card is a container widget that can hold other widgets or content.
type Card struct {
	title       string
	subtitle    string
	tools       []template.HTML
	content     template.HTML
	footer      template.HTML
	class       string
	id          string
	collapsible bool
	collapsed   bool
	removable   bool
	loading     bool
	fullHeight  bool
}

// New creates a new card widget.
func New() *Card {
	return &Card{}
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

// Class adds CSS classes to the card.
func (c *Card) Class(class string) *Card {
	c.class = class
	return c
}

// ID sets the card ID.
func (c *Card) ID(id string) *Card {
	c.id = id
	return c
}

// Content sets the card content.
func (c *Card) Content(content template.HTML) *Card {
	c.content = content
	return c
}

// Body is an alias for Content.
func (c *Card) Body(content template.HTML) *Card {
	return c.Content(content)
}

// Tool adds a tool button/link to the card header.
func (c *Card) Tool(tool template.HTML) *Card {
	c.tools = append(c.tools, tool)
	return c
}

// Footer sets the card footer content.
func (c *Card) Footer(footer template.HTML) *Card {
	c.footer = footer
	return c
}

// Collapsible makes the card collapsible.
func (c *Card) Collapsible(collapsible bool) *Card {
	c.collapsible = collapsible
	return c
}

// Collapsed sets the initial collapsed state.
func (c *Card) Collapsed(collapsed bool) *Card {
	c.collapsed = collapsed
	return c
}

// Removable makes the card removable.
func (c *Card) Removable(removable bool) *Card {
	c.removable = removable
	return c
}

// Loading shows a loading state on the card.
func (c *Card) Loading(loading bool) *Card {
	c.loading = loading
	return c
}

// FullHeight makes the card take full available height.
func (c *Card) FullHeight(fullHeight bool) *Card {
	c.fullHeight = fullHeight
	return c
}

// Title returns the card title.
func (c *Card) GetTitle() string {
	return c.title
}

// Subtitle returns the card subtitle.
func (c *Card) GetSubtitle() string {
	return c.subtitle
}

// Tools returns the card tools.
func (c *Card) GetTools() []template.HTML {
	return c.tools
}

// Content returns the card content.
func (c *Card) GetContent() template.HTML {
	return c.content
}

// Footer returns the card footer.
func (c *Card) GetFooter() template.HTML {
	return c.footer
}

// Class returns the card CSS class.
func (c *Card) GetClass() string {
	return c.class
}

// ID returns the card ID.
func (c *Card) GetID() string {
	return c.id
}

// IsCollapsible returns whether the card is collapsible.
func (c *Card) IsCollapsible() bool {
	return c.collapsible
}

// IsCollapsed returns whether the card is collapsed.
func (c *Card) IsCollapsed() bool {
	return c.collapsed
}

// IsRemovable returns whether the card is removable.
func (c *Card) IsRemovable() bool {
	return c.removable
}

// IsLoading returns whether the card is in loading state.
func (c *Card) IsLoading() bool {
	return c.loading
}

// IsFullHeight returns whether the card has full height.
func (c *Card) IsFullHeight() bool {
	return c.fullHeight
}

// RenderContext provides data for rendering the card.
type RenderContext struct {
	Title       string
	Subtitle    string
	Tools       []template.HTML
	Content     template.HTML
	Footer      template.HTML
	Class       string
	ID          string
	Collapsible bool
	Collapsed   bool
	Removable   bool
	Loading     bool
	FullHeight  bool
}

// Render prepares the card for rendering.
func (c *Card) Render() *RenderContext {
	return &RenderContext{
		Title:       c.title,
		Subtitle:    c.subtitle,
		Tools:       c.tools,
		Content:     c.content,
		Footer:      c.footer,
		Class:       c.class,
		ID:          c.id,
		Collapsible: c.collapsible,
		Collapsed:   c.collapsed,
		Removable:   c.removable,
		Loading:     c.loading,
		FullHeight:  c.fullHeight,
	}
}
