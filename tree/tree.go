package tree

// Builder defines the tree page rendering preferences.
type Builder struct {
	Title            string
	Description      string
	EmptyText        string
	TitleField       string
	DescriptionField string
}

// New creates a tree builder.
func New() *Builder {
	return &Builder{
		EmptyText:        "No tree data yet.",
		TitleField:       "Title",
		DescriptionField: "Description",
	}
}

// Branch defines which fields should be exposed in the tree cards.
func (b *Builder) Branch(titleField, descriptionField string) {
	if titleField != "" {
		b.TitleField = titleField
	}
	if descriptionField != "" {
		b.DescriptionField = descriptionField
	}
}
