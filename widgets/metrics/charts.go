package metrics

import (
	"encoding/json"
	"fmt"
	"time"
)

// DonutChart represents a donut/pie chart for metrics
type DonutChart struct {
	title    string
	subtitle string
	labels   []string
	data     []float64
	colors   []string
	height   string
	id       string
	showLegend bool
}

// Donut creates a new donut chart
func Donut() *DonutChart {
	return &DonutChart{
		height:     "350px",
		showLegend: true,
		colors:     []string{"#3B82F6", "#10B981", "#F59E0B", "#EF4444", "#8B5CF6", "#EC4899"},
	}
}

// Title sets the chart title
func (d *DonutChart) Title(title string) *DonutChart {
	d.title = title
	return d
}

// Subtitle sets the chart subtitle
func (d *DonutChart) Subtitle(subtitle string) *DonutChart {
	d.subtitle = subtitle
	return d
}

// Height sets the chart height
func (d *DonutChart) Height(height string) *DonutChart {
	d.height = height
	return d
}

// ID sets a custom ID
func (d *DonutChart) ID(id string) *DonutChart {
	d.id = id
	return d
}

// Colors sets custom colors
func (d *DonutChart) Colors(colors ...string) *DonutChart {
	d.colors = colors
	return d
}

// Data adds data points
func (d *DonutChart) Data(label string, value float64) *DonutChart {
	d.labels = append(d.labels, label)
	d.data = append(d.data, value)
	return d
}

// ShowLegend enables/disables legend
func (d *DonutChart) ShowLegend(show bool) *DonutChart {
	d.showLegend = show
	return d
}

// DonutRenderContext provides data for rendering
type DonutRenderContext struct {
	ID         string
	Title      string
	Subtitle   string
	Height     string
	DataJSON   string
	LabelsJSON string
	ColorsJSON string
	ShowLegend bool
}

// Render prepares the donut chart for rendering
func (d *DonutChart) Render() *DonutRenderContext {
	id := d.id
	if id == "" {
		id = "donut-chart-" + generateID()
	}

	dataJSON, _ := json.Marshal(d.data)
	labelsJSON, _ := json.Marshal(d.labels)
	colorsJSON, _ := json.Marshal(d.colors)

	return &DonutRenderContext{
		ID:         id,
		Title:      d.title,
		Subtitle:   d.subtitle,
		Height:     d.height,
		DataJSON:   string(dataJSON),
		LabelsJSON: string(labelsJSON),
		ColorsJSON: string(colorsJSON),
		ShowLegend: d.showLegend,
	}
}

// LineChart represents a line chart for metrics
type LineChart struct {
	title      string
	subtitle   string
	xAxis      []string
	series     []LineSeries
	height     string
	id         string
	colors     []string
	showArea   bool
	showPoints bool
}

// LineSeries represents a data series for line charts
type LineSeries struct {
	Name string
	Data []float64
}

// Line creates a new line chart
func Line() *LineChart {
	return &LineChart{
		height:     "350px",
		colors:     []string{"#3B82F6", "#10B981", "#F59E0B", "#EF4444", "#8B5CF6"},
		showPoints: true,
	}
}

// Title sets the chart title
func (l *LineChart) Title(title string) *LineChart {
	l.title = title
	return l
}

// Subtitle sets the chart subtitle
func (l *LineChart) Subtitle(subtitle string) *LineChart {
	l.subtitle = subtitle
	return l
}

// Height sets the chart height
func (l *LineChart) Height(height string) *LineChart {
	l.height = height
	return l
}

// ID sets a custom ID
func (l *LineChart) ID(id string) *LineChart {
	l.id = id
	return l
}

// XAxis sets the x-axis labels
func (l *LineChart) XAxis(labels ...string) *LineChart {
	l.xAxis = labels
	return l
}

// Series adds a data series
func (l *LineChart) Series(name string, data []float64) *LineChart {
	l.series = append(l.series, LineSeries{
		Name: name,
		Data: data,
	})
	return l
}

// ShowArea enables area fill under the line
func (l *LineChart) ShowArea(show bool) *LineChart {
	l.showArea = show
	return l
}

// ShowPoints enables point markers
func (l *LineChart) ShowPoints(show bool) *LineChart {
	l.showPoints = show
	return l
}

// LineRenderContext provides data for rendering
type LineRenderContext struct {
	ID         string
	Title      string
	Subtitle   string
	Height     string
	XAxisJSON  string
	SeriesJSON string
	ColorsJSON string
	ShowArea   bool
	ShowPoints bool
}

// Render prepares the line chart for rendering
func (l *LineChart) Render() *LineRenderContext {
	id := l.id
	if id == "" {
		id = "line-chart-" + generateID()
	}

	xAxisJSON, _ := json.Marshal(l.xAxis)
	seriesJSON, _ := json.Marshal(l.series)
	colorsJSON, _ := json.Marshal(l.colors)

	return &LineRenderContext{
		ID:         id,
		Title:      l.title,
		Subtitle:   l.subtitle,
		Height:     l.height,
		XAxisJSON:  string(xAxisJSON),
		SeriesJSON: string(seriesJSON),
		ColorsJSON: string(colorsJSON),
		ShowArea:   l.showArea,
		ShowPoints: l.showPoints,
	}
}

// BarChart represents a bar chart for metrics
type BarChart struct {
	title      string
	subtitle   string
	xAxis      []string
	series     []BarSeries
	height     string
	id         string
	colors     []string
	horizontal bool
	stacked    bool
}

// BarSeries represents a data series for bar charts
type BarSeries struct {
	Name string
	Data []float64
}

// Bar creates a new bar chart
func Bar() *BarChart {
	return &BarChart{
		height: "350px",
		colors: []string{"#3B82F6", "#10B981", "#F59E0B", "#EF4444", "#8B5CF6"},
	}
}

// Title sets the chart title
func (b *BarChart) Title(title string) *BarChart {
	b.title = title
	return b
}

// Subtitle sets the chart subtitle
func (b *BarChart) Subtitle(subtitle string) *BarChart {
	b.subtitle = subtitle
	return b
}

// Height sets the chart height
func (b *BarChart) Height(height string) *BarChart {
	b.height = height
	return b
}

// ID sets a custom ID
func (b *BarChart) ID(id string) *BarChart {
	b.id = id
	return b
}

// XAxis sets the x-axis labels
func (b *BarChart) XAxis(labels ...string) *BarChart {
	b.xAxis = labels
	return b
}

// Series adds a data series
func (b *BarChart) Series(name string, data []float64) *BarChart {
	b.series = append(b.series, BarSeries{
		Name: name,
		Data: data,
	})
	return b
}

// Horizontal sets the bar orientation
func (b *BarChart) Horizontal(horizontal bool) *BarChart {
	b.horizontal = horizontal
	return b
}

// Stacked enables stacked bars
func (b *BarChart) Stacked(stacked bool) *BarChart {
	b.stacked = stacked
	return b
}

// BarRenderContext provides data for rendering
type BarRenderContext struct {
	ID         string
	Title      string
	Subtitle   string
	Height     string
	XAxisJSON  string
	SeriesJSON string
	ColorsJSON string
	Horizontal bool
	Stacked    bool
}

// Render prepares the bar chart for rendering
func (b *BarChart) Render() *BarRenderContext {
	id := b.id
	if id == "" {
		id = "bar-chart-" + generateID()
	}

	xAxisJSON, _ := json.Marshal(b.xAxis)
	seriesJSON, _ := json.Marshal(b.series)
	colorsJSON, _ := json.Marshal(b.colors)

	return &BarRenderContext{
		ID:         id,
		Title:      b.title,
		Subtitle:   b.subtitle,
		Height:     b.height,
		XAxisJSON:  string(xAxisJSON),
		SeriesJSON: string(seriesJSON),
		ColorsJSON: string(colorsJSON),
		Horizontal: b.horizontal,
		Stacked:    b.stacked,
	}
}

// generateID creates a simple random ID
func generateID2() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
