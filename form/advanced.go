package form

import "time"

// Extended Field Types
const (
	FieldMonth       FieldType = "month"
	FieldYear        FieldType = "year"
	FieldTimezone    FieldType = "timezone"
	FieldTel         FieldType = "tel"
	FieldCaptcha     FieldType = "captcha"
	FieldArray       FieldType = "array"
	FieldCascade     FieldType = "cascade"
	FieldFieldset    FieldType = "fieldset"
	FieldPlainInput  FieldType = "plaininput"
	FieldWebUploader FieldType = "webuploader"
	FieldNullable    FieldType = "nullable"
)

// ==================== Month Field ====================

// Month creates a month picker field
func Month(name, label string) *Field {
	return &Field{
		Name:  name,
		Label: label,
		Type:  FieldMonth,
	}
}

// ==================== Year Field ====================

// Year creates a year picker field
func Year(name, label string) *Field {
	return &Field{
		Name:  name,
		Label: label,
		Type:  FieldYear,
	}
}

// YearRange sets the selectable year range
func (f *Field) YearRange(start, end int) *Field {
	// Store range in Options for rendering
	f.Options = []Option{
		{Value: "start", Label: string(rune(start))},
		{Value: "end", Label: string(rune(end))},
	}
	return f
}

// ==================== Timezone Field ====================

// Timezone creates a timezone selector field
func Timezone(name, label string) *Field {
	field := &Field{
		Name:  name,
		Label: label,
		Type:  FieldTimezone,
	}
	// Populate with common timezones
	field.Options = getTimezoneOptions()
	return field
}

// getTimezoneOptions returns common timezone options
func getTimezoneOptions() []Option {
	zones := []string{
		"UTC",
		"Asia/Shanghai",
		"Asia/Tokyo",
		"Asia/Seoul",
		"Asia/Singapore",
		"Asia/Hong_Kong",
		"Asia/Taipei",
		"Asia/Bangkok",
		"Asia/Dubai",
		"Asia/Kolkata",
		"Europe/London",
		"Europe/Paris",
		"Europe/Berlin",
		"Europe/Moscow",
		"America/New_York",
		"America/Chicago",
		"America/Denver",
		"America/Los_Angeles",
		"America/Toronto",
		"America/Sao_Paulo",
		"Australia/Sydney",
		"Australia/Melbourne",
		"Pacific/Auckland",
	}

	opts := make([]Option, len(zones))
	for i, z := range zones {
		opts[i] = Option{Value: z, Label: z}
	}
	return opts
}

// ==================== Tel Field ====================

// Tel creates a telephone input field
func Tel(name, label string) *Field {
	return &Field{
		Name:  name,
		Label: label,
		Type:  FieldTel,
	}
}

// Pattern sets the validation pattern for tel
func (f *Field) Pattern(pattern string) *Field {
	// Store in placeholder temporarily, or add new field to Field struct
	f.MaskFormat = pattern
	return f
}

// ==================== Captcha Field ====================

// CaptchaConfig holds captcha configuration
type CaptchaConfig struct {
	Type     string // image, math, text
	Length   int
	Width    int
	Height   int
	TTL      time.Duration
}

// Captcha creates a captcha verification field
func Captcha(name, label string) *Field {
	return &Field{
		Name:  name,
		Label: label,
		Type:  FieldCaptcha,
	}
}

// CaptchaType sets the captcha type (image, math, text)
func (f *Field) CaptchaType(t string) *Field {
	// Store in Options
	f.Options = append(f.Options, Option{Value: "type", Label: t})
	return f
}

// CaptchaLength sets the captcha length
func (f *Field) CaptchaLength(length int) *Field {
	// Store min/max for length
	f.MinVal = float64Ptr(float64(length))
	return f
}

// Refreshable enables captcha refresh
func (f *Field) Refreshable() *Field {
	f.Options = append(f.Options, Option{Value: "refreshable", Label: "true"})
	return f
}

// ==================== Array Field ====================

// ArrayField creates an array input field (dynamic list of values)
func ArrayField(name, label string) *Field {
	return &Field{
		Name:  name,
		Label: label,
		Type:  FieldArray,
	}
}

// ArrayItemType sets the type of array items
func (f *Field) ArrayItemType(itemType string) *Field {
	f.Options = append(f.Options, Option{Value: "itemtype", Label: itemType})
	return f
}

// ArrayMinItems sets minimum number of items
func (f *Field) ArrayMinItems(min int) *Field {
	f.RepeaterMinRows = min
	return f
}

// ArrayMaxItems sets maximum number of items
func (f *Field) ArrayMaxItems(max int) *Field {
	// Store in MaxVal
	f.MaxVal = float64Ptr(float64(max))
	return f
}

// ==================== Cascade Field ====================

// CascadeGroup creates a cascading group field
func CascadeGroup(name, label string) *Field {
	return &Field{
		Name:  name,
		Label: label,
		Type:  FieldCascade,
	}
}

// CascadeLevels sets the cascade levels
func (f *Field) CascadeLevels(levels ...string) *Field {
	f.Options = []Option{}
	for _, level := range levels {
		f.Options = append(f.Options, Option{Value: level, Label: level})
	}
	return f
}

// CascadeData sets the hierarchical data for cascade
func (f *Field) CascadeData(data map[string][]Option) *Field {
	// Store cascade data - will need custom handling in template
	f.Type = FieldCascade
	return f
}

// ==================== Fieldset Field ====================

// Fieldset creates a fieldset/fieldset grouping
func Fieldset(title string, fn func(*Builder)) *Field {
	builder := &Builder{}
	fn(builder)

	return &Field{
		Name:           title,
		Label:          title,
		Type:           FieldFieldset,
		RepeaterFields: builder.Fields,
	}
}

// ==================== PlainInput Field ====================

// PlainInput creates a plain text input (no styling)
func PlainInput(name, label string) *Field {
	return &Field{
		Name:  name,
		Label: label,
		Type:  FieldPlainInput,
	}
}

// ==================== WebUploader Field ====================

// WebUploader creates an advanced web uploader field
func WebUploader(name, label string) *Field {
	return &Field{
		Name:  name,
		Label: label,
		Type:  FieldWebUploader,
	}
}

// Chunked enables chunked upload for large files
func (f *Field) Chunked(chunkSize int64) *Field {
	f.MaxFileSize = chunkSize
	f.Options = append(f.Options, Option{Value: "chunked", Label: "true"})
	return f
}

// Multiple enables multiple file upload
func (f *Field) WebUploaderMultiple() *Field {
	f.Multiple = true
	return f
}

// DragDrop enables drag and drop upload
func (f *Field) DragDrop() *Field {
	f.Options = append(f.Options, Option{Value: "dragdrop", Label: "true"})
	return f
}

// Preview enables file preview
func (f *Field) Preview() *Field {
	f.Options = append(f.Options, Option{Value: "preview", Label: "true"})
	return f
}

// ==================== Nullable Field ====================

// Nullable creates a nullable wrapper for other fields
func Nullable(field *Field) *Field {
	return &Field{
		Name:  field.Name,
		Label: field.Label,
		Type:  FieldNullable,
		// Store the actual field type in Options
		Options: []Option{{Value: "basetype", Label: string(field.Type)}},
	}
}

// NullableCheckboxLabel sets the checkbox label for nullable
func (f *Field) NullableCheckboxLabel(label string) *Field {
	f.Options = append(f.Options, Option{Value: "checkboxlabel", Label: label})
	return f
}

// NullableDefault sets the default null state
func (f *Field) NullableDefault(null bool) *Field {
	if null {
		f.Options = append(f.Options, Option{Value: "defaultnull", Label: "true"})
	}
	return f
}

// Helper function
func float64Ptr(f float64) *float64 {
	return &f
}
