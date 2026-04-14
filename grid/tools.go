package grid

import (
	"fmt"
	"html/template"
	"net/url"
	"strings"
)

// ==================== Export Button Tool ====================

// ExportFormat describes supported export formats
type ExportFormat string

const (
	ExportExcel ExportFormat = "xlsx"
	ExportCSV   ExportFormat = "csv"
	ExportPDF   ExportFormat = "pdf"
)

// ExportButtonTool provides data export functionality
type ExportButtonTool struct {
	label    string
	formats  []ExportFormat
	filename string
}

// ExportButton creates a new export button tool
func ExportButton() *ExportButtonTool {
	return &ExportButtonTool{
		label:    "Export",
		formats:  []ExportFormat{ExportExcel, ExportCSV},
		filename: "export",
	}
}

// Label sets the button label
func (e *ExportButtonTool) Label(label string) *ExportButtonTool {
	e.label = label
	return e
}

// Formats sets available export formats
func (e *ExportButtonTool) Formats(formats ...ExportFormat) *ExportButtonTool {
	e.formats = formats
	return e
}

// Filename sets the export filename
func (e *ExportButtonTool) Filename(name string) *ExportButtonTool {
	e.filename = name
	return e
}

// Render implements Tool interface
func (e *ExportButtonTool) Render() template.HTML {
	if len(e.formats) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(`<div class="relative inline-block text-left" x-data="{ open: false }">`)
	sb.WriteString(fmt.Sprintf(
		`<button @click="open = !open" type="button" class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">`+
			`<svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">`+
			`<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"></path>`+
			`</svg>%s<svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>`+
			`</button>`,
		e.label))

	// Dropdown menu
	sb.WriteString(`<div x-show="open" @click.away="open = false" x-transition class="absolute right-0 mt-2 w-48 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 z-50" style="display: none;">`)
	sb.WriteString(`<div class="py-1">`)

	for _, format := range e.formats {
		sb.WriteString(fmt.Sprintf(
			`<a href="?export=%s&filename=%s" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 flex items-center">`+
				`<span class="w-6 h-6 mr-2 flex items-center justify-center bg-gray-100 rounded text-xs font-bold text-gray-600">%s</span>%s Export`+
				`</a>`,
			format, url.QueryEscape(e.filename), strings.ToUpper(string(format)), strings.ToUpper(string(format))))
	}

	sb.WriteString(`</div></div></div>`)

	return template.HTML(sb.String())
}

// Script implements Tool interface
func (e *ExportButtonTool) Script() string {
	return ""
}

// ==================== Column Selector Tool ====================

// ColumnSelectorTool allows users to show/hide columns
type ColumnSelectorTool struct {
	label   string
	columns []ColumnOption
}

// ColumnOption represents a column that can be toggled
type ColumnOption struct {
	Name    string
	Label   string
	Visible bool
}

// ColumnSelector creates a new column selector tool
func ColumnSelector() *ColumnSelectorTool {
	return &ColumnSelectorTool{
		label: "Columns",
	}
}

// Label sets the button label
func (c *ColumnSelectorTool) Label(label string) *ColumnSelectorTool {
	c.label = label
	return c
}

// Columns sets available columns
func (c *ColumnSelectorTool) Columns(columns ...ColumnOption) *ColumnSelectorTool {
	c.columns = columns
	return c
}

// Render implements Tool interface
func (c *ColumnSelectorTool) Render() template.HTML {
	if len(c.columns) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(`<div class="relative inline-block text-left" x-data="{ open: false, columns: [`)

	// Build Alpine.js data
	for i, col := range c.columns {
		if i > 0 {
			sb.WriteString(", ")
		}
		visible := "false"
		if col.Visible {
			visible = "true"
		}
		sb.WriteString(fmt.Sprintf(`{name: '%s', label: '%s', visible: %s}`,
			template.JSEscapeString(col.Name),
			template.JSEscapeString(col.Label),
			visible))
	}
	sb.WriteString(`] }">`)

	sb.WriteString(fmt.Sprintf(
		`<button @click="open = !open" type="button" class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">`+
			`<svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">`+
			`<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17V7m0 10a2 2 0 01-2 2H5a2 2 0 01-2-2V7a2 2 0 012-2h2a2 2 0 012 2m0 10a2 2 0 002 2h2a2 2 0 002-2M9 7a2 2 0 012-2h2a2 2 0 012 2m0 10V7m0 10a2 2 0 002 2h2a2 2 0 002-2V7a2 2 0 00-2-2h-2a2 2 0 00-2 2"></path>`+
			`</svg>%s<svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>`+
			`</button>`,
		c.label))

	// Dropdown menu
	sb.WriteString(`<div x-show="open" @click.away="open = false" x-transition class="absolute right-0 mt-2 w-64 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 z-50 p-4" style="display: none;">`)
	sb.WriteString(`<template x-for="col in columns" :key="col.name">`)
	sb.WriteString(`<div class="flex items-center mb-2">`)
	sb.WriteString(`<input type="checkbox" x-model="col.visible" class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded">`)
	sb.WriteString(`<span x-text="col.label" class="ml-2 text-sm text-gray-700"></span>`)
	sb.WriteString(`</div>`)
	sb.WriteString(`</template>`)
	sb.WriteString(`<div class="mt-3 pt-3 border-t border-gray-200 flex justify-between">`)
	sb.WriteString(`<button @click="columns.forEach(c => c.visible = true)" class="text-xs text-blue-600 hover:text-blue-800">Select All</button>`)
	sb.WriteString(`<button @click="columns.forEach(c => c.visible = false)" class="text-xs text-gray-500 hover:text-gray-700">Select None</button>`)
	sb.WriteString(`</div></div></div>`)

	return template.HTML(sb.String())
}

// Script implements Tool interface
func (c *ColumnSelectorTool) Script() string {
	return ""
}

// ==================== Refresh Button Tool ====================

// RefreshButtonTool provides grid refresh functionality
type RefreshButtonTool struct {
	label string
}

// RefreshButton creates a new refresh button tool
func RefreshButton() *RefreshButtonTool {
	return &RefreshButtonTool{
		label: "Refresh",
	}
}

// Label sets the button label
func (r *RefreshButtonTool) Label(label string) *RefreshButtonTool {
	r.label = label
	return r
}

// Render implements Tool interface
func (r *RefreshButtonTool) Render() template.HTML {
	return template.HTML(fmt.Sprintf(
		`<button type="button" onclick="window.location.reload()" class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">`+
			`<svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">`+
			`<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>`+
			`</svg>%s`+
			`</button>`,
		r.label))
}

// Script implements Tool interface
func (r *RefreshButtonTool) Script() string {
	return ""
}

// ==================== Per Page Selector Tool ====================

// PerPageSelectorTool allows users to select items per page
type PerPageSelectorTool struct {
	label   string
	options []int
	current int
}

// PerPageSelector creates a new per page selector tool
func PerPageSelector() *PerPageSelectorTool {
	return &PerPageSelectorTool{
		label:   "Per Page",
		options: []int{10, 20, 50, 100, 200},
		current: 20,
	}
}

// Label sets the label
func (p *PerPageSelectorTool) Label(label string) *PerPageSelectorTool {
	p.label = label
	return p
}

// Options sets available options
func (p *PerPageSelectorTool) Options(options ...int) *PerPageSelectorTool {
	p.options = options
	return p
}

// Current sets the current value
func (p *PerPageSelectorTool) Current(n int) *PerPageSelectorTool {
	p.current = n
	return p
}

// Render implements Tool interface
func (p *PerPageSelectorTool) Render() template.HTML {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(
		`<div class="inline-flex items-center"><span class="text-sm text-gray-600 mr-2">%s:</span>`+
			`<select onchange="window.location.search = updateQueryParam(window.location.search, 'per_page', this.value)" class="block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md">`,
		p.label))

	for _, opt := range p.options {
		selected := ""
		if opt == p.current {
			selected = " selected"
		}
		sb.WriteString(fmt.Sprintf(`<option value="%d"%s>%d</option>`, opt, selected, opt))
	}

	sb.WriteString(`</select></div>`)

	return template.HTML(sb.String())
}

// Script implements Tool interface
func (p *PerPageSelectorTool) Script() string {
	return `
function updateQueryParam(search, key, value) {
	const params = new URLSearchParams(search);
	params.set(key, value);
	params.delete('page'); // Reset to first page
	return params.toString();
}
`
}

// ==================== Quick Search Tool ====================

// QuickSearchTool provides a quick search input
type QuickSearchTool struct {
	placeholder string
	fields      []string
}

// QuickSearch creates a new quick search tool
func QuickSearch() *QuickSearchTool {
	return &QuickSearchTool{
		placeholder: "Quick Search...",
		fields:      []string{},
	}
}

// Placeholder sets the input placeholder
func (q *QuickSearchTool) Placeholder(text string) *QuickSearchTool {
	q.placeholder = text
	return q
}

// Fields sets the fields to search
func (q *QuickSearchTool) Fields(fields ...string) *QuickSearchTool {
	q.fields = fields
	return q
}

// Render implements Tool interface
func (q *QuickSearchTool) Render() template.HTML {
	return template.HTML(fmt.Sprintf(
		`<div class="relative">
			<div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
				<svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
				</svg>
			</div>
			<input type="text" name="quick_search" placeholder="%s"
				class="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md leading-5 bg-white placeholder-gray-500 focus:outline-none focus:placeholder-gray-400 focus:ring-1 focus:ring-blue-500 focus:border-blue-500 sm:text-sm">
		</div>`,
		q.placeholder))
}

// Script implements Tool interface
func (q *QuickSearchTool) Script() string {
	return ""
}

// ==================== Filter Button Tool ====================

// FilterButtonTool toggles the filter panel
type FilterButtonTool struct {
	label string
}

// FilterButton creates a new filter button tool
func FilterButton() *FilterButtonTool {
	return &FilterButtonTool{
		label: "Filter",
	}
}

// Label sets the button label
func (f *FilterButtonTool) Label(label string) *FilterButtonTool {
	f.label = label
	return f
}

// Render implements Tool interface
func (f *FilterButtonTool) Render() template.HTML {
	return template.HTML(fmt.Sprintf(
		`<button type="button" @click="showFilters = !showFilters" class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">`+
			`<svg class="w-4 h-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">`+
			`<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z"></path>`+
			`</svg>%s`+
			`</button>`,
		f.label))
}

// Script implements Tool interface
func (f *FilterButtonTool) Script() string {
	return ""
}

// ==================== Quick Create Button Tool ====================

// QuickCreateButtonTool provides quick create functionality
type QuickCreateButtonTool struct {
	label string
	url   string
}

// QuickCreateButton creates a new quick create button tool
func QuickCreateButton() *QuickCreateButtonTool {
	return &QuickCreateButtonTool{
		label: "Quick Create",
		url:   "create",
	}
}

// Label sets the button label
func (q *QuickCreateButtonTool) Label(label string) *QuickCreateButtonTool {
	q.label = label
	return q
}

// URL sets the create URL
func (q *QuickCreateButtonTool) URL(url string) *QuickCreateButtonTool {
	q.url = url
	return q
}

// Render implements Tool interface
func (q *QuickCreateButtonTool) Render() template.HTML {
	return template.HTML(fmt.Sprintf(
		`<a href="%s" class="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">`+
			`<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">`+
			`<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>`+
			`</svg>%s`+
			`</a>`,
		q.url, q.label))
}

// Script implements Tool interface
func (q *QuickCreateButtonTool) Script() string {
	return ""
}

// ==================== Batch Actions Tool ====================

// BatchActionsTool provides batch action buttons
type BatchActionsTool struct {
	label   string
	actions []BatchActionItem
}

// BatchActionItem represents a batch action
type BatchActionItem struct {
	Label   string
	URL     string
	Style   string
	Confirm string
}

// BatchActions creates a new batch actions tool
func BatchActions() *BatchActionsTool {
	return &BatchActionsTool{
		label: "Actions",
	}
}

// Label sets the label
func (b *BatchActionsTool) Label(label string) *BatchActionsTool {
	b.label = label
	return b
}

// Action adds a batch action
func (b *BatchActionsTool) Action(label, url string) *BatchActionsTool {
	b.actions = append(b.actions, BatchActionItem{
		Label: label,
		URL:   url,
		Style: "default",
	})
	return b
}

// ActionWithConfirm adds a batch action with confirmation
func (b *BatchActionsTool) ActionWithConfirm(label, url, confirm string) *BatchActionsTool {
	b.actions = append(b.actions, BatchActionItem{
		Label:   label,
		URL:     url,
		Style:   "danger",
		Confirm: confirm,
	})
	return b
}

// Render implements Tool interface
func (b *BatchActionsTool) Render() template.HTML {
	if len(b.actions) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(`<div class="relative inline-block text-left" x-data="{ open: false, selected: [] }">`)
	sb.WriteString(fmt.Sprintf(
		`<button @click="if (selected.length === 0) { alert('Please select items first'); return; } open = !open" type="button" class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">`+
			`%s<span class="ml-1 text-gray-400" x-text="selected.length > 0 ? '(' + selected.length + ')' : ''"></span>`+
			`<svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>`+
			`</button>`,
		b.label))

	// Dropdown menu
	sb.WriteString(`<div x-show="open" @click.away="open = false" x-transition class="absolute left-0 mt-2 w-48 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 z-50" style="display: none;">`)
	sb.WriteString(`<div class="py-1">`)

	for _, action := range b.actions {
		class := "block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
		if action.Style == "danger" {
			class = "block w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50"
		}

		confirm := ""
		if action.Confirm != "" {
			confirm = fmt.Sprintf(` onclick="if (!confirm('%s')) return false;"`, template.JSEscapeString(action.Confirm))
		}

		sb.WriteString(fmt.Sprintf(
			`<form method="POST" action="%s" class="m-0">`+
				`<input type="hidden" name="ids" :value="selected.join(',')">`+
				`<button type="submit" class="%s"%s>%s</button>`+
				`</form>`,
			action.URL, class, confirm, action.Label))
	}

	sb.WriteString(`</div></div></div>`)

	return template.HTML(sb.String())
}

// Script implements Tool interface
func (b *BatchActionsTool) Script() string {
	return ""
}
