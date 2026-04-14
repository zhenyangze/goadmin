package grid

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"
	"strings"
)

// Displayer is the interface for grid cell displayers
type Displayer interface {
	// Display returns the HTML representation of the value
	Display(record any, value any) template.HTML
}

// DisplayerFunc is a function type that implements Displayer
type DisplayerFunc func(record any, value any) template.HTML

// Display implements the Displayer interface
func (f DisplayerFunc) Display(record any, value any) template.HTML {
	return f(record, value)
}

// ToFormatter converts a Displayer to a Formatter
func ToFormatter(d Displayer) Formatter {
	return func(record any, value any) template.HTML {
		return d.Display(record, value)
	}
}

// Display sets a displayer for the column
func (c *Column) Displayer(d Displayer) *Column {
	c.Formatter = ToFormatter(d)
	return c
}

// ==================== Badge Displayer ====================

// Badge styles
const (
	BadgeStyleDefault BadgeStyle = "default"
	BadgeStylePrimary BadgeStyle = "primary"
	BadgeStyleSuccess BadgeStyle = "success"
	BadgeStyleWarning BadgeStyle = "warning"
	BadgeStyleDanger  BadgeStyle = "danger"
	BadgeStyleInfo    BadgeStyle = "info"
)

// BadgeStyle describes the visual style of a badge
type BadgeStyle string

// BadgeDisplayer displays values as badges
type BadgeDisplayer struct {
	style     BadgeStyle
	colorFunc func(record any, value any) BadgeStyle
}

// Badge creates a new badge displayer
func Badge() *BadgeDisplayer {
	return &BadgeDisplayer{
		style: BadgeStyleDefault,
	}
}

// Style sets the badge style
func (b *BadgeDisplayer) Style(style BadgeStyle) *BadgeDisplayer {
	b.style = style
	return b
}

// Color sets a dynamic color function
func (b *BadgeDisplayer) Color(fn func(record any, value any) BadgeStyle) *BadgeDisplayer {
	b.colorFunc = fn
	return b
}

// Display implements the Displayer interface
func (b *BadgeDisplayer) Display(record any, value any) template.HTML {
	style := b.style
	if b.colorFunc != nil {
		style = b.colorFunc(record, value)
	}

	class := fmt.Sprintf("inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-%s-100 text-%s-800",
		style, style)

	return template.HTML(fmt.Sprintf(`<span class="%s">%v</span>`, class, value))
}

// ==================== Label Displayer ====================

// Label styles
const (
	LabelStyleDefault LabelStyle = "default"
	LabelStylePrimary LabelStyle = "primary"
	LabelStyleSuccess LabelStyle = "success"
	LabelStyleWarning LabelStyle = "warning"
	LabelStyleDanger  LabelStyle = "danger"
	LabelStyleInfo    LabelStyle = "info"
)

// LabelStyle describes the visual style of a label
type LabelStyle string

// LabelDisplayer displays values as labels/tags
type LabelDisplayer struct {
	style     LabelStyle
	colorFunc func(record any, value any) LabelStyle
}

// Label creates a new label displayer
func Label() *LabelDisplayer {
	return &LabelDisplayer{
		style: LabelStyleDefault,
	}
}

// Style sets the label style
func (l *LabelDisplayer) Style(style LabelStyle) *LabelDisplayer {
	l.style = style
	return l
}

// Color sets a dynamic color function
func (l *LabelDisplayer) Color(fn func(record any, value any) LabelStyle) *LabelDisplayer {
	l.colorFunc = fn
	return l
}

// Display implements the Displayer interface
func (l *LabelDisplayer) Display(record any, value any) template.HTML {
	style := l.style
	if l.colorFunc != nil {
		style = l.colorFunc(record, value)
	}

	class := fmt.Sprintf("inline-flex items-center px-2 py-0.5 rounded text-sm font-medium bg-%s-100 text-%s-800",
		style, style)

	return template.HTML(fmt.Sprintf(`<span class="%s">%v</span>`, class, value))
}

// ==================== Image Displayer ====================

// ImageDisplayer displays image thumbnails
type ImageDisplayer struct {
	width     string
	height    string
	preview   bool
	imageFunc func(record any, value any) string
}

// Image creates a new image displayer
func Image() *ImageDisplayer {
	return &ImageDisplayer{
		width:   "50",
		height:  "50",
		preview: true,
	}
}

// Size sets the image size
func (i *ImageDisplayer) Size(width, height int) *ImageDisplayer {
	i.width = fmt.Sprintf("%d", width)
	i.height = fmt.Sprintf("%d", height)
	return i
}

// Preview enables/disables image preview on click
func (i *ImageDisplayer) Preview(enabled bool) *ImageDisplayer {
	i.preview = enabled
	return i
}

// Src sets a custom image source function
func (i *ImageDisplayer) Src(fn func(record any, value any) string) *ImageDisplayer {
	i.imageFunc = fn
	return i
}

// Display implements the Displayer interface
func (i *ImageDisplayer) Display(record any, value any) template.HTML {
	src := fmt.Sprintf("%v", value)
	if i.imageFunc != nil {
		src = i.imageFunc(record, value)
	}

	if src == "" || src == "<nil>" {
		return template.HTML(`<span class="text-gray-400">-</span>`)
	}

	imgHTML := fmt.Sprintf(
		`<img src="%s" class="rounded object-cover" style="width:%spx;height:%spx;" alt="">`,
		src, i.width, i.height)

	if i.preview {
		return template.HTML(fmt.Sprintf(
			`<a href="%s" target="_blank" class="inline-block hover:opacity-80 transition-opacity">%s</a>`,
			src, imgHTML))
	}

	return template.HTML(imgHTML)
}

// ==================== Link Displayer ====================

// LinkDisplayer displays values as clickable links
type LinkDisplayer struct {
	urlFunc   func(record any, value any) string
	target    string
	maxLength int
}

// Link creates a new link displayer
func Link() *LinkDisplayer {
	return &LinkDisplayer{
		target:    "_self",
		maxLength: 0,
	}
}

// URL sets the link URL function
func (l *LinkDisplayer) URL(fn func(record any, value any) string) *LinkDisplayer {
	l.urlFunc = fn
	return l
}

// Target sets the link target
func (l *LinkDisplayer) Target(target string) *LinkDisplayer {
	l.target = target
	return l
}

// MaxLength sets the maximum display length (0 = no limit)
func (l *LinkDisplayer) MaxLength(length int) *LinkDisplayer {
	l.maxLength = length
	return l
}

// Display implements the Displayer interface
func (l *LinkDisplayer) Display(record any, value any) template.HTML {
	displayValue := fmt.Sprintf("%v", value)
	if displayValue == "<nil>" {
		return template.HTML(`<span class="text-gray-400">-</span>`)
	}

	// Truncate if needed
	originalValue := displayValue
	if l.maxLength > 0 && len(displayValue) > l.maxLength {
		displayValue = displayValue[:l.maxLength] + "..."
	}

	href := "#"
	if l.urlFunc != nil {
		href = l.urlFunc(record, value)
	} else {
		// Default: use the value as URL
		href = originalValue
	}

	return template.HTML(fmt.Sprintf(
		`<a href="%s" target="%s" class="text-blue-600 hover:text-blue-800 hover:underline">%s</a>`,
		href, l.target, displayValue))
}

// ==================== ProgressBar Displayer ====================

// ProgressBarDisplayer displays values as progress bars
type ProgressBarDisplayer struct {
	min       float64
	max       float64
	colorFunc func(record any, value any) string
	showText  bool
}

// ProgressBar creates a new progress bar displayer
func ProgressBar() *ProgressBarDisplayer {
	return &ProgressBarDisplayer{
		min:      0,
		max:      100,
		showText: true,
	}
}

// Range sets the min/max range
func (p *ProgressBarDisplayer) Range(min, max float64) *ProgressBarDisplayer {
	p.min = min
	p.max = max
	return p
}

// Color sets a dynamic color function
func (p *ProgressBarDisplayer) Color(fn func(record any, value any) string) *ProgressBarDisplayer {
	p.colorFunc = fn
	return p
}

// ShowText enables/disables text display
func (p *ProgressBarDisplayer) ShowText(show bool) *ProgressBarDisplayer {
	p.showText = show
	return p
}

// Display implements the Displayer interface
func (p *ProgressBarDisplayer) Display(record any, value any) template.HTML {
	var numValue float64
	switch v := value.(type) {
	case float64:
		numValue = v
	case float32:
		numValue = float64(v)
	case int:
		numValue = float64(v)
	case int64:
		numValue = float64(v)
	default:
		numValue = 0
	}

	// Calculate percentage
	percentage := 0.0
	if p.max > p.min {
		percentage = ((numValue - p.min) / (p.max - p.min)) * 100
	}
	if percentage < 0 {
		percentage = 0
	}
	if percentage > 100 {
		percentage = 100
	}

	// Determine color
	color := "blue"
	if p.colorFunc != nil {
		color = p.colorFunc(record, value)
	} else {
		// Auto color based on percentage
		if percentage < 30 {
			color = "red"
		} else if percentage < 70 {
			color = "yellow"
		} else {
			color = "green"
		}
	}

	text := ""
	if p.showText {
		text = fmt.Sprintf(`<span class="text-xs text-gray-600 ml-2">%.1f%%</span>`, percentage)
	}

	return template.HTML(fmt.Sprintf(
		`<div class="flex items-center w-full max-w-xs">`+
			`<div class="flex-1 h-2 bg-gray-200 rounded-full overflow-hidden">`+
			`<div class="h-full bg-%s-500 rounded-full" style="width:%.1f%%"></div>`+
			`</div>%s`+
			`</div>`,
		color, percentage, text))
}

// ==================== SwitchDisplay Displayer ====================

// SwitchDisplayer displays boolean values as switches
type SwitchDisplayer struct {
	onText  string
	offText string
	onColor string
	offColor string
}

// SwitchDisplay creates a new switch displayer
func SwitchDisplay() *SwitchDisplayer {
	return &SwitchDisplayer{
		onText:   "ON",
		offText:  "OFF",
		onColor:  "green",
		offColor: "gray",
	}
}

// Text sets the on/off text
func (s *SwitchDisplayer) Text(on, off string) *SwitchDisplayer {
	s.onText = on
	s.offText = off
	return s
}

// Color sets the on/off colors
func (s *SwitchDisplayer) Color(on, off string) *SwitchDisplayer {
	s.onColor = on
	s.offColor = off
	return s
}

// Display implements the Displayer interface
func (s *SwitchDisplayer) Display(record any, value any) template.HTML {
	isOn := false
	switch v := value.(type) {
	case bool:
		isOn = v
	case int, int64:
		isOn = v != 0
	case string:
		isOn = v == "1" || v == "true" || v == "yes" || v == "on"
	}

	if isOn {
		return template.HTML(fmt.Sprintf(
			`<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-%s-100 text-%s-800">`+
				`<span class="w-2 h-2 bg-%s-500 rounded-full mr-1.5"></span>%s`+
				`</span>`,
			s.onColor, s.onColor, s.onColor, s.onText))
	}

	return template.HTML(fmt.Sprintf(
		`<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-%s-100 text-%s-800">`+
			`<span class="w-2 h-2 bg-%s-500 rounded-full mr-1.5"></span>%s`+
			`</span>`,
		s.offColor, s.offColor, s.offColor, s.offText))
}

// ==================== QRCode Displayer ====================

// QRCodeDisplayer displays QR codes
type QRCodeDisplayer struct {
	size  int
	textFunc func(record any, value any) string
}

// QRCode creates a new QR code displayer
func QRCode() *QRCodeDisplayer {
	return &QRCodeDisplayer{
		size: 64,
	}
}

// Size sets the QR code size
func (q *QRCodeDisplayer) Size(size int) *QRCodeDisplayer {
	q.size = size
	return q
}

// Text sets a custom text function for the QR code
func (q *QRCodeDisplayer) Text(fn func(record any, value any) string) *QRCodeDisplayer {
	q.textFunc = fn
	return q
}

// Display implements the Displayer interface
func (q *QRCodeDisplayer) Display(record any, value any) template.HTML {
	text := fmt.Sprintf("%v", value)
	if q.textFunc != nil {
		text = q.textFunc(record, value)
	}

	// Use a simple QR code generation service or inline SVG
	// For now, use a data URI with a placeholder that can be replaced with actual QR generation
	escapedText := url.QueryEscape(text)
	imgURL := fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=%dx%d&data=%s",
		q.size, q.size, escapedText)

	return template.HTML(fmt.Sprintf(
		`<img src="%s" width="%d" height="%d" alt="QR Code" class="rounded" title="%s">`,
		imgURL, q.size, q.size, template.HTMLEscapeString(text)))
}

// ==================== Copyable Displayer ====================

// CopyableDisplayer displays values with a copy button
type CopyableDisplayer struct {
	maxLength int
	showIcon  bool
}

// Copyable creates a new copyable displayer
func Copyable() *CopyableDisplayer {
	return &CopyableDisplayer{
		maxLength: 50,
		showIcon:  true,
	}
}

// MaxLength sets the maximum display length
func (c *CopyableDisplayer) MaxLength(length int) *CopyableDisplayer {
	c.maxLength = length
	return c
}

// ShowIcon enables/disables the copy icon
func (c *CopyableDisplayer) ShowIcon(show bool) *CopyableDisplayer {
	c.showIcon = show
	return c
}

// Display implements the Displayer interface
func (c *CopyableDisplayer) Display(record any, value any) template.HTML {
	text := fmt.Sprintf("%v", value)
	if text == "<nil>" {
		return template.HTML(`<span class="text-gray-400">-</span>`)
	}

	displayText := text
	if c.maxLength > 0 && len(text) > c.maxLength {
		displayText = text[:c.maxLength] + "..."
	}

	iconHTML := ""
	if c.showIcon {
		iconHTML = `<svg class="w-4 h-4 ml-1 opacity-0 group-hover:opacity-100 transition-opacity" fill="none" stroke="currentColor" viewBox="0 0 24 24">` +
			`<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"></path>` +
			`</svg>`
	}

	return template.HTML(fmt.Sprintf(
		`<span class="group inline-flex items-center cursor-pointer hover:text-blue-600" `+
			`onclick="navigator.clipboard.writeText('%s');this.classList.add('text-green-600');setTimeout(()=>this.classList.remove('text-green-600'),500)">`+
			`%s%s`+
			`</span>`,
		template.JSEscapeString(text), displayText, iconHTML))
}

// ==================== Limit Displayer ====================

// LimitDisplayer truncates text with tooltip
type LimitDisplayer struct {
	limit     int
	tooltip   bool
	replace   string
}

// Limit creates a new limit displayer
func Limit(length int) *LimitDisplayer {
	return &LimitDisplayer{
		limit:   length,
		tooltip: true,
		replace: "...",
	}
}

// Tooltip enables/disables tooltip
func (l *LimitDisplayer) Tooltip(enabled bool) *LimitDisplayer {
	l.tooltip = enabled
	return l
}

// Replace sets the replacement string
func (l *LimitDisplayer) Replace(replacement string) *LimitDisplayer {
	l.replace = replacement
	return l
}

// Display implements the Displayer interface
func (l *LimitDisplayer) Display(record any, value any) template.HTML {
	text := fmt.Sprintf("%v", value)
	if text == "<nil>" {
		return template.HTML(`<span class="text-gray-400">-</span>`)
	}

	if len(text) <= l.limit {
		return template.HTML(template.HTMLEscapeString(text))
	}

	truncated := text[:l.limit] + l.replace

	if l.tooltip {
		return template.HTML(fmt.Sprintf(
			`<span title="%s" class="cursor-help border-b border-dotted border-gray-400">%s</span>`,
			template.HTMLEscapeString(text), template.HTMLEscapeString(truncated)))
	}

	return template.HTML(template.HTMLEscapeString(truncated))
}

// ==================== Table Displayer (for nested data) ====================

// TableDisplayer displays nested data as a small table
type TableDisplayer struct {
	columns   []string
	keyFunc   func(record any) string
	dataFunc  func(record any) []map[string]any
	maxRows   int
}

// Table creates a new table displayer
func Table() *TableDisplayer {
	return &TableDisplayer{
		maxRows: 5,
	}
}

// Columns sets the column keys to display
func (t *TableDisplayer) Columns(keys ...string) *TableDisplayer {
	t.columns = keys
	return t
}

// Data sets the data function
func (t *TableDisplayer) Data(fn func(record any) []map[string]any) *TableDisplayer {
	t.dataFunc = fn
	return t
}

// MaxRows sets the maximum number of rows to display
func (t *TableDisplayer) MaxRows(n int) *TableDisplayer {
	t.maxRows = n
	return t
}

// Display implements the Displayer interface
func (t *TableDisplayer) Display(record any, value any) template.HTML {
	var data []map[string]any

	// Try to parse value as JSON array
	if str, ok := value.(string); ok && str != "" {
		if err := json.Unmarshal([]byte(str), &data); err != nil {
			// Try as single object
			var single map[string]any
			if err := json.Unmarshal([]byte(str), &single); err == nil {
				data = []map[string]any{single}
			}
		}
	} else if arr, ok := value.([]map[string]any); ok {
		data = arr
	} else if t.dataFunc != nil {
		data = t.dataFunc(record)
	}

	if len(data) == 0 {
		return template.HTML(`<span class="text-gray-400">-</span>`)
	}

	// Determine columns if not set
	cols := t.columns
	if len(cols) == 0 && len(data) > 0 {
		for k := range data[0] {
			cols = append(cols, k)
		}
	}

	// Build table HTML
	var sb strings.Builder
	sb.WriteString(`<table class="min-w-full text-xs border border-gray-200">`)

	// Header
	sb.WriteString(`<thead class="bg-gray-50"><tr>`)
	for _, col := range cols {
		sb.WriteString(fmt.Sprintf(`<th class="px-2 py-1 text-left font-medium text-gray-500 border-b">%s</th>`, col))
	}
	sb.WriteString(`</tr></thead>`)

	// Body
	sb.WriteString(`<tbody>`)
	for i, row := range data {
		if i >= t.maxRows {
			sb.WriteString(fmt.Sprintf(`<tr><td colspan="%d" class="px-2 py-1 text-gray-500 italic">... %d more</td></tr>`,
				len(cols), len(data)-t.maxRows))
			break
		}
		sb.WriteString(`<tr class="border-b border-gray-100">`)
		for _, col := range cols {
			val := ""
			if v, ok := row[col]; ok {
				val = fmt.Sprintf("%v", v)
			}
			sb.WriteString(fmt.Sprintf(`<td class="px-2 py-1 text-gray-700">%s</td>`, template.HTMLEscapeString(val)))
		}
		sb.WriteString(`</tr>`)
	}
	sb.WriteString(`</tbody></table>`)

	return template.HTML(sb.String())
}

// ==================== Editable Displayer ====================

// EditableDisplayer allows inline editing
type EditableDisplayer struct {
	fieldName string
	displayType string // text, select, textarea
	options   []Option
	saveURL   func(record any) string
}

// Editable creates a new editable displayer
func Editable(fieldName string) *EditableDisplayer {
	return &EditableDisplayer{
		fieldName:   fieldName,
		displayType: "text",
	}
}

// Type sets the input type
func (e *EditableDisplayer) Type(t string) *EditableDisplayer {
	e.displayType = t
	return e
}

// Options sets options for select type
func (e *EditableDisplayer) Options(options ...Option) *EditableDisplayer {
	e.options = options
	return e
}

// SaveURL sets the save URL function
func (e *EditableDisplayer) SaveURL(fn func(record any) string) *EditableDisplayer {
	e.saveURL = fn
	return e
}

// Display implements the Displayer interface
func (e *EditableDisplayer) Display(record any, value any) template.HTML {
	displayValue := fmt.Sprintf("%v", value)
	if displayValue == "<nil>" {
		displayValue = ""
	}

	var inputHTML string
	switch e.displayType {
	case "select":
		var opts strings.Builder
		opts.WriteString(`<select class="editable-input border rounded px-2 py-1 text-sm w-full">`)
		for _, opt := range e.options {
			selected := ""
			if opt.Value == displayValue {
				selected = " selected"
			}
			opts.WriteString(fmt.Sprintf(`<option value="%s"%s>%s</option>`, opt.Value, selected, opt.Label))
		}
		opts.WriteString(`</select>`)
		inputHTML = opts.String()
	case "textarea":
		inputHTML = fmt.Sprintf(`<textarea class="editable-input border rounded px-2 py-1 text-sm w-full" rows="3">%s</textarea>`,
			template.HTMLEscapeString(displayValue))
	default:
		inputHTML = fmt.Sprintf(`<input type="text" class="editable-input border rounded px-2 py-1 text-sm w-full" value="%s">`,
			template.HTMLEscapeString(displayValue))
	}

	return template.HTML(fmt.Sprintf(
		`<div class="editable-cell group cursor-pointer hover:bg-gray-50 p-1 rounded" data-field="%s">`+
			`<span class="editable-display">%s</span>`+
			`<span class="editable-edit hidden">%s</span>`+
			`</div>`,
		e.fieldName, template.HTMLEscapeString(displayValue), inputHTML))
}

// ==================== Button Displayer ====================

// ButtonDisplayer displays a button
type ButtonDisplayer struct {
	label     string
	style     ActionStyle
	urlFunc   func(record any, value any) string
	onClick   func(record any, value any) string
}

// Button creates a new button displayer
func Button(label string) *ButtonDisplayer {
	return &ButtonDisplayer{
		label: label,
		style: ActionDefault,
	}
}

// Style sets the button style
func (b *ButtonDisplayer) Style(style ActionStyle) *ButtonDisplayer {
	b.style = style
	return b
}

// URL sets the button URL
func (b *ButtonDisplayer) URL(fn func(record any, value any) string) *ButtonDisplayer {
	b.urlFunc = fn
	return b
}

// OnClick sets the onclick handler
func (b *ButtonDisplayer) OnClick(fn func(record any, value any) string) *ButtonDisplayer {
	b.onClick = fn
	return b
}

// Display implements the Displayer interface
func (b *ButtonDisplayer) Display(record any, value any) template.HTML {
	styleClasses := map[ActionStyle]string{
		ActionDefault: "bg-white border border-gray-300 text-gray-700 hover:bg-gray-50",
		ActionPrimary: "bg-blue-600 text-white hover:bg-blue-700",
		ActionGhost:   "bg-transparent text-gray-600 hover:bg-gray-100",
		ActionDanger:  "bg-red-600 text-white hover:bg-red-700",
	}

	class := styleClasses[b.style]
	if class == "" {
		class = styleClasses[ActionDefault]
	}

	href := "#"
	if b.urlFunc != nil {
		href = b.urlFunc(record, value)
	}

	onclick := ""
	if b.onClick != nil {
		onclick = fmt.Sprintf(` onclick="%s"`, b.onClick(record, value))
	}

	label := b.label
	if label == "" {
		label = fmt.Sprintf("%v", value)
	}

	return template.HTML(fmt.Sprintf(
		`<a href="%s" class="inline-flex items-center px-3 py-1.5 text-sm font-medium rounded %s%s">%s</a>`,
		href, class, onclick, label))
}

// ==================== Checkbox Displayer ====================

// CheckboxDisplayer displays boolean values as checkboxes
type CheckboxDisplayer struct {
	checkedFunc func(record any, value any) bool
	disabled    bool
}

// Checkbox creates a new checkbox displayer
func Checkbox() *CheckboxDisplayer {
	return &CheckboxDisplayer{}
}

// Checked sets a function to determine if checked
func (c *CheckboxDisplayer) Checked(fn func(record any, value any) bool) *CheckboxDisplayer {
	c.checkedFunc = fn
	return c
}

// Disabled sets the disabled state
func (c *CheckboxDisplayer) Disabled(disabled bool) *CheckboxDisplayer {
	c.disabled = disabled
	return c
}

// Display implements the Displayer interface
func (c *CheckboxDisplayer) Display(record any, value any) template.HTML {
	isChecked := false
	if c.checkedFunc != nil {
		isChecked = c.checkedFunc(record, value)
	} else {
		// Default check logic
		switch v := value.(type) {
		case bool:
			isChecked = v
		case int, int64:
			isChecked = v != 0
		case string:
			isChecked = v == "1" || v == "true" || v == "yes" || v == "on"
		}
	}

	checkedAttr := ""
	if isChecked {
		checkedAttr = " checked"
	}

	disabledAttr := ""
	if c.disabled {
		disabledAttr = " disabled"
	}

	return template.HTML(fmt.Sprintf(
		`<input type="checkbox" class="w-4 h-4 text-blue-600 rounded border-gray-300 focus:ring-blue-500"%s%s onclick="return false;">`,
		checkedAttr, disabledAttr))
}

// ==================== Radio Displayer ====================

// RadioDisplayer displays values as radio buttons
type RadioDisplayer struct {
	options     []Option
	selectedFunc func(record any, value any) string
	disabled    bool
	name        string
}

// Radio creates a new radio displayer
func Radio() *RadioDisplayer {
	return &RadioDisplayer{}
}

// Options sets the radio options
func (r *RadioDisplayer) Options(options ...Option) *RadioDisplayer {
	r.options = options
	return r
}

// OptionsMap sets options from a map
func (r *RadioDisplayer) OptionsMap(m map[string]string) *RadioDisplayer {
	for k, v := range m {
		r.options = append(r.options, Option{Value: k, Label: v})
	}
	return r
}

// Selected sets a function to determine selected value
func (r *RadioDisplayer) Selected(fn func(record any, value any) string) *RadioDisplayer {
	r.selectedFunc = fn
	return r
}

// Disabled sets the disabled state
func (r *RadioDisplayer) Disabled(disabled bool) *RadioDisplayer {
	r.disabled = disabled
	return r
}

// Name sets the radio group name
func (r *RadioDisplayer) Name(name string) *RadioDisplayer {
	r.name = name
	return r
}

// Display implements the Displayer interface
func (r *RadioDisplayer) Display(record any, value any) template.HTML {
	if len(r.options) == 0 {
		return template.HTML(`<span class="text-gray-400">-</span>`)
	}

	selectedValue := ""
	if r.selectedFunc != nil {
		selectedValue = r.selectedFunc(record, value)
	} else {
		selectedValue = fmt.Sprintf("%v", value)
	}

	disabledAttr := ""
	if r.disabled {
		disabledAttr = " disabled"
	}

	name := r.name
	if name == "" {
		name = "radio-group"
	}

	var sb strings.Builder
	sb.WriteString(`<div class="flex flex-wrap gap-3">`)
	for _, opt := range r.options {
		checkedAttr := ""
		if opt.Value == selectedValue {
			checkedAttr = " checked"
		}
		sb.WriteString(fmt.Sprintf(
			`<label class="inline-flex items-center cursor-pointer">`+
				`<input type="radio" name="%s" value="%s" class="w-4 h-4 text-blue-600 border-gray-300 focus:ring-blue-500"%s%s onclick="return false;">`+
				`<span class="ml-2 text-sm text-gray-700">%s</span>`+
				`</label>`,
			name, opt.Value, checkedAttr, disabledAttr, opt.Label))
	}
	sb.WriteString(`</div>`)

	return template.HTML(sb.String())
}

// ==================== Select Displayer ====================

// SelectDisplayer displays values as styled select options
type SelectDisplayer struct {
	options    []Option
	colorMap   map[string]string
	defaultColor string
}

// Select creates a new select displayer
func SelectDisplay() *SelectDisplayer {
	return &SelectDisplayer{
		colorMap: make(map[string]string),
		defaultColor: "gray",
	}
}

// Options sets the select options
func (s *SelectDisplayer) Options(options ...Option) *SelectDisplayer {
	s.options = options
	return s
}

// OptionsMap sets options from a map
func (s *SelectDisplayer) OptionsMap(m map[string]string) *SelectDisplayer {
	for k, v := range m {
		s.options = append(s.options, Option{Value: k, Label: v})
	}
	return s
}

// Color sets the color for a specific value
func (s *SelectDisplayer) Color(value, color string) *SelectDisplayer {
	s.colorMap[value] = color
	return s
}

// DefaultColor sets the default color
func (s *SelectDisplayer) DefaultColor(color string) *SelectDisplayer {
	s.defaultColor = color
	return s
}

// Display implements the Displayer interface
func (s *SelectDisplayer) Display(record any, value any) template.HTML {
	valueStr := fmt.Sprintf("%v", value)
	if valueStr == "<nil>" {
		return template.HTML(`<span class="text-gray-400">-</span>`)
	}

	// Find label for value
	label := valueStr
	for _, opt := range s.options {
		if opt.Value == valueStr {
			label = opt.Label
			break
		}
	}

	color := s.defaultColor
	if c, ok := s.colorMap[valueStr]; ok {
		color = c
	}

	return template.HTML(fmt.Sprintf(
		`<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-%s-100 text-%s-800">%s</span>`,
		color, color, template.HTMLEscapeString(label)))
}

// ==================== SwitchGroup Displayer ====================

// SwitchGroupDisplayer displays multiple switches for array values
type SwitchGroupDisplayer struct {
	options     []Option
	separator   string
}

// SwitchGroup creates a new switch group displayer
func SwitchGroup() *SwitchGroupDisplayer {
	return &SwitchGroupDisplayer{
		separator: ",",
	}
}

// Options sets the switch options
func (s *SwitchGroupDisplayer) Options(options ...Option) *SwitchGroupDisplayer {
	s.options = options
	return s
}

// Separator sets the value separator
func (s *SwitchGroupDisplayer) Separator(sep string) *SwitchGroupDisplayer {
	s.separator = sep
	return s
}

// Display implements the Displayer interface
func (s *SwitchGroupDisplayer) Display(record any, value any) template.HTML {
	// Parse value as array
	var values []string
	switch v := value.(type) {
	case []string:
		values = v
	case []any:
		for _, item := range v {
			values = append(values, fmt.Sprintf("%v", item))
		}
	case string:
		if v != "" {
			values = strings.Split(v, s.separator)
		}
	default:
		values = []string{fmt.Sprintf("%v", value)}
	}

	if len(s.options) == 0 {
		return template.HTML(`<span class="text-gray-400">-</span>`)
	}

	// Build value set for quick lookup
	valueSet := make(map[string]bool)
	for _, v := range values {
		valueSet[strings.TrimSpace(v)] = true
	}

	var sb strings.Builder
	sb.WriteString(`<div class="flex flex-wrap gap-2">`)
	for _, opt := range s.options {
		isOn := valueSet[opt.Value]
		bgClass := "bg-gray-200"
		dotClass := "translate-x-0.5"
		if isOn {
			bgClass = "bg-green-500"
			dotClass = "translate-x-5"
		}
		sb.WriteString(fmt.Sprintf(
			`<div class="flex items-center space-x-2">`+
				`<div class="w-10 h-5 %s rounded-full relative transition-colors">`+
				`<div class="w-4 h-4 bg-white rounded-full absolute top-0.5 left-0.5 %s transition-transform"></div>`+
				`</div>`+
				`<span class="text-sm text-gray-700">%s</span>`+
				`</div>`,
			bgClass, dotClass, opt.Label))
	}
	sb.WriteString(`</div>`)

	return template.HTML(sb.String())
}

// ==================== Input Displayer ====================

// InputDisplayer displays values as styled input fields
type InputDisplayer struct {
	placeholder string
	readonly    bool
	maxLength   int
	type_       string // text, number, email, etc.
}

// Input creates a new input displayer
func Input() *InputDisplayer {
	return &InputDisplayer{
		type_: "text",
	}
}

// Placeholder sets the placeholder
func (i *InputDisplayer) Placeholder(placeholder string) *InputDisplayer {
	i.placeholder = placeholder
	return i
}

// Readonly sets readonly state
func (i *InputDisplayer) Readonly(readonly bool) *InputDisplayer {
	i.readonly = readonly
	return i
}

// MaxLength sets max length
func (i *InputDisplayer) MaxLength(length int) *InputDisplayer {
	i.maxLength = length
	return i
}

// Type sets the input type
func (i *InputDisplayer) Type(t string) *InputDisplayer {
	i.type_ = t
	return i
}

// Display implements the Displayer interface
func (i *InputDisplayer) Display(record any, value any) template.HTML {
	valueStr := fmt.Sprintf("%v", value)
	if valueStr == "<nil>" {
		valueStr = ""
	}

	readonlyAttr := ""
	if i.readonly {
		readonlyAttr = " readonly"
	}

	maxLengthAttr := ""
	if i.maxLength > 0 {
		maxLengthAttr = fmt.Sprintf(` maxlength="%d"`, i.maxLength)
	}

	placeholderAttr := ""
	if i.placeholder != "" {
		placeholderAttr = fmt.Sprintf(` placeholder="%s"`, i.placeholder)
	}

	return template.HTML(fmt.Sprintf(
		`<input type="%s" value="%s" class="w-full px-3 py-1.5 text-sm border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"%s%s%s onclick="return false;">`,
		i.type_, template.HTMLEscapeString(valueStr), readonlyAttr, maxLengthAttr, placeholderAttr))
}

// ==================== Textarea Displayer ====================

// TextareaDisplayer displays values as styled textareas
type TextareaDisplayer struct {
	rows     int
	cols     int
	readonly bool
	maxLength int
}

// Textarea creates a new textarea displayer
func Textarea() *TextareaDisplayer {
	return &TextareaDisplayer{
		rows: 3,
		cols: 30,
	}
}

// Rows sets the number of rows
func (t *TextareaDisplayer) Rows(rows int) *TextareaDisplayer {
	t.rows = rows
	return t
}

// Cols sets the number of columns
func (t *TextareaDisplayer) Cols(cols int) *TextareaDisplayer {
	t.cols = cols
	return t
}

// Readonly sets readonly state
func (t *TextareaDisplayer) Readonly(readonly bool) *TextareaDisplayer {
	t.readonly = readonly
	return t
}

// MaxLength sets max length
func (t *TextareaDisplayer) MaxLength(length int) *TextareaDisplayer {
	t.maxLength = length
	return t
}

// Display implements the Displayer interface
func (t *TextareaDisplayer) Display(record any, value any) template.HTML {
	valueStr := fmt.Sprintf("%v", value)
	if valueStr == "<nil>" {
		valueStr = ""
	}

	readonlyAttr := ""
	if t.readonly {
		readonlyAttr = " readonly"
	}

	maxLengthAttr := ""
	if t.maxLength > 0 {
		maxLengthAttr = fmt.Sprintf(` maxlength="%d"`, t.maxLength)
	}

	return template.HTML(fmt.Sprintf(
		`<textarea rows="%d" cols="%d" class="w-full px-3 py-1.5 text-sm border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500 resize-none"%s%s onclick="return false;">%s</textarea>`,
		t.rows, t.cols, readonlyAttr, maxLengthAttr, template.HTMLEscapeString(valueStr)))
}

// ==================== Expand Displayer ====================

// ExpandDisplayer displays expandable content
type ExpandDisplayer struct {
	summary    string
	contentFunc func(record any, value any) template.HTML
	expanded   bool
}

// Expand creates a new expand displayer
func Expand() *ExpandDisplayer {
	return &ExpandDisplayer{}
}

// Summary sets the summary text
func (e *ExpandDisplayer) Summary(summary string) *ExpandDisplayer {
	e.summary = summary
	return e
}

// Content sets the content function
func (e *ExpandDisplayer) Content(fn func(record any, value any) template.HTML) *ExpandDisplayer {
	e.contentFunc = fn
	return e
}

// Expanded sets initial expanded state
func (e *ExpandDisplayer) Expanded(expanded bool) *ExpandDisplayer {
	e.expanded = expanded
	return e
}

// Display implements the Displayer interface
func (e *ExpandDisplayer) Display(record any, value any) template.HTML {
	summary := e.summary
	if summary == "" {
		summary = fmt.Sprintf("%v", value)
		if len(summary) > 50 {
			summary = summary[:50] + "..."
		}
	}

	openAttr := ""
	if e.expanded {
		openAttr = " open"
	}

	content := template.HTML("")
	if e.contentFunc != nil {
		content = e.contentFunc(record, value)
	} else {
		content = template.HTML(fmt.Sprintf("<div class='p-3 bg-gray-50 rounded'>%v</div>", value))
	}

	return template.HTML(fmt.Sprintf(
		`<details class="group"%s>`+
			`<summary class="cursor-pointer text-blue-600 hover:text-blue-800 font-medium">%s</summary>`+
			`<div class="mt-2 text-sm text-gray-600">%s</div>`+
			`</details>`,
		openAttr, template.HTMLEscapeString(summary), content))
}

// ==================== Modal Displayer ====================

// ModalDisplayer displays a modal trigger button
type ModalDisplayer struct {
	triggerLabel string
	triggerStyle ActionStyle
	title        string
	contentFunc  func(record any, value any) template.HTML
	size         string // sm, md, lg, xl
}

// Modal creates a new modal displayer
func Modal(label string) *ModalDisplayer {
	return &ModalDisplayer{
		triggerLabel: label,
		triggerStyle: ActionDefault,
		size:         "md",
	}
}

// Title sets the modal title
func (m *ModalDisplayer) Title(title string) *ModalDisplayer {
	m.title = title
	return m
}

// Content sets the modal content function
func (m *ModalDisplayer) Content(fn func(record any, value any) template.HTML) *ModalDisplayer {
	m.contentFunc = fn
	return m
}

// Size sets the modal size
func (m *ModalDisplayer) Size(size string) *ModalDisplayer {
	m.size = size
	return m
}

// TriggerStyle sets the trigger button style
func (m *ModalDisplayer) TriggerStyle(style ActionStyle) *ModalDisplayer {
	m.triggerStyle = style
	return m
}

// Display implements the Displayer interface
func (m *ModalDisplayer) Display(record any, value any) template.HTML {
	styleClasses := map[ActionStyle]string{
		ActionDefault: "bg-white border border-gray-300 text-gray-700 hover:bg-gray-50",
		ActionPrimary: "bg-blue-600 text-white hover:bg-blue-700",
		ActionGhost:   "bg-transparent text-gray-600 hover:bg-gray-100",
		ActionDanger:  "bg-red-600 text-white hover:bg-red-700",
	}

	class := styleClasses[m.triggerStyle]
	if class == "" {
		class = styleClasses[ActionDefault]
	}

	sizeClasses := map[string]string{
		"sm": "max-w-sm",
		"md": "max-w-md",
		"lg": "max-w-lg",
		"xl": "max-w-xl",
	}
	sizeClass := sizeClasses[m.size]
	if sizeClass == "" {
		sizeClass = sizeClasses["md"]
	}

	title := m.title
	if title == "" {
		title = "Details"
	}

	content := template.HTML("")
	if m.contentFunc != nil {
		content = m.contentFunc(record, value)
	} else {
		content = template.HTML(fmt.Sprintf("<p>%v</p>", value))
	}

	return template.HTML(fmt.Sprintf(
		`<span x-data="{ open: false }">`+
			`<button @click="open = true" class="inline-flex items-center px-3 py-1.5 text-sm font-medium rounded %s">%s</button>`+
			`<div x-show="open" class="fixed inset-0 z-50 overflow-y-auto" style="display: none;">`+
			`<div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">`+
			`<div x-show="open" @click="open = false" class="fixed inset-0 transition-opacity bg-gray-500 bg-opacity-75"></div>`+
			`<span class="hidden sm:inline-block sm:align-middle sm:h-screen">&#8203;</span>`+
			`<div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle %s">`+
			`<div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">`+
			`<div class="sm:flex sm:items-start">`+
			`<div class="mt-3 text-center sm:mt-0 sm:text-left w-full">`+
			`<h3 class="text-lg leading-6 font-medium text-gray-900 mb-4">%s</h3>`+
			`<div class="mt-2">%s</div>`+
			`</div>`+
			`</div>`+
			`</div>`+
			`<div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">`+
			`<button @click="open = false" type="button" class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-blue-600 text-base font-medium text-white hover:bg-blue-700 focus:outline-none sm:ml-3 sm:w-auto sm:text-sm">Close</button>`+
			`</div>`+
			`</div>`+
			`</div>`+
			`</div>`+
			`</span>`,
		class, m.triggerLabel, sizeClass, title, content))
}

// ==================== Downloadable Displayer ====================

// DownloadableDisplayer displays values as downloadable links
type DownloadableDisplayer struct {
	urlFunc    func(record any, value any) string
	filename   func(record any, value any) string
	text       string
	icon       bool
}

// Downloadable creates a new downloadable displayer
func Downloadable() *DownloadableDisplayer {
	return &DownloadableDisplayer{
		icon: true,
	}
}

// URL sets the download URL function
func (d *DownloadableDisplayer) URL(fn func(record any, value any) string) *DownloadableDisplayer {
	d.urlFunc = fn
	return d
}

// Filename sets the filename function
func (d *DownloadableDisplayer) Filename(fn func(record any, value any) string) *DownloadableDisplayer {
	d.filename = fn
	return d
}

// Text sets the link text
func (d *DownloadableDisplayer) Text(text string) *DownloadableDisplayer {
	d.text = text
	return d
}

// ShowIcon shows/hides the download icon
func (d *DownloadableDisplayer) ShowIcon(show bool) *DownloadableDisplayer {
	d.icon = show
	return d
}

// Display implements the Displayer interface
func (d *DownloadableDisplayer) Display(record any, value any) template.HTML {
	valueStr := fmt.Sprintf("%v", value)
	if valueStr == "<nil>" || valueStr == "" {
		return template.HTML(`<span class="text-gray-400">-</span>`)
	}

	url := valueStr
	if d.urlFunc != nil {
		url = d.urlFunc(record, value)
	}

	filename := ""
	if d.filename != nil {
		filename = d.filename(record, value)
	} else {
		filename = fmt.Sprintf("download-%v", value)
	}

	text := d.text
	if text == "" {
		text = valueStr
		if len(text) > 30 {
			text = text[:30] + "..."
		}
	}

	iconHTML := ""
	if d.icon {
		iconHTML = `<svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"></path></svg>`
	}

	return template.HTML(fmt.Sprintf(
		`<a href="%s" download="%s" class="inline-flex items-center text-blue-600 hover:text-blue-800 hover:underline">`+
			`%s%s`+
			`</a>`,
		url, filename, iconHTML, template.HTMLEscapeString(text)))
}

// ==================== Orderable Displayer ====================

// OrderableDisplayer displays sortable/orderable indicators
type OrderableDisplayer struct {
	orderFunc  func(record any, value any) int
	showArrows bool
}

// Orderable creates a new orderable displayer
func Orderable() *OrderableDisplayer {
	return &OrderableDisplayer{
		showArrows: true,
	}
}

// Order sets the order function
func (o *OrderableDisplayer) Order(fn func(record any, value any) int) *OrderableDisplayer {
	o.orderFunc = fn
	return o
}

// ShowArrows shows/hides the order arrows
func (o *OrderableDisplayer) ShowArrows(show bool) *OrderableDisplayer {
	o.showArrows = show
	return o
}

// Display implements the Displayer interface
func (o *OrderableDisplayer) Display(record any, value any) template.HTML {
	order := 0
	if o.orderFunc != nil {
		order = o.orderFunc(record, value)
	} else {
		switch v := value.(type) {
		case int:
			order = v
		case int64:
			order = int(v)
		case float64:
			order = int(v)
		default:
			order = 0
		}
	}

	arrowHTML := ""
	if o.showArrows {
		arrowHTML = `<div class="flex flex-col ml-2 text-gray-400">` +
			`<svg class="w-3 h-3 hover:text-gray-600 cursor-pointer" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M14.707 12.707a1 1 0 01-1.414 0L10 9.414l-3.293 3.293a1 1 0 01-1.414-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 010 1.414z" clip-rule="evenodd"></path></svg>` +
			`<svg class="w-3 h-3 hover:text-gray-600 cursor-pointer" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd"></path></svg>` +
			`</div>`
	}

	return template.HTML(fmt.Sprintf(
		`<div class="flex items-center justify-center">`+
			`<span class="inline-flex items-center justify-center w-8 h-8 bg-blue-100 text-blue-800 text-sm font-medium rounded-full">%d</span>`+
			`%s`+
			`</div>`,
		order, arrowHTML))
}

// ==================== Tree Displayer ====================

// TreeDisplayer displays hierarchical tree data
type TreeDisplayer struct {
	labelFunc   func(item map[string]any) string
	childrenKey string
	maxDepth    int
	expanded    bool
}

// Tree creates a new tree displayer
func Tree() *TreeDisplayer {
	return &TreeDisplayer{
		childrenKey: "children",
		maxDepth:    3,
	}
}

// Label sets the label function
func (t *TreeDisplayer) Label(fn func(item map[string]any) string) *TreeDisplayer {
	t.labelFunc = fn
	return t
}

// ChildrenKey sets the children key
func (t *TreeDisplayer) ChildrenKey(key string) *TreeDisplayer {
	t.childrenKey = key
	return t
}

// MaxDepth sets the maximum depth to display
func (t *TreeDisplayer) MaxDepth(depth int) *TreeDisplayer {
	t.maxDepth = depth
	return t
}

// Expanded sets initial expanded state
func (t *TreeDisplayer) Expanded(expanded bool) *TreeDisplayer {
	t.expanded = expanded
	return t
}

// Display implements the Displayer interface
func (t *TreeDisplayer) Display(record any, value any) template.HTML {
	var data []map[string]any

	switch v := value.(type) {
	case []map[string]any:
		data = v
	case string:
		if v != "" {
			json.Unmarshal([]byte(v), &data)
		}
	case []any:
		for _, item := range v {
			if m, ok := item.(map[string]any); ok {
				data = append(data, m)
			}
		}
	}

	if len(data) == 0 {
		return template.HTML(`<span class="text-gray-400">-</span>`)
	}

	treeHTML := t.buildTree(data, 0)
	return template.HTML(treeHTML)
}

func (t *TreeDisplayer) buildTree(items []map[string]any, depth int) string {
	if depth >= t.maxDepth {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(`<ul class="space-y-1">`)
	for _, item := range items {
		label := ""
		if t.labelFunc != nil {
			label = t.labelFunc(item)
		} else if name, ok := item["name"].(string); ok {
			label = name
		} else if title, ok := item["title"].(string); ok {
			label = title
		} else {
			label = fmt.Sprintf("%v", item)
		}

		sb.WriteString(`<li class="flex items-center">`)
		sb.WriteString(fmt.Sprintf(
			`<span class="inline-flex items-center text-sm text-gray-700">`+
				`<svg class="w-4 h-4 mr-1.5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"></path></svg>`+
				`%s</span>`, template.HTMLEscapeString(label)))

		// Process children
		if children, ok := item[t.childrenKey].([]any); ok && len(children) > 0 {
			childMaps := make([]map[string]any, 0, len(children))
			for _, child := range children {
				if m, ok := child.(map[string]any); ok {
					childMaps = append(childMaps, m)
				}
			}
			if len(childMaps) > 0 {
				sb.WriteString(`<ul class="ml-6 mt-1 space-y-1">`)
				sb.WriteString(t.buildTree(childMaps, depth+1))
				sb.WriteString(`</ul>`)
			}
		}
		sb.WriteString(`</li>`)
	}
	sb.WriteString(`</ul>`)

	return sb.String()
}

// ==================== DialogTree Displayer ====================

// DialogTreeDisplayer displays a tree in a dialog
type DialogTreeDisplayer struct {
	triggerLabel string
	title        string
	treeDataFunc func(record any, value any) []map[string]any
	size         string
}

// DialogTree creates a new dialog tree displayer
func DialogTree(label string) *DialogTreeDisplayer {
	return &DialogTreeDisplayer{
		triggerLabel: label,
		size:         "lg",
	}
}

// Title sets the dialog title
func (d *DialogTreeDisplayer) Title(title string) *DialogTreeDisplayer {
	d.title = title
	return d
}

// TreeData sets the tree data function
func (d *DialogTreeDisplayer) TreeData(fn func(record any, value any) []map[string]any) *DialogTreeDisplayer {
	d.treeDataFunc = fn
	return d
}

// Size sets the dialog size
func (d *DialogTreeDisplayer) Size(size string) *DialogTreeDisplayer {
	d.size = size
	return d
}

// Display implements the Displayer interface
func (d *DialogTreeDisplayer) Display(record any, value any) template.HTML {
	title := d.title
	if title == "" {
		title = "Tree View"
	}

	sizeClasses := map[string]string{
		"sm":  "max-w-sm",
		"md":  "max-w-md",
		"lg":  "max-w-lg",
		"xl":  "max-w-xl",
		"2xl": "max-w-2xl",
		"full": "max-w-full mx-4",
	}
	sizeClass := sizeClasses[d.size]
	if sizeClass == "" {
		sizeClass = sizeClasses["lg"]
	}

	var treeData []map[string]any
	if d.treeDataFunc != nil {
		treeData = d.treeDataFunc(record, value)
	} else {
		// Try to parse value as tree data
		switch v := value.(type) {
		case []map[string]any:
			treeData = v
		case string:
			if v != "" {
				json.Unmarshal([]byte(v), &treeData)
			}
		}
	}

	tree := &TreeDisplayer{childrenKey: "children", maxDepth: 10}
	treeHTML := tree.buildTree(treeData, 0)

	return template.HTML(fmt.Sprintf(
		`<span x-data="{ open: false }">`+
			`<button @click="open = true" class="inline-flex items-center px-3 py-1.5 text-sm font-medium bg-white border border-gray-300 rounded hover:bg-gray-50 text-gray-700">`+
			`<svg class="w-4 h-4 mr-1.5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"></path></svg>`+
			`%s</button>`+
			`<div x-show="open" class="fixed inset-0 z-50 overflow-y-auto" style="display: none;">`+
			`<div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">`+
			`<div x-show="open" @click="open = false" class="fixed inset-0 transition-opacity bg-gray-500 bg-opacity-75"></div>`+
			`<span class="hidden sm:inline-block sm:align-middle sm:h-screen">&#8203;</span>`+
			`<div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle %s">`+
			`<div class="bg-white px-4 pt-5 pb-4 sm:p-6">`+
			`<div class="flex justify-between items-center mb-4">`+
			`<h3 class="text-lg font-medium text-gray-900">%s</h3>`+
			`<button @click="open = false" class="text-gray-400 hover:text-gray-500">`+
			`<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>`+
			`</button>`+
			`</div>`+
			`<div class="max-h-96 overflow-y-auto">%s</div>`+
			`</div>`+
			`</div>`+
			`</div>`+
			`</div>`+
			`</span>`,
		d.triggerLabel, sizeClass, title, treeHTML))
}

// DropdownAction represents a single action in the dropdown
type DropdownAction struct {
	Label   string
	URL     func(record any, value any) string
	Style   ActionStyle
	Confirm string
}

// DropdownActionsDisplayer displays a dropdown menu of actions
type DropdownActionsDisplayer struct {
	label   string
	actions []DropdownAction
}

// DropdownActions creates a new dropdown actions displayer
func DropdownActions(label string) *DropdownActionsDisplayer {
	return &DropdownActionsDisplayer{
		label:   label,
		actions: []DropdownAction{},
	}
}

// Action adds an action to the dropdown
func (d *DropdownActionsDisplayer) Action(label string, url func(record any, value any) string) *DropdownActionsDisplayer {
	d.actions = append(d.actions, DropdownAction{
		Label: label,
		URL:   url,
		Style: ActionDefault,
	})
	return d
}

// ActionWithConfirm adds an action with confirmation
func (d *DropdownActionsDisplayer) ActionWithConfirm(label string, url func(record any, value any) string, confirm string) *DropdownActionsDisplayer {
	d.actions = append(d.actions, DropdownAction{
		Label:   label,
		URL:     url,
		Style:   ActionDanger,
		Confirm: confirm,
	})
	return d
}

// Display implements the Displayer interface
func (d *DropdownActionsDisplayer) Display(record any, value any) template.HTML {
	if len(d.actions) == 0 {
		return template.HTML(`<span class="text-gray-400">-</span>`)
	}

	var items strings.Builder
	for _, action := range d.actions {
		url := "#"
		if action.URL != nil {
			url = action.URL(record, value)
		}

		confirm := ""
		if action.Confirm != "" {
			confirm = fmt.Sprintf(` onclick="return confirm('%s')"`, template.JSEscapeString(action.Confirm))
		}

		class := "block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
		if action.Style == ActionDanger {
			class = "block px-4 py-2 text-sm text-red-600 hover:bg-red-50"
		}

		items.WriteString(fmt.Sprintf(
			`<a href="%s" class="%s"%s>%s</a>`,
			url, class, confirm, action.Label))
	}

	return template.HTML(fmt.Sprintf(
		`<div class="relative inline-block text-left" x-data="{ open: false }">`+
			`<button @click="open = !open" class="inline-flex items-center px-3 py-1.5 text-sm font-medium bg-white border border-gray-300 rounded hover:bg-gray-50">`+
			`%s <svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>`+
			`</button>`+
			`<div x-show="open" @click.away="open = false" class="absolute right-0 z-10 mt-2 w-48 bg-white rounded-md shadow-lg ring-1 ring-black ring-opacity-5" style="display: none;">`+
			`<div class="py-1">%s</div>`+
			`</div>`+
			`</div>`,
		d.label, items.String()))
}
