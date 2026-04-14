// Package tab provides tab components for organizing content.
package tab

import (
	"fmt"
	"html/template"
	"time"
)

// Tab represents a single tab panel.
type Tab struct {
	ID      string
	Title   string
	Content template.HTML
	Icon    string
	Badge   string
}

// Tabs is a collection of tabs.
type Tabs struct {
	tabs    []*Tab
	id      string
	class   string
	stacked bool // Vertical tabs on left side
}

// New creates a new tab container.
func New() *Tabs {
	return &Tabs{
		tabs: make([]*Tab, 0),
	}
}

// ID sets the container ID.
func (t *Tabs) ID(id string) *Tabs {
	t.id = id
	return t
}

// Class adds CSS classes.
func (t *Tabs) Class(class string) *Tabs {
	t.class = class
	return t
}

// Stacked enables vertical tab layout.
func (t *Tabs) Stacked(stacked bool) *Tabs {
	t.stacked = stacked
	return t
}

// Add adds a new tab.
func (t *Tabs) Add(title string, content template.HTML) *Tab {
	tab := &Tab{
		ID:      "tab-" + generateID(),
		Title:   title,
		Content: content,
	}
	t.tabs = append(t.tabs, tab)
	return tab
}

// AddTab adds an existing tab.
func (t *Tabs) AddTab(tab *Tab) *Tabs {
	if tab.ID == "" {
		tab.ID = "tab-" + generateID()
	}
	t.tabs = append(t.tabs, tab)
	return t
}

// WithIcon sets the tab icon.
func (t *Tab) WithIcon(icon string) *Tab {
	t.Icon = icon
	return t
}

// WithBadge sets the tab badge.
func (t *Tab) WithBadge(badge string) *Tab {
	t.Badge = badge
	return t
}

// RenderContext provides data for rendering.
type RenderContext struct {
	ID      string
	Class   string
	Stacked bool
	Tabs    []*Tab
}

// Render prepares tabs for rendering.
func (t *Tabs) Render() *RenderContext {
	return &RenderContext{
		ID:      t.id,
		Class:   t.class,
		Stacked: t.stacked,
		Tabs:    t.tabs,
	}
}

// generateID creates a unique ID.
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
