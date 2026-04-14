package grid

import (
	"context"
	"html/template"
	"strings"
)

// Formatter customizes grid cell rendering.
type Formatter func(record any, value any) template.HTML

// FilterKind describes a supported filter input.
type FilterKind string

const (
	FilterText   FilterKind = "text"
	FilterSelect FilterKind = "select"
)

// Option is used by select-like inputs.
type Option struct {
	Value string
	Label string
}

// ActionStyle describes the visual intent of an action.
type ActionStyle string

const (
	ActionDefault ActionStyle = "default"
	ActionPrimary ActionStyle = "primary"
	ActionGhost   ActionStyle = "ghost"
	ActionDanger  ActionStyle = "danger"
)

// PageAction describes a toolbar action above the grid.
type PageAction struct {
	Label   string
	URL     string
	Style   ActionStyle
	Method  string
	Confirm string
}

// RowActionURL builds a row action target from the current record.
type RowActionURL func(record any) string

// RowAction describes an action rendered for each row.
type RowAction struct {
	Label   string
	URL     RowActionURL
	Style   ActionStyle
	Method  string
	Confirm string
}

// BatchAction describes a batch action for selected rows.
type BatchAction struct {
	Label   string
	URL     string
	Style   ActionStyle
	Method  string
	Confirm string
}

// Filter describes a grid filter field.
type Filter struct {
	Name    string
	Label   string
	Kind    FilterKind
	Options []Option
}

// Column describes one grid column.
type Column struct {
	Name      string
	Label     string
	Sortable  bool
	Formatter Formatter
}

// SortableColumn marks the column as sortable.
func (c *Column) SortableColumn() *Column {
	c.Sortable = true
	return c
}

// Display sets a custom formatter.
func (c *Column) Display(fn Formatter) *Column {
	c.Formatter = fn
	return c
}

// Tool is the interface for custom toolbar tools (AbstractTool equivalent).
type Tool interface {
	// Render returns the HTML for the tool
	Render() template.HTML
	// Script returns JavaScript code for the tool (optional)
	Script() string
}

// ToolFunc is a function that implements the Tool interface.
type ToolFunc struct {
	RenderFn func() template.HTML
	ScriptFn func() string
}

// Render returns the HTML for the tool.
func (t ToolFunc) Render() template.HTML {
	if t.RenderFn != nil {
		return t.RenderFn()
	}
	return ""
}

// Script returns JavaScript code for the tool.
func (t ToolFunc) Script() string {
	if t.ScriptFn != nil {
		return t.ScriptFn()
	}
	return ""
}

// BatchHandler is the callback for processing batch actions.
type BatchHandler func(ctx context.Context, ids []string) error

// BatchActionWithHandler extends BatchAction with server-side handler.
type BatchActionWithHandler struct {
	Label   string
	Style   ActionStyle
	Method  string
	Confirm string
	Handler BatchHandler
}

// Builder defines the grid page.
type Builder struct {
	Title              string
	Description        string
	QuickSearches      []string
	Columns            []*Column
	Filters            []*Filter
	CreateLabel        string
	DisableCreate      bool
	DisableView        bool
	DisableEdit        bool
	DisableDelete      bool
	DisableBatchDelete bool
	DisableRowSelector bool
	PageActions        []*PageAction
	RowActions         []*RowAction
	BatchActions       []*BatchAction
	BatchActionHandlers []*BatchActionWithHandler
	// Toolbar tools
	Tools           []Tool
	ToolsWithOutline bool
	DisableRefresh   bool
	// Dialog form settings
	EnableDialogCreate bool
	EnableDialogEdit   bool
	DialogWidth        string
	DialogHeight       string
	// Pagination settings
	DisablePagination bool
	SimplePaginate    bool
	PerPageOptions    []int
	// Display settings
	ScrollbarX      bool
	TableClasses    []string
	// Row selector settings
	RowSelectorTitleColumn string
	RowSelectorIDColumn    string
	RowSelectorChecked     func(record any) bool
	RowSelectorDisabled    func(record any) bool
	RowSelectorClick       bool
	// Inline edit
	EnableQuickEdit bool
}

// New creates a grid builder.
func New() *Builder {
	return &Builder{CreateLabel: "Create"}
}

// Column adds a generic column.
func (b *Builder) Column(name, label string) *Column {
	col := &Column{Name: name, Label: label}
	b.Columns = append(b.Columns, col)
	return col
}

// ID adds an ID column.
func (b *Builder) ID(label string) *Column {
	return b.Column("ID", label)
}

// QuickSearch enables a single keyword search over fields.
func (b *Builder) QuickSearch(fields ...string) {
	b.QuickSearches = append(b.QuickSearches, fields...)
}

// Filter adds a filter input.
func (b *Builder) Filter(name, label string, kind FilterKind) *Filter {
	filter := &Filter{Name: name, Label: label, Kind: kind}
	b.Filters = append(b.Filters, filter)
	return filter
}

// PageAction adds a top-level grid action.
func (b *Builder) PageAction(label, url string) *PageAction {
	action := &PageAction{
		Label:  label,
		URL:    url,
		Style:  ActionGhost,
		Method: "GET",
	}
	b.PageActions = append(b.PageActions, action)
	return action
}

// RowAction adds a per-row action.
func (b *Builder) RowAction(label string, url RowActionURL) *RowAction {
	action := &RowAction{
		Label:  label,
		URL:    url,
		Style:  ActionGhost,
		Method: "GET",
	}
	b.RowActions = append(b.RowActions, action)
	return action
}

// BatchAction adds a batch action for selected rows.
func (b *Builder) BatchAction(label, url string) *BatchAction {
	action := &BatchAction{
		Label:  label,
		URL:    url,
		Style:  ActionGhost,
		Method: "POST",
	}
	b.BatchActions = append(b.BatchActions, action)
	return action
}

// EnableDialogCreate enables dialog form for creating.
func (b *Builder) UseDialogCreate() *Builder {
	b.EnableDialogCreate = true
	return b
}

// EnableDialogEdit enables dialog form for editing.
func (b *Builder) UseDialogEdit() *Builder {
	b.EnableDialogEdit = true
	return b
}

// SetDialogFormDimensions sets the dialog form dimensions.
func (b *Builder) SetDialogFormDimensions(width, height string) *Builder {
	b.DialogWidth = width
	b.DialogHeight = height
	return b
}

// DisableBatchDelete disables batch delete button.
func (b *Builder) HideBatchDelete() *Builder {
	b.DisableBatchDelete = true
	return b
}

// DisableRowSelector disables row selector checkbox.
func (b *Builder) HideRowSelector() *Builder {
	b.DisableRowSelector = true
	return b
}

// DisablePagination disables pagination.
func (b *Builder) HidePagination() *Builder {
	b.DisablePagination = true
	return b
}

// SimplePaginate enables simple paginate mode.
func (b *Builder) UseSimplePaginate() *Builder {
	b.SimplePaginate = true
	return b
}

// PerPages sets the per page options.
func (b *Builder) SetPerPages(options []int) *Builder {
	b.PerPageOptions = options
	return b
}

// ScrollbarX enables horizontal scrollbar.
func (b *Builder) UseScrollbarX() *Builder {
	b.ScrollbarX = true
	return b
}

// AddTableClass adds CSS classes to the table.
func (b *Builder) AddTableClass(classes ...string) *Builder {
	b.TableClasses = append(b.TableClasses, classes...)
	return b
}

// RowSelector configures row selector.
func (b *Builder) RowSelector() *RowSelector {
	return &RowSelector{builder: b}
}

// RowSelector provides row selector configuration.
type RowSelector struct {
	builder *Builder
}

// TitleColumn sets the title column for row selector.
func (r *RowSelector) TitleColumn(column string) *RowSelector {
	r.builder.RowSelectorTitleColumn = column
	return r
}

// IDColumn sets the ID column for row selector.
func (r *RowSelector) IDColumn(column string) *RowSelector {
	r.builder.RowSelectorIDColumn = column
	return r
}

// Checked sets the checked callback for row selector.
func (r *RowSelector) Checked(fn func(record any) bool) *RowSelector {
	r.builder.RowSelectorChecked = fn
	return r
}

// Disable sets the disabled callback for row selector.
func (r *RowSelector) Disable(fn func(record any) bool) *RowSelector {
	r.builder.RowSelectorDisabled = fn
	return r
}

// Click enables click to select row.
func (r *RowSelector) Click() *RowSelector {
	r.builder.RowSelectorClick = true
	return r
}

// EnableQuickEdit enables quick inline edit.
func (b *Builder) UseQuickEdit() *Builder {
	b.EnableQuickEdit = true
	return b
}

// AddTool adds a custom toolbar tool.
func (b *Builder) AddTool(tool Tool) *Builder {
	b.Tools = append(b.Tools, tool)
	return b
}

// UseToolsWithOutline enables outline style for toolbar buttons.
func (b *Builder) UseToolsWithOutline() *Builder {
	b.ToolsWithOutline = true
	return b
}

// HideRefresh hides the built-in refresh button.
func (b *Builder) HideRefresh() *Builder {
	b.DisableRefresh = true
	return b
}

// BatchActionWithHandler adds a batch action with server-side handler.
func (b *Builder) BatchActionWithHandler(label string, handler BatchHandler) *BatchActionWithHandler {
	action := &BatchActionWithHandler{
		Label:   label,
		Style:   ActionGhost,
		Method:  "POST",
		Handler: handler,
	}
	b.BatchActionHandlers = append(b.BatchActionHandlers, action)
	return action
}

// WithStyle overrides the action style.
func (a *BatchActionWithHandler) WithStyle(style ActionStyle) *BatchActionWithHandler {
	a.Style = style
	return a
}

// WithMethod overrides the HTTP method.
func (a *BatchActionWithHandler) WithMethod(method string) *BatchActionWithHandler {
	if method != "" {
		a.Method = strings.ToUpper(method)
	}
	return a
}

// WithConfirm adds a confirmation prompt.
func (a *BatchActionWithHandler) WithConfirm(message string) *BatchActionWithHandler {
	a.Confirm = message
	return a
}

// WithOptions sets select options.
func (f *Filter) WithOptions(options ...Option) *Filter {
	f.Options = append(f.Options, options...)
	return f
}

// WithStyle overrides the action style.
func (a *PageAction) WithStyle(style ActionStyle) *PageAction {
	a.Style = style
	return a
}

// WithMethod overrides the HTTP method.
func (a *PageAction) WithMethod(method string) *PageAction {
	if method != "" {
		a.Method = strings.ToUpper(method)
	}
	return a
}

// WithConfirm adds a confirmation prompt.
func (a *PageAction) WithConfirm(message string) *PageAction {
	a.Confirm = message
	return a
}

// WithStyle overrides the action style.
func (a *RowAction) WithStyle(style ActionStyle) *RowAction {
	a.Style = style
	return a
}

// WithMethod overrides the HTTP method.
func (a *RowAction) WithMethod(method string) *RowAction {
	if method != "" {
		a.Method = strings.ToUpper(method)
	}
	return a
}

// WithConfirm adds a confirmation prompt.
func (a *RowAction) WithConfirm(message string) *RowAction {
	a.Confirm = message
	return a
}

// WithStyle overrides the action style for batch action.
func (a *BatchAction) WithStyle(style ActionStyle) *BatchAction {
	a.Style = style
	return a
}

// WithMethod overrides the HTTP method for batch action.
func (a *BatchAction) WithMethod(method string) *BatchAction {
	if method != "" {
		a.Method = strings.ToUpper(method)
	}
	return a
}

// WithConfirm adds a confirmation prompt for batch action.
func (a *BatchAction) WithConfirm(message string) *BatchAction {
	a.Confirm = message
	return a
}
