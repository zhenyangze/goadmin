// Package metrics provides metrics widgets for dashboard data visualization.
package metrics

import (
	"html/template"
)

// MetricCard displays a single metric value with optional trend/chart
type MetricCard struct {
	title       string
	subtitle    string
	value       string
	icon        string
	color       string
	trend       float64
	trendLabel  string
	chart       template.HTML
	subtext     string
	link        string
	linkText    string
	id          string
}

// Card creates a new metric card widget
func Card() *MetricCard {
	return &MetricCard{
		color: "blue",
		icon:  "chart",
	}
}

// Title sets the card title
func (m *MetricCard) Title(title string) *MetricCard {
	m.title = title
	return m
}

// Subtitle sets the card subtitle
func (m *MetricCard) Subtitle(subtitle string) *MetricCard {
	m.subtitle = subtitle
	return m
}

// Value sets the metric value
func (m *MetricCard) Value(value string) *MetricCard {
	m.value = value
	return m
}

// Icon sets the icon (chart, users, dollar, shopping-cart, etc.)
func (m *MetricCard) Icon(icon string) *MetricCard {
	m.icon = icon
	return m
}

// Color sets the accent color (blue, green, red, yellow, purple, pink, indigo)
func (m *MetricCard) Color(color string) *MetricCard {
	m.color = color
	return m
}

// Trend sets the trend percentage and label
func (m *MetricCard) Trend(percent float64, label string) *MetricCard {
	m.trend = percent
	m.trendLabel = label
	return m
}

// Subtext sets additional text below the value
func (m *MetricCard) Subtext(text string) *MetricCard {
	m.subtext = text
	return m
}

// Link sets a link for the card
func (m *MetricCard) Link(url, text string) *MetricCard {
	m.link = url
	m.linkText = text
	return m
}

// ID sets a custom ID
func (m *MetricCard) ID(id string) *MetricCard {
	m.id = id
	return m
}

// RenderContext provides data for rendering
type CardRenderContext struct {
	ID         string
	Title      string
	Subtitle   string
	Value      string
	Icon       string
	Color      string
	ColorClass string
	Trend      float64
	TrendUp    bool
	TrendLabel string
	ShowTrend  bool
	Subtext    string
	Link       string
	LinkText   string
}

// Render prepares the metric card for rendering
func (m *MetricCard) Render() *CardRenderContext {
	id := m.id
	if id == "" {
		id = "metric-card-" + generateID()
	}

	colorMap := map[string]string{
		"blue":    "bg-blue-500",
		"green":   "bg-green-500",
		"red":     "bg-red-500",
		"yellow":  "bg-yellow-500",
		"purple":  "bg-purple-500",
		"pink":    "bg-pink-500",
		"indigo":  "bg-indigo-500",
		"orange":  "bg-orange-500",
		"cyan":    "bg-cyan-500",
		"teal":    "bg-teal-500",
	}

	return &CardRenderContext{
		ID:         id,
		Title:      m.title,
		Subtitle:   m.subtitle,
		Value:      m.value,
		Icon:       m.icon,
		Color:      m.color,
		ColorClass: colorMap[m.color],
		Trend:      m.trend,
		TrendUp:    m.trend >= 0,
		TrendLabel: m.trendLabel,
		ShowTrend:  m.trend != 0 || m.trendLabel != "",
		Subtext:    m.subtext,
		Link:       m.link,
		LinkText:   m.linkText,
	}
}

// getIconSVG returns the SVG for the icon
func getIconSVG(icon string) template.HTML {
	icons := map[string]string{
		"chart": `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path></svg>`,
		"users": `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"></path></svg>`,
		"dollar": `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>`,
		"shopping-cart": `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z"></path></svg>`,
		"eye": `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path></svg>`,
		"heart": `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"></path></svg>`,
		"star": `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z"></path></svg>`,
		"trending-up": `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path></svg>`,
		"trending-down": `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 17h8m0 0V9m0 8l-8-8-4 4-6-6"></path></svg>`,
		"box": `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"></path></svg>`,
	}

	svg, ok := icons[icon]
	if !ok {
		return template.HTML(icons["chart"])
	}
	return template.HTML(svg)
}

// GetIcon returns the icon SVG
func (m *MetricCard) GetIcon() template.HTML {
	return getIconSVG(m.icon)
}
