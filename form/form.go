package form

// FieldType is the rendered form input variant.
type FieldType string

const (
	FieldHidden    FieldType = "hidden"
	FieldDisplay   FieldType = "display"
	FieldText      FieldType = "text"
	FieldNumber    FieldType = "number"
	FieldEmail     FieldType = "email"
	FieldURL       FieldType = "url"
	FieldIP        FieldType = "ip"
	FieldMobile    FieldType = "mobile"
	FieldPassword  FieldType = "password"
	FieldTextarea  FieldType = "textarea"
	FieldSelect    FieldType = "select"
	FieldMulti     FieldType = "multiselect"
	FieldRadio     FieldType = "radio"
	FieldCheckbox  FieldType = "checkbox"
	FieldSwitch    FieldType = "switch"
	FieldTags      FieldType = "tags"
	FieldDate      FieldType = "date"
	FieldTime      FieldType = "time"
	FieldDatetime  FieldType = "datetime-local"
	FieldDateRange FieldType = "daterange"
	FieldTimeRange FieldType = "timerange"
	FieldDateTimeRange FieldType = "datetimerange"
	FieldUpload    FieldType = "file"
	FieldImage     FieldType = "image"
	FieldMultipleImage FieldType = "multipleimage"
	FieldMultipleFile  FieldType = "multiplefile"
	FieldRepeater  FieldType = "repeater"
	FieldEditor    FieldType = "editor"
	FieldMarkdown  FieldType = "markdown"
	FieldColor     FieldType = "color"
	FieldIcon      FieldType = "icon"
	FieldRange     FieldType = "range"
	FieldRate      FieldType = "rate"
	FieldCurrency  FieldType = "currency"
	FieldKeyValue  FieldType = "keyvalue"
	FieldDivider   FieldType = "divider"
	FieldHtml      FieldType = "html"
	FieldSlider    FieldType = "slider"
	FieldAutocomplete FieldType = "autocomplete"
	FieldListbox   FieldType = "listbox"
	FieldMap       FieldType = "map"
	FieldTree      FieldType = "tree"
	FieldSelectTable FieldType = "selecttable"
	FieldTable     FieldType = "table"
	FieldHasMany   FieldType = "hasmany"
	FieldEmbeds    FieldType = "embeds"
)

// Option represents one select option.
type Option struct {
	Value string
	Label string
}

// Field defines a form input.
type Field struct {
	Name              string
	SecondName        string
	Label             string
	Type              FieldType
	ValuePath         string
	SecondValuePath   string
	Multiple          bool
	MaxFileSize       int64
	AllowedExtensions []string
	Required          bool
	Readonly          bool
	Disabled          bool
	Help              string
	Placeholder       string
	SecondPlaceholder string
	Options           []Option
	RepeaterFields    []*Field
	RepeaterMinRows   int
	// Number/Range fields
	MinVal   *float64
	MaxVal   *float64
	StepVal  *float64
	// Display options
	IsInline bool
	// Editor height
	EditorHeight int
	// Currency symbol
	CurrencySymbol string
	// KeyValue placeholders
	KeyPlaceholder   string
	ValuePlaceholder string
	// Mobile mask format
	MaskFormat string
	// Slider options
	SliderMin     *float64
	SliderMax     *float64
	SliderStep    *float64
	SliderPostfix string
	// Autocomplete AJAX URL
	AutocompleteURL string
	// Image/Upload options
	UploadDir       string
	UploadMaxCount  int
	IsSortable      bool
	// Html content
	HtmlContent string
	// Map provider (tencent, amap, baidu, google)
	MapProvider string
	// Tree nodes data
	TreeNodes []TreeNode
	TreeIDColumn string
	TreeTitleColumn string
	TreeParentColumn string
	TreeExpand bool
	TreeAllowParentSelect bool
	// SelectTable config
	SelectTableURL string
	SelectTableTitle string
	SelectTableDialogWidth string
	SelectTableDisplayField string
	SelectTableValueField string
	// Table/HasMany/Embeds fields
	NestedFields []*Field
	// HasMany options
	HasManyLabel string
	HasManyTableMode bool
}

// TreeNode represents a node in tree field
type TreeNode struct {
	ID       string      `json:"id"`
	Title    string      `json:"title"`
	ParentID string      `json:"parent_id,omitempty"`
	Children []TreeNode  `json:"children,omitempty"`
}

// Builder defines one resource form.
type Builder struct {
	Title         string
	Description   string
	Fields        []*Field
	HideDelete    bool
	SubmitLabel   string
	DeleteLabel   string
	CancelBackURL string
}

// New creates a form builder.
func New() *Builder {
	return &Builder{
		SubmitLabel: "Save",
		DeleteLabel: "Delete",
	}
}

// Field adds a field.
func (b *Builder) Field(name, label string, typ FieldType) *Field {
	field := &Field{Name: name, Label: label, Type: typ}
	b.Fields = append(b.Fields, field)
	return field
}

// Hidden adds a hidden field.
func (b *Builder) Hidden(name string) *Field {
	return b.Field(name, "", FieldHidden)
}

// Display adds a read-only field.
func (b *Builder) Display(name, label string) *Field {
	field := b.Field(name, label, FieldDisplay)
	field.Readonly = true
	return field
}

// Text adds a text field.
func (b *Builder) Text(name, label string) *Field {
	return b.Field(name, label, FieldText)
}

// Number adds a number input field.
func (b *Builder) Number(name, label string) *Field {
	return b.Field(name, label, FieldNumber)
}

// Email adds an email field.
func (b *Builder) Email(name, label string) *Field {
	return b.Field(name, label, FieldEmail)
}

// URL adds a URL field.
func (b *Builder) URL(name, label string) *Field {
	return b.Field(name, label, FieldURL)
}

// Password adds a password field.
func (b *Builder) Password(name, label string) *Field {
	return b.Field(name, label, FieldPassword)
}

// Textarea adds a textarea field.
func (b *Builder) Textarea(name, label string) *Field {
	return b.Field(name, label, FieldTextarea)
}

// Editor adds a rich text editor field.
func (b *Builder) Editor(name, label string) *Field {
	return b.Field(name, label, FieldEditor)
}

// Select adds a select field.
func (b *Builder) Select(name, label string, options ...Option) *Field {
	field := b.Field(name, label, FieldSelect)
	field.Options = append(field.Options, options...)
	return field
}

// MultiSelect adds a multiple select field.
func (b *Builder) MultiSelect(name, label string, options ...Option) *Field {
	field := b.Field(name, label, FieldMulti)
	field.Options = append(field.Options, options...)
	return field
}

// Radio adds a radio button group field.
func (b *Builder) Radio(name, label string, options ...Option) *Field {
	field := b.Field(name, label, FieldRadio)
	field.Options = append(field.Options, options...)
	return field
}

// Checkbox adds a checkbox group field.
func (b *Builder) Checkbox(name, label string, options ...Option) *Field {
	field := b.Field(name, label, FieldCheckbox)
	field.Options = append(field.Options, options...)
	return field
}

// Switch adds a boolean checkbox field.
func (b *Builder) Switch(name, label string) *Field {
	return b.Field(name, label, FieldSwitch)
}

// Tags adds a comma-separated tags field.
func (b *Builder) Tags(name, label string) *Field {
	return b.Field(name, label, FieldTags)
}

// Date adds a date field.
func (b *Builder) Date(name, label string) *Field {
	return b.Field(name, label, FieldDate)
}

// Time adds a time field.
func (b *Builder) Time(name, label string) *Field {
	return b.Field(name, label, FieldTime)
}

// Datetime adds a datetime-local field.
func (b *Builder) Datetime(name, label string) *Field {
	return b.Field(name, label, FieldDatetime)
}

// DateRange adds a paired date range field.
func (b *Builder) DateRange(startName, endName, label string) *Field {
	field := b.Field(startName, label, FieldDateRange)
	field.SecondName = endName
	return field
}

// Color adds a color picker field.
func (b *Builder) Color(name, label string) *Field {
	return b.Field(name, label, FieldColor)
}

// Icon adds an icon picker field.
func (b *Builder) Icon(name, label string) *Field {
	return b.Field(name, label, FieldIcon)
}

// Range adds a range slider field.
func (b *Builder) Range(name, label string) *Field {
	return b.Field(name, label, FieldRange)
}

// Rate adds a star rating field.
func (b *Builder) Rate(name, label string) *Field {
	return b.Field(name, label, FieldRate)
}

// Currency adds a currency input field.
func (b *Builder) Currency(name, label string) *Field {
	return b.Field(name, label, FieldCurrency)
}

// KeyValue adds a dynamic key-value pairs field.
func (b *Builder) KeyValue(name, label string) *Field {
	return b.Field(name, label, FieldKeyValue)
}

// IP adds an IP address input field.
func (b *Builder) IP(name, label string) *Field {
	return b.Field(name, label, FieldIP)
}

// Mobile adds a mobile phone input field with optional mask.
func (b *Builder) Mobile(name, label string) *Field {
	return b.Field(name, label, FieldMobile)
}

// TimeRange adds a paired time range field.
func (b *Builder) TimeRange(startName, endName, label string) *Field {
	field := b.Field(startName, label, FieldTimeRange)
	field.SecondName = endName
	return field
}

// DateTimeRange adds a paired datetime range field.
func (b *Builder) DateTimeRange(startName, endName, label string) *Field {
	field := b.Field(startName, label, FieldDateTimeRange)
	field.SecondName = endName
	return field
}

// Image adds an image upload field.
func (b *Builder) Image(name, label string) *Field {
	return b.Field(name, label, FieldImage)
}

// MultipleImage adds a multiple image upload field.
func (b *Builder) MultipleImage(name, label string) *Field {
	field := b.Field(name, label, FieldMultipleImage)
	field.Multiple = true
	return field
}

// MultipleFile adds a multiple file upload field.
func (b *Builder) MultipleFile(name, label string) *Field {
	field := b.Field(name, label, FieldMultipleFile)
	field.Multiple = true
	return field
}

// Markdown adds a markdown editor field.
func (b *Builder) Markdown(name, label string) *Field {
	return b.Field(name, label, FieldMarkdown)
}

// Slider adds a slider input field.
func (b *Builder) Slider(name, label string) *Field {
	return b.Field(name, label, FieldSlider)
}

// Autocomplete adds an autocomplete input field.
func (b *Builder) Autocomplete(name, label string) *Field {
	return b.Field(name, label, FieldAutocomplete)
}

// Listbox adds a listbox multi-select field.
func (b *Builder) Listbox(name, label string, options ...Option) *Field {
	field := b.Field(name, label, FieldListbox)
	field.Options = append(field.Options, options...)
	return field
}

// Divider adds a visual divider field.
func (b *Builder) Divider() *Field {
	return b.Field("_divider", "", FieldDivider)
}

// Html adds a custom HTML content field.
func (b *Builder) Html(content string) *Field {
	field := b.Field("_html", "", FieldHtml)
	field.HtmlContent = content
	return field
}

// Map adds a map picker field for latitude/longitude selection.
func (b *Builder) Map(latitudeField, longitudeField, label string) *Field {
	field := b.Field(latitudeField, label, FieldMap)
	field.SecondName = longitudeField
	field.MapProvider = "tencent" // default provider
	return field
}

// Tree adds a tree selector field.
func (b *Builder) Tree(name, label string) *Field {
	return b.Field(name, label, FieldTree)
}

// SelectTable adds a table selector field with dialog.
func (b *Builder) SelectTable(name, label string) *Field {
	return b.Field(name, label, FieldSelectTable)
}

// Table adds a table form field for editing JSON 2D array.
func (b *Builder) Table(name, label string, fn func(*TableBuilder)) *Field {
	field := b.Field(name, label, FieldTable)
	tb := &TableBuilder{field: field}
	if fn != nil {
		fn(tb)
	}
	return field
}

// HasMany adds a has-many relation form field.
func (b *Builder) HasMany(name, label string, fn func(*NestedFormBuilder)) *Field {
	field := b.Field(name, label, FieldHasMany)
	field.HasManyLabel = label
	nb := &NestedFormBuilder{field: field}
	if fn != nil {
		fn(nb)
	}
	return field
}

// Embeds adds an embeds form field for JSON object editing.
func (b *Builder) Embeds(name, label string, fn func(*NestedFormBuilder)) *Field {
	field := b.Field(name, label, FieldEmbeds)
	nb := &NestedFormBuilder{field: field}
	if fn != nil {
		fn(nb)
	}
	return field
}

// Upload adds a file upload field.
func (b *Builder) Upload(name, label string) *Field {
	return b.Field(name, label, FieldUpload)
}

// Repeater adds a has-many style nested editor backed by a JSON string field.
func (b *Builder) Repeater(name, label string) *RepeaterBuilder {
	field := b.Field(name, label, FieldRepeater)
	field.RepeaterMinRows = 3
	return &RepeaterBuilder{field: field}
}

// MultiUpload adds a multiple-file upload field.
func (b *Builder) MultiUpload(name, label string) *Field {
	field := b.Field(name, label, FieldUpload)
	field.Multiple = true
	return field
}

// MarkRequired makes validation intent visible in the UI.
func (f *Field) MarkRequired() *Field {
	f.Required = true
	return f
}

// WithHelp sets contextual help text.
func (f *Field) WithHelp(help string) *Field {
	f.Help = help
	return f
}

// WithPlaceholder sets the placeholder text.
func (f *Field) WithPlaceholder(placeholder string) *Field {
	f.Placeholder = placeholder
	return f
}

// ValueFrom overrides the record field used to populate the form value.
func (f *Field) ValueFrom(path string) *Field {
	f.ValuePath = path
	return f
}

// SecondValueFrom overrides the record field used to populate the second value.
func (f *Field) SecondValueFrom(path string) *Field {
	f.SecondValuePath = path
	return f
}

// WithSecondPlaceholder sets the placeholder for the second input in paired fields.
func (f *Field) WithSecondPlaceholder(placeholder string) *Field {
	f.SecondPlaceholder = placeholder
	return f
}

// AllowExtensions restricts upload file extensions.
func (f *Field) AllowExtensions(exts ...string) *Field {
	f.AllowedExtensions = append(f.AllowedExtensions, exts...)
	return f
}

// MaxSize limits upload file size in bytes.
func (f *Field) MaxSize(size int64) *Field {
	f.MaxFileSize = size
	return f
}

// Min sets the minimum value for number/range fields.
func (f *Field) Min(min float64) *Field {
	f.MinVal = &min
	return f
}

// Max sets the maximum value for number/range fields.
func (f *Field) Max(max float64) *Field {
	f.MaxVal = &max
	return f
}

// Step sets the step value for number/range fields.
func (f *Field) Step(step float64) *Field {
	f.StepVal = &step
	return f
}

// Inline sets the inline display mode for radio/checkbox fields.
func (f *Field) Inline(inline bool) *Field {
	f.IsInline = inline
	return f
}

// Height sets the editor height in pixels.
func (f *Field) Height(height int) *Field {
	f.EditorHeight = height
	return f
}

// Symbol sets the currency symbol.
func (f *Field) Symbol(symbol string) *Field {
	f.CurrencySymbol = symbol
	return f
}

// WithKeyPlaceholder sets the key placeholder for key-value fields.
func (f *Field) WithKeyPlaceholder(placeholder string) *Field {
	f.KeyPlaceholder = placeholder
	return f
}

// WithValuePlaceholder sets the value placeholder for key-value fields.
func (f *Field) WithValuePlaceholder(placeholder string) *Field {
	f.ValuePlaceholder = placeholder
	return f
}

// Disabled marks the field as disabled.
func (f *Field) MarkDisabled() *Field {
	f.Disabled = true
	return f
}

// WithMask sets the mask format for mobile fields.
func (f *Field) WithMask(mask string) *Field {
	f.MaskFormat = mask
	return f
}

// WithSliderOptions sets the slider options.
func (f *Field) WithSliderOptions(min, max float64, step float64, postfix string) *Field {
	f.SliderMin = &min
	f.SliderMax = &max
	f.SliderStep = &step
	f.SliderPostfix = postfix
	return f
}

// WithAutocompleteURL sets the AJAX URL for autocomplete fields.
func (f *Field) WithAutocompleteURL(url string) *Field {
	f.AutocompleteURL = url
	return f
}

// Dir sets the upload directory for file/image fields.
func (f *Field) Dir(dir string) *Field {
	f.UploadDir = dir
	return f
}

// Limit sets the maximum upload count for multiple file/image fields.
func (f *Field) Limit(count int) *Field {
	f.UploadMaxCount = count
	return f
}

// Sortable enables sorting for multiple file/image fields.
func (f *Field) Sortable() *Field {
	f.IsSortable = true
	return f
}

// MapProvider sets the map provider (tencent, amap, baidu, google).
func (f *Field) Provider(provider string) *Field {
	f.MapProvider = provider
	return f
}

// Nodes sets the tree nodes data.
func (f *Field) Nodes(nodes []TreeNode) *Field {
	f.TreeNodes = nodes
	return f
}

// SetIDColumn sets the tree node ID column name.
func (f *Field) SetIDColumn(column string) *Field {
	f.TreeIDColumn = column
	return f
}

// SetTitleColumn sets the tree node title column name.
func (f *Field) SetTitleColumn(column string) *Field {
	f.TreeTitleColumn = column
	return f
}

// SetParentColumn sets the tree node parent column name.
func (f *Field) SetParentColumn(column string) *Field {
	f.TreeParentColumn = column
	return f
}

// Expand sets whether tree nodes are expanded by default.
func (f *Field) Expand(expand bool) *Field {
	f.TreeExpand = expand
	return f
}

// AllowParentNode allows selecting parent nodes in tree.
func (f *Field) AllowParentNode() *Field {
	f.TreeAllowParentSelect = true
	return f
}

// From sets the select table URL or renderable.
func (f *Field) From(url string) *Field {
	f.SelectTableURL = url
	return f
}

// DialogWidth sets the select table dialog width.
func (f *Field) DialogWidth(width string) *Field {
	f.SelectTableDialogWidth = width
	return f
}

// Pluck sets the display and value fields for select table.
func (f *Field) Pluck(displayField, valueField string) *Field {
	f.SelectTableDisplayField = displayField
	f.SelectTableValueField = valueField
	return f
}

// UseTable sets has-many to use table mode display.
func (f *Field) UseTable() *Field {
	f.HasManyTableMode = true
	return f
}

// RepeaterBuilder configures nested sub-fields.
type RepeaterBuilder struct {
	field *Field
}

// TableBuilder configures table form fields.
type TableBuilder struct {
	field *Field
}

// NestedFormBuilder configures nested form fields for HasMany and Embeds.
type NestedFormBuilder struct {
	field *Field
}

func (b *RepeaterBuilder) add(name, label string, typ FieldType) *Field {
	child := &Field{Name: name, Label: label, Type: typ}
	b.field.RepeaterFields = append(b.field.RepeaterFields, child)
	return child
}

func (b *RepeaterBuilder) Hidden(name string) *Field {
	return b.add(name, "", FieldHidden)
}

func (b *RepeaterBuilder) Text(name, label string) *Field {
	return b.add(name, label, FieldText)
}

func (b *RepeaterBuilder) Number(name, label string) *Field {
	return b.add(name, label, FieldNumber)
}

func (b *RepeaterBuilder) Textarea(name, label string) *Field {
	return b.add(name, label, FieldTextarea)
}

func (b *RepeaterBuilder) Editor(name, label string) *Field {
	return b.add(name, label, FieldEditor)
}

func (b *RepeaterBuilder) Date(name, label string) *Field {
	return b.add(name, label, FieldDate)
}

func (b *RepeaterBuilder) Time(name, label string) *Field {
	return b.add(name, label, FieldTime)
}

func (b *RepeaterBuilder) Datetime(name, label string) *Field {
	return b.add(name, label, FieldDatetime)
}

func (b *RepeaterBuilder) Select(name, label string, options ...Option) *Field {
	child := b.add(name, label, FieldSelect)
	child.Options = append(child.Options, options...)
	return child
}

func (b *RepeaterBuilder) Radio(name, label string, options ...Option) *Field {
	child := b.add(name, label, FieldRadio)
	child.Options = append(child.Options, options...)
	return child
}

func (b *RepeaterBuilder) Checkbox(name, label string, options ...Option) *Field {
	child := b.add(name, label, FieldCheckbox)
	child.Options = append(child.Options, options...)
	return child
}

func (b *RepeaterBuilder) Switch(name, label string) *Field {
	return b.add(name, label, FieldSwitch)
}

func (b *RepeaterBuilder) Tags(name, label string) *Field {
	return b.add(name, label, FieldTags)
}

func (b *RepeaterBuilder) Color(name, label string) *Field {
	return b.add(name, label, FieldColor)
}

func (b *RepeaterBuilder) Icon(name, label string) *Field {
	return b.add(name, label, FieldIcon)
}

func (b *RepeaterBuilder) Range(name, label string) *Field {
	return b.add(name, label, FieldRange)
}

func (b *RepeaterBuilder) Rate(name, label string) *Field {
	return b.add(name, label, FieldRate)
}

func (b *RepeaterBuilder) Currency(name, label string) *Field {
	return b.add(name, label, FieldCurrency)
}

func (b *RepeaterBuilder) IP(name, label string) *Field {
	return b.add(name, label, FieldIP)
}

func (b *RepeaterBuilder) Mobile(name, label string) *Field {
	return b.add(name, label, FieldMobile)
}

func (b *RepeaterBuilder) Markdown(name, label string) *Field {
	return b.add(name, label, FieldMarkdown)
}

func (b *RepeaterBuilder) Slider(name, label string) *Field {
	return b.add(name, label, FieldSlider)
}

func (b *RepeaterBuilder) Autocomplete(name, label string) *Field {
	return b.add(name, label, FieldAutocomplete)
}

func (b *RepeaterBuilder) Listbox(name, label string, options ...Option) *Field {
	child := b.add(name, label, FieldListbox)
	child.Options = append(child.Options, options...)
	return child
}

// TableBuilder methods
func (b *TableBuilder) Text(name, label string) *Field {
	child := &Field{Name: name, Label: label, Type: FieldText}
	b.field.NestedFields = append(b.field.NestedFields, child)
	return child
}

func (b *TableBuilder) Number(name, label string) *Field {
	child := &Field{Name: name, Label: label, Type: FieldNumber}
	b.field.NestedFields = append(b.field.NestedFields, child)
	return child
}

func (b *TableBuilder) Select(name, label string, options ...Option) *Field {
	child := &Field{Name: name, Label: label, Type: FieldSelect}
	child.Options = append(child.Options, options...)
	b.field.NestedFields = append(b.field.NestedFields, child)
	return child
}

func (b *TableBuilder) Date(name, label string) *Field {
	child := &Field{Name: name, Label: label, Type: FieldDate}
	b.field.NestedFields = append(b.field.NestedFields, child)
	return child
}

// NestedFormBuilder methods
func (b *NestedFormBuilder) Text(name, label string) *Field {
	child := &Field{Name: name, Label: label, Type: FieldText}
	b.field.NestedFields = append(b.field.NestedFields, child)
	return child
}

func (b *NestedFormBuilder) Number(name, label string) *Field {
	child := &Field{Name: name, Label: label, Type: FieldNumber}
	b.field.NestedFields = append(b.field.NestedFields, child)
	return child
}

func (b *NestedFormBuilder) Textarea(name, label string) *Field {
	child := &Field{Name: name, Label: label, Type: FieldTextarea}
	b.field.NestedFields = append(b.field.NestedFields, child)
	return child
}

func (b *NestedFormBuilder) Select(name, label string, options ...Option) *Field {
	child := &Field{Name: name, Label: label, Type: FieldSelect}
	child.Options = append(child.Options, options...)
	b.field.NestedFields = append(b.field.NestedFields, child)
	return child
}

func (b *NestedFormBuilder) Date(name, label string) *Field {
	child := &Field{Name: name, Label: label, Type: FieldDate}
	b.field.NestedFields = append(b.field.NestedFields, child)
	return child
}

func (b *NestedFormBuilder) Datetime(name, label string) *Field {
	child := &Field{Name: name, Label: label, Type: FieldDatetime}
	b.field.NestedFields = append(b.field.NestedFields, child)
	return child
}

func (b *NestedFormBuilder) Switch(name, label string) *Field {
	child := &Field{Name: name, Label: label, Type: FieldSwitch}
	b.field.NestedFields = append(b.field.NestedFields, child)
	return child
}

func (b *NestedFormBuilder) Image(name, label string) *Field {
	child := &Field{Name: name, Label: label, Type: FieldImage}
	b.field.NestedFields = append(b.field.NestedFields, child)
	return child
}

func (b *NestedFormBuilder) Hidden(name string) *Field {
	child := &Field{Name: name, Label: "", Type: FieldHidden}
	b.field.NestedFields = append(b.field.NestedFields, child)
	return child
}

func (b *RepeaterBuilder) MinRows(rows int) *RepeaterBuilder {
	if rows > 0 {
		b.field.RepeaterMinRows = rows
	}
	return b
}

func (b *RepeaterBuilder) ValueFrom(path string) *RepeaterBuilder {
	b.field.ValuePath = path
	return b
}

func (b *RepeaterBuilder) WithHelp(help string) *RepeaterBuilder {
	b.field.Help = help
	return b
}

func (b *RepeaterBuilder) MarkRequired() *RepeaterBuilder {
	b.field.Required = true
	return b
}
