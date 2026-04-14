// Package tooltip provides tooltip widgets.
package tooltip

import (
	"fmt"
	"html/template"
)

// Tooltip shows a tooltip on hover
type Tooltip struct {
	content   template.HTML
	text      string
	position  string // top, bottom, left, right
	trigger   string // hover, click, focus
	delay     int    // ms
}

// New creates a new tooltip
func New() *Tooltip {
	return &Tooltip{
		position: "top",
		trigger:  "hover",
		delay:    0,
	}
}

// Content sets the tooltip content
func (t *Tooltip) Content(content template.HTML) *Tooltip {
	t.content = content
	return t
}

// Text sets the tooltip text (simple version)
func (t *Tooltip) Text(text string) *Tooltip {
	t.text = text
	return t
}

// Position sets the tooltip position
func (t *Tooltip) Position(pos string) *Tooltip {
	t.position = pos
	return t
}

// Trigger sets the trigger event
func (t *Tooltip) Trigger(trigger string) *Tooltip {
	t.trigger = trigger
	return t
}

// Delay sets the show delay
func (t *Tooltip) Delay(delay int) *Tooltip {
	t.delay = delay
	return t
}

// Wrap wraps content with tooltip
func (t *Tooltip) Wrap(content template.HTML) template.HTML {
	tipContent := t.text
	if t.content != "" {
		tipContent = string(t.content)
	}

	return template.HTML(fmt.Sprintf(
		`<span class="tooltip" data-tooltip="%s" data-position="%s" data-trigger="%s">%s</span>`,
		template.HTMLEscapeString(tipContent),
		t.position,
		t.trigger,
		content))
}
