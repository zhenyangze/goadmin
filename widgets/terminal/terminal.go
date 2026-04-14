// Package terminal provides terminal/console output widget
package terminal

import (
	"fmt"
	"html/template"
	"strings"
)

// Line represents a terminal line
type Line struct {
	Text  string
	Type  string // info, success, warning, error, command
	Timestamp string
}

// Terminal is a widget that displays terminal/console output
type Terminal struct {
	id       string
	lines    []Line
	prompt   string
	theme    string // dark, light
	height   string
	autoScroll bool
	copyable   bool
}

// New creates a new terminal widget
func New() *Terminal {
	return &Terminal{
		prompt:     "$ ",
		theme:      "dark",
		height:     "300px",
		autoScroll: true,
		copyable:   true,
	}
}

// ID sets the terminal ID
func (t *Terminal) ID(id string) *Terminal {
	t.id = id
	return t
}

// Prompt sets the prompt character
func (t *Terminal) Prompt(prompt string) *Terminal {
	t.prompt = prompt
	return t
}

// Theme sets the theme (dark or light)
func (t *Terminal) Theme(theme string) *Terminal {
	t.theme = theme
	return t
}

// Height sets the height
func (t *Terminal) Height(height string) *Terminal {
	t.height = height
	return t
}

// AutoScroll enables/disables auto scroll
func (t *Terminal) AutoScroll(auto bool) *Terminal {
	t.autoScroll = auto
	return t
}

// Copyable enables/disables copy button
func (t *Terminal) Copyable(copyable bool) *Terminal {
	t.copyable = copyable
	return t
}

// AddLine adds a line
func (t *Terminal) AddLine(text string, lineType string) *Terminal {
	t.lines = append(t.lines, Line{Text: text, Type: lineType})
	return t
}

// AddInfo adds an info line
func (t *Terminal) AddInfo(text string) *Terminal {
	return t.AddLine(text, "info")
}

// AddSuccess adds a success line
func (t *Terminal) AddSuccess(text string) *Terminal {
	return t.AddLine(text, "success")
}

// AddWarning adds a warning line
func (t *Terminal) AddWarning(text string) *Terminal {
	return t.AddLine(text, "warning")
}

// AddError adds an error line
func (t *Terminal) AddError(text string) *Terminal {
	return t.AddLine(text, "error")
}

// AddCommand adds a command line
func (t *Terminal) AddCommand(text string) *Terminal {
	return t.AddLine(text, "command")
}

// Lines sets all lines
func (t *Terminal) Lines(lines []Line) *Terminal {
	t.lines = lines
	return t
}

// Clear clears all lines
func (t *Terminal) Clear() *Terminal {
	t.lines = []Line{}
	return t
}

// RenderContext provides data for rendering
type RenderContext struct {
	ID         string
	Lines      []Line
	Prompt     string
	Theme      string
	Height     string
	AutoScroll bool
	Copyable   bool
}

// Render generates rendering context
func (t *Terminal) Render() *RenderContext {
	return &RenderContext{
		ID:         t.id,
		Lines:      t.lines,
		Prompt:     t.prompt,
		Theme:      t.theme,
		Height:     t.height,
		AutoScroll: t.autoScroll,
		Copyable:   t.copyable,
	}
}

// RenderTo generates HTML
func (t *Terminal) RenderTo() template.HTML {
	ctx := t.Render()
	return ctx.ToHTML()
}

// ToHTML converts context to HTML
func (ctx *RenderContext) ToHTML() template.HTML {
	var sb strings.Builder

	// Theme classes
	bgClass := "bg-gray-900"
	textClass := "text-green-400"
	borderClass := "border-gray-800"

	if ctx.Theme == "light" {
		bgClass = "bg-gray-100"
		textClass = "text-gray-800"
		borderClass = "border-gray-300"
	}

	idAttr := ""
	if ctx.ID != "" {
		idAttr = fmt.Sprintf(` id="%s"`, ctx.ID)
	}

	// Container
	sb.WriteString(fmt.Sprintf(`<div%s class="%s %s rounded-lg border %s overflow-hidden font-mono text-sm">`,
		idAttr, bgClass, textClass, borderClass))

	// Header with copy button
	if ctx.Copyable {
		sb.WriteString(fmt.Sprintf(`<div class="flex justify-between items-center px-4 py-2 %s border-b %s">`, bgClass, borderClass))
		sb.WriteString(`<span class="text-xs text-gray-500">Terminal</span>`)

		// Copy button
		allText := ""
		for _, line := range ctx.Lines {
			allText += ctx.Prompt + line.Text + "\n"
		}

		sb.WriteString(fmt.Sprintf(`<button class="text-xs text-gray-500 hover:text-gray-300 flex items-center" onclick="navigator.clipboard.writeText('%s')">`,
			template.JSEscapeString(strings.TrimSpace(allText))))
		sb.WriteString(`<svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"></path></svg>`)
		sb.WriteString(`Copy</button>`)
		sb.WriteString(`</div>`)
	}

	// Terminal content
	contentClass := "p-4 overflow-y-auto"
	if ctx.AutoScroll {
		contentClass += " auto-scroll"
	}
	sb.WriteString(fmt.Sprintf(`<div class="%s" style="height:%s;">`, contentClass, ctx.Height))

	// Lines
	for _, line := range ctx.Lines {
		lineClass := ""
		icon := ""

		switch line.Type {
		case "success":
			lineClass = "text-green-400"
			icon = "✓ "
		case "warning":
			lineClass = "text-yellow-400"
			icon = "⚠ "
		case "error":
			lineClass = "text-red-400"
			icon = "✗ "
		case "command":
			lineClass = "text-blue-400"
			icon = ctx.Prompt
		default:
			lineClass = "text-gray-300"
		}

		if ctx.Theme == "light" {
			switch line.Type {
			case "success":
				lineClass = "text-green-600"
			case "warning":
				lineClass = "text-yellow-600"
			case "error":
				lineClass = "text-red-600"
			case "command":
				lineClass = "text-blue-600"
			default:
				lineClass = "text-gray-700"
			}
		}

		sb.WriteString(fmt.Sprintf(`<div class="%s whitespace-pre-wrap break-words">`, lineClass))
		sb.WriteString(template.HTMLEscapeString(icon + line.Text))
		sb.WriteString(`</div>`)
	}

	if len(ctx.Lines) == 0 {
		sb.WriteString(`<div class="text-gray-500 italic">No output...</div>`)
	}

	sb.WriteString(`</div>`)
	sb.WriteString(`</div>`)

	// Auto scroll script
	if ctx.AutoScroll && ctx.ID != "" {
		sb.WriteString(fmt.Sprintf(`<script>
(function() {
	const terminal = document.getElementById('%s');
	if (terminal) {
		const content = terminal.querySelector('.overflow-y-auto');
		if (content) content.scrollTop = content.scrollHeight;
	}
})();
</script>`, ctx.ID))
	}

	return template.HTML(sb.String())
}

// Static creates a terminal with static content
func Static(content string) *Terminal {
	t := New()
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			if strings.HasPrefix(line, "$") {
				t.AddCommand(strings.TrimPrefix(line, "$"))
			} else if strings.HasPrefix(line, "Error:") || strings.HasPrefix(line, "[ERROR]") {
				t.AddError(line)
			} else if strings.HasPrefix(line, "Success:") || strings.HasPrefix(line, "[OK]") {
				t.AddSuccess(line)
			} else {
				t.AddInfo(line)
			}
		}
	}
	return t
}
