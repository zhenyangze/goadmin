package goadmin

import (
	"html/template"
	"net/url"
)

// DialogConfig configures a dialog/modal
type DialogConfig struct {
	ID       string
	Title    string
	Width    string
	Height   string
	Content  template.HTML
	Footer   template.HTML
	ShowClose bool
}

// DialogManager manages dialog instances
type DialogManager struct {
	dialogs []DialogConfig
}

// NewDialogManager creates a new dialog manager
func NewDialogManager() *DialogManager {
	return &DialogManager{
		dialogs: []DialogConfig{},
	}
}

// AddDialog adds a dialog configuration
func (dm *DialogManager) AddDialog(config DialogConfig) string {
	if config.ID == "" {
		config.ID = generateDialogID()
	}
	if config.Width == "" {
		config.Width = "600px"
	}
	if config.Height == "" {
		config.Height = "auto"
	}
	if config.ShowClose == false {
		config.ShowClose = true // default to true
	}
	dm.dialogs = append(dm.dialogs, config)
	return config.ID
}

// GetDialogs returns all dialog configurations
func (dm *DialogManager) GetDialogs() []DialogConfig {
	return dm.dialogs
}

// RenderDialog generates the HTML for a dialog
func RenderDialog(config DialogConfig) template.HTML {
	closeBtn := ""
	if config.ShowClose {
		closeBtn = `
			<button class="modal-close" onclick="closeDialog('` + config.ID + `')" aria-label="Close">
				<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<line x1="18" y1="6" x2="6" y2="18"/>
					<line x1="6" y1="6" x2="18" y2="18"/>
				</svg>
			</button>`
	}

	footer := ""
	if config.Footer != "" {
		footer = `<div class="modal-footer">` + string(config.Footer) + `</div>`
	}

	height := config.Height
	if height != "auto" {
		height = config.Height
	}

	return template.HTML(`
<div id="` + config.ID + `" class="modal-wrapper" style="display: none;">
	<div class="modal-overlay" onclick="closeDialog('` + config.ID + `')"></div>
	<div class="modal-dialog" style="width: ` + config.Width + `; height: ` + height + `;">
		<div class="modal-header">
			<h4 class="modal-title">` + template.HTMLEscapeString(config.Title) + `</h4>
			` + closeBtn + `
		</div>
		<div class="modal-body">
			` + string(config.Content) + `
		</div>
		` + footer + `
	</div>
</div>`)
}

// DialogTrigger generates a button/link to open a dialog
func DialogTrigger(dialogID string, label string, style string) template.HTML {
	class := "btn btn-primary"
	switch style {
	case "secondary":
		class = "btn btn-secondary"
	case "danger":
		class = "btn btn-danger"
	case "ghost":
		class = "btn btn-ghost"
	case "link":
		class = "text-blue-600 hover:underline"
	}

	if style == "link" {
		return template.HTML(`<a href="#" onclick="openDialog('` + dialogID + `'); return false;" class="` + class + `">` + label + `</a>`)
	}

	return template.HTML(`<button type="button" onclick="openDialog('` + dialogID + `')" class="` + class + `">` + label + `</button>`)
}

// generateDialogID generates a unique dialog ID
var dialogCounter int

func generateDialogID() string {
	dialogCounter++
	return "dialog_" + string(rune(dialogCounter))
}

// ==================== Dialog Table Component ====================

// DialogTableConfig configures a dialog table
type DialogTableConfig struct {
	ID             string
	Title          string
	URL            string
	Width          string
	Height         string
	Columns        []DialogTableColumn
	MultiSelect    bool
	SelectCallback string // JavaScript function name to call on select
}

// DialogTableColumn defines a column in the dialog table
type DialogTableColumn struct {
	Name  string
	Label string
}

// DialogTableSelection represents a selected item
type DialogTableSelection struct {
	ID          string            `json:"id"`
	DisplayText string            `json:"display_text"`
	Data        map[string]string `json:"data"`
}

// RenderDialogTable generates HTML for a dialog table
func RenderDialogTable(config DialogTableConfig) template.HTML {
	if config.ID == "" {
		config.ID = "dialog_table_" + string(rune(dialogCounter))
		dialogCounter++
	}
	if config.Width == "" {
		config.Width = "800px"
	}
	if config.Height == "" {
		config.Height = "600px"
	}

	var columnsHTML string
	for _, col := range config.Columns {
		columnsHTML += `<th>` + template.HTMLEscapeString(col.Label) + `</th>`
	}

	multiSelectAttr := ""
	if config.MultiSelect {
		multiSelectAttr = " multiple"
	}

	content := template.HTML(`
<div class="dialog-table-wrapper" data-url="` + template.HTMLEscapeString(config.URL) + `"` + multiSelectAttr + `>
	<div class="dialog-table-search mb-4">
		<input type="text" class="form-input" placeholder="Search..." onkeyup="searchDialogTable(this)">
	</div>
	<div class="dialog-table-container overflow-auto" style="max-height: 400px;">
		<table class="data-table">
			<thead>
				<tr>` + columnsHTML + `</tr>
			</thead>
			<tbody class="dialog-table-body">
				<!-- Data loaded via AJAX -->
			</tbody>
		</table>
	</div>
	<div class="dialog-table-pagination mt-4 flex justify-between items-center">
		<span class="text-sm text-gray-600">Loading...</span>
		<div class="pagination-buttons"></div>
	</div>
</div>`)

	footer := template.HTML(`
<button type="button" class="btn btn-primary" onclick="confirmDialogTableSelection('` + config.ID + `')">Select</button>
<button type="button" class="btn btn-ghost" onclick="closeDialog('` + config.ID + `')">Cancel</button>`)

	dialogConfig := DialogConfig{
		ID:        config.ID,
		Title:     config.Title,
		Width:     config.Width,
		Height:    config.Height,
		Content:   content,
		Footer:    footer,
		ShowClose: true,
	}

	return RenderDialog(dialogConfig)
}

// DialogTableTrigger creates a button to open the dialog table
func DialogTableTrigger(dialogID string, label string) template.HTML {
	return template.HTML(`
<button type="button" onclick="loadDialogTableData('` + dialogID + `'); openDialog('` + dialogID + `')" class="btn btn-secondary">
	<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="mr-2">
		<rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
		<line x1="3" y1="9" x2="21" y2="9"/>
		<line x1="9" y1="21" x2="9" y2="9"/>
	</svg>
	` + label + `
</button>`)
}

// ==================== Dialog JavaScript ====================

// DialogJavaScript returns the JavaScript needed for dialogs
func DialogJavaScript() template.HTML {
	return template.HTML(`
<script>
// Global dialog functions
window.openDialog = function(dialogId) {
	const dialog = document.getElementById(dialogId);
	if (dialog) {
		dialog.style.display = 'block';
		document.body.style.overflow = 'hidden';

		// Dispatch custom event
		dialog.dispatchEvent(new CustomEvent('dialog:open', { bubbles: true }));
	}
};

window.closeDialog = function(dialogId) {
	const dialog = typeof dialogId === 'string' ? document.getElementById(dialogId) : dialogId;
	if (dialog) {
		dialog.style.display = 'none';
		document.body.style.overflow = '';

		// Dispatch custom event
		dialog.dispatchEvent(new CustomEvent('dialog:close', { bubbles: true }));
	}
};

// Close dialog on Escape key
document.addEventListener('keydown', function(e) {
	if (e.key === 'Escape') {
		const openDialogs = document.querySelectorAll('.modal-wrapper[style*="block"]');
		openDialogs.forEach(dialog => closeDialog(dialog.id));
	}
});

// Dialog Table Functions
window.loadDialogTableData = function(dialogId) {
	const dialog = document.getElementById(dialogId);
	if (!dialog) return;

	const wrapper = dialog.querySelector('.dialog-table-wrapper');
	const url = wrapper.dataset.url;
	const tbody = wrapper.querySelector('.dialog-table-body');

	fetch(url)
		.then(response => response.json())
		.then(data => {
			renderDialogTableRows(tbody, data.rows || []);
			updateDialogTablePagination(wrapper, data.pagination || {});
		})
		.catch(err => {
			tbody.innerHTML = '<tr><td colspan="100" class="text-center text-red-500">Failed to load data</td></tr>';
		});
};

window.renderDialogTableRows = function(tbody, rows) {
	if (!rows || rows.length === 0) {
		tbody.innerHTML = '<tr><td colspan="100" class="text-center text-gray-500">No data found</td></tr>';
		return;
	}

	const isMultiSelect = tbody.closest('.dialog-table-wrapper').hasAttribute('multiple');
	const inputType = isMultiSelect ? 'checkbox' : 'radio';
	const inputName = isMultiSelect ? 'dialog_table_selection[]' : 'dialog_table_selection';

	tbody.innerHTML = rows.map(row =>
		'<tr onclick="selectDialogTableRow(this, event)">' +
		'<td><input type="' + inputType + '" name="' + inputName + '" value="' + (row.id || '') + '" data-display="' + (row.display_text || '') + '"></td>' +
		Object.keys(row).filter(k => k !== 'id' && k !== 'display_text').map(key =>
			'<td>' + (row[key] || '') + '</td>'
		).join('') +
		'</tr>'
	).join('');
};

window.selectDialogTableRow = function(row, event) {
	const input = row.querySelector('input[type="checkbox"], input[type="radio"]');
	if (!input) return;

	// Don't toggle if clicking directly on the input
	if (event.target === input) return;

	input.checked = !input.checked;

	if (input.type === 'radio') {
		// Deselect other rows for radio buttons
		row.parentElement.querySelectorAll('tr').forEach(r => r.classList.remove('selected'));
	}

	if (input.checked) {
		row.classList.add('selected');
	} else {
		row.classList.remove('selected');
	}
};

window.searchDialogTable = function(input) {
	const term = input.value.toLowerCase();
	const tbody = input.closest('.dialog-table-wrapper').querySelector('.dialog-table-body');
	const rows = tbody.querySelectorAll('tr');

	rows.forEach(row => {
		const text = row.textContent.toLowerCase();
		row.style.display = text.includes(term) ? '' : 'none';
	});
};

window.confirmDialogTableSelection = function(dialogId) {
	const dialog = document.getElementById(dialogId);
	const selected = dialog.querySelectorAll('input[name^="dialog_table_selection"]:checked');

	if (selected.length === 0) {
		alert('Please select at least one item');
		return;
	}

	const selections = Array.from(selected).map(input => ({
		id: input.value,
		display_text: input.dataset.display
	}));

	// Dispatch event with selection data
	dialog.dispatchEvent(new CustomEvent('dialogtable:select', {
		detail: { selections: selections },
		bubbles: true
	}));

	closeDialog(dialogId);
};
</script>`)
}

// ==================== Dialog Form Component ====================

// DialogFormConfig configures a dialog form
type DialogFormConfig struct {
	ID     string
	Title  string
	URL    string
	Width  string
	Height string
	Fields template.HTML
}

// RenderDialogForm generates HTML for a dialog form
func RenderDialogForm(config DialogFormConfig) template.HTML {
	if config.ID == "" {
		config.ID = "dialog_form_" + url.QueryEscape(config.Title)
	}
	if config.Width == "" {
		config.Width = "700px"
	}
	if config.Height == "" {
		config.Height = "auto"
	}

	content := template.HTML(`
<form id="` + config.ID + `_form" method="POST" action="` + template.HTMLEscapeString(config.URL) + `" data-dialog-form="true">` + `
	` + string(config.Fields) + `
</form>`)

	footer := template.HTML(`
<button type="submit" form="` + config.ID + `_form" class="btn btn-primary">Save</button>
<button type="button" class="btn btn-ghost" onclick="closeDialog('` + config.ID + `')">Cancel</button>`)

	dialogConfig := DialogConfig{
		ID:        config.ID,
		Title:     config.Title,
		Width:     config.Width,
		Height:    config.Height,
		Content:   content,
		Footer:    footer,
		ShowClose: true,
	}

	return RenderDialog(dialogConfig)
}

// DialogFormScript returns JavaScript for handling dialog form submissions
func DialogFormScript() template.HTML {
	return template.HTML(`
<script>
// Handle dialog form submissions
document.addEventListener('submit', function(e) {
	const form = e.target;
	if (!form.hasAttribute('data-dialog-form')) return;

	e.preventDefault();

	const formData = new FormData(form);
	const submitBtn = form.querySelector('button[type="submit"]');

	if (submitBtn) {
		submitBtn.disabled = true;
		submitBtn.textContent = 'Saving...';
	}

	fetch(form.action, {
		method: 'POST',
		body: formData,
		headers: {
			'X-Requested-With': 'XMLHttpRequest'
		}
	})
	.then(response => response.json())
	.then(result => {
		if (result.success) {
			// Close the dialog
			const dialog = form.closest('.modal-wrapper');
			if (dialog) closeDialog(dialog.id);

			// Reload page or update grid
			if (result.reload) {
				location.reload();
			} else if (result.script) {
				eval(result.script);
			}

			// Show success message
			if (window.showSuccessMessage) {
				window.showSuccessMessage(result.message || 'Saved successfully');
			}
		} else {
			// Show validation errors
			if (result.errors) {
				Object.keys(result.errors).forEach(field => {
					const input = form.querySelector('[name="' + field + '"]');
					if (input) {
						input.classList.add('error');
						const errorEl = document.createElement('span');
						errorEl.className = 'error-message';
						errorEl.textContent = result.errors[field];
						input.parentNode.appendChild(errorEl);
					}
				});
			} else if (result.message) {
				alert(result.message);
			}
		}
	})
	.catch(err => {
		alert('An error occurred. Please try again.');
		console.error(err);
	})
	.finally(() => {
		if (submitBtn) {
			submitBtn.disabled = false;
			submitBtn.textContent = 'Save';
		}
	});
});
</script>`)
}
