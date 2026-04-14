// Package chart provides chart widgets for dashboard and data visualization.
package chart

import (
	"encoding/json"
	"fmt"
	"time"
)

// ChartType describes the type of chart.
type ChartType string

const (
	TypeLine   ChartType = "line"
	TypeBar    ChartType = "bar"
	TypePie    ChartType = "pie"
	TypeDoughnut ChartType = "doughnut"
	TypeRadar  ChartType = "radar"
	TypeArea   ChartType = "area"
)

// Dataset represents a data series in the chart.
type Dataset struct {
	Label       string   `json:"label"`
	Data        []float64 `json:"data"`
	BackgroundColor interface{} `json:"backgroundColor,omitempty"`
	BorderColor     interface{} `json:"borderColor,omitempty"`
	BorderWidth     int         `json:"borderWidth,omitempty"`
	Fill            bool        `json:"fill,omitempty"`
	Tension         float64     `json:"tension,omitempty"`
}

// Chart is a chart widget for data visualization.
type Chart struct {
	chartType    ChartType
	title        string
	subtitle     string
	labels       []string
	datasets     []*Dataset
	height       string
	width        string
	options      map[string]interface{}
	responsive   bool
	aspectRatio  *float64
	id           string
}

// New creates a new chart widget.
func New(chartType ChartType) *Chart {
	return &Chart{
		chartType:  chartType,
		height:     "300px",
		responsive: true,
		options:    make(map[string]interface{}),
	}
}

// Line creates a new line chart.
func Line() *Chart {
	return New(TypeLine)
}

// Bar creates a new bar chart.
func Bar() *Chart {
	return New(TypeBar)
}

// Pie creates a new pie chart.
func Pie() *Chart {
	return New(TypePie)
}

// Doughnut creates a new doughnut chart.
func Doughnut() *Chart {
	return New(TypeDoughnut)
}

// Radar creates a new radar chart.
func Radar() *Chart {
	return New(TypeRadar)
}

// Area creates a new area chart (line with fill).
func Area() *Chart {
	c := New(TypeLine)
	c.chartType = TypeArea
	return c
}

// Title sets the chart title.
func (c *Chart) Title(title string) *Chart {
	c.title = title
	return c
}

// Subtitle sets the chart subtitle.
func (c *Chart) Subtitle(subtitle string) *Chart {
	c.subtitle = subtitle
	return c
}

// Labels sets the X-axis labels or pie chart labels.
func (c *Chart) Labels(labels ...string) *Chart {
	c.labels = labels
	return c
}

// Dataset adds a data series to the chart.
func (c *Chart) Dataset(label string, data []float64) *Chart {
	ds := &Dataset{
		Label: label,
		Data:  data,
	}
	c.datasets = append(c.datasets, ds)
	return c
}

// Colors sets the colors for datasets (for pie/doughnut charts).
func (c *Chart) Colors(colors ...string) *Chart {
	// Apply colors to the last dataset
	if len(c.datasets) > 0 {
		ds := c.datasets[len(c.datasets)-1]
		ds.BackgroundColor = colors
		ds.BorderColor = colors
	}
	return c
}

// Height sets the chart height.
func (c *Chart) Height(height string) *Chart {
	c.height = height
	return c
}

// Width sets the chart width.
func (c *Chart) Width(width string) *Chart {
	c.width = width
	return c
}

// Responsive enables/disables responsive resizing.
func (c *Chart) Responsive(responsive bool) *Chart {
	c.responsive = responsive
	return c
}

// AspectRatio sets the chart aspect ratio.
func (c *Chart) AspectRatio(ratio float64) *Chart {
	c.aspectRatio = &ratio
	return c
}

// ID sets a custom ID for the chart.
func (c *Chart) ID(id string) *Chart {
	c.id = id
	return c
}

// Option sets a custom chart option.
func (c *Chart) Option(key string, value interface{}) *Chart {
	c.options[key] = value
	return c
}

// SmoothLine makes line charts smooth (sets tension).
func (c *Chart) SmoothLine() *Chart {
	for _, ds := range c.datasets {
		ds.Tension = 0.4
	}
	return c
}

// FillArea enables area fill for line charts.
func (c *Chart) FillArea() *Chart {
	for _, ds := range c.datasets {
		ds.Fill = true
	}
	return c
}

// Data returns the chart data as JSON.
func (c *Chart) Data() string {
	data := map[string]interface{}{
		"labels":   c.labels,
		"datasets": c.datasets,
	}
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

// Options returns the chart options as JSON.
func (c *Chart) Options() string {
	opts := map[string]interface{}{
		"responsive": c.responsive,
		"plugins": map[string]interface{}{
			"title": map[string]interface{}{
				"display": c.title != "",
				"text":    c.title,
			},
			"legend": map[string]interface{}{
				"display": true,
				"position": "top",
			},
		},
	}
	if c.aspectRatio != nil {
		opts["aspectRatio"] = *c.aspectRatio
	}
	for k, v := range c.options {
		opts[k] = v
	}
	jsonOpts, _ := json.Marshal(opts)
	return string(jsonOpts)
}

// ChartType returns the chart type string for Chart.js.
func (c *Chart) ChartType() string {
	if c.chartType == TypeArea {
		return "line" // Area is a line chart with fill
	}
	return string(c.chartType)
}

// HasFill returns true if the chart should fill the area.
func (c *Chart) HasFill() bool {
	return c.chartType == TypeArea
}

// RenderContext provides data for rendering the chart.
type RenderContext struct {
	ID       string
	Type     string
	Data     string
	Options  string
	Height   string
	Width    string
	Title    string
	Subtitle string
	HasFill  bool
}

// Render prepares the chart for rendering.
func (c *Chart) Render() *RenderContext {
	id := c.id
	if id == "" {
		// Generate a random ID
		id = "chart-" + generateID()
	}
	return &RenderContext{
		ID:       id,
		Type:     c.ChartType(),
		Data:     c.Data(),
		Options:  c.Options(),
		Height:   c.height,
		Width:    c.width,
		Title:    c.title,
		Subtitle: c.subtitle,
		HasFill:  c.HasFill(),
	}
}

// generateID creates a simple random ID.
func generateID() string {
	// Simple timestamp-based ID
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
