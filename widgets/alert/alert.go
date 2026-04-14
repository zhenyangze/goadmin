// Package alert provides alert and notification components.
package alert

import (
	"encoding/json"
	"html/template"
)

// AlertType represents the type of alert.
type AlertType string

const (
	TypeSuccess AlertType = "success"
	TypeError   AlertType = "error"
	TypeWarning AlertType = "warning"
	TypeInfo    AlertType = "info"
)

// Alert is an alert component.
type Alert struct {
	typ         AlertType
	title       string
	message     string
	dismissible bool
	class       string
}

// New creates a new alert.
func New(alertType AlertType, message string) *Alert {
	return &Alert{
		typ:     alertType,
		message: message,
	}
}

// Success creates a success alert.
func Success(message string) *Alert {
	return New(TypeSuccess, message)
}

// Error creates an error alert.
func Error(message string) *Alert {
	return New(TypeError, message)
}

// Warning creates a warning alert.
func Warning(message string) *Alert {
	return New(TypeWarning, message)
}

// Info creates an info alert.
func Info(message string) *Alert {
	return New(TypeInfo, message)
}

// Title sets the alert title.
func (a *Alert) Title(title string) *Alert {
	a.title = title
	return a
}

// Dismissible makes the alert dismissible.
func (a *Alert) Dismissible(dismissible bool) *Alert {
	a.dismissible = dismissible
	return a
}

// Class adds CSS classes.
func (a *Alert) Class(class string) *Alert {
	a.class = class
	return a
}

// RenderContext provides data for rendering.
type RenderContext struct {
	Type        string
	Title       string
	Message     string
	Dismissible bool
	Class       string
}

// Render prepares the alert for rendering.
func (a *Alert) Render() *RenderContext {
	return &RenderContext{
		Type:        string(a.typ),
		Title:       a.title,
		Message:     a.message,
		Dismissible: a.dismissible,
		Class:       a.class,
	}
}

// Toast is a toast notification.
type Toast struct {
	typ     AlertType
	title   string
	message string
}

// NewToast creates a new toast.
func NewToast(alertType AlertType, message string) *Toast {
	return &Toast{
		typ:     alertType,
		message: message,
	}
}

// SuccessToast creates a success toast.
func SuccessToast(message string) *Toast {
	return NewToast(TypeSuccess, message)
}

// ErrorToast creates an error toast.
func ErrorToast(message string) *Toast {
	return NewToast(TypeError, message)
}

// WarningToast creates a warning toast.
func WarningToast(message string) *Toast {
	return NewToast(TypeWarning, message)
}

// InfoToast creates an info toast.
func InfoToast(message string) *Toast {
	return NewToast(TypeInfo, message)
}

// Title sets the toast title.
func (t *Toast) Title(title string) *Toast {
	t.title = title
	return t
}

// RenderContext returns the toast render context.
func (t *Toast) Render() map[string]interface{} {
	return map[string]interface{}{
		"type":    string(t.typ),
		"title":   t.title,
		"message": t.message,
	}
}

// FlashMessage is a flash message for session storage.
type FlashMessage struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

// FlashContainer holds flash messages.
type FlashContainer struct {
	messages []FlashMessage
}

// NewFlashContainer creates a new flash container.
func NewFlashContainer() *FlashContainer {
	return &FlashContainer{
		messages: make([]FlashMessage, 0),
	}
}

// Add adds a flash message.
func (f *FlashContainer) Add(alertType AlertType, message string) *FlashContainer {
	return f.AddWithTitle(alertType, "", message)
}

// AddWithTitle adds a flash message with title.
func (f *FlashContainer) AddWithTitle(alertType AlertType, title, message string) *FlashContainer {
	f.messages = append(f.messages, FlashMessage{
		Type:    string(alertType),
		Title:   title,
		Message: message,
	})
	return f
}

// Success adds a success flash.
func (f *FlashContainer) Success(message string) *FlashContainer {
	return f.Add(TypeSuccess, message)
}

// Error adds an error flash.
func (f *FlashContainer) Error(message string) *FlashContainer {
	return f.Add(TypeError, message)
}

// Warning adds a warning flash.
func (f *FlashContainer) Warning(message string) *FlashContainer {
	return f.Add(TypeWarning, message)
}

// Info adds an info flash.
func (f *FlashContainer) Info(message string) *FlashContainer {
	return f.Add(TypeInfo, message)
}

// Messages returns all flash messages.
func (f *FlashContainer) Messages() []FlashMessage {
	return f.messages
}

// JSON returns flash messages as JSON.
func (f *FlashContainer) JSON() string {
	if len(f.messages) == 0 {
		return "[]"
	}
	// Simple JSON encoding
	data, _ := json.Marshal(f.messages)
	return string(data)
}

// HTML returns the toast container HTML with flash messages.
func (f *FlashContainer) HTML() template.HTML {
	if len(f.messages) == 0 {
		return template.HTML(`<div data-flash-messages="[]"></div>`)
	}
	return template.HTML(`<div data-flash-messages="` + template.HTMLEscapeString(f.JSON()) + `"></div>`)
}
