package theme

// Theme contains the visual tokens used by the shared templates.
type Theme struct {
	Accent     string
	Sidebar    string
	SidebarInk string
	Body       string
	Surface    string
}

// Default returns a warm neutral Tailwind-friendly palette.
func Default() Theme {
	return Theme{
		Accent:     "#b45309",
		Sidebar:    "#111827",
		SidebarInk: "#f9fafb",
		Body:       "#f7f4ee",
		Surface:    "#ffffff",
	}
}
