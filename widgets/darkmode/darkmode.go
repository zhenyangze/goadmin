// Package darkmode provides dark mode switcher widget.
package darkmode

import "html/template"

// DarkModeSwitcher toggles between light and dark themes
type DarkModeSwitcher struct {
	label       string
	defaultDark bool
	auto        bool // auto detect system preference
}

// New creates a new dark mode switcher
func New() *DarkModeSwitcher {
	return &DarkModeSwitcher{
		label: "Dark Mode",
		auto:  true,
	}
}

// Label sets the switcher label
func (d *DarkModeSwitcher) Label(label string) *DarkModeSwitcher {
	d.label = label
	return d
}

// DefaultDark sets the default state
func (d *DarkModeSwitcher) DefaultDark(dark bool) *DarkModeSwitcher {
	d.defaultDark = dark
	return d
}

// AutoDetect enables system preference detection
func (d *DarkModeSwitcher) AutoDetect(auto bool) *DarkModeSwitcher {
	d.auto = auto
	return d
}

// RenderContext provides data for rendering
type RenderContext struct {
	Label       string
	DefaultDark bool
	Auto        bool
}

// Render prepares the switcher for rendering
func (d *DarkModeSwitcher) Render() *RenderContext {
	return &RenderContext{
		Label:       d.label,
		DefaultDark: d.defaultDark,
		Auto:        d.auto,
	}
}

// JavaScript returns the required JavaScript
func JavaScript() template.HTML {
	return template.HTML(`
<script>
(function() {
	// Check for saved preference or system preference
	const savedTheme = localStorage.getItem('theme');
	const systemDark = window.matchMedia('(prefers-color-scheme: dark)').matches;

	if (savedTheme === 'dark' || (!savedTheme && systemDark)) {
		document.documentElement.classList.add('dark');
	}

	window.toggleDarkMode = function() {
		const isDark = document.documentElement.classList.toggle('dark');
		localStorage.setItem('theme', isDark ? 'dark' : 'light');
	};
})();
</script>`)
}
