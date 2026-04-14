// Package radio provides radio button widget for forms and displays.
package radio

import (
	"fmt"
	"html/template"
	"strings"
)

// Option represents a radio option
type Option struct {
	Value    string
	Label    string
	Disabled bool
}

// Radio is a widget that displays radio buttons
type Radio struct {
	name      string
	options   []Option
	selected  string
	inline    bool
	class     string
	disabled  bool
	required  bool
}

// New creates a new radio widget
func New(name string) *Radio {
	return &Radio{
		name:    name,
		options: []Option{},
	}
}

// Option adds a single option
func (r *Radio) Option(value, label string) *Radio {
	r.options = append(r.options, Option{Value: value, Label: label})
	return r
}

// Options adds multiple options
func (r *Radio) Options(opts ...Option) *Radio {
	r.options = append(r.options, opts...)
	return r
}

// OptionsMap adds options from a map
func (r *Radio) OptionsMap(m map[string]string) *Radio {
	for k, v := range m {
		r.options = append(r.options, Option{Value: k, Label: v})
	}
	return r
}

// Selected sets the selected value
func (r *Radio) Selected(value string) *Radio {
	r.selected = value
	return r
}

// Inline sets whether radios are displayed inline
func (r *Radio) Inline(inline bool) *Radio {
	r.inline = inline
	return r
}

// Class sets CSS classes
func (r *Radio) Class(class string) *Radio {
	r.class = class
	return r
}

// Disabled sets disabled state
func (r *Radio) Disabled(disabled bool) *Radio {
	r.disabled = disabled
	return r
}

// Required sets required state
func (r *Radio) Required(required bool) *Radio {
	r.required = required
	return r
}

// RenderContext provides data for rendering
type RenderContext struct {
	Name     string
	Options  []Option
	Selected string
	Inline   bool
	Class    string
	Disabled bool
	Required bool
}

// Render generates HTML rendering context
func (r *Radio) Render() *RenderContext {
	return &RenderContext{
		Name:     r.name,
		Options:  r.options,
		Selected: r.selected,
		Inline:   r.inline,
		Class:    r.class,
		Disabled: r.disabled,
		Required: r.required,
	}
}

// RenderTo generates HTML template
func (r *Radio) RenderTo() template.HTML {
	ctx := r.Render()
	return ctx.ToHTML()
}

// ToHTML converts context to HTML
func (ctx *RenderContext) ToHTML() template.HTML {
	var sb strings.Builder

	containerClass := "space-y-2"
	if ctx.Inline {
		containerClass = "flex flex-wrap gap-4"
	}
	if ctx.Class != "" {
		containerClass += " " + ctx.Class
	}

	sb.WriteString(fmt.Sprintf(`<div class="%s">`, containerClass))

	for _, opt := range ctx.Options {
		checkedAttr := ""
		if opt.Value == ctx.Selected {
			checkedAttr = " checked"
		}

		disabledAttr := ""
		if opt.Disabled || ctx.Disabled {
			disabledAttr = " disabled"
		}

		sb.WriteString(fmt.Sprintf(
			`<label class="inline-flex items-center cursor-pointer">`+
				`<input type="radio" name="%s" value="%s" class="w-4 h-4 text-blue-600 border-gray-300 focus:ring-blue-500"%s%s>`+
				`<span class="ml-2 text-sm text-gray-700">%s</span>`+
				`</label>`,
			ctx.Name, opt.Value, checkedAttr, disabledAttr, opt.Label))
	}

	sb.WriteString(`</div>`)

	return template.HTML(sb.String())
}
