package metrics

import (
	"encoding/json"
	"fmt"
	"time"
)

// RadialBar represents a radial/progress bar chart
type RadialBar struct {
	title    string
	subtitle string
	series   []RadialBarSeries
	height   string
	id       string
	colors   []string
}

// RadialBarSeries represents a single series in the radial bar
type RadialBarSeries struct {
	Name  string
	Value float64
	Max   float64
}

// NewRadialBar creates a new radial bar widget
func NewRadialBar() *RadialBar {
	return &RadialBar{
		height: "350px",
		colors: []string{"#3B82F6", "#10B981", "#F59E0B", "#EF4444", "#8B5CF6"},
	}
}

// Title sets the chart title
func (r *RadialBar) Title(title string) *RadialBar {
	r.title = title
	return r
}

// Subtitle sets the chart subtitle
func (r *RadialBar) Subtitle(subtitle string) *RadialBar {
	r.subtitle = subtitle
	return r
}

// Height sets the chart height
func (r *RadialBar) Height(height string) *RadialBar {
	r.height = height
	return r
}

// ID sets a custom ID
func (r *RadialBar) ID(id string) *RadialBar {
	r.id = id
	return r
}

// Colors sets custom colors
func (r *RadialBar) Colors(colors ...string) *RadialBar {
	r.colors = colors
	return r
}

// Series adds a data series
func (r *RadialBar) Series(name string, value, max float64) *RadialBar {
	r.series = append(r.series, RadialBarSeries{
		Name:  name,
		Value: value,
		Max:   max,
	})
	return r
}

// RenderContext provides data for rendering
type RadialBarRenderContext struct {
	ID       string
	Title    string
	Subtitle string
	Height   string
	DataJSON string
	Options  string
}

// Render prepares the radial bar for rendering
func (r *RadialBar) Render() *RadialBarRenderContext {
	id := r.id
	if id == "" {
		id = "radial-bar-" + generateID()
	}

	// Prepare data for ApexCharts
	labels := make([]string, len(r.series))
	data := make([]float64, len(r.series))
	for i, s := range r.series {
		labels[i] = s.Name
		if s.Max > 0 {
			data[i] = (s.Value / s.Max) * 100
		} else {
			data[i] = s.Value
		}
	}

	seriesData := []map[string]interface{}{
		{"data": data},
	}

	dataJSON, _ := json.Marshal(seriesData)

	options := map[string]interface{}{
		"chart": map[string]interface{}{
			"type": "radialBar",
		},
		"plotOptions": map[string]interface{}{
			"radialBar": map[string]interface{}{
				"dataLabels": map[string]interface{}{
					"name": map[string]interface{}{
						"fontSize": "16px",
					},
					"value": map[string]interface{}{
						"fontSize": "14px",
					},
				},
			},
		},
		"labels": labels,
		"colors": r.colors,
	}

	optionsJSON, _ := json.Marshal(options)

	return &RadialBarRenderContext{
		ID:       id,
		Title:    r.title,
		Subtitle: r.subtitle,
		Height:   r.height,
		DataJSON: string(dataJSON),
		Options:  string(optionsJSON),
	}
}

// RoundChart represents a circular progress chart
type RoundChart struct {
	title     string
	subtitle  string
	value     float64
	max       float64
	label     string
	size      string
	thickness string
	color     string
	trackColor string
	id        string
}

// Round creates a new round/circular progress chart
func Round() *RoundChart {
	return &RoundChart{
		value:      0,
		max:        100,
		size:       "150",
		thickness:  "10",
		color:      "#3B82F6",
		trackColor: "#E5E7EB",
	}
}

// Title sets the chart title
func (r *RoundChart) Title(title string) *RoundChart {
	r.title = title
	return r
}

// Subtitle sets the chart subtitle
func (r *RoundChart) Subtitle(subtitle string) *RoundChart {
	r.subtitle = subtitle
	return r
}

// Value sets the current value
func (r *RoundChart) Value(value float64) *RoundChart {
	r.value = value
	return r
}

// Max sets the maximum value
func (r *RoundChart) Max(max float64) *RoundChart {
	r.max = max
	return r
}

// Label sets the center label
func (r *RoundChart) Label(label string) *RoundChart {
	r.label = label
	return r
}

// Size sets the chart size in pixels
func (r *RoundChart) Size(size int) *RoundChart {
	r.size = fmt.Sprintf("%d", size)
	return r
}

// Thickness sets the stroke thickness
func (r *RoundChart) Thickness(thickness int) *RoundChart {
	r.thickness = fmt.Sprintf("%d", thickness)
	return r
}

// Color sets the progress color
func (r *RoundChart) Color(color string) *RoundChart {
	r.color = color
	return r
}

// TrackColor sets the background track color
func (r *RoundChart) TrackColor(color string) *RoundChart {
	r.trackColor = color
	return r
}

// ID sets a custom ID
func (r *RoundChart) ID(id string) *RoundChart {
	r.id = id
	return r
}

// RenderContext provides data for rendering
type RoundRenderContext struct {
	ID         string
	Title      string
	Subtitle   string
	Value      float64
	Max        float64
	Percentage float64
	Label      string
	Size       string
	Thickness  string
	Color      string
	TrackColor string
}

// Render prepares the round chart for rendering
func (r *RoundChart) Render() *RoundRenderContext {
	id := r.id
	if id == "" {
		id = "round-chart-" + generateID()
	}

	percentage := 0.0
	if r.max > 0 {
		percentage = (r.value / r.max) * 100
	}

	label := r.label
	if label == "" {
		label = fmt.Sprintf("%.0f%%", percentage)
	}

	return &RoundRenderContext{
		ID:         id,
		Title:      r.title,
		Subtitle:   r.subtitle,
		Value:      r.value,
		Max:        r.max,
		Percentage: percentage,
		Label:      label,
		Size:       r.size,
		Thickness:  r.thickness,
		Color:      r.color,
		TrackColor: r.trackColor,
	}
}

// SingleRound represents a single circular progress indicator
type SingleRound struct {
	title      string
	subtitle   string
	value      float64
	max        float64
	label      string
	size       string
	color      string
	trackColor string
	showValue  bool
	id         string
}

// SingleRound creates a new single round indicator
func SingleRoundChart() *SingleRound {
	return &SingleRound{
		value:      0,
		max:        100,
		size:       "120",
		color:      "#3B82F6",
		trackColor: "#E5E7EB",
		showValue:  true,
	}
}

// Title sets the chart title
func (s *SingleRound) Title(title string) *SingleRound {
	s.title = title
	return s
}

// Subtitle sets the chart subtitle
func (s *SingleRound) Subtitle(subtitle string) *SingleRound {
	s.subtitle = subtitle
	return s
}

// Value sets the current value
func (s *SingleRound) Value(value float64) *SingleRound {
	s.value = value
	return s
}

// Max sets the maximum value
func (s *SingleRound) Max(max float64) *SingleRound {
	s.max = max
	return s
}

// Label sets the center label
func (s *SingleRound) Label(label string) *SingleRound {
	s.label = label
	return s
}

// Size sets the chart size
func (s *SingleRound) Size(size int) *SingleRound {
	s.size = fmt.Sprintf("%d", size)
	return s
}

// Color sets the progress color
func (s *SingleRound) Color(color string) *SingleRound {
	s.color = color
	return s
}

// TrackColor sets the track color
func (s *SingleRound) TrackColor(color string) *SingleRound {
	s.trackColor = color
	return s
}

// ShowValue enables/disables value display
func (s *SingleRound) ShowValue(show bool) *SingleRound {
	s.showValue = show
	return s
}

// ID sets a custom ID
func (s *SingleRound) ID(id string) *SingleRound {
	s.id = id
	return s
}

// RenderContext provides data for rendering
type SingleRoundRenderContext struct {
	ID         string
	Title      string
	Subtitle   string
	Value      float64
	Max        float64
	Percentage float64
	Label      string
	Size       string
	Color      string
	TrackColor string
	ShowValue  bool
}

// Render prepares the single round chart for rendering
func (s *SingleRound) Render() *SingleRoundRenderContext {
	id := s.id
	if id == "" {
		id = "single-round-" + generateID()
	}

	percentage := 0.0
	if s.max > 0 {
		percentage = (s.value / s.max) * 100
	}

	label := s.label
	if label == "" && s.showValue {
		label = fmt.Sprintf("%.0f", s.value)
	}

	return &SingleRoundRenderContext{
		ID:         id,
		Title:      s.title,
		Subtitle:   s.subtitle,
		Value:      s.value,
		Max:        s.max,
		Percentage: percentage,
		Label:      label,
		Size:       s.size,
		Color:      s.color,
		TrackColor: s.trackColor,
		ShowValue:  s.showValue,
	}
}

// generateID creates a simple random ID
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
