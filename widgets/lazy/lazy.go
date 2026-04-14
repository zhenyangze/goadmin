// Package lazy provides lazy loading widgets.
package lazy

import "html/template"

// Lazy loads content asynchronously
type Lazy struct {
	url         string
	placeholder template.HTML
	event       string // scroll, click, visible
	id          string
}

// New creates a new lazy loading widget
func New() *Lazy {
	return &Lazy{
		event: "visible",
	}
}

// URL sets the content URL
func (l *Lazy) URL(url string) *Lazy {
	l.url = url
	return l
}

// Placeholder sets the loading placeholder
func (l *Lazy) Placeholder(placeholder template.HTML) *Lazy {
	l.placeholder = placeholder
	return l
}

// Event sets the trigger event
func (l *Lazy) Event(event string) *Lazy {
	l.event = event
	return l
}

// ID sets a custom ID
func (l *Lazy) ID(id string) *Lazy {
	l.id = id
	return l
}

// RenderContext provides data for rendering
type RenderContext struct {
	ID          string
	URL         string
	Placeholder template.HTML
	Event       string
}

// Render prepares the lazy widget for rendering
func (l *Lazy) Render() *RenderContext {
	return &RenderContext{
		ID:          l.id,
		URL:         l.url,
		Placeholder: l.placeholder,
		Event:       l.event,
	}
}

// JavaScript returns the required JavaScript
func JavaScript() template.HTML {
	return template.HTML(`
<script>
document.addEventListener('DOMContentLoaded', function() {
	// Intersection Observer for lazy loading
	const observer = new IntersectionObserver((entries) => {
		entries.forEach(entry => {
			if (entry.isIntersecting) {
				const el = entry.target;
				const url = el.dataset.url;
				if (url) {
					fetch(url)
						.then(r => r.text())
						.then(html => {
							el.innerHTML = html;
							el.classList.add('lazy-loaded');
						});
				}
				observer.unobserve(el);
			}
		});
	});

	document.querySelectorAll('[data-lazy="true"]').forEach(el => {
		observer.observe(el);
	});
});
</script>`)
}
