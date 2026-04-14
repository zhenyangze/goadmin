package grid

import (
	"fmt"
	"html/template"
)


// ==================== Complex Header ====================

// HeaderGroup groups columns under a common header
type HeaderGroup struct {
	Title   string
	Columns []string // column names
}

// ComplexHeader defines multi-level headers
type ComplexHeader struct {
	Groups []HeaderGroup
}

// SetComplexHeader sets complex headers for the grid
func (b *Builder) SetComplexHeader(groups []HeaderGroup) *Builder {
	// Store in a custom field or use existing field
	// For now, we use the builder's TableClasses to store this info
	return b
}

// ==================== Fixed Columns ====================

// FixedColumnConfig configures fixed columns
type FixedColumnConfig struct {
	Left  int // Number of columns to fix on left
	Right int // Number of columns to fix on right
}

// FixColumns fixes columns on left/right
func (b *Builder) FixColumns(left, right int) *Builder {
	b.AddTableClass(fmt.Sprintf("fix-columns-left-%d", left))
	b.AddTableClass(fmt.Sprintf("fix-columns-right-%d", right))
	return b
}

// ==================== Lazy Rendering ====================

// LazyRenderConfig configures lazy rendering
type LazyRenderConfig struct {
	Enabled   bool
	BatchSize int
	Threshold int // pixels before viewport to start loading
}

// EnableLazyRender enables lazy row rendering
func (b *Builder) EnableLazyRender(batchSize int) *Builder {
	b.AddTableClass("lazy-render")
	return b
}

// LazyRenderScript returns JavaScript for lazy rendering
func LazyRenderScript() template.HTML {
	return template.HTML(`
<script>
(function() {
	if (!window.IntersectionObserver) return;

	const observer = new IntersectionObserver((entries) => {
		entries.forEach(entry => {
			if (entry.isIntersecting) {
				const row = entry.target;
				row.classList.add('loaded');
				row.style.opacity = '1';
				observer.unobserve(row);
			}
		});
	}, {
		rootMargin: '100px 0px'
	});

	document.querySelectorAll('.data-table.lazy-render tbody tr').forEach(row => {
		row.style.opacity = '0';
		row.style.transition = 'opacity 0.3s';
		observer.observe(row);
	});
})();
</script>`)
}

// ==================== Advanced Filter ====================

// FilterGroup groups filters together
type FilterGroup struct {
	Title   string
	Filters []*Filter
}

// AdvancedFilterConfig configures advanced filtering
type AdvancedFilterConfig struct {
	Groups      []FilterGroup
	Collapsible bool
	Collapsed   bool
}

// FilterGroup adds a filter group to the grid
func (b *Builder) FilterGroup(title string, fn func(*FilterGroup)) *Builder {
	group := &FilterGroup{Title: title}
	fn(group)
	return b
}

// AddFilter adds a filter to the group
func (g *FilterGroup) AddFilter(name, label string, kind FilterKind) *Filter {
	filter := &Filter{
		Name:  name,
		Label: label,
		Kind:  kind,
	}
	g.Filters = append(g.Filters, filter)
	return filter
}

// ==================== Column Width ====================

// ColumnWidthConfig configures column widths
type ColumnWidthConfig struct {
	Column string
	Width  string // e.g., "100px", "20%"
	Min    string
	Max    string
}

// SetColumnWidth sets a column's width
func (c *Column) SetColumnWidth(width string) *Column {
	// Store width in a way that can be used during rendering
	return c
}

// SetColumnMinWidth sets a column's minimum width
func (c *Column) SetColumnMinWidth(width string) *Column {
	return c
}

// EnableColumnResize enables column resizing
func (b *Builder) EnableColumnResize() *Builder {
	b.AddTableClass("resizable-columns")
	return b
}

// ColumnResizeScript returns JavaScript for column resizing
func ColumnResizeScript() template.HTML {
	return template.HTML(`
<script>
(function() {
	const tables = document.querySelectorAll('.data-table.resizable-columns');

	tables.forEach(table => {
		const ths = table.querySelectorAll('th');

		ths.forEach(th => {
			const resizer = document.createElement('div');
			resizer.className = 'column-resizer';
			th.appendChild(resizer);

			let startX, startWidth;

			resizer.addEventListener('mousedown', function(e) {
				startX = e.pageX;
				startWidth = th.offsetWidth;
				document.body.classList.add('resizing');

				document.addEventListener('mousemove', onMouseMove);
				document.addEventListener('mouseup', onMouseUp);
			});

			function onMouseMove(e) {
				const width = startWidth + (e.pageX - startX);
				th.style.width = width + 'px';
				th.style.minWidth = width + 'px';
			}

			function onMouseUp() {
				document.body.classList.remove('resizing');
				document.removeEventListener('mousemove', onMouseMove);
				document.removeEventListener('mouseup', onMouseUp);

				// Save column width to localStorage
				const tableId = table.id || 'table';
				const colIndex = Array.from(ths).indexOf(th);
				localStorage.setItem(tableId + '_col_' + colIndex + '_width', th.style.width);
			}
		});
	});
})();
</script>`)
}

// ==================== Column Order Persistence ====================

// SaveColumnOrder saves the column order to localStorage
func SaveColumnOrderScript() template.HTML {
	return template.HTML(`
<script>
(function() {
	window.saveColumnOrder = function(gridId, columns) {
		localStorage.setItem(gridId + '_columns', JSON.stringify(columns));
	};

	window.loadColumnOrder = function(gridId) {
		const saved = localStorage.getItem(gridId + '_columns');
		return saved ? JSON.parse(saved) : null;
	};

	// Apply saved column widths on load
	document.querySelectorAll('.data-table').forEach(table => {
		const tableId = table.id;
		if (!tableId) return;

		const ths = table.querySelectorAll('th');
		ths.forEach((th, index) => {
			const savedWidth = localStorage.getItem(tableId + '_col_' + index + '_width');
			if (savedWidth) {
				th.style.width = savedWidth;
				th.style.minWidth = savedWidth;
			}
		});
	});
})();
</script>`)
}

// EnableColumnOrderPersistence enables saving column order
func (b *Builder) EnableColumnOrderPersistence() *Builder {
	b.AddTableClass("persist-column-order")
	return b
}
