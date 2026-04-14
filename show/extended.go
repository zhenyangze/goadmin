package show

import (
	"fmt"
	"html/template"
	"strings"
)

// Extended Item Types (defined alongside existing types)
const (
	ItemHtml     ItemType = "html"
	ItemNewline  ItemType = "newline"
	ItemPanel    ItemType = "panel"
	ItemRow      ItemType = "row"
	ItemRelation ItemType = "relation"
)

// ==================== HTML Display ====================

// Html adds raw HTML content to the show page
func (b *Builder) Html(content template.HTML) *Builder {
	b.Items = append(b.Items, &Item{
		Type: ItemHtml,
		Formatter: func(record any, value any) template.HTML {
			return content
		},
	})
	return b
}

// HtmlWithLabel adds HTML content with a label
func (b *Builder) HtmlWithLabel(label string, content template.HTML) *Item {
	item := &Item{
		Type:  ItemHtml,
		Label: label,
		Formatter: func(record any, value any) template.HTML {
			return content
		},
	}
	b.Items = append(b.Items, item)
	return item
}

// ==================== Newline ====================

// Newline adds a line break
func (b *Builder) Newline() *Builder {
	b.Items = append(b.Items, &Item{
		Type: ItemNewline,
	})
	return b
}

// Space adds vertical spacing
func (b *Builder) Space(height string) *Builder {
	b.Items = append(b.Items, &Item{
		Type: ItemNewline,
		Formatter: func(record any, value any) template.HTML {
			return template.HTML(fmt.Sprintf(`<div style="height:%s"></div>`, height))
		},
	})
	return b
}

// ==================== Panel ====================

// Panel creates a panel container for grouping fields
func (b *Builder) Panel(title string, fn func(*PanelBuilder)) *Builder {
	panelBuilder := &PanelBuilder{
		title: title,
	}
	fn(panelBuilder)

	b.Items = append(b.Items, &Item{
		Type:  ItemPanel,
		Title: title,
		Formatter: func(record any, value any) template.HTML {
			return panelBuilder.Render(record)
		},
	})
	return b
}

// PanelBuilder builds a panel with fields
type PanelBuilder struct {
	title       string
	subtitle    string
	icon        string
	collapsible bool
	collapsed   bool
	tools       []template.HTML
	items       []*Item
	footer      template.HTML
}

// Title sets the panel title
func (p *PanelBuilder) Title(title string) *PanelBuilder {
	p.title = title
	return p
}

// Subtitle sets the panel subtitle
func (p *PanelBuilder) Subtitle(subtitle string) *PanelBuilder {
	p.subtitle = subtitle
	return p
}

// Icon sets the panel icon
func (p *PanelBuilder) Icon(icon string) *PanelBuilder {
	p.icon = icon
	return p
}

// Collapsible makes the panel collapsible
func (p *PanelBuilder) Collapsible(collapsed bool) *PanelBuilder {
	p.collapsible = true
	p.collapsed = collapsed
	return p
}

// Tool adds a tool to the panel header
func (p *PanelBuilder) Tool(tool template.HTML) *PanelBuilder {
	p.tools = append(p.tools, tool)
	return p
}

// Field adds a field to the panel
func (p *PanelBuilder) Field(name, label string) *Item {
	item := &Item{
		Type:  ItemField,
		Name:  name,
		Label: label,
	}
	p.items = append(p.items, item)
	return item
}

// Html adds HTML content to the panel
func (p *PanelBuilder) Html(content template.HTML) *PanelBuilder {
	p.items = append(p.items, &Item{
		Type: ItemHtml,
		Formatter: func(record any, value any) template.HTML {
			return content
		},
	})
	return p
}

// Divider adds a divider to the panel
func (p *PanelBuilder) Divider(title string) *PanelBuilder {
	p.items = append(p.items, &Item{
		Type:  ItemDivider,
		Title: title,
	})
	return p
}

// Footer sets the panel footer
func (p *PanelBuilder) Footer(footer template.HTML) *PanelBuilder {
	p.footer = footer
	return p
}

// Render generates the panel HTML
func (p *PanelBuilder) Render(record any) template.HTML {
	var content strings.Builder

	// Build content
	for _, item := range p.items {
		switch item.Type {
		case ItemField:
			label := item.Label
			if label == "" {
				label = item.Name
			}
			value := ""
			if item.Formatter != nil {
				value = string(item.Formatter(record, nil))
			}
			content.WriteString(fmt.Sprintf(
				`<div class="field-row"><label class="field-label">%s</label><div class="field-value">%s</div></div>`,
				template.HTMLEscapeString(label), value))
		case ItemDivider:
			if item.Title != "" {
				content.WriteString(fmt.Sprintf(`<h4 class="panel-divider">%s</h4>`, template.HTMLEscapeString(item.Title)))
			} else {
				content.WriteString(`<hr class="panel-divider-line">`)
			}
		case ItemHtml:
			if item.Formatter != nil {
				content.WriteString(string(item.Formatter(record, nil)))
			}
		}
	}

	// Panel HTML
	collapsedClass := ""
	if p.collapsed {
		collapsedClass = " collapsed"
	}
	collapsibleAttr := ""
	if p.collapsible {
		collapsibleAttr = fmt.Sprintf(` data-collapsible="true" data-collapsed="%t"`, p.collapsed)
	}

	toolsHTML := ""
	if len(p.tools) > 0 {
		toolsHTML = `<div class="panel-tools">`
		for _, tool := range p.tools {
			toolsHTML += string(tool)
		}
		toolsHTML += `</div>`
	}

	iconHTML := ""
	if p.icon != "" {
		iconHTML = fmt.Sprintf(`<span class="panel-icon">%s</span>`, p.icon)
	}

	subtitleHTML := ""
	if p.subtitle != "" {
		subtitleHTML = fmt.Sprintf(`<span class="panel-subtitle">%s</span>`, template.HTMLEscapeString(p.subtitle))
	}

	footerHTML := ""
	if p.footer != "" {
		footerHTML = fmt.Sprintf(`<div class="panel-footer">%s</div>`, p.footer)
	}

	return template.HTML(fmt.Sprintf(
		`<div class="show-panel%s"%s>`+
			`<div class="panel-header">`+
			`<h3 class="panel-title">%s%s%s</h3>%s`+
			`</div>`+
			`<div class="panel-body">%s</div>%s`+
			`</div>`,
		collapsedClass, collapsibleAttr,
		iconHTML, template.HTMLEscapeString(p.title), subtitleHTML,
		toolsHTML,
		content.String(),
		footerHTML))
}

// ==================== Row Layout ====================

// Row creates a row layout for organizing fields horizontally
func (b *Builder) Row(fn func(*RowBuilder)) *Builder {
	rowBuilder := &RowBuilder{}
	fn(rowBuilder)

	b.Items = append(b.Items, &Item{
		Type: ItemRow,
		Formatter: func(record any, value any) template.HTML {
			return rowBuilder.Render(record)
		},
	})
	return b
}

// RowBuilder builds a row with columns
type RowBuilder struct {
	columns []*RowColumn
}

// RowColumn represents a column in a row
type RowColumn struct {
	Width int // 1-12, using grid system
	Items []*Item
}

// Column adds a column to the row
func (r *RowBuilder) Column(width int, fn func(*RowColumn)) *RowBuilder {
	col := &RowColumn{Width: width}
	fn(col)
	r.columns = append(r.columns, col)
	return r
}

// Field adds a field to the column
func (c *RowColumn) Field(name, label string) *Item {
	item := &Item{
		Type:  ItemField,
		Name:  name,
		Label: label,
	}
	c.Items = append(c.Items, item)
	return item
}

// Html adds HTML content to the column
func (c *RowColumn) Html(content template.HTML) *RowColumn {
	c.Items = append(c.Items, &Item{
		Type: ItemHtml,
		Formatter: func(record any, value any) template.HTML {
			return content
		},
	})
	return c
}

// Render generates the row HTML
func (r *RowBuilder) Render(record any) template.HTML {
	var columnsHTML strings.Builder

	for _, col := range r.columns {
		var colContent strings.Builder

		for _, item := range col.Items {
			switch item.Type {
			case ItemField:
				label := item.Label
				if label == "" {
					label = item.Name
				}
				value := ""
				if item.Formatter != nil {
					value = string(item.Formatter(record, nil))
				}
				colContent.WriteString(fmt.Sprintf(
					`<div class="field-row"><label class="field-label">%s</label><div class="field-value">%s</div></div>`,
					template.HTMLEscapeString(label), value))
			case ItemHtml:
				if item.Formatter != nil {
					colContent.WriteString(string(item.Formatter(record, nil)))
				}
			}
		}

		columnsHTML.WriteString(fmt.Sprintf(
			`<div class="row-column col-%d">%s</div>`,
			col.Width, colContent.String()))
	}

	return template.HTML(fmt.Sprintf(
		`<div class="show-row">%s</div>`,
		columnsHTML.String()))
}

// ==================== Relation Display ====================

// Relation displays related model data
func (b *Builder) Relation(name, label string, config *RelationConfig) *Item {
	item := &Item{
		Type:  ItemRelation,
		Name:  name,
		Label: label,
		Formatter: func(record any, value any) template.HTML {
			return config.Render(record, name)
		},
	}
	b.Items = append(b.Items, item)
	return item
}

// RelationConfig configures relation display
type RelationConfig struct {
	Title       string
	Fields      []string
	DisplayFunc func(record any) template.HTML
	GridConfig  *RelationGridConfig
}

// RelationGridConfig configures a grid display for relations
type RelationGridConfig struct {
	Columns    []string
	Actions    bool
	Pagination bool
}

// Render generates the relation HTML
func (r *RelationConfig) Render(record any, fieldName string) template.HTML {
	if r.DisplayFunc != nil {
		return r.DisplayFunc(record)
	}

	// Default relation display
	title := r.Title
	if title == "" {
		title = "Related Data"
	}

	return template.HTML(fmt.Sprintf(
		`<div class="show-relation" data-field="%s">`+
			`<div class="relation-header"><h4>%s</h4></div>`+
			`<div class="relation-content">Relation data loading...</div>`+
			`</div>`,
		fieldName, template.HTMLEscapeString(title)))
}

// HasOne configures a has-one relation display
func HasOne(title string) *RelationConfig {
	return &RelationConfig{
		Title:  title,
		Fields: []string{},
	}
}

// HasMany configures a has-many relation display
func HasMany(title string) *RelationConfig {
	return &RelationConfig{
		Title:  title,
		Fields: []string{},
		GridConfig: &RelationGridConfig{
			Columns:    []string{},
			Actions:    false,
			Pagination: false,
		},
	}
}

// WithFields sets the fields to display
func (r *RelationConfig) WithFields(fields ...string) *RelationConfig {
	r.Fields = fields
	return r
}

// Display sets a custom display function
func (r *RelationConfig) Display(fn func(record any) template.HTML) *RelationConfig {
	r.DisplayFunc = fn
	return r
}

// Grid configures grid display for has-many relations
func (r *RelationConfig) Grid(columns ...string) *RelationConfig {
	r.GridConfig = &RelationGridConfig{
		Columns:    columns,
		Actions:    true,
		Pagination: true,
	}
	return r
}
