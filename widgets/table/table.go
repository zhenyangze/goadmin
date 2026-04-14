// Package table provides data table widget
package table

import (
	"fmt"
	"html/template"
	"strings"
)

// Column represents a table column
type Column struct {
	Name     string
	Label    string
	Sortable bool
	Width    string
}

// Row represents a table row
type Row struct {
	Values map[string]template.HTML
	Class  string
}

// Table is a widget that displays data in a table
type Table struct {
	id        string
	columns   []Column
	rows      []Row
	emptyText string
	striped   bool
	bordered  bool
	hover     bool
	compact   bool
	class     string
}

// New creates a new table widget
func New() *Table {
	return &Table{
		emptyText: "No data available",
		striped:   true,
		hover:     true,
	}
}

// ID sets the table ID
func (t *Table) ID(id string) *Table {
	t.id = id
	return t
}

// Column adds a column
func (t *Table) Column(name, label string) *Table {
	t.columns = append(t.columns, Column{Name: name, Label: label})
	return t
}

// ColumnWithWidth adds a column with width
func (t *Table) ColumnWithWidth(name, label, width string) *Table {
	t.columns = append(t.columns, Column{Name: name, Label: label, Width: width})
	return t
}

// Row adds a data row
func (t *Table) Row(values map[string]template.HTML) *Table {
	t.rows = append(t.rows, Row{Values: values})
	return t
}

// Rows sets all rows
func (t *Table) Rows(rows []Row) *Table {
	t.rows = rows
	return t
}

// EmptyText sets the empty text
func (t *Table) EmptyText(text string) *Table {
	t.emptyText = text
	return t
}

// Striped sets striped rows
func (t *Table) Striped(striped bool) *Table {
	t.striped = striped
	return t
}

// Bordered sets bordered style
func (t *Table) Bordered(bordered bool) *Table {
	t.bordered = bordered
	return t
}

// Hover sets hover effect
func (t *Table) Hover(hover bool) *Table {
	t.hover = hover
	return t
}

// Compact sets compact style
func (t *Table) Compact(compact bool) *Table {
	t.compact = compact
	return t
}

// Class sets CSS classes
func (t *Table) Class(class string) *Table {
	t.class = class
	return t
}

// RenderContext provides data for rendering
type RenderContext struct {
	ID        string
	Columns   []Column
	Rows      []Row
	EmptyText string
	Striped   bool
	Bordered  bool
	Hover     bool
	Compact   bool
	Class     string
}

// Render generates HTML rendering context
func (t *Table) Render() *RenderContext {
	return &RenderContext{
		ID:        t.id,
		Columns:   t.columns,
		Rows:      t.rows,
		EmptyText: t.emptyText,
		Striped:   t.striped,
		Bordered:  t.bordered,
		Hover:     t.hover,
		Compact:   t.compact,
		Class:     t.class,
	}
}

// RenderTo generates HTML
func (t *Table) RenderTo() template.HTML {
	ctx := t.Render()
	return ctx.ToHTML()
}

// ToHTML converts context to HTML
func (ctx *RenderContext) ToHTML() template.HTML {
	var sb strings.Builder

	idAttr := ""
	if ctx.ID != "" {
		idAttr = fmt.Sprintf(` id="%s"`, ctx.ID)
	}

	tableClass := "min-w-full divide-y divide-gray-200"
	if ctx.Bordered {
		tableClass += " border border-gray-200"
	}
	if ctx.Class != "" {
		tableClass += " " + ctx.Class
	}

	sb.WriteString(fmt.Sprintf(`<div class="overflow-x-auto">`))
	sb.WriteString(fmt.Sprintf(`<table%s class="%s">`, idAttr, tableClass))

	// Header
	sb.WriteString(`<thead class="bg-gray-50"><tr>`)
	for _, col := range ctx.Columns {
		widthAttr := ""
		if col.Width != "" {
			widthAttr = fmt.Sprintf(` style="width:%s"`, col.Width)
		}
		thClass := "px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
		if ctx.Compact {
			thClass = "px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
		}
		sb.WriteString(fmt.Sprintf(`<th class="%s"%s>%s</th>`, thClass, widthAttr, col.Label))
	}
	sb.WriteString(`</tr></thead>`)

	// Body
	sb.WriteString(`<tbody class="bg-white divide-y divide-gray-200">`)

	if len(ctx.Rows) == 0 {
		sb.WriteString(fmt.Sprintf(`<tr><td colspan="%d" class="px-6 py-8 text-center text-gray-500">%s</td></tr>`,
			len(ctx.Columns), ctx.EmptyText))
	} else {
		for i, row := range ctx.Rows {
			trClass := ""
			if ctx.Striped && i%2 == 1 {
				trClass += " bg-gray-50"
			}
			if ctx.Hover {
				trClass += " hover:bg-gray-100"
			}
			if row.Class != "" {
				trClass += " " + row.Class
			}
			sb.WriteString(fmt.Sprintf(`<tr class="%s">`, strings.TrimSpace(trClass)))

			for _, col := range ctx.Columns {
				val := row.Values[col.Name]
				if val == "" {
					val = template.HTML(`<span class="text-gray-400">-</span>`)
				}
				tdClass := "px-6 py-4 whitespace-nowrap text-sm text-gray-900"
				if ctx.Compact {
					tdClass = "px-4 py-2 whitespace-nowrap text-sm text-gray-900"
				}
				sb.WriteString(fmt.Sprintf(`<td class="%s">%s</td>`, tdClass, val))
			}
			sb.WriteString(`</tr>`)
		}
	}

	sb.WriteString(`</tbody>`)
	sb.WriteString(`</table>`)
	sb.WriteString(`</div>`)

	return template.HTML(sb.String())
}

// SimpleTable creates a simple table from data
func SimpleTable(headers []string, data [][]string) *Table {
	t := New()
	for _, h := range headers {
		t.Column(h, h)
	}
	for _, row := range data {
		values := make(map[string]template.HTML)
		for i, val := range row {
			if i < len(headers) {
				values[headers[i]] = template.HTML(template.HTMLEscapeString(val))
			}
		}
		t.Row(values)
	}
	return t
}
