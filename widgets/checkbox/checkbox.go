// Package checkbox provides checkbox widget for forms and displays.
package checkbox

import (
	"fmt"
	"html/template"
	"strings"
)

// Option represents a checkbox option
type Option struct {
	Value    string
	Label    string
	Checked  bool
	Disabled bool
}

// Checkbox is a widget that displays checkboxes
type Checkbox struct {
	name      string
	options   []Option
	inline    bool
	class     string
	disabled  bool
	required  bool
}

// New creates a new checkbox widget
func New(name string) *Checkbox {
	return &Checkbox{
		name:    name,
		options: []Option{},
	}
}

// Option adds a single option
func (c *Checkbox) Option(value, label string) *Checkbox {
	c.options = append(c.options, Option{Value: value, Label: label})
	return c
}

// Options adds multiple options
func (c *Checkbox) Options(opts ...Option) *Checkbox {
	c.options = append(c.options, opts...)
	return c
}

// OptionsMap adds options from a map
func (c *Checkbox) OptionsMap(m map[string]string) *Checkbox {
	for k, v := range m {
		c.options = append(c.options, Option{Value: k, Label: v})
	}
	return c
}

// Inline sets whether checkboxes are displayed inline
func (c *Checkbox) Inline(inline bool) *Checkbox {
	c.inline = inline
	return c
}

// Class sets CSS classes
func (c *Checkbox) Class(class string) *Checkbox {
	c.class = class
	return c
}

// Disabled sets disabled state
func (c *Checkbox) Disabled(disabled bool) *Checkbox {
	c.disabled = disabled
	return c
}

// Required sets required state
func (c *Checkbox) Required(required bool) *Checkbox {
	c.required = required
	return c
}

// SetChecked sets checked values
func (c *Checkbox) SetChecked(values []string) *Checkbox {
	valueSet := make(map[string]bool)
	for _, v := range values {
		valueSet[v] = true
	}
	for i := range c.options {
		c.options[i].Checked = valueSet[c.options[i].Value]
	}
	return c
}

// RenderContext provides data for rendering
type RenderContext struct {
	Name     string
	Options  []Option
	Inline   bool
	Class    string
	Disabled bool
	Required bool
}

// Render generates HTML rendering context
func (c *Checkbox) Render() *RenderContext {
	return &RenderContext{
		Name:     c.name,
		Options:  c.options,
		Inline:   c.inline,
		Class:    c.class,
		Disabled: c.disabled,
		Required: c.required,
	}
}

// RenderTo generates HTML template
func (c *Checkbox) RenderTo() template.HTML {
	ctx := c.Render()
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
		if opt.Checked {
			checkedAttr = " checked"
		}

		disabledAttr := ""
		if opt.Disabled || ctx.Disabled {
			disabledAttr = " disabled"
		}

		sb.WriteString(fmt.Sprintf(
			`<label class="inline-flex items-center cursor-pointer">`+
				`<input type="checkbox" name="%s[]" value="%s" class="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"%s%s>`+
				`<span class="ml-2 text-sm text-gray-700">%s</span>`+
				`</label>`,
			ctx.Name, opt.Value, checkedAttr, disabledAttr, opt.Label))
	}

	sb.WriteString(`</div>`)

	return template.HTML(sb.String())
}
