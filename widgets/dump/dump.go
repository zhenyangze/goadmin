// Package dump provides data dumping/debugging widgets.
package dump

import (
	"encoding/json"
	"fmt"
	"html/template"
)

// Dump displays variable data for debugging
type Dump struct {
	data  interface{}
	label string
	depth int
}

// New creates a new dump widget
func New(data interface{}) *Dump {
	return &Dump{
		data:  data,
		depth: 5,
	}
}

// Label sets the dump label
func (d *Dump) Label(label string) *Dump {
	d.label = label
	return d
}

// Depth sets the recursion depth
func (d *Dump) Depth(depth int) *Dump {
	d.depth = depth
	return d
}

// RenderContext provides data for rendering
type RenderContext struct {
	Label    string
	Content  template.HTML
	JSONData string
}

// Render prepares the dump for rendering
func (d *Dump) Render() *RenderContext {
	jsonData, _ := json.MarshalIndent(d.data, "", "  ")

	return &RenderContext{
		Label:    d.label,
		Content:  formatDump(d.data, 0, d.depth),
		JSONData: string(jsonData),
	}
}

// formatDump formats the data as HTML
func formatDump(data interface{}, level, maxDepth int) template.HTML {
	if level >= maxDepth {
		return template.HTML(`<span class="dump-max-depth">... (max depth)</span>`)
	}

	switch v := data.(type) {
	case nil:
		return template.HTML(`<span class="dump-null">null</span>`)
	case bool:
		return template.HTML(fmt.Sprintf(`<span class="dump-bool">%t</span>`, v))
	case float64:
		return template.HTML(fmt.Sprintf(`<span class="dump-number">%v</span>`, v))
	case int:
		return template.HTML(fmt.Sprintf(`<span class="dump-number">%d</span>`, v))
	case int64:
		return template.HTML(fmt.Sprintf(`<span class="dump-number">%d</span>`, v))
	case string:
		return template.HTML(fmt.Sprintf(`<span class="dump-string">"%s"</span>`, template.HTMLEscapeString(v)))
	case []interface{}:
		return formatArray(v, level, maxDepth)
	case map[string]interface{}:
		return formatMap(v, level, maxDepth)
	default:
		return template.HTML(fmt.Sprintf(`<span class="dump-other">%v</span>`, v))
	}
}

func formatArray(arr []interface{}, level, maxDepth int) template.HTML {
	if len(arr) == 0 {
		return template.HTML(`<span class="dump-array">[]</span>`)
	}

	var html string
	html += fmt.Sprintf(`<span class="dump-array">array(%d) [</span>`, len(arr))
	html += `<div class="dump-indent">`

	for i, item := range arr {
		html += fmt.Sprintf(`<div class="dump-item"><span class="dump-key">[%d]</span> => %s</div>`,
			i, formatDump(item, level+1, maxDepth))
	}

	html += `</div>`
	html += `<span class="dump-array">]</span>`

	return template.HTML(html)
}

func formatMap(m map[string]interface{}, level, maxDepth int) template.HTML {
	if len(m) == 0 {
		return template.HTML(`<span class="dump-object">{}</span>`)
	}

	var html string
	html += fmt.Sprintf(`<span class="dump-object">object(%d) {</span>`, len(m))
	html += `<div class="dump-indent">`

	for key, val := range m {
		html += fmt.Sprintf(`<div class="dump-item"><span class="dump-key">["%s"]</span> => %s</div>`,
			template.HTMLEscapeString(key), formatDump(val, level+1, maxDepth))
	}

	html += `</div>`
	html += `<span class="dump-object">}</span>`

	return template.HTML(html)
}

// Var dumps a variable (alias for New)
func Var(data interface{}) *Dump {
	return New(data)
}
