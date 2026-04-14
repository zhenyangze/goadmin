// Package dropdown provides dropdown menu components.
package dropdown

import (
	"fmt"
	"html/template"
	"time"
)

// Item represents a dropdown menu item.
type Item struct {
	Label    string
	URL      string
	Icon     template.HTML
	Class    string
	Confirm  string
	Divider  bool
	Disabled bool
}

// Dropdown is a dropdown menu component.
type Dropdown struct {
	id       string
	label    string
	class    string
	items    []*Item
	position string // left, right
}

// New creates a new dropdown.
func New(label string) *Dropdown {
	return &Dropdown{
		label:    label,
		class:    "btn-secondary",
		items:    make([]*Item, 0),
		position: "left",
	}
}

// ID sets the dropdown ID.
func (d *Dropdown) ID(id string) *Dropdown {
	d.id = id
	return d
}

// Class sets the button class.
func (d *Dropdown) Class(class string) *Dropdown {
	d.class = class
	return d
}

// Position sets the menu position (left or right).
func (d *Dropdown) Position(position string) *Dropdown {
	d.position = position
	return d
}

// Button adds a regular menu item.
func (d *Dropdown) Button(label, url string) *Item {
	item := &Item{
		Label: label,
		URL:   url,
	}
	d.items = append(d.items, item)
	return item
}

// Link is an alias for Button.
func (d *Dropdown) Link(label, url string) *Item {
	return d.Button(label, url)
}

// Divider adds a divider.
func (d *Dropdown) Divider() *Dropdown {
	d.items = append(d.items, &Item{Divider: true})
	return d
}

// Icon sets the item icon.
func (i *Item) WithIcon(icon template.HTML) *Item {
	i.Icon = icon
	return i
}

// Confirm adds a confirmation dialog.
func (i *Item) WithConfirm(message string) *Item {
	i.Confirm = message
	return i
}

// Disabled disables the item.
func (i *Item) WithDisabled(disabled bool) *Item {
	i.Disabled = disabled
	return i
}

// Danger styles the item as dangerous.
func (i *Item) Danger() *Item {
	i.Class = "text-danger"
	return i
}

// RenderContext provides data for rendering.
type RenderContext struct {
	ID       string
	Label    string
	Class    string
	Position string
	Items    []*Item
}

// Render prepares the dropdown for rendering.
func (d *Dropdown) Render() *RenderContext {
	id := d.id
	if id == "" {
		id = "dropdown-" + fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return &RenderContext{
		ID:       id,
		Label:    d.label,
		Class:    d.class,
		Position: d.position,
		Items:    d.items,
	}
}
