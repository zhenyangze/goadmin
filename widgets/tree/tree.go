// Package tree provides tree view widgets.
package tree

import "html/template"

// TreeWidget displays a tree structure
type TreeWidget struct {
	id         string
	nodes      []*TreeNode
	expandable bool
	selectable bool
	multiSelect bool
	showIcon   bool
	onSelect   string // JS function name
}

// TreeNode represents a node in the tree
type TreeNode struct {
	ID       string
	Label    string
	Icon     string
	Children []*TreeNode
	Expanded bool
	Selected bool
	Disabled bool
	Data     map[string]string
}

// New creates a new tree widget
func New() *TreeWidget {
	return &TreeWidget{
		nodes:      []*TreeNode{},
		expandable: true,
		selectable: false,
		showIcon:   true,
	}
}

// ID sets the tree ID
func (t *TreeWidget) ID(id string) *TreeWidget {
	t.id = id
	return t
}

// Nodes sets the tree nodes
func (t *TreeWidget) Nodes(nodes ...*TreeNode) *TreeWidget {
	t.nodes = nodes
	return t
}

// Expandable makes nodes expandable
func (t *TreeWidget) Expandable(enable bool) *TreeWidget {
	t.expandable = enable
	return t
}

// Selectable enables node selection
func (t *TreeWidget) Selectable(enable bool) *TreeWidget {
	t.selectable = enable
	return t
}

// MultiSelect enables multiple selection
func (t *TreeWidget) MultiSelect(enable bool) *TreeWidget {
	t.multiSelect = enable
	return t
}

// ShowIcon shows/hides node icons
func (t *TreeWidget) ShowIcon(show bool) *TreeWidget {
	t.showIcon = show
	return t
}

// OnSelect sets the selection callback
func (t *TreeWidget) OnSelect(fn string) *TreeWidget {
	t.onSelect = fn
	return t
}

// AddNode adds a node to the tree
func (t *TreeWidget) AddNode(node *TreeNode) *TreeWidget {
	t.nodes = append(t.nodes, node)
	return t
}

// TreeNode creates a new tree node
func Node(id, label string) *TreeNode {
	return &TreeNode{
		ID:    id,
		Label: label,
		Data:  make(map[string]string),
	}
}

// WithIcon sets the node icon
func (n *TreeNode) WithIcon(icon string) *TreeNode {
	n.Icon = icon
	return n
}

// WithChildren sets child nodes
func (n *TreeNode) WithChildren(children ...*TreeNode) *TreeNode {
	n.Children = children
	return n
}

// WithData sets node data
func (n *TreeNode) WithData(key, value string) *TreeNode {
	n.Data[key] = value
	return n
}

// RenderContext provides data for rendering
type RenderContext struct {
	ID          string
	Nodes       []*TreeNode
	Expandable  bool
	Selectable  bool
	MultiSelect bool
	ShowIcon    bool
	OnSelect    string
}

// Render prepares the tree for rendering
func (t *TreeWidget) Render() *RenderContext {
	return &RenderContext{
		ID:          t.id,
		Nodes:       t.nodes,
		Expandable:  t.expandable,
		Selectable:  t.selectable,
		MultiSelect: t.multiSelect,
		ShowIcon:    t.showIcon,
		OnSelect:    t.onSelect,
	}
}

// JavaScript returns the tree JavaScript
func JavaScript() template.HTML {
	return template.HTML(`
<script>
(function() {
	// Tree toggle functionality
	document.querySelectorAll('.tree-toggle').forEach(toggle => {
		toggle.addEventListener('click', function() {
			const node = this.closest('.tree-node');
			node.classList.toggle('expanded');
			node.classList.toggle('collapsed');
		});
	});

	// Tree selection
	document.querySelectorAll('.tree-selectable .tree-label').forEach(label => {
		label.addEventListener('click', function() {
			const node = this.closest('.tree-node');
			const tree = this.closest('.tree-widget');

			if (tree.dataset.multiSelect !== 'true') {
				tree.querySelectorAll('.tree-node').forEach(n => n.classList.remove('selected'));
			}

			node.classList.toggle('selected');

			// Trigger callback
			const onSelect = tree.dataset.onSelect;
			if (onSelect && window[onSelect]) {
				window[onSelect](node.dataset.id, node.classList.contains('selected'));
			}
		});
	});
})();
</script>`)
}
