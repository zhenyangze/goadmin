package grid

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
)

// ==================== Quick Edit Action ====================

// QuickEditConfig configures inline quick edit
type QuickEditConfig struct {
	FieldName   string
	FieldType   string // text, select, textarea, number
	Options     []Option // for select type
	SaveURL     string
	SaveMethod  string
	Placeholder string
}

// QuickEditAction creates a quick edit row action
func QuickEditAction(fieldName string) *RowAction {
	return &RowAction{
		Label: "Quick Edit",
		URL: func(record any) string {
			return fmt.Sprintf("javascript:openQuickEdit('%s', %s)", fieldName, toJSON(record))
		},
		Style:  ActionGhost,
		Method: "GET",
	}
}

// QuickEditButton adds a quick edit button to grid
func (b *Builder) QuickEditButton(fieldName string, config *QuickEditConfig) *Builder {
	// Add the quick edit action
	b.RowAction("Quick Edit", func(record any) string {
		return fmt.Sprintf("#quick-edit-%s-%v", fieldName, getRecordID(record))
	}).WithMethod("POST")

	return b
}

// JavaScript for quick edit
func QuickEditScript() template.HTML {
	return template.HTML(`
<script>
(function() {
	window.openQuickEdit = function(fieldName, record) {
		const cell = document.querySelector('[data-field="' + fieldName + '][data-id="' + record.id + '"]');
		if (!cell) return;

		const originalValue = cell.textContent.trim();
		const input = document.createElement('input');
		input.type = 'text';
		input.value = originalValue;
		input.className = 'quick-edit-input';

		cell.innerHTML = '';
		cell.appendChild(input);
		input.focus();
		input.select();

		const save = function() {
			const newValue = input.value;
			// Save via AJAX
			fetch(cell.dataset.saveUrl || window.location.href + '/quick-edit', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'X-Requested-With': 'XMLHttpRequest'
				},
				body: JSON.stringify({
					id: record.id,
					field: fieldName,
					value: newValue
				})
			})
			.then(r => r.json())
			.then(result => {
				if (result.success) {
					cell.textContent = newValue;
					cell.classList.add('quick-edit-saved');
					setTimeout(() => cell.classList.remove('quick-edit-saved'), 1000);
				} else {
					alert(result.message || 'Save failed');
					cell.textContent = originalValue;
				}
			})
			.catch(() => {
				cell.textContent = originalValue;
			});
		};

		const cancel = function() {
			cell.textContent = originalValue;
		};

		input.addEventListener('blur', save);
		input.addEventListener('keydown', function(e) {
			if (e.key === 'Enter') {
				save();
			} else if (e.key === 'Escape') {
				cancel();
			}
		});
	};
})();
</script>`)
}

// ==================== Context Menu Actions ====================

// ContextMenuAction represents an action in context menu
type ContextMenuAction struct {
	Label   string
	Icon    string
	URL     RowActionURL
	Style   ActionStyle
	Method  string
	Confirm string
	Divider bool // Add divider after this action
}

// ContextMenuActions groups actions in a context menu
type ContextMenuActions struct {
	label   string
	actions []ContextMenuAction
}

// ContextMenu creates a new context menu
func ContextMenu(label string) *ContextMenuActions {
	return &ContextMenuActions{
		label:   label,
		actions: []ContextMenuAction{},
	}
}

// Action adds an action to the context menu
func (c *ContextMenuActions) Action(label string, url RowActionURL) *ContextMenuActions {
	c.actions = append(c.actions, ContextMenuAction{
		Label:  label,
		URL:    url,
		Style:  ActionDefault,
		Method: "GET",
	})
	return c
}

// ActionWithStyle adds an action with custom style
func (c *ContextMenuActions) ActionWithStyle(label string, url RowActionURL, style ActionStyle) *ContextMenuActions {
	c.actions = append(c.actions, ContextMenuAction{
		Label:  label,
		URL:    url,
		Style:  style,
		Method: "GET",
	})
	return c
}

// ActionWithConfirm adds an action with confirmation
func (c *ContextMenuActions) ActionWithConfirm(label string, url RowActionURL, confirm string) *ContextMenuActions {
	c.actions = append(c.actions, ContextMenuAction{
		Label:   label,
		URL:     url,
		Style:   ActionDanger,
		Method:  "POST",
		Confirm: confirm,
	})
	return c
}

// Divider adds a divider
func (c *ContextMenuActions) Divider() *ContextMenuActions {
	if len(c.actions) > 0 {
		c.actions[len(c.actions)-1].Divider = true
	}
	return c
}

// Render generates the context menu HTML
func (c *ContextMenuActions) Render(record any) template.HTML {
	if len(c.actions) == 0 {
		return ""
	}

	var itemsHTML string
	for _, action := range c.actions {
		url := action.URL(record)
		confirm := ""
		if action.Confirm != "" {
			confirm = fmt.Sprintf(` onclick="return confirm('%s')"`, template.JSEscapeString(action.Confirm))
		}

		class := "context-menu-item"
		if action.Style == ActionDanger {
			class += " context-menu-item-danger"
		}

		iconHTML := ""
		if action.Icon != "" {
			iconHTML = fmt.Sprintf(`<span class="context-menu-icon">%s</span>`, action.Icon)
		}

		itemsHTML += fmt.Sprintf(
			`<a href="%s" class="%s"%s>%s%s</a>`,
			url, class, confirm, iconHTML, action.Label)

		if action.Divider {
			itemsHTML += `<div class="context-menu-divider"></div>`
		}
	}

	return template.HTML(fmt.Sprintf(`
<div class="context-menu-wrapper">
	<button type="button" class="context-menu-trigger" onclick="toggleContextMenu(this)">
		<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
			<circle cx="12" cy="5" r="1"/>
			<circle cx="12" cy="12" r="1"/>
			<circle cx="12" cy="19" r="1"/>
		</svg>
	</button>
	<div class="context-menu-dropdown" style="display: none;">
		%s
	</div>
</div>`, itemsHTML))
}

// ContextMenuScript returns the JavaScript for context menus
func ContextMenuScript() template.HTML {
	return template.HTML(`
<script>
(function() {
	window.toggleContextMenu = function(trigger) {
		const dropdown = trigger.nextElementSibling;
		const isOpen = dropdown.style.display !== 'none';

		// Close all other menus
		document.querySelectorAll('.context-menu-dropdown').forEach(d => {
			d.style.display = 'none';
		});

		if (!isOpen) {
			dropdown.style.display = 'block';
			// Position dropdown
			const rect = trigger.getBoundingClientRect();
			dropdown.style.top = rect.bottom + 'px';
			dropdown.style.right = (window.innerWidth - rect.right) + 'px';
		}
	};

	// Close on click outside
	document.addEventListener('click', function(e) {
		if (!e.target.closest('.context-menu-wrapper')) {
			document.querySelectorAll('.context-menu-dropdown').forEach(d => {
				d.style.display = 'none';
			});
		}
	});
})();
</script>`)
}

// ==================== Batch Action Confirmation ====================

// BatchConfirmConfig configures batch action confirmation dialog
type BatchConfirmConfig struct {
	Title       string
	Message     string
	ConfirmText string
	CancelText  string
	Dangerous   bool
}

// DefaultBatchConfirm provides default confirmation config
func DefaultBatchConfirm() *BatchConfirmConfig {
	return &BatchConfirmConfig{
		Title:       "Confirm Action",
		Message:     "Are you sure you want to perform this action on the selected items?",
		ConfirmText: "Confirm",
		CancelText:  "Cancel",
		Dangerous:   false,
	}
}

// WithConfirmation enables confirmation for batch action
func (b *BatchAction) WithConfirmation(config *BatchConfirmConfig) *BatchAction {
	if config == nil {
		config = DefaultBatchConfirm()
	}
	b.Confirm = config.Message
	return b
}

// WithConfirmation enables confirmation for batch action with handler
func (b *BatchActionWithHandler) WithConfirmation(config *BatchConfirmConfig) *BatchActionWithHandler {
	if config == nil {
		config = DefaultBatchConfirm()
	}
	b.Confirm = config.Message
	return b
}

// BatchConfirmDialogScript returns JavaScript for batch confirmation
func BatchConfirmDialogScript() template.HTML {
	return template.HTML(`
<script>
(function() {
	window.showBatchConfirm = function(title, message, onConfirm, onCancel) {
		const modal = document.createElement('div');
		modal.className = 'modal-wrapper is-open';
		modal.innerHTML = '<div class="modal-overlay" onclick="this.parentElement.remove()"></div>' +
			'<div class="modal-dialog" style="width: 400px;">' +
				'<div class="modal-header">' +
					'<h4 class="modal-title">' + title + '</h4>' +
					'<button class="modal-close" onclick="this.closest(\'.modal-wrapper\').remove()">' +
						'<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">' +
							'<line x1="18" y1="6" x2="6" y2="18"/>' +
							'<line x1="6" y1="6" x2="18" y2="18"/>' +
						'</svg>' +
					'</button>' +
				'</div>' +
				'<div class="modal-body">' +
					'<p>' + message + '</p>' +
					'<p class="batch-selected-count"></p>' +
				'</div>' +
				'<div class="modal-footer">' +
					'<button type="button" class="btn btn-primary" id="batch-confirm-btn">Confirm</button>' +
					'<button type="button" class="btn btn-ghost" onclick="this.closest(\'.modal-wrapper\').remove()">Cancel</button>' +
				'</div>' +
			'</div>';

		document.body.appendChild(modal);

		modal.querySelector('#batch-confirm-btn').addEventListener('click', function() {
			if (onConfirm) onConfirm();
			modal.remove();
		});
	};
})();
</script>`)
}

// ==================== Action Permission Control ====================

// PermissionChecker checks if user has permission for action
type PermissionChecker func(ctx context.Context, action string, record any) bool

// ActionPermission controls action visibility based on permissions
type ActionPermission struct {
	checker PermissionChecker
	actions map[string]string // action name -> permission key
}

// NewActionPermission creates a new permission controller
func NewActionPermission(checker PermissionChecker) *ActionPermission {
	return &ActionPermission{
		checker: checker,
		actions: make(map[string]string),
	}
}

// RegisterAction registers an action with permission key
func (ap *ActionPermission) RegisterAction(actionName, permissionKey string) *ActionPermission {
	ap.actions[actionName] = permissionKey
	return ap
}

// Can checks if user can perform action
func (ap *ActionPermission) Can(ctx context.Context, actionName string, record any) bool {
	permissionKey, ok := ap.actions[actionName]
	if !ok {
		return true // No permission required
	}
	if ap.checker == nil {
		return true
	}
	return ap.checker(ctx, permissionKey, record)
}

// FilterRowActions filters row actions based on permissions
func (ap *ActionPermission) FilterRowActions(ctx context.Context, actions []*RowAction, record any) []*RowAction {
	if ap.checker == nil {
		return actions
	}

	var filtered []*RowAction
	for _, action := range actions {
		permissionKey := ap.actions[action.Label]
		if permissionKey == "" || ap.checker(ctx, permissionKey, record) {
			filtered = append(filtered, action)
		}
	}
	return filtered
}

// FilterPageActions filters page actions based on permissions
func (ap *ActionPermission) FilterPageActions(ctx context.Context, actions []*PageAction) []*PageAction {
	if ap.checker == nil {
		return actions
	}

	var filtered []*PageAction
	for _, action := range actions {
		permissionKey := ap.actions[action.Label]
		if permissionKey == "" || ap.checker(ctx, permissionKey, nil) {
			filtered = append(filtered, action)
		}
	}
	return filtered
}

// FilterBatchActions filters batch actions based on permissions
func (ap *ActionPermission) FilterBatchActions(ctx context.Context, actions []*BatchAction) []*BatchAction {
	if ap.checker == nil {
		return actions
	}

	var filtered []*BatchAction
	for _, action := range actions {
		permissionKey := ap.actions[action.Label]
		if permissionKey == "" || ap.checker(ctx, permissionKey, nil) {
			filtered = append(filtered, action)
		}
	}
	return filtered
}

// Helper functions
func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func getRecordID(record interface{}) string {
	// Try to extract ID from record
	if r, ok := record.(map[string]interface{}); ok {
		if id, ok := r["id"]; ok {
			return fmt.Sprintf("%v", id)
		}
	}
	return "0"
}
