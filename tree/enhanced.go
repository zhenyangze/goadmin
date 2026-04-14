package tree

import (
	"context"
	"fmt"
	"html/template"
)

// ActionStyle describes the visual style of a tree action
type ActionStyle string

const (
	ActionDefault ActionStyle = "default"
	ActionPrimary ActionStyle = "primary"
	ActionDanger  ActionStyle = "danger"
	ActionGhost   ActionStyle = "ghost"
)

// TreeAction represents an action on a tree node
type TreeAction struct {
	Label    string
	Icon     string
	URL      string
	Style    ActionStyle
	Confirm  string
	Ajax     bool
	Visible  func(node any) bool
}

// TreeTool represents a toolbar tool for the tree
type TreeTool struct {
	Label   string
	Icon    string
	Handler func(ctx context.Context) error
	URL     string
}

// TreeNode represents a node in the tree with enhanced features
type TreeNode struct {
	ID          string
	ParentID    string
	Label       string
	Description string
	Icon        string
	Expanded    bool
	Selected    bool
	Disabled    bool
	Order       int
	Data        map[string]interface{}
	Children    []*TreeNode
	Actions     []*TreeAction
}

// TreeRowAction represents a row action for tree items
type TreeRowAction struct {
	Label   string
	Icon    string
	URL     func(node *TreeNode) string
	Style   ActionStyle
	Confirm string
}

// DragDropConfig configures drag and drop functionality
type DragDropConfig struct {
	Enabled   bool
	SaveURL   string
	OnDragEnd string // JavaScript callback
}

// BatchActionConfig configures batch actions for tree
type BatchActionConfig struct {
	Label    string
	Icon     string
	URL      string
	Confirm  string
	Handler  func(ctx context.Context, ids []string) error
}

// EnhancedBuilder extends the tree builder with actions and tools
type EnhancedBuilder struct {
	*Builder
	Actions        []*TreeAction
	Tools          []*TreeTool
	RowActions     []*TreeRowAction
	DragDrop       *DragDropConfig
	BatchActions   []*BatchActionConfig
	Selectable     bool
	MultiSelect    bool
	ShowCheckbox   bool
	ShowIcon       bool
	ShowActions    bool
	OnSelect       string // JavaScript callback
}

// NewEnhanced creates an enhanced tree builder
func NewEnhanced() *EnhancedBuilder {
	return &EnhancedBuilder{
		Builder:     New(),
		Actions:     []*TreeAction{},
		Tools:       []*TreeTool{},
		RowActions:  []*TreeRowAction{},
		DragDrop:    &DragDropConfig{Enabled: false},
		BatchActions: []*BatchActionConfig{},
		Selectable:  false,
		MultiSelect: false,
		ShowCheckbox: false,
		ShowIcon:    true,
		ShowActions: true,
	}
}

// AddAction adds an action to tree nodes
func (b *EnhancedBuilder) AddAction(action *TreeAction) *EnhancedBuilder {
	b.Actions = append(b.Actions, action)
	return b
}

// AddTool adds a toolbar tool
func (b *EnhancedBuilder) AddTool(tool *TreeTool) *EnhancedBuilder {
	b.Tools = append(b.Tools, tool)
	return b
}

// AddRowAction adds a row action
func (b *EnhancedBuilder) AddRowAction(action *TreeRowAction) *EnhancedBuilder {
	b.RowActions = append(b.RowActions, action)
	return b
}

// EnableDragDrop enables drag and drop functionality
func (b *EnhancedBuilder) EnableDragDrop(saveURL string) *EnhancedBuilder {
	b.DragDrop.Enabled = true
	b.DragDrop.SaveURL = saveURL
	return b
}

// EnableSelection enables node selection
func (b *EnhancedBuilder) EnableSelection(multi bool) *EnhancedBuilder {
	b.Selectable = true
	b.MultiSelect = multi
	b.ShowCheckbox = true
	return b
}

// AddBatchAction adds a batch action
func (b *EnhancedBuilder) AddBatchAction(action *BatchActionConfig) *EnhancedBuilder {
	b.BatchActions = append(b.BatchActions, action)
	return b
}

// SetOnSelect sets the selection callback
func (b *EnhancedBuilder) SetOnSelect(callback string) *EnhancedBuilder {
	b.OnSelect = callback
	return b
}

// TreeAction creates a tree action
func Action(label string) *TreeAction {
	return &TreeAction{
		Label:   label,
		Style:   ActionDefault,
		Visible: func(node any) bool { return true },
	}
}

// WithIcon sets the action icon
func (a *TreeAction) WithIcon(icon string) *TreeAction {
	a.Icon = icon
	return a
}

// WithURL sets the action URL
func (a *TreeAction) WithURL(url string) *TreeAction {
	a.URL = url
	return a
}

// WithStyle sets the action style
func (a *TreeAction) WithStyle(style ActionStyle) *TreeAction {
	a.Style = style
	return a
}

// WithConfirm sets the confirmation message
func (a *TreeAction) WithConfirm(msg string) *TreeAction {
	a.Confirm = msg
	return a
}

// WithAjax makes the action use AJAX
func (a *TreeAction) WithAjax() *TreeAction {
	a.Ajax = true
	return a
}

// WithVisible sets the visibility condition
func (a *TreeAction) WithVisible(fn func(node any) bool) *TreeAction {
	a.Visible = fn
	return a
}

// TreeTool creates a tree tool
func Tool(label string) *TreeTool {
	return &TreeTool{Label: label}
}

// WithIcon sets the tool icon
func (t *TreeTool) WithIcon(icon string) *TreeTool {
	t.Icon = icon
	return t
}

// WithHandler sets the tool handler
func (t *TreeTool) WithHandler(fn func(ctx context.Context) error) *TreeTool {
	t.Handler = fn
	return t
}

// WithURL sets the tool URL
func (t *TreeTool) WithURL(url string) *TreeTool {
	t.URL = url
	return t
}

// TreeRowAction creates a tree row action
func RowAction(label string) *TreeRowAction {
	return &TreeRowAction{
		Label: label,
		Style: ActionDefault,
		URL:   func(node *TreeNode) string { return "#" },
	}
}

// SetIcon sets the row action icon
func (a *TreeRowAction) SetIcon(icon string) *TreeRowAction {
	a.Icon = icon
	return a
}

// SetURL sets the URL function
func (a *TreeRowAction) SetURL(fn func(node *TreeNode) string) *TreeRowAction {
	a.URL = fn
	return a
}

// SetStyle sets the row action style
func (a *TreeRowAction) SetStyle(style ActionStyle) *TreeRowAction {
	a.Style = style
	return a
}

// SetConfirm sets the confirmation message
func (a *TreeRowAction) SetConfirm(msg string) *TreeRowAction {
	a.Confirm = msg
	return a
}

// JavaScript returns the JavaScript for tree enhancements
func JavaScript() template.HTML {
	return template.HTML(`
<script>
(function() {
	// Tree drag and drop
	window.initTreeDragDrop = function(treeId, saveURL) {
		const tree = document.getElementById(treeId);
		if (!tree) return;

		tree.querySelectorAll('.tree-node').forEach(node => {
			node.setAttribute('draggable', 'true');

			node.addEventListener('dragstart', function(e) {
				e.dataTransfer.setData('text/plain', node.dataset.id);
				node.classList.add('dragging');
			});

			node.addEventListener('dragend', function() {
				node.classList.remove('dragging');
			});

			node.addEventListener('dragover', function(e) {
				e.preventDefault();
				node.classList.add('drag-over');
			});

			node.addEventListener('dragleave', function() {
				node.classList.remove('drag-over');
			});

			node.addEventListener('drop', function(e) {
				e.preventDefault();
				node.classList.remove('drag-over');

				const draggedId = e.dataTransfer.getData('text/plain');
				const targetId = node.dataset.id;

				if (draggedId && draggedId !== targetId) {
					// Send to server
					fetch(saveURL, {
						method: 'POST',
						headers: {
							'Content-Type': 'application/json',
							'X-Requested-With': 'XMLHttpRequest'
						},
						body: JSON.stringify({
							draggedId: draggedId,
							targetId: targetId
						})
					})
					.then(r => r.json())
					.then(result => {
						if (result.success) {
							location.reload();
						} else {
							alert(result.message || 'Move failed');
						}
					});
				}
			});
		});
	};

	// Tree batch operations
	window.getSelectedTreeNodes = function(treeId) {
		const tree = document.getElementById(treeId);
		if (!tree) return [];

		return Array.from(tree.querySelectorAll('.tree-node.selected'))
			.map(node => node.dataset.id);
	};

	// Tree toggle expand/collapse
	window.toggleTreeNode = function(nodeId) {
		const node = document.querySelector('[data-id="' + nodeId + '"]');
		if (!node) return;

		const children = node.querySelector('.tree-children');
		if (children) {
			children.style.display = children.style.display === 'none' ? 'block' : 'none';
			node.classList.toggle('expanded');
			node.classList.toggle('collapsed');
		}
	};
})();
</script>`)
}

// RenderTreeActions generates HTML for tree actions
func RenderTreeActions(actions []*TreeAction, node *TreeNode) template.HTML {
	if len(actions) == 0 {
		return ""
	}

	var html string
	for _, action := range actions {
		if action.Visible != nil && !action.Visible(node) {
			continue
		}

		class := "tree-action"
		switch action.Style {
		case ActionPrimary:
			class += " tree-action-primary"
		case ActionDanger:
			class += " tree-action-danger"
		case ActionGhost:
			class += " tree-action-ghost"
		}

		iconHTML := ""
		if action.Icon != "" {
			iconHTML = fmt.Sprintf(`<span class="action-icon">%s</span>`, action.Icon)
		}

		confirm := ""
		if action.Confirm != "" {
			confirm = fmt.Sprintf(` onclick="return confirm('%s')"`, template.JSEscapeString(action.Confirm))
		}

		ajax := ""
		if action.Ajax {
			ajax = ` data-ajax="true"`
		}

		html += fmt.Sprintf(
			`<a href="%s" class="%s"%s%s>%s%s</a>`,
			action.URL, class, confirm, ajax, iconHTML, action.Label)
	}

	return template.HTML(html)
}

// RenderTreeTools generates HTML for tree tools
func RenderTreeTools(tools []*TreeTool) template.HTML {
	if len(tools) == 0 {
		return ""
	}

	var html string
	for _, tool := range tools {
		iconHTML := ""
		if tool.Icon != "" {
			iconHTML = fmt.Sprintf(`<span class="tool-icon">%s</span>`, tool.Icon)
		}

		html += fmt.Sprintf(
			`<button type="button" class="tree-tool" onclick="window.location.href='%s'">%s%s</button>`,
			tool.URL, iconHTML, tool.Label)
	}

	return template.HTML(html)
}

// RenderTreeCheckbox generates HTML for tree checkbox
func RenderTreeCheckbox(node *TreeNode, multiSelect bool) template.HTML {
	type_attr := "radio"
	if multiSelect {
		type_attr = "checkbox"
	}

	checked := ""
	if node.Selected {
		checked = " checked"
	}

	disabled := ""
	if node.Disabled {
		disabled = " disabled"
	}

	return template.HTML(fmt.Sprintf(
		`<input type="%s" name="tree_selection" value="%s" class="tree-checkbox"%s%s>`,
		type_attr, template.HTMLEscapeString(node.ID), checked, disabled))
}
