// Package form provides standalone tool forms for admin pages.
// Tool forms can be embedded in pages or displayed in modals without registering routes.
package form

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	formbuilder "github.com/zhenyangze/goadmin/form"
)

// Handler is the callback for processing form submissions.
type Handler func(ctx context.Context, values url.Values) error

// ToolForm is a standalone form widget that can be embedded in pages or displayed in modals.
type ToolForm struct {
	form        *formbuilder.Builder
	handler     Handler
	defaults    map[string]interface{}
	successMsg  string
	errorMsg    string
	savedScript string
	errorScript string
	action      string
	method      string
	ajax        bool
	layout      string // horizontal, vertical, inline
	title       string
	description string
	async       bool // Load form asynchronously via AJAX
	asyncURL    string
}

// New creates a new tool form widget.
func New() *ToolForm {
	return &ToolForm{
		form:       formbuilder.New(),
		defaults:   make(map[string]interface{}),
		successMsg: "操作成功",
		errorMsg:   "操作失败",
		method:     "POST",
		layout:     "vertical",
	}
}

// Title sets the form title.
func (f *ToolForm) Title(title string) *ToolForm {
	f.title = title
	return f
}

// Description sets the form description.
func (f *ToolForm) Description(desc string) *ToolForm {
	f.description = desc
	return f
}

// Action sets the form action URL.
func (f *ToolForm) Action(action string) *ToolForm {
	f.action = action
	return f
}

// Method sets the form HTTP method.
func (f *ToolForm) Method(method string) *ToolForm {
	f.method = method
	return f
}

// Ajax enables/disables AJAX submission.
func (f *ToolForm) Ajax(enable bool) *ToolForm {
	f.ajax = enable
	return f
}

// Layout sets the form layout (horizontal, vertical, inline).
func (f *ToolForm) Layout(layout string) *ToolForm {
	f.layout = layout
	return f
}

// Field adds a field to the form (returns *ToolForm for chaining).
func (f *ToolForm) Field(name, label string, typ formbuilder.FieldType) *ToolForm {
	f.form.Field(name, label, typ)
	return f
}

// Hidden adds a hidden field.
func (f *ToolForm) Hidden(name string) *ToolForm {
	f.form.Hidden(name)
	return f
}

// Text adds a text field.
func (f *ToolForm) Text(name, label string) *ToolForm {
	f.form.Text(name, label)
	return f
}

// Number adds a number field.
func (f *ToolForm) Number(name, label string) *ToolForm {
	f.form.Number(name, label)
	return f
}

// Email adds an email field.
func (f *ToolForm) Email(name, label string) *ToolForm {
	f.form.Email(name, label)
	return f
}

// Password adds a password field.
func (f *ToolForm) Password(name, label string) *ToolForm {
	f.form.Password(name, label)
	return f
}

// Textarea adds a textarea field.
func (f *ToolForm) Textarea(name, label string) *ToolForm {
	f.form.Textarea(name, label)
	return f
}

// Editor adds a rich text editor field.
func (f *ToolForm) Editor(name, label string) *ToolForm {
	f.form.Editor(name, label)
	return f
}

// Select adds a select field.
func (f *ToolForm) Select(name, label string, options ...formbuilder.Option) *ToolForm {
	f.form.Select(name, label, options...)
	return f
}

// MultiSelect adds a multiple select field.
func (f *ToolForm) MultiSelect(name, label string, options ...formbuilder.Option) *ToolForm {
	f.form.MultiSelect(name, label, options...)
	return f
}

// Radio adds a radio button group field.
func (f *ToolForm) Radio(name, label string, options ...formbuilder.Option) *ToolForm {
	f.form.Radio(name, label, options...)
	return f
}

// Checkbox adds a checkbox group field.
func (f *ToolForm) Checkbox(name, label string, options ...formbuilder.Option) *ToolForm {
	f.form.Checkbox(name, label, options...)
	return f
}

// Switch adds a boolean checkbox field.
func (f *ToolForm) Switch(name, label string) *ToolForm {
	f.form.Switch(name, label)
	return f
}

// Tags adds a comma-separated tags field.
func (f *ToolForm) Tags(name, label string) *ToolForm {
	f.form.Tags(name, label)
	return f
}

// Date adds a date field.
func (f *ToolForm) Date(name, label string) *ToolForm {
	f.form.Date(name, label)
	return f
}

// Time adds a time field.
func (f *ToolForm) Time(name, label string) *ToolForm {
	f.form.Time(name, label)
	return f
}

// Datetime adds a datetime-local field.
func (f *ToolForm) Datetime(name, label string) *ToolForm {
	f.form.Datetime(name, label)
	return f
}

// DateRange adds a paired date range field.
func (f *ToolForm) DateRange(startName, endName, label string) *ToolForm {
	f.form.DateRange(startName, endName, label)
	return f
}

// Color adds a color picker field.
func (f *ToolForm) Color(name, label string) *ToolForm {
	f.form.Color(name, label)
	return f
}

// Icon adds an icon picker field.
func (f *ToolForm) Icon(name, label string) *ToolForm {
	f.form.Icon(name, label)
	return f
}

// Range adds a range slider field.
func (f *ToolForm) Range(name, label string) *ToolForm {
	f.form.Range(name, label)
	return f
}

// Rate adds a star rating field.
func (f *ToolForm) Rate(name, label string) *ToolForm {
	f.form.Rate(name, label)
	return f
}

// Currency adds a currency input field.
func (f *ToolForm) Currency(name, label string) *ToolForm {
	f.form.Currency(name, label)
	return f
}

// KeyValue adds a dynamic key-value pairs field.
func (f *ToolForm) KeyValue(name, label string) *ToolForm {
	f.form.KeyValue(name, label)
	return f
}

// Divider adds a divider field.
func (f *ToolForm) Divider() *ToolForm {
	f.form.Divider()
	return f
}

// Html adds raw HTML content.
func (f *ToolForm) Html(content string) *ToolForm {
	f.form.Html(content)
	return f
}

// IP adds an IP address field.
func (f *ToolForm) IP(name, label string) *ToolForm {
	f.form.IP(name, label)
	return f
}

// Mobile adds a mobile phone field.
func (f *ToolForm) Mobile(name, label string) *ToolForm {
	f.form.Mobile(name, label)
	return f
}

// URL adds a URL field.
func (f *ToolForm) URL(name, label string) *ToolForm {
	f.form.URL(name, label)
	return f
}

// TimeRange adds a time range field.
func (f *ToolForm) TimeRange(startName, endName, label string) *ToolForm {
	f.form.TimeRange(startName, endName, label)
	return f
}

// DateTimeRange adds a datetime range field.
func (f *ToolForm) DateTimeRange(startName, endName, label string) *ToolForm {
	f.form.DateTimeRange(startName, endName, label)
	return f
}

// Upload adds a file upload field.
func (f *ToolForm) Upload(name, label string) *ToolForm {
	f.form.Upload(name, label)
	return f
}

// Image adds an image upload field.
func (f *ToolForm) Image(name, label string) *ToolForm {
	f.form.Image(name, label)
	return f
}

// MultipleImage adds a multiple image upload field.
func (f *ToolForm) MultipleImage(name, label string) *ToolForm {
	f.form.MultipleImage(name, label)
	return f
}

// MultipleFile adds a multiple file upload field.
func (f *ToolForm) MultipleFile(name, label string) *ToolForm {
	f.form.MultipleFile(name, label)
	return f
}

// Repeater adds a repeater field and returns the builder for further configuration.
func (f *ToolForm) Repeater(name, label string) *formbuilder.RepeaterBuilder {
	return f.form.Repeater(name, label)
}

// Markdown adds a markdown editor field.
func (f *ToolForm) Markdown(name, label string) *ToolForm {
	f.form.Markdown(name, label)
	return f
}

// Slider adds a slider field.
func (f *ToolForm) Slider(name, label string) *ToolForm {
	f.form.Slider(name, label)
	return f
}

// Autocomplete adds an autocomplete field.
func (f *ToolForm) Autocomplete(name, label string) *ToolForm {
	f.form.Autocomplete(name, label)
	return f
}

// Listbox adds a listbox field.
func (f *ToolForm) Listbox(name, label string, options ...formbuilder.Option) *ToolForm {
	f.form.Listbox(name, label, options...)
	return f
}

// Map adds a map coordinate picker field.
func (f *ToolForm) Map(latitudeField, longitudeField, label string) *ToolForm {
	f.form.Map(latitudeField, longitudeField, label)
	return f
}

// Tree adds a tree selection field.
func (f *ToolForm) Tree(name, label string) *ToolForm {
	f.form.Tree(name, label)
	return f
}

// SelectTable adds a table selection field.
func (f *ToolForm) SelectTable(name, label string) *ToolForm {
	f.form.SelectTable(name, label)
	return f
}

// Table adds a table form field.
func (f *ToolForm) Table(name, label string, fn func(*formbuilder.TableBuilder)) *ToolForm {
	f.form.Table(name, label, fn)
	return f
}

// HasMany adds a has-many relation field.
func (f *ToolForm) HasMany(name, label string, fn func(*formbuilder.NestedFormBuilder)) *ToolForm {
	f.form.HasMany(name, label, fn)
	return f
}

// Embeds adds an embeds field.
func (f *ToolForm) Embeds(name, label string, fn func(*formbuilder.NestedFormBuilder)) *ToolForm {
	f.form.Embeds(name, label, fn)
	return f
}

// Handle sets the form submission handler.
func (f *ToolForm) Handle(handler Handler) *ToolForm {
	f.handler = handler
	return f
}

// Default sets default values for form fields.
func (f *ToolForm) Default(defaults map[string]interface{}) *ToolForm {
	f.defaults = defaults
	return f
}

// DefaultValue sets a single default value for a form field.
func (f *ToolForm) DefaultValue(name string, value interface{}) *ToolForm {
	f.defaults[name] = value
	return f
}

// Success sets the success message and optional JavaScript callback.
func (f *ToolForm) Success(message string, script ...string) *ToolForm {
	f.successMsg = message
	if len(script) > 0 {
		f.savedScript = script[0]
	}
	return f
}

// Error sets the error message and optional JavaScript callback.
func (f *ToolForm) Error(message string, script ...string) *ToolForm {
	f.errorMsg = message
	if len(script) > 0 {
		f.errorScript = script[0]
	}
	return f
}

// Async enables async loading of the form via AJAX.
func (f *ToolForm) Async(url string) *ToolForm {
	f.async = true
	f.asyncURL = url
	return f
}

// IsAsync returns whether the form is loaded asynchronously.
func (f *ToolForm) IsAsync() bool {
	return f.async
}

// AsyncURL returns the URL for async loading.
func (f *ToolForm) AsyncURL() string {
	return f.asyncURL
}

// Builder returns the underlying form builder.
func (f *ToolForm) Builder() *formbuilder.Builder {
	return f.form
}

// Fields returns all form fields.
func (f *ToolForm) Fields() []*formbuilder.Field {
	return f.form.Fields
}

// Handler returns the form submission handler.
func (f *ToolForm) GetHandler() Handler {
	return f.handler
}

// Defaults returns the default values.
func (f *ToolForm) Defaults() map[string]interface{} {
	return f.defaults
}

// RenderContext provides data for rendering the form.
type RenderContext struct {
	Title       string
	Description string
	Action      string
	Method      string
	Ajax        bool
	Layout      string
	Fields      []*formbuilder.Field
	Defaults    map[string]interface{}
	CSRF        string
	Prefix      string
}

// Render prepares the form for rendering.
func (f *ToolForm) Render(csrf, prefix string) *RenderContext {
	return &RenderContext{
		Title:       f.title,
		Description: f.description,
		Action:      f.action,
		Method:      f.method,
		Ajax:        f.ajax,
		Layout:      f.layout,
		Fields:      f.form.Fields,
		Defaults:    f.defaults,
		CSRF:        csrf,
		Prefix:      prefix,
	}
}

// Process handles form submission.
func (f *ToolForm) Process(ctx context.Context, values url.Values) (*Response, error) {
	if f.handler == nil {
		return nil, fmt.Errorf("no handler registered for form")
	}

	err := f.handler(ctx, values)
	if err != nil {
		return &Response{
			Success: false,
			Message: f.errorMsg,
			Script:  f.errorScript,
			Error:   err,
		}, nil
	}

	return &Response{
		Success: true,
		Message: f.successMsg,
		Script:  f.savedScript,
	}, nil
}

// Response is the form processing response.
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Script  string      `json:"script,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   error       `json:"-"`
}

// JSON returns the response as JSON.
func (r *Response) JSON() ([]byte, error) {
	return json.Marshal(r)
}

// HTML returns an HTML response for the form submission.
func (r *Response) HTML() template.HTML {
	if r.Success {
		if r.Script != "" {
			return template.HTML(fmt.Sprintf(`<script>%s</script>`, r.Script))
		}
		return template.HTML(`<div class="alert alert-success">` + r.Message + `</div>`)
	}

	if r.Script != "" {
		return template.HTML(fmt.Sprintf(`<script>%s</script>`, r.Script))
	}
	return template.HTML(`<div class="alert alert-danger">` + r.Message + `</div>`)
}

// LazyRenderable is the interface for async loaded forms.
type LazyRenderable interface {
	Render() template.HTML
}

// ModalForm wraps a ToolForm for modal display.
type ModalForm struct {
	form   *ToolForm
	title  string
	width  string // small, default, large, full
	height string
}

// NewModal creates a new modal form wrapper.
func NewModal(form *ToolForm, title string) *ModalForm {
	return &ModalForm{
		form:   form,
		title:  title,
		width:  "default",
		height: "auto",
	}
}

// Width sets the modal width.
func (m *ModalForm) Width(width string) *ModalForm {
	m.width = width
	return m
}

// Height sets the modal height.
func (m *ModalForm) Height(height string) *ModalForm {
	m.height = height
	return m
}

// Render returns the modal form data.
func (m *ModalForm) Render(csrf, prefix string) map[string]interface{} {
	return map[string]interface{}{
		"Title":  m.title,
		"Width":  m.width,
		"Height": m.height,
		"Form":   m.form.Render(csrf, prefix),
	}
}

// HTTPHandler creates an HTTP handler for the tool form.
// It handles both GET (render) and POST (submit) requests.
func HTTPHandler(form *ToolForm, tmpl *template.Template, getCSRF func() string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		prefix := "/admin" // This should be configurable

		switch r.Method {
		case http.MethodGet:
			// Render form
			renderCtx := form.Render(getCSRF(), prefix)
			// Execute template with render context
			data := map[string]interface{}{
				"Form": renderCtx,
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if err := tmpl.ExecuteTemplate(w, "toolform.tmpl", data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		case http.MethodPost:
			// Process form submission
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			resp, err := form.Process(ctx, r.PostForm)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Check if AJAX request
			if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)
				return
			}

			// HTML response
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if resp.Success {
				if resp.Script != "" {
					fmt.Fprintf(w, "<script>%s</script>", resp.Script)
				} else {
					fmt.Fprintf(w, "<div class=\"alert alert-success\">%s</div>", resp.Message)
				}
			} else {
				if resp.Script != "" {
					fmt.Fprintf(w, "<script>%s</script>", resp.Script)
				} else {
					fmt.Fprintf(w, "<div class=\"alert alert-danger\">%s</div>", resp.Message)
				}
			}

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
