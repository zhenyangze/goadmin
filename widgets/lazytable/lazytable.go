// Package lazytable provides lazy loading table widget
package lazytable

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

// LazyTable is a widget that displays a table with lazy loading
type LazyTable struct {
	id        string
	loadURL   string
	columns   []Column
	emptyText string
	pageSize  int
	striped   bool
	bordered  bool
	header    bool
	footer    bool
	class     string
}

// New creates a new lazy table widget
func New(id string) *LazyTable {
	return &LazyTable{
		id:        id,
		emptyText: "No data available",
		pageSize:  10,
		striped:   true,
		header:    true,
		footer:    false,
	}
}

// ID sets the table ID
func (t *LazyTable) ID(id string) *LazyTable {
	t.id = id
	return t
}

// LoadURL sets the data source URL
func (t *LazyTable) LoadURL(url string) *LazyTable {
	t.loadURL = url
	return t
}

// Column adds a column
func (t *LazyTable) Column(name, label string) *LazyTable {
	t.columns = append(t.columns, Column{Name: name, Label: label})
	return t
}

// EmptyText sets the empty text
func (t *LazyTable) EmptyText(text string) *LazyTable {
	t.emptyText = text
	return t
}

// PageSize sets the number of items per page
func (t *LazyTable) PageSize(size int) *LazyTable {
	t.pageSize = size
	return t
}

// Striped sets striped rows
func (t *LazyTable) Striped(striped bool) *LazyTable {
	t.striped = striped
	return t
}

// Bordered sets bordered style
func (t *LazyTable) Bordered(bordered bool) *LazyTable {
	t.bordered = bordered
	return t
}

// Class sets CSS classes
func (t *LazyTable) Class(class string) *LazyTable {
	t.class = class
	return t
}

// RenderContext provides data for rendering
type RenderContext struct {
	ID        string
	LoadURL   string
	Columns   []Column
	EmptyText string
	PageSize  int
	Striped   bool
	Bordered  bool
	Class     string
}

// Render generates HTML rendering context
func (t *LazyTable) Render() *RenderContext {
	return &RenderContext{
		ID:        t.id,
		LoadURL:   t.loadURL,
		Columns:   t.columns,
		EmptyText: t.emptyText,
		PageSize:  t.pageSize,
		Striped:   t.striped,
		Bordered:  t.bordered,
		Class:     t.class,
	}
}

// RenderTo generates HTML
func (t *LazyTable) RenderTo() template.HTML {
	ctx := t.Render()
	return ctx.ToHTML()
}

// ToHTML converts context to HTML
func (ctx *RenderContext) ToHTML() template.HTML {
	var sb strings.Builder

	tableClass := "min-w-full divide-y divide-gray-200"
	if ctx.Striped {
		tableClass += " striped"
	}
	if ctx.Bordered {
		tableClass += " border border-gray-200"
	}
	if ctx.Class != "" {
		tableClass += " " + ctx.Class
	}

	sb.WriteString(fmt.Sprintf(`<div id="%s" class="lazy-table" data-load-url="%s" data-page-size="%d">`,
		ctx.ID, ctx.LoadURL, ctx.PageSize))

	// Table skeleton for lazy loading
	sb.WriteString(fmt.Sprintf(`<table class="%s">`, tableClass))

	// Header
	sb.WriteString(`<thead class="bg-gray-50"><tr>`)
	for _, col := range ctx.Columns {
		widthAttr := ""
		if col.Width != "" {
			widthAttr = fmt.Sprintf(` style="width:%s"`, col.Width)
		}
		sortClass := ""
		if col.Sortable {
			sortClass = " cursor-pointer hover:bg-gray-100"
		}
		sb.WriteString(fmt.Sprintf(`<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider%s"%s>%s</th>`,
			sortClass, widthAttr, col.Label))
	}
	sb.WriteString(`</tr></thead>`)

	// Body with loading state
	sb.WriteString(`<tbody class="bg-white divide-y divide-gray-200">`)
	sb.WriteString(fmt.Sprintf(`<tr class="loading-row"><td colspan="%d" class="px-6 py-8 text-center text-gray-500">`, len(ctx.Columns)))
	sb.WriteString(`<div class="flex items-center justify-center space-x-2">`)
	sb.WriteString(`<div class="w-4 h-4 border-2 border-blue-600 border-t-transparent rounded-full animate-spin"></div>`)
	sb.WriteString(`<span>Loading...</span>`)
	sb.WriteString(`</div>`)
	sb.WriteString(`</td></tr>`)
	sb.WriteString(`</tbody>`)

	sb.WriteString(`</table>`)

	// Empty state (hidden by default)
	sb.WriteString(fmt.Sprintf(`<div class="empty-state hidden text-center py-8 text-gray-500">%s</div>`, ctx.EmptyText))

	sb.WriteString(`</div>`)

	// JavaScript for lazy loading
	sb.WriteString(fmt.Sprintf(`<script>
(function() {
	const table = document.getElementById('%s');
	const loadURL = table.dataset.loadUrl;
	const pageSize = parseInt(table.dataset.pageSize);

	function loadData() {
		fetch(loadURL + '?pageSize=' + pageSize)
			.then(r => r.json())
			.then(data => {
				const tbody = table.querySelector('tbody');
				const loadingRow = tbody.querySelector('.loading-row');
				if (loadingRow) loadingRow.remove();

				if (!data || data.length === 0) {
					table.querySelector('.empty-state').classList.remove('hidden');
					return;
				}

				data.forEach(row => {
					const tr = document.createElement('tr');
					tr.className = 'hover:bg-gray-50';
					tbody.appendChild(tr);
				});
			})
			.catch(err => {
				const tbody = table.querySelector('tbody');
				tbody.innerHTML = '<tr><td colspan="%d" class="px-6 py-4 text-center text-red-500">Failed to load data</td></tr>';
			});
	}

	setTimeout(loadData, 100);
})();
</script>`, len(ctx.Columns)))

	return template.HTML(sb.String())
}
