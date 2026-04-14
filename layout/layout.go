// Package layout provides layout components for organizing content.
package layout

import (
	"fmt"
	"html/template"
)

// Column represents a column in a grid layout
type Column struct {
	Width   int // 1-12
	Offset  int // 0-11
	Content template.HTML
	Class   string
}

// NewColumn creates a new column with specified width (1-12)
func NewColumn(width int, content template.HTML) *Column {
	return &Column{
		Width:   width,
		Content: content,
	}
}

// WithOffset sets the column offset
func (c *Column) WithOffset(offset int) *Column {
	c.Offset = offset
	return c
}

// WithClass adds CSS classes
func (c *Column) WithClass(class string) *Column {
	c.Class = class
	return c
}

// Render generates the column HTML
func (c *Column) Render() template.HTML {
	offsetClass := ""
	if c.Offset > 0 {
		offsetClass = fmt.Sprintf(" col-offset-%d", c.Offset)
	}

	class := fmt.Sprintf("col-%d%s", c.Width, offsetClass)
	if c.Class != "" {
		class += " " + c.Class
	}

	return template.HTML(fmt.Sprintf(
		`<div class="%s">%s</div>`,
		class, c.Content))
}

// Row represents a row containing columns
type Row struct {
	Columns []*Column
	Class   string
	Gutter  string // small, medium, large
}

// NewRow creates a new row
func NewRow() *Row {
	return &Row{
		Columns: []*Column{},
	}
}

// AddColumn adds a column to the row
func (r *Row) AddColumn(col *Column) *Row {
	r.Columns = append(r.Columns, col)
	return r
}

// WithClass adds CSS classes
func (r *Row) WithClass(class string) *Row {
	r.Class = class
	return r
}

// WithGutter sets the gutter size
func (r *Row) WithGutter(size string) *Row {
	r.Gutter = size
	return r
}

// Render generates the row HTML
func (r *Row) Render() template.HTML {
	class := "row"
	if r.Class != "" {
		class += " " + r.Class
	}
	if r.Gutter != "" {
		class += " gutter-" + r.Gutter
	}

	var colsHTML string
	for _, col := range r.Columns {
		colsHTML += string(col.Render())
	}

	return template.HTML(fmt.Sprintf(
		`<div class="%s">%s</div>`,
		class, colsHTML))
}

// Content wraps content with padding and optional background
type Content struct {
	Content  template.HTML
	Padding  string // small, medium, large
	BgColor  string
	Bordered bool
	Shadow   bool
	Class    string
}

// NewContent creates new content wrapper
func NewContent(content template.HTML) *Content {
	return &Content{
		Content: content,
		Padding: "medium",
	}
}

// WithPadding sets padding size
func (c *Content) WithPadding(size string) *Content {
	c.Padding = size
	return c
}

// WithBgColor sets background color
func (c *Content) WithBgColor(color string) *Content {
	c.BgColor = color
	return c
}

// WithBorder adds border
func (c *Content) WithBorder() *Content {
	c.Bordered = true
	return c
}

// WithShadow adds shadow
func (c *Content) WithShadow() *Content {
	c.Shadow = true
	return c
}

// WithClass adds CSS classes
func (c *Content) WithClass(class string) *Content {
	c.Class = class
	return c
}

// Render generates the content HTML
func (c *Content) Render() template.HTML {
	class := "content-wrapper"
	if c.Padding != "" {
		class += " padding-" + c.Padding
	}
	if c.BgColor != "" {
		class += " bg-" + c.BgColor
	}
	if c.Bordered {
		class += " bordered"
	}
	if c.Shadow {
		class += " shadow"
	}
	if c.Class != "" {
		class += " " + c.Class
	}

	return template.HTML(fmt.Sprintf(
		`<div class="%s">%s</div>`,
		class, c.Content))
}

// Section represents a content section with header
type Section struct {
	Title       string
	Subtitle    string
	Content     template.HTML
	Actions     []template.HTML
	Collapsible bool
	Collapsed   bool
	Bordered    bool
	Class       string
}

// NewSection creates a new section
func NewSection(title string) *Section {
	return &Section{
		Title:   title,
		Actions: []template.HTML{},
	}
}

// WithSubtitle sets the subtitle
func (s *Section) WithSubtitle(subtitle string) *Section {
	s.Subtitle = subtitle
	return s
}

// WithContent sets the content
func (s *Section) WithContent(content template.HTML) *Section {
	s.Content = content
	return s
}

// AddAction adds an action button
func (s *Section) AddAction(action template.HTML) *Section {
	s.Actions = append(s.Actions, action)
	return s
}

// SetCollapsible makes the section collapsible
func (s *Section) SetCollapsible(collapsed bool) *Section {
	s.Collapsible = true
	s.Collapsed = collapsed
	return s
}

// WithBorder adds border
func (s *Section) WithBorder() *Section {
	s.Bordered = true
	return s
}

// WithClass adds CSS classes
func (s *Section) WithClass(class string) *Section {
	s.Class = class
	return s
}

// Render generates the section HTML
func (s *Section) Render() template.HTML {
	class := "section"
	if s.Bordered {
		class += " section-bordered"
	}
	if s.Collapsible {
		class += " collapsible"
		if s.Collapsed {
			class += " collapsed"
		}
	}
	if s.Class != "" {
		class += " " + s.Class
	}

	var actionsHTML string
	if len(s.Actions) > 0 {
		actionsHTML = `<div class="section-actions">`
		for _, action := range s.Actions {
			actionsHTML += string(action)
		}
		actionsHTML += `</div>`
	}

	subtitleHTML := ""
	if s.Subtitle != "" {
		subtitleHTML = fmt.Sprintf(`<p class="section-subtitle">%s</p>`, template.HTMLEscapeString(s.Subtitle))
	}

	toggleHTML := ""
	if s.Collapsible {
		toggleHTML = `<button type="button" class="section-toggle" onclick="toggleSection(this)">
			<svg class="icon" viewBox="0 0 24 24"><path d="M19 9l-7 7-7-7"/></svg>
		</button>`
	}

	return template.HTML(fmt.Sprintf(
		`<section class="%s">`+
			`<div class="section-header">`+
			`<div class="section-title-wrapper">`+
			`<h3 class="section-title">%s</h3>%s%s`+
			`</div>%s`+
			`</div>`+
			`<div class="section-body">%s</div>`+
			`</section>`,
		class,
		template.HTMLEscapeString(s.Title),
		subtitleHTML,
		toggleHTML,
		actionsHTML,
		s.Content))
}

// Responsive breakpoints
const (
	BreakpointSM = 640  // Small devices
	BreakpointMD = 768  // Medium devices
	BreakpointLG = 1024 // Large devices
	BreakpointXL = 1280 // Extra large devices
)

// ResponsiveColumn represents a responsive column
type ResponsiveColumn struct {
	XS      int // Extra small (default)
	SM      int // Small
	MD      int // Medium
	LG      int // Large
	XL      int // Extra large
	Offset  map[string]int
	Content template.HTML
	Class   string
}

// ResponsiveCol creates a responsive column
func ResponsiveCol(content template.HTML) *ResponsiveColumn {
	return &ResponsiveColumn{
		XS:      12,
		Content: content,
		Offset:  make(map[string]int),
	}
}

// SetSM sets small breakpoint
func (c *ResponsiveColumn) SetSM(width int) *ResponsiveColumn {
	c.SM = width
	return c
}

// SetMD sets medium breakpoint
func (c *ResponsiveColumn) SetMD(width int) *ResponsiveColumn {
	c.MD = width
	return c
}

// SetLG sets large breakpoint
func (c *ResponsiveColumn) SetLG(width int) *ResponsiveColumn {
	c.LG = width
	return c
}

// SetXL sets extra large breakpoint
func (c *ResponsiveColumn) SetXL(width int) *ResponsiveColumn {
	c.XL = width
	return c
}

// SetOffset sets offset for a breakpoint
func (c *ResponsiveColumn) SetOffset(bp string, offset int) *ResponsiveColumn {
	c.Offset[bp] = offset
	return c
}

// Render generates responsive column HTML
func (c *ResponsiveColumn) Render() template.HTML {
	class := fmt.Sprintf("col-%d", c.XS)

	if c.SM > 0 {
		class += fmt.Sprintf(" col-sm-%d", c.SM)
	}
	if c.MD > 0 {
		class += fmt.Sprintf(" col-md-%d", c.MD)
	}
	if c.LG > 0 {
		class += fmt.Sprintf(" col-lg-%d", c.LG)
	}
	if c.XL > 0 {
		class += fmt.Sprintf(" col-xl-%d", c.XL)
	}

	for bp, offset := range c.Offset {
		if offset > 0 {
			class += fmt.Sprintf(" %s-offset-%d", bp, offset)
		}
	}

	if c.Class != "" {
		class += " " + c.Class
	}

	return template.HTML(fmt.Sprintf(
		`<div class="%s">%s</div>`,
		class, c.Content))
}

// Layout JavaScript
func JavaScript() template.HTML {
	return template.HTML(`
<script>
(function() {
	window.toggleSection = function(btn) {
		const section = btn.closest('.section');
		section.classList.toggle('collapsed');
	};
})();
</script>`)
}
