package goadmin

import (
	"fmt"
	"html/template"
	"net/url"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func normalizePath(path string) string {
	path = "/" + strings.Trim(path, "/")
	if path == "/" {
		return path
	}
	return strings.TrimSuffix(path, "/")
}

func joinURL(parts ...string) string {
	var cleaned []string
	for i, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if i == 0 {
			cleaned = append(cleaned, strings.TrimRight(part, "/"))
			continue
		}
		cleaned = append(cleaned, strings.Trim(part, "/"))
	}
	out := strings.Join(cleaned, "/")
	if out == "" {
		return "/"
	}
	if !strings.HasPrefix(out, "/") {
		out = "/" + out
	}
	return out
}

func buildURL(base string, values url.Values) string {
	encoded := values.Encode()
	if encoded == "" {
		return base
	}
	return base + "?" + encoded
}

func valueFromPath(record any, path string) any {
	current := reflect.ValueOf(record)
	for _, part := range strings.Split(path, ".") {
		if !current.IsValid() {
			return nil
		}
		for current.Kind() == reflect.Pointer {
			if current.IsNil() {
				return nil
			}
			current = current.Elem()
		}
		switch current.Kind() {
		case reflect.Struct:
			current = current.FieldByName(part)
		case reflect.Map:
			current = current.MapIndex(reflect.ValueOf(part))
		default:
			return nil
		}
	}
	if !current.IsValid() {
		return nil
	}
	return current.Interface()
}

func valueListFromPath(record any, path string) []string {
	value := valueFromPath(record, path)
	if value == nil {
		return nil
	}
	rv := reflect.ValueOf(value)
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		out := make([]string, 0, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			item := rv.Index(i).Interface()
			id := valueFromPath(item, "ID")
			if id != nil {
				out = append(out, formatValue(id))
				continue
			}
			out = append(out, formatValue(item))
		}
		return out
	default:
		text := formatValue(value)
		if text == "" {
			return nil
		}
		return []string{text}
	}
}

func formatValue(value any) string {
	if value == nil {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return typed
	case fmt.Stringer:
		return typed.String()
	case time.Time:
		if typed.IsZero() {
			return ""
		}
		return typed.Format("2006-01-02 15:04")
	case []string:
		return strings.Join(typed, ", ")
	}

	rv := reflect.ValueOf(value)
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return ""
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		parts := make([]string, 0, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			item := rv.Index(i).Interface()
			name := valueFromPath(item, "Name")
			if text := formatValue(name); text != "" {
				parts = append(parts, text)
				continue
			}
			parts = append(parts, formatValue(item))
		}
		return strings.Join(parts, ", ")
	case reflect.Struct:
		if field := rv.FieldByName("Name"); field.IsValid() {
			return formatValue(field.Interface())
		}
		if field := rv.FieldByName("Title"); field.IsValid() {
			return formatValue(field.Interface())
		}
	}
	return fmt.Sprint(value)
}

func formatInputValue(value any, fieldType string) string {
	if value == nil {
		return ""
	}
	switch typed := value.(type) {
	case time.Time:
		if typed.IsZero() {
			return ""
		}
		if fieldType == "date" || fieldType == "daterange" {
			return typed.Format("2006-01-02")
		}
		if fieldType == "datetime-local" {
			return typed.Format("2006-01-02T15:04")
		}
		return typed.Format("2006-01-02 15:04")
	}
	return formatValue(value)
}

func isImagePath(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp":
		return true
	default:
		return false
	}
}

func splitCommaSeparated(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func toHTML(value string) template.HTML {
	return template.HTML(template.HTMLEscapeString(value))
}

func parseUint(raw string) uint {
	value, _ := strconv.ParseUint(raw, 10, 64)
	return uint(value)
}
