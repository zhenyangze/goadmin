package show

import "html/template"

// Formatter customizes field rendering.
type Formatter func(record any, value any) template.HTML

// ItemType describes the row kind in the detail page.
type ItemType string

const (
	ItemField   ItemType = "field"
	ItemDivider ItemType = "divider"
)

// Item is a field row or divider.
type Item struct {
	Type      ItemType
	Name      string
	Label     string
	Title     string
	Formatter Formatter
}

// Builder defines a show/detail page.
type Builder struct {
	Title       string
	Description string
	Items       []*Item
}

// New creates a show builder.
func New() *Builder {
	return &Builder{}
}

// Field adds a detail field.
func (b *Builder) Field(name, label string) *Item {
	item := &Item{Type: ItemField, Name: name, Label: label}
	b.Items = append(b.Items, item)
	return item
}

// Divider adds a visual section divider.
func (b *Builder) Divider(title string) {
	b.Items = append(b.Items, &Item{Type: ItemDivider, Title: title})
}

// Display sets a custom renderer.
func (i *Item) Display(fn Formatter) *Item {
	i.Formatter = fn
	return i
}
