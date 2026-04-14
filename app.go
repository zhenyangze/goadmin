package goadmin

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/zhenyangze/goadmin/form"
	"github.com/zhenyangze/goadmin/grid"
	"github.com/zhenyangze/goadmin/show"
	"github.com/zhenyangze/goadmin/tree"
	widgetform "github.com/zhenyangze/goadmin/widgets/form"
)

//go:embed assets/templates/*.tmpl assets/styles/admin.css
var assetFS embed.FS

type dashboardFunc func(context.Context) (DashboardData, error)

// App is the reusable admin HTTP handler.
type App struct {
	cfg            Config
	auth           AuthService
	templates      *template.Template
	dashboard      dashboardFunc
	dashboardPages map[string]DashboardPage
	routes         map[string]Route
	resourceMap    map[string]Resource
	resources      []Resource
	css            []byte
	uploadForm     *http.ServeMux
	toolForms      map[string]*widgetform.ToolForm
}

// New creates a new admin application.
func New(cfg Config, authService AuthService) (*App, error) {
	cfg = cfg.WithDefaults()
	if cfg.SessionSecret == "" {
		return nil, errors.New("session secret is required")
	}

	css, err := assetFS.ReadFile("assets/styles/admin.css")
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("admin").Funcs(template.FuncMap{
		"safeHTML": func(v any) template.HTML {
			switch typed := v.(type) {
			case template.HTML:
				return typed
			default:
				return template.HTML(template.HTMLEscapeString(fmt.Sprint(v)))
			}
		},
		"substr": func(s string, start, length int) string {
			if start < 0 || start >= len(s) {
				return ""
			}
			end := start + length
			if end > len(s) {
				end = len(s)
			}
			return s[start:end]
		},
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"add": func(a, b int) int {
			return a + b
		},
	}).ParseFS(assetFS, "assets/templates/*.tmpl")
	if err != nil {
		return nil, err
	}

	return &App{
		cfg:            cfg,
		auth:           authService,
		templates:      tmpl,
		dashboardPages: map[string]DashboardPage{},
		routes:         map[string]Route{},
		resourceMap:    map[string]Resource{},
		css:            css,
		toolForms:      map[string]*widgetform.ToolForm{},
	}, nil
}

// Register adds a CRUD resource to the app.
func (a *App) Register(resource Resource) {
	resource.Path = strings.Trim(resource.Path, "/")
	a.resourceMap[resource.Path] = resource
	a.resources = append(a.resources, resource)
}

// SetDashboard configures a custom dashboard callback.
func (a *App) SetDashboard(fn func(context.Context) (DashboardData, error)) {
	a.dashboard = fn
}

// RegisterDashboardPage adds a custom dashboard-style page under the admin prefix.
func (a *App) RegisterDashboardPage(page DashboardPage) {
	page.Path = strings.Trim(page.Path, "/")
	a.dashboardPages[page.Path] = page
}

// RegisterRoute adds a custom non-resource route under the admin prefix.
func (a *App) RegisterRoute(route Route) {
	route.Path = strings.Trim(route.Path, "/")
	a.routes[route.Path] = route
}

// RegisterToolForm registers a tool form under the admin prefix.
// The form can be accessed at /{prefix}/form/{path}.
func (a *App) RegisterToolForm(path string, form *widgetform.ToolForm) {
	path = strings.Trim(path, "/")
	// Ensure path doesn't conflict with resources
	if _, ok := a.resourceMap[path]; ok {
		panic(fmt.Sprintf("tool form path %q conflicts with registered resource", path))
	}
	// Ensure path doesn't conflict with dashboard pages
	if _, ok := a.dashboardPages[path]; ok {
		panic(fmt.Sprintf("tool form path %q conflicts with dashboard page", path))
	}
	// Ensure path doesn't conflict with routes
	if _, ok := a.routes[path]; ok {
		panic(fmt.Sprintf("tool form path %q conflicts with registered route", path))
	}
	a.toolForms[path] = form
}

// Handler returns the reusable HTTP handler.
func (a *App) Handler() http.Handler {
	return a
}

// ServeHTTP implements http.Handler.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	prefix := normalizePath(a.cfg.Prefix)
	if !strings.HasPrefix(normalizePath(r.URL.Path), prefix) && normalizePath(r.URL.Path) != prefix {
		http.NotFound(w, r)
		return
	}

	if normalizePath(r.URL.Path) == joinURL(prefix, "assets", "admin.css") {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		_, _ = w.Write(a.css)
		return
	}
	if a.serveUploads(w, r) {
		return
	}

	path := strings.TrimPrefix(r.URL.Path, prefix)
	path = strings.Trim(path, "/")

	switch path {
	case "login":
		a.handleLogin(w, r)
		return
	case "logout":
		a.handleLogout(w, r)
		return
	}

	state, identity := a.requireIdentity(w, r)
	if identity == nil {
		return
	}

	if path == "" {
		a.handleDashboard(w, r, state, identity)
		return
	}

	if page, ok := a.dashboardPages[path]; ok && currentMethod(r) == http.MethodGet {
		allowed, err := a.auth.Authorize(r.Context(), identity, page.Permission)
		if err != nil {
			a.writeError(w, http.StatusInternalServerError, err)
			return
		}
		if !allowed {
			http.NotFound(w, r)
			return
		}
		data, err := page.Build(r.Context(), r, identity)
		if err != nil {
			a.writeError(w, http.StatusInternalServerError, err)
			return
		}
		a.renderShell(w, r, state, identity, "dashboard", pageData{
			PageTitle:       fallback(data.Title, page.Title),
			PageDescription: fallback(data.Description, page.Description),
			Content:         a.buildDashboardView(data),
		})
		return
	}

	if route, ok := a.routes[path]; ok {
		allowed, err := a.auth.Authorize(r.Context(), identity, route.Permission)
		if err != nil {
			a.writeError(w, http.StatusInternalServerError, err)
			return
		}
		if !allowed {
			http.NotFound(w, r)
			return
		}
		if !routeAllows(route, currentMethod(r)) {
			http.NotFound(w, r)
			return
		}
		if err := route.Handler(w, r, identity); err != nil {
			a.writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	// Handle tool forms
	if form, ok := a.toolForms[path]; ok {
		a.handleToolForm(w, r, state, identity, form)
		return
	}

	parts := strings.Split(path, "/")
	resource, ok := a.resourceMap[parts[0]]
	if !ok {
		http.NotFound(w, r)
		return
	}
	allowed, err := a.auth.Authorize(r.Context(), identity, resource.Permission)
	if err != nil {
		a.writeError(w, http.StatusInternalServerError, err)
		return
	}
	if !allowed {
		a.renderShell(w, r, state, identity, "dashboard", pageData{
			PageTitle:       "Forbidden",
			PageDescription: "You do not have permission to access this module.",
			Message:         "Forbidden",
			Content: dashboardView{
				Cards: []DashboardCard{{Title: "Access denied", Value: "403", Hint: "Permission check failed"}},
			},
		})
		return
	}

	method := currentMethod(r)
	switch len(parts) {
	case 1:
		if method == http.MethodGet {
			a.handleGrid(w, r, state, identity, resource)
			return
		}
		if method == http.MethodPost {
			a.handleCreate(w, r, state, identity, resource)
			return
		}
	case 2:
		if parts[1] == "new" && method == http.MethodGet {
			a.handleNewForm(w, r, state, identity, resource)
			return
		}
		if parts[1] == "tree" && method == http.MethodGet {
			a.handleTree(w, r, state, identity, resource)
			return
		}
		if method == http.MethodGet {
			a.handleShow(w, r, state, identity, resource, parts[1])
			return
		}
		if method == http.MethodPut {
			a.handleUpdate(w, r, state, identity, resource, parts[1])
			return
		}
	case 3:
		if parts[2] == "edit" && method == http.MethodGet {
			a.handleEditForm(w, r, state, identity, resource, parts[1])
			return
		}
		if parts[2] == "delete" && method == http.MethodPost {
			a.handleDelete(w, r, state, identity, resource, parts[1])
			return
		}
		if parts[2] == "batch-delete" && method == http.MethodPost {
			a.handleBatchDelete(w, r, state, identity, resource)
			return
		}
		// Handle custom batch actions
		if strings.HasPrefix(parts[2], "batch-action") && method == http.MethodPost {
			a.handleCustomBatchAction(w, r, state, identity, resource, parts[2])
			return
		}
	}

	http.NotFound(w, r)
}

func (a *App) requireIdentity(w http.ResponseWriter, r *http.Request) (*sessionState, *Identity) {
	state, err := a.readSession(r)
	if err != nil {
		http.Redirect(w, r, joinURL(a.cfg.Prefix, "login"), http.StatusFound)
		return nil, nil
	}
	identity, err := a.auth.FindIdentity(r.Context(), state.UserID)
	if err != nil {
		a.clearSession(w)
		http.Redirect(w, r, joinURL(a.cfg.Prefix, "login"), http.StatusFound)
		return nil, nil
	}
	return state, identity
}

func (a *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		_ = a.templates.ExecuteTemplate(w, "login", map[string]any{
			"AppName": a.cfg.AppName,
			"Title":   a.cfg.Title,
			"Prefix":  normalizePath(a.cfg.Prefix),
			"Error":   r.URL.Query().Get("error"),
			"Theme":   a.cfg.Theme,
		})
		return
	}
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	if err := r.ParseForm(); err != nil {
		a.writeError(w, http.StatusBadRequest, err)
		return
	}
	identity, err := a.auth.Authenticate(r.Context(), r.FormValue("username"), r.FormValue("password"))
	if err != nil {
		http.Redirect(w, r, joinURL(a.cfg.Prefix, "login")+"?error="+url.QueryEscape(err.Error()), http.StatusFound)
		return
	}
	if err := a.writeSession(w, identity.ID); err != nil {
		a.writeError(w, http.StatusInternalServerError, err)
		return
	}
	http.Redirect(w, r, normalizePath(a.cfg.Prefix), http.StatusFound)
}

func (a *App) handleLogout(w http.ResponseWriter, r *http.Request) {
	if currentMethod(r) != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	state, _ := a.readSession(r)
	if state == nil || r.FormValue("_csrf") != state.CSRF {
		a.writeError(w, http.StatusBadRequest, errors.New("invalid csrf token"))
		return
	}
	a.clearSession(w)
	http.Redirect(w, r, joinURL(a.cfg.Prefix, "login"), http.StatusFound)
}

func (a *App) handleDashboard(w http.ResponseWriter, r *http.Request, state *sessionState, identity *Identity) {
	data := DashboardData{
		Title:       "Dashboard",
		Description: "Go implementation inspired by Dcat Admin.",
		Cards: []DashboardCard{
			{Title: "Resources", Value: strconv.Itoa(len(a.resources)), Hint: "Registered admin modules"},
		},
	}
	if a.dashboard != nil {
		if custom, err := a.dashboard(r.Context()); err == nil {
			data = custom
		}
	}
	a.renderShell(w, r, state, identity, "dashboard", pageData{
		PageTitle:       data.Title,
		PageDescription: data.Description,
		Content:         a.buildDashboardView(data),
	})
}

func (a *App) handleGrid(w http.ResponseWriter, r *http.Request, state *sessionState, identity *Identity, resource Resource) {
	builder := grid.New()
	if resource.BuildGrid != nil {
		resource.BuildGrid(builder)
	}
	query := ListQuery{
		Page:      intFromQuery(r, "page", 1),
		PerPage:   intFromQuery(r, "per_page", 10),
		Search:    r.URL.Query().Get("q"),
		Sort:      r.URL.Query().Get("sort"),
		Direction: r.URL.Query().Get("direction"),
		Filters:   map[string]string{},
	}
	for _, filter := range builder.Filters {
		query.Filters[filter.Name] = r.URL.Query().Get("f_" + filter.Name)
	}

	result, err := resource.Repository.List(r.Context(), query)
	if err != nil {
		a.writeError(w, http.StatusInternalServerError, err)
		return
	}

	baseURL := joinURL(a.cfg.Prefix, resource.Path)
	colSpan := len(builder.Columns) + 1 // +1 for Actions column
	if !builder.DisableRowSelector {
		colSpan++ // +1 for checkbox column
	}
	view := gridView{
		Title:              fallback(builder.Title, resource.Title),
		Description:        fallback(builder.Description, resource.Description),
		Actions:            a.buildGridPageActions(baseURL, builder, resource),
		BatchActions:       a.buildGridBatchActions(baseURL, builder),
		Tools:              a.buildGridTools(builder),
		Columns:            a.buildGridColumns(r, baseURL, builder, query),
		Rows:               a.buildGridRows(baseURL, builder, result.Items),
		Filters:            a.buildGridFilters(builder, query),
		QuickSearch:        query.Search,
		CurrentPath:        baseURL,
		Pagination:         buildPagination(baseURL, query, result.Total, result.Page, result.PerPage),
		ResultSummary:      fmt.Sprintf("Total %d records", result.Total),
		EmptyText:          fallback(resource.EmptyText, "No records yet."),
		ColSpan:            colSpan,
		CSRF:               state.CSRF,
		// New features
		DisableRowSelector: builder.DisableRowSelector,
		RowSelector: gridRowSelectorView{
			Enabled:     !builder.DisableRowSelector,
			TitleColumn: builder.RowSelectorTitleColumn,
			IDColumn:    builder.RowSelectorIDColumn,
			Click:       builder.RowSelectorClick,
		},
		EnableDialogCreate: builder.EnableDialogCreate,
		EnableDialogEdit:   builder.EnableDialogEdit,
		DialogWidth:        fallback(builder.DialogWidth, "700px"),
		DialogHeight:       fallback(builder.DialogHeight, "670px"),
		ScrollbarX:         builder.ScrollbarX,
		TableClasses:       builder.TableClasses,
		PerPageOptions:     builder.PerPageOptions,
		CurrentPerPage:     query.PerPage,
		ToolsWithOutline:   builder.ToolsWithOutline,
		DisableRefresh:     builder.DisableRefresh,
	}

	a.renderShell(w, r, state, identity, "grid", pageData{
		PageTitle:       view.Title,
		PageDescription: view.Description,
		Helper:          a.resourceHelper(resource),
		Content:         view,
	})
}

func (a *App) handleNewForm(w http.ResponseWriter, r *http.Request, state *sessionState, identity *Identity, resource Resource) {
	builder := form.New()
	if resource.BuildForm == nil {
		http.NotFound(w, r)
		return
	}
	resource.BuildForm(builder)
	view := a.buildFormView(resource, builder, nil, state.CSRF, "")
	a.renderShell(w, r, state, identity, "form", pageData{
		PageTitle:       fallback(builder.Title, "Create "+resource.Title),
		PageDescription: fallback(builder.Description, resource.Description),
		Helper:          a.resourceHelper(resource),
		Content:         view,
	})
}

func (a *App) handleEditForm(w http.ResponseWriter, r *http.Request, state *sessionState, identity *Identity, resource Resource, id string) {
	builder := form.New()
	if resource.BuildForm == nil {
		http.NotFound(w, r)
		return
	}
	resource.BuildForm(builder)
	record, err := resource.Repository.Get(r.Context(), id)
	if err != nil {
		a.writeError(w, http.StatusNotFound, err)
		return
	}
	view := a.buildFormView(resource, builder, record, state.CSRF, id)
	a.renderShell(w, r, state, identity, "form", pageData{
		PageTitle:       fallback(builder.Title, "Edit "+resource.Title),
		PageDescription: fallback(builder.Description, resource.Description),
		Helper:          a.resourceHelper(resource),
		Content:         view,
	})
}

func (a *App) handleShow(w http.ResponseWriter, r *http.Request, state *sessionState, identity *Identity, resource Resource, id string) {
	if resource.BuildShow == nil {
		http.NotFound(w, r)
		return
	}
	builder := show.New()
	resource.BuildShow(builder)
	record, err := resource.Repository.Get(r.Context(), id)
	if err != nil {
		a.writeError(w, http.StatusNotFound, err)
		return
	}
	view := a.buildShowView(resource, builder, record, id)
	a.renderShell(w, r, state, identity, "show", pageData{
		PageTitle:       fallback(builder.Title, resource.Title+" detail"),
		PageDescription: fallback(builder.Description, resource.Description),
		Helper:          a.resourceHelper(resource),
		Content:         view,
	})
}

func (a *App) handleTree(w http.ResponseWriter, r *http.Request, state *sessionState, identity *Identity, resource Resource) {
	provider, ok := resource.Repository.(TreeProvider)
	if !ok || resource.BuildTree == nil {
		http.NotFound(w, r)
		return
	}
	builder := tree.New()
	resource.BuildTree(builder)
	nodes, err := provider.Tree(r.Context())
	if err != nil {
		a.writeError(w, http.StatusInternalServerError, err)
		return
	}
	baseURL := joinURL(a.cfg.Prefix, resource.Path)
	viewNodes := a.buildTreeNodeViews(nodes, baseURL)
	a.renderShell(w, r, state, identity, "tree", pageData{
		PageTitle:       fallback(builder.Title, resource.Title+" tree"),
		PageDescription: fallback(builder.Description, resource.Description),
		Helper:          a.resourceHelper(resource),
		Content: treeView{
			Title:       fallback(builder.Title, resource.Title+" tree"),
			Description: fallback(builder.Description, resource.Description),
			EmptyText:   builder.EmptyText,
			Nodes:       viewNodes,
		},
	})
}

func (a *App) handleCreate(w http.ResponseWriter, r *http.Request, state *sessionState, identity *Identity, resource Resource) {
	if resource.BuildForm == nil {
		http.NotFound(w, r)
		return
	}
	builder := form.New()
	resource.BuildForm(builder)
	if err := a.parseFormRequest(resource, r); err != nil {
		a.writeError(w, http.StatusBadRequest, err)
		return
	}
	if r.FormValue("_csrf") != state.CSRF {
		a.writeError(w, http.StatusBadRequest, errors.New("invalid csrf token"))
		return
	}
	values := a.formValues(resource)
	submitted, err := values(r)
	if err != nil {
		a.writeError(w, http.StatusBadRequest, err)
		return
	}
	if fieldErrors := a.validateForm(builder, submitted, nil); len(fieldErrors) > 0 {
		a.cleanupSubmittedUploads(r.Context(), resource, submitted)
		a.renderFormError(w, r, state, identity, resource, builder, nil, "", submitted, fieldErrors, "Please correct the highlighted fields.")
		return
	}
	// 将 identity 添加到 context
	ctx := context.WithValue(r.Context(), "identity", identity)
	if err := resource.Repository.Create(ctx, submitted); err != nil {
		a.cleanupSubmittedUploads(r.Context(), resource, submitted)
		a.renderFormError(w, r, state, identity, resource, builder, nil, "", submitted, nil, err.Error())
		return
	}
	http.Redirect(w, r, joinURL(a.cfg.Prefix, resource.Path)+"?flash=created", http.StatusFound)
}

func (a *App) handleUpdate(w http.ResponseWriter, r *http.Request, state *sessionState, identity *Identity, resource Resource, id string) {
	if resource.BuildForm == nil {
		http.NotFound(w, r)
		return
	}
	builder := form.New()
	resource.BuildForm(builder)
	existing, err := resource.Repository.Get(r.Context(), id)
	if err != nil {
		a.writeError(w, http.StatusNotFound, err)
		return
	}
	uploadFields := a.uploadFields(resource)
	if err := a.parseFormRequest(resource, r); err != nil {
		a.writeError(w, http.StatusBadRequest, err)
		return
	}
	if r.FormValue("_csrf") != state.CSRF {
		a.writeError(w, http.StatusBadRequest, errors.New("invalid csrf token"))
		return
	}
	values := a.formValues(resource)
	submitted, err := values(r)
	if err != nil {
		a.writeError(w, http.StatusBadRequest, err)
		return
	}
	if fieldErrors := a.validateForm(builder, submitted, existing); len(fieldErrors) > 0 {
		a.cleanupSubmittedUploads(r.Context(), resource, submitted)
		a.renderFormError(w, r, state, identity, resource, builder, existing, id, submitted, fieldErrors, "Please correct the highlighted fields.")
		return
	}
	// 将 identity 添加到 context
	ctx := context.WithValue(r.Context(), "identity", identity)
	if err := resource.Repository.Update(ctx, id, submitted); err != nil {
		a.cleanupSubmittedUploads(r.Context(), resource, submitted)
		a.renderFormError(w, r, state, identity, resource, builder, existing, id, submitted, nil, err.Error())
		return
	}
	a.cleanupReplacedUploads(r.Context(), existing, uploadFields, submitted)
	http.Redirect(w, r, joinURL(a.cfg.Prefix, resource.Path, id)+"?flash=updated", http.StatusFound)
}

func (a *App) handleDelete(w http.ResponseWriter, r *http.Request, state *sessionState, identity *Identity, resource Resource, id string) {
	if err := r.ParseForm(); err != nil {
		a.writeError(w, http.StatusBadRequest, err)
		return
	}
	if r.FormValue("_csrf") != state.CSRF {
		a.writeError(w, http.StatusBadRequest, errors.New("invalid csrf token"))
		return
	}
	var existing any
	var uploadFields []*form.Field
	if uploadFields = a.uploadFields(resource); len(uploadFields) > 0 {
		record, err := resource.Repository.Get(r.Context(), id)
		if err == nil {
			existing = record
		}
	}
	// 将 identity 添加到 context
	ctx := context.WithValue(r.Context(), "identity", identity)
	if err := resource.Repository.Delete(ctx, id); err != nil {
		a.writeError(w, http.StatusBadRequest, err)
		return
	}
	a.cleanupDeletedUploads(r.Context(), existing, uploadFields)
	http.Redirect(w, r, joinURL(a.cfg.Prefix, resource.Path)+"?flash=deleted", http.StatusFound)
}

func (a *App) handleBatchDelete(w http.ResponseWriter, r *http.Request, state *sessionState, _ *Identity, resource Resource) {
	if err := r.ParseForm(); err != nil {
		a.writeError(w, http.StatusBadRequest, err)
		return
	}
	if r.FormValue("_csrf") != state.CSRF {
		a.writeError(w, http.StatusBadRequest, errors.New("invalid csrf token"))
		return
	}
	ids := r.FormValue("_ids")
	if ids == "" {
		http.Redirect(w, r, joinURL(a.cfg.Prefix, resource.Path), http.StatusFound)
		return
	}
	idList := strings.Split(ids, ",")
	for _, id := range idList {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		var existing any
		var uploadFields []*form.Field
		if uploadFields = a.uploadFields(resource); len(uploadFields) > 0 {
			record, err := resource.Repository.Get(r.Context(), id)
			if err == nil {
				existing = record
			}
		}
		if err := resource.Repository.Delete(r.Context(), id); err != nil {
			a.writeError(w, http.StatusBadRequest, err)
			return
		}
		a.cleanupDeletedUploads(r.Context(), existing, uploadFields)
	}
	http.Redirect(w, r, joinURL(a.cfg.Prefix, resource.Path)+"?flash=batch-deleted", http.StatusFound)
}

func (a *App) handleCustomBatchAction(w http.ResponseWriter, r *http.Request, state *sessionState, _ *Identity, resource Resource, actionPath string) {
	if err := r.ParseForm(); err != nil {
		a.writeError(w, http.StatusBadRequest, err)
		return
	}
	if r.FormValue("_csrf") != state.CSRF {
		a.writeError(w, http.StatusBadRequest, errors.New("invalid csrf token"))
		return
	}

	// Extract action name from path (e.g., "batch-action:Export" -> "Export")
	parts := strings.SplitN(actionPath, ":", 2)
	actionName := ""
	if len(parts) == 2 {
		actionName = parts[1]
	} else {
		// Try URL-decoded version
		actionName = strings.TrimPrefix(actionPath, "batch-action")
		actionName = strings.TrimPrefix(actionName, "/")
		actionName, _ = url.QueryUnescape(actionName)
	}

	ids := r.FormValue("_ids")
	if ids == "" {
		if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"success": false,
				"message": "No items selected",
			})
			return
		}
		http.Redirect(w, r, joinURL(a.cfg.Prefix, resource.Path), http.StatusFound)
		return
	}

	idList := strings.Split(ids, ",")
	filteredIDs := make([]string, 0, len(idList))
	for _, id := range idList {
		id = strings.TrimSpace(id)
		if id != "" {
			filteredIDs = append(filteredIDs, id)
		}
	}

	// Find the handler from the builder
	builder := grid.New()
	if resource.BuildGrid != nil {
		resource.BuildGrid(builder)
	}

	var handler grid.BatchHandler
	for _, action := range builder.BatchActionHandlers {
		if action != nil && action.Label == actionName {
			handler = action.Handler
			break
		}
	}

	if handler == nil {
		if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"success": false,
				"message": "Action not found: " + actionName,
			})
			return
		}
		http.NotFound(w, r)
		return
	}

	// Execute the handler
	err := handler(r.Context(), filteredIDs)
	if err != nil {
		if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		a.writeError(w, http.StatusInternalServerError, err)
		return
	}

	// Return AJAX response or redirect
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"message": "Operation completed successfully",
		})
		return
	}
	http.Redirect(w, r, joinURL(a.cfg.Prefix, resource.Path)+"?flash=batch-action-completed", http.StatusFound)
}

func (a *App) formValues(resource Resource) func(*http.Request) (Values, error) {
	builder := form.New()
	resource.BuildForm(builder)
	return func(r *http.Request) (Values, error) {
		values := Values{}
		for _, field := range builder.Fields {
			if field.Type == form.FieldDisplay {
				continue
			}
			if field.Type == form.FieldRepeater {
				encoded, err := collectRepeaterValue(r.Form, field)
				if err != nil {
					return nil, err
				}
				values[field.Name] = []string{encoded}
				continue
			}
			if field.Type == form.FieldUpload {
				files := uploadHeaders(r, field.Name)
				locations := make([]string, 0, len(files))
				if !field.Multiple && len(files) == 0 {
					file, header, err := r.FormFile(field.Name)
					if err == nil {
						files = []*multipart.FileHeader{header}
						location, saveErr := func() (string, error) {
							defer file.Close()
							if err := a.validateUpload(field, header); err != nil {
								return "", err
							}
							return a.saveUpload(r.Context(), file, header)
						}()
						if saveErr != nil {
							return nil, saveErr
						}
						if strings.TrimSpace(location) != "" {
							locations = append(locations, location)
						}
					}
				}
				if len(files) == 0 {
					continue
				}
				if len(locations) == 0 {
					for _, header := range files {
						if err := a.validateUpload(field, header); err != nil {
							return nil, err
						}
						file, err := header.Open()
						if err != nil {
							return nil, err
						}
						location, saveErr := a.saveUpload(r.Context(), file, header)
						file.Close()
						if saveErr != nil {
							return nil, saveErr
						}
						if strings.TrimSpace(location) != "" {
							locations = append(locations, location)
						}
					}
				}
				if len(locations) > 0 {
					values[field.Name] = locations
				}
				continue
			}
			collectFieldValues(values, r.Form, field.Name, field.Type)
			if field.SecondName != "" {
				collectFieldValues(values, r.Form, field.SecondName, form.FieldDate)
			}
		}
		return values, nil
	}
}

func uploadHeaders(r *http.Request, name string) []*multipart.FileHeader {
	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		return nil
	}
	return r.MultipartForm.File[name]
}

func (a *App) parseFormRequest(resource Resource, r *http.Request) error {
	builder := form.New()
	resource.BuildForm(builder)
	for _, field := range builder.Fields {
		if field.Type == form.FieldUpload {
			return r.ParseMultipartForm(32 << 20)
		}
	}
	return r.ParseForm()
}

func (a *App) uploadFields(resource Resource) []*form.Field {
	if resource.BuildForm == nil {
		return nil
	}
	builder := form.New()
	resource.BuildForm(builder)
	fields := make([]*form.Field, 0)
	for _, field := range builder.Fields {
		if field.Type == form.FieldUpload {
			fields = append(fields, field)
		}
	}
	return fields
}

func (a *App) cleanupReplacedUploads(ctx context.Context, record any, fields []*form.Field, submitted Values) {
	if record == nil {
		return
	}
	for _, field := range fields {
		newValues, ok := submitted[field.Name]
		if !ok || len(newValues) == 0 {
			continue
		}
		oldValues := a.uploadValuesFromRecord(record, field)
		for _, oldValue := range oldValues {
			keep := false
			for _, newValue := range newValues {
				if oldValue == newValue {
					keep = true
					break
				}
			}
			if !keep {
				_ = a.deleteUpload(ctx, oldValue)
			}
		}
	}
}

func (a *App) cleanupSubmittedUploads(ctx context.Context, resource Resource, submitted Values) {
	for _, field := range a.uploadFields(resource) {
		for _, value := range submitted.All(field.Name) {
			_ = a.deleteUpload(ctx, value)
		}
	}
}

func (a *App) validateForm(builder *form.Builder, submitted Values, record any) map[string]string {
	fieldErrors := map[string]string{}
	for _, field := range builder.Fields {
		if field.Type == form.FieldDisplay || field.Type == form.FieldHidden {
			continue
		}
		currentValue := strings.TrimSpace(submitted.First(field.Name))
		existingValues := a.uploadValuesFromRecord(record, field)
		switch field.Type {
		case form.FieldPassword:
			if field.Required && record == nil && currentValue == "" {
				fieldErrors[field.Name] = "This field is required."
			}
		case form.FieldMulti:
			if field.Required && len(submitted.All(field.Name)) == 0 {
				fieldErrors[field.Name] = "Select at least one option."
			}
		case form.FieldUpload:
			if field.Required && len(submitted.All(field.Name)) == 0 && len(existingValues) == 0 {
				fieldErrors[field.Name] = "Please upload a file."
			}
		case form.FieldRepeater:
			if field.Required {
				raw := strings.TrimSpace(submitted.First(field.Name))
				if strings.TrimSpace(raw) == "" {
					fieldErrors[field.Name] = "Please add at least one item."
				}
			}
		case form.FieldDateRange:
			start := strings.TrimSpace(submitted.First(field.Name))
			end := strings.TrimSpace(submitted.First(field.SecondName))
			if field.Required && (start == "" || end == "") {
				fieldErrors[field.Name] = "Please provide both start and end dates."
				continue
			}
			if (start == "") != (end == "") {
				fieldErrors[field.Name] = "Start and end dates must be provided together."
				continue
			}
			if start != "" && end != "" && start > end {
				fieldErrors[field.Name] = "Start date must be before or equal to end date."
			}
		default:
			if field.Required && currentValue == "" {
				fieldErrors[field.Name] = "This field is required."
			}
		}
	}
	return fieldErrors
}

func (a *App) renderFormError(w http.ResponseWriter, r *http.Request, state *sessionState, identity *Identity, resource Resource, builder *form.Builder, record any, id string, submitted Values, fieldErrors map[string]string, message string) {
	w.WriteHeader(http.StatusBadRequest)
	view := a.buildFormViewState(resource, builder, record, state.CSRF, id, submitted, fieldErrors, message)
	a.renderShell(w, r, state, identity, "form", pageData{
		PageTitle:       view.Title,
		PageDescription: view.Description,
		Content:         view,
	})
}

func (a *App) cleanupDeletedUploads(ctx context.Context, record any, fields []*form.Field) {
	if record == nil {
		return
	}
	for _, field := range fields {
		for _, value := range a.uploadValuesFromRecord(record, field) {
			_ = a.deleteUpload(ctx, value)
		}
	}
}

func (a *App) uploadValuesFromRecord(record any, field *form.Field) []string {
	path := field.Name
	if field.ValuePath != "" {
		path = field.ValuePath
	}
	value := valueFromPath(record, path)
	if value == nil {
		return nil
	}
	if field.Multiple {
		switch typed := value.(type) {
		case string:
			return splitCommaSeparated(typed)
		case []string:
			return typed
		default:
			return splitCommaSeparated(formatValue(value))
		}
	}
	return []string{formatValue(value)}
}

func collectFieldValues(values Values, formValues map[string][]string, name string, fieldType form.FieldType) {
	submitted, ok := formValues[name]
	if !ok {
		if fieldType == form.FieldMulti {
			values[name] = nil
		} else if fieldType == form.FieldSwitch {
			values[name] = []string{"0"}
		}
		return
	}
	raw := append([]string(nil), submitted...)
	if fieldType == form.FieldPassword && len(raw) > 0 && strings.TrimSpace(raw[0]) == "" {
		return
	}
	if fieldType == form.FieldMulti {
		cleaned := make([]string, 0, len(raw))
		for _, value := range raw {
			if trimmed := strings.TrimSpace(value); trimmed != "" {
				cleaned = append(cleaned, trimmed)
			}
		}
		values[name] = cleaned
		return
	}
	if len(raw) == 0 || strings.TrimSpace(raw[0]) == "" {
		return
	}
	values[name] = raw[:1]
}

func (a *App) renderShell(w http.ResponseWriter, r *http.Request, state *sessionState, identity *Identity, contentTemplate string, data pageData) {
	navigation, err := a.auth.Navigation(r.Context(), identity)
	if err != nil {
		a.writeError(w, http.StatusInternalServerError, err)
		return
	}
	page := layoutData{
		AppName:         a.cfg.AppName,
		Title:           a.cfg.Title,
		PageTitle:       data.PageTitle,
		PageDescription: data.PageDescription,
		CurrentPath:     normalizePath(r.URL.Path),
		Prefix:          normalizePath(a.cfg.Prefix),
		Flash:           r.URL.Query().Get("flash"),
		FlashMessages:   r.URL.Query().Get("flash-messages"),
		Theme:           a.cfg.Theme,
		User:            identity,
		Menu:            a.menuView(navigation, normalizePath(r.URL.Path)),
		CSRF:            state.CSRF,
		ContentHTML:     a.renderPartial(contentTemplate, data.Content),
		Message:         data.Message,
		Helper:          data.Helper,
	}
	if err := a.templates.ExecuteTemplate(w, "layout", page); err != nil {
		a.writeError(w, http.StatusInternalServerError, err)
	}
}

func (a *App) renderPartial(name string, data any) template.HTML {
	var buf bytes.Buffer
	if err := a.templates.ExecuteTemplate(&buf, name, data); err != nil {
		return template.HTML(template.HTMLEscapeString(err.Error()))
	}
	return template.HTML(buf.String())
}

// FlashRedirect redirects with a flash message.
func (a *App) FlashRedirect(w http.ResponseWriter, r *http.Request, redirectURL string, flashType, message string) {
	http.Redirect(w, r, redirectURL+"?flash="+url.QueryEscape(message), http.StatusFound)
}

// FlashSuccess redirects with a success flash message.
func (a *App) FlashSuccess(w http.ResponseWriter, r *http.Request, url, message string) {
	a.FlashRedirect(w, r, url, "success", message)
}

// FlashError redirects with an error flash message.
func (a *App) FlashError(w http.ResponseWriter, r *http.Request, url, message string) {
	a.FlashRedirect(w, r, url, "error", message)
}

func (a *App) writeError(w http.ResponseWriter, code int, err error) {
	http.Error(w, err.Error(), code)
}

func currentMethod(r *http.Request) string {
	if r.Method != http.MethodPost {
		return r.Method
	}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		_ = r.ParseMultipartForm(32 << 20)
	} else {
		_ = r.ParseForm()
	}
	if override := r.FormValue("_method"); override != "" {
		return strings.ToUpper(override)
	}
	return r.Method
}

func intFromQuery(r *http.Request, key string, fallbackValue int) int {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return fallbackValue
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < 1 {
		return fallbackValue
	}
	return value
}

func fallback(value, other string) string {
	if value != "" {
		return value
	}
	return other
}

type pageData struct {
	PageTitle       string
	PageDescription string
	Message         string
	Helper          *helperPanelView
	Content         any
}

type layoutData struct {
	AppName         string
	Title           string
	PageTitle       string
	PageDescription string
	CurrentPath     string
	Prefix          string
	Flash           string
	FlashMessages   string
	Theme           any
	User            *Identity
	Menu            []menuItemView
	CSRF            string
	ContentHTML     template.HTML
	Message         string
	Helper          *helperPanelView
}

type menuItemView struct {
	Title    string
	Icon     string
	URL      string
	Active   bool
	Expanded bool
	Children []menuItemView
}

type dashboardView struct {
	Cards  []DashboardCard
	Panels []dashboardPanelView
}

type dashboardPanelView struct {
	Title       string
	Description string
	EmptyText   string
	Actions     []dashboardActionView
	Items       []dashboardPanelItemView
}

type dashboardActionView struct {
	Label string
	URL   string
	Class string
}

type dashboardPanelItemView struct {
	Title       string
	Description string
	Value       string
	URL         string
	Tags        []string
}

type helperPanelView struct {
	Title string
	Tags  []string
	Items []string
}

type gridView struct {
	Title              string
	Description        string
	QuickSearch        string
	Filters            []gridFilterView
	Columns            []gridColumnView
	Rows               []gridRowView
	Actions            []gridActionView
	BatchActions       []gridActionView
	Tools              []gridToolView
	Pagination         []paginationLink
	CurrentPath        string
	ResultSummary      string
	EmptyText          string
	ColSpan            int
	CSRF               string
	// New features
	DisableRowSelector bool
	RowSelector        gridRowSelectorView
	EnableDialogCreate bool
	EnableDialogEdit   bool
	DialogWidth        string
	DialogHeight       string
	ScrollbarX         bool
	TableClasses       []string
	PerPageOptions     []int
	CurrentPerPage     int
	ToolsWithOutline   bool
	DisableRefresh     bool
}

type gridToolView struct {
	HTML   template.HTML
	Script string
}

type gridRowSelectorView struct {
	Enabled    bool
	TitleColumn string
	IDColumn   string
	Click      bool
}

type gridFilterView struct {
	Name     string
	Label    string
	Kind     string
	Value    string
	Options  []gridOptionView
	InputKey string
}

type gridOptionView struct {
	Label    string
	Value    string
	Selected bool
}

type gridColumnView struct {
	Label    string
	Sortable bool
	SortURL  string
}

type gridRowView struct {
	ID       string
	Cells    []template.HTML
	Actions  []gridActionView
	Checked  bool
	Disabled bool
}

type gridActionView struct {
	Label   string
	URL     string
	Method  string
	Confirm string
	Class   string
}

type paginationLink struct {
	Label   string
	URL     string
	Active  bool
	Current bool
}

type formView struct {
	Title        string
	Description  string
	Message      string
	Action       string
	Method       string
	Enctype      string
	DeleteAction string
	ShowDelete   bool
	BackURL      string
	SubmitLabel  string
	DeleteLabel  string
	CSRF         string
	Fields       []formFieldView
}

type formFieldView struct {
	Name              string
	SecondName        string
	Label             string
	Type              string
	Value             string
	SecondValue       string
	Values            []string
	Checked           bool
	Multiple          bool
	FileURL           string
	FileURLs          []string
	IsImage           bool
	RepeaterFields    []repeaterChildView
	RepeaterRows      []repeaterRowView
	KeyValueRows      []keyValueRowView
	KeyPlaceholder    string
	ValuePlaceholder  string
	Error             string
	Help              string
	Required          bool
	Readonly          bool
	Disabled          bool
	Placeholder       string
	SecondPlaceholder string
	Options           []gridOptionView
	Min               *float64
	Max               *float64
	Step              *float64
	Inline            bool
	CurrencySymbol    string
	EditorHeight      int
	MaskFormat        string
	SliderMin         *float64
	SliderMax         *float64
	SliderStep        *float64
	SliderPostfix     string
	AutocompleteURL   string
	UploadDir         string
	UploadMaxCount    int
	IsSortable        bool
	HtmlContent       string
	// Complex fields
	MapProvider           string
	TreeNodes             []form.TreeNode
	TreeIDColumn          string
	TreeTitleColumn       string
	TreeParentColumn      string
	TreeExpand            bool
	TreeAllowParentSelect bool
	SelectTableURL        string
	SelectTableTitle      string
	SelectTableDialogWidth string
	SelectTableDisplayField string
	SelectTableValueField string
	NestedFields          []*form.Field
	HasManyLabel          string
	HasManyTableMode      bool
	// Dynamic data for complex fields
	TableRows             []tableRowView
	HasManyRows           []hasManyRowView
	EmbedsValues          map[string]string
}

type tableRowView struct {
	RowIndex  int
	FieldName string
	RowValues map[string]string
}

type hasManyRowView struct {
	ItemIndex  int
	FieldName  string
	ItemLabel  string
	ItemValues map[string]string
}

type keyValueRowView struct {
	FieldName        string
	Index            int
	Key              string
	Value            string
	KeyPlaceholder   string
	ValuePlaceholder string
}

type repeaterChildView struct {
	Name        string
	Label       string
	Type        string
	Placeholder string
	Options     []gridOptionView
}

type repeaterRowView struct {
	Index  int
	Values map[string]string
}

type showView struct {
	Title       string
	Description string
	Items       []showItemView
	EditURL     string
	BackURL     string
}

type showItemView struct {
	Type  string
	Label string
	Title string
	Value template.HTML
}

type treeView struct {
	Title       string
	Description string
	EmptyText   string
	Nodes       []treeNodeView
}

type treeNodeView struct {
	ID          string
	ParentID    string
	Title       string
	Description string
	URL         string
	Children    []treeNodeView
}

func (a *App) buildGridColumns(r *http.Request, baseURL string, builder *grid.Builder, query ListQuery) []gridColumnView {
	var columns []gridColumnView
	for _, column := range builder.Columns {
		link := ""
		if column.Sortable {
			values := cloneQuery(r.URL.Query())
			values.Set("sort", column.Name)
			if query.Sort == column.Name && strings.EqualFold(query.Direction, "asc") {
				values.Set("direction", "desc")
			} else {
				values.Set("direction", "asc")
			}
			link = buildURL(baseURL, values)
		}
		columns = append(columns, gridColumnView{
			Label:    column.Label,
			Sortable: column.Sortable,
			SortURL:  link,
		})
	}
	return columns
}

func (a *App) buildTreeNodeViews(nodes []TreeNode, baseURL string) []treeNodeView {
	result := make([]treeNodeView, 0, len(nodes))
	for _, node := range nodes {
		viewNode := treeNodeView{
			ID:          node.ID,
			ParentID:    node.ParentID,
			Title:       node.Title,
			Description: node.Description,
			URL:         joinURL(baseURL, node.ID),
			Children:    a.buildTreeNodeViews(node.Children, baseURL),
		}
		result = append(result, viewNode)
	}
	return result
}

func (a *App) buildDashboardView(data DashboardData) dashboardView {
	view := dashboardView{Cards: data.Cards}
	for _, panel := range data.Panels {
		panelView := dashboardPanelView{
			Title:       panel.Title,
			Description: panel.Description,
			EmptyText:   fallback(panel.EmptyText, "No items yet."),
		}
		for _, action := range panel.Actions {
			if strings.TrimSpace(action.URL) == "" {
				continue
			}
			panelView.Actions = append(panelView.Actions, dashboardActionView{
				Label: action.Label,
				URL:   action.URL,
				Class: dashboardActionClass(action.Style),
			})
		}
		for _, item := range panel.Items {
			panelView.Items = append(panelView.Items, dashboardPanelItemView{
				Title:       item.Title,
				Description: item.Description,
				Value:       item.Value,
				URL:         item.URL,
				Tags:        append([]string(nil), item.Tags...),
			})
		}
		view.Panels = append(view.Panels, panelView)
	}
	return view
}

func (a *App) resourceHelper(resource Resource) *helperPanelView {
	return nil
}

func dashboardActionClass(style string) string {
	switch strings.ToLower(strings.TrimSpace(style)) {
	case "primary":
		return "btn btn-primary"
	case "danger":
		return "btn btn-danger"
	default:
		return "btn btn-ghost"
	}
}

func (a *App) buildGridRows(baseURL string, builder *grid.Builder, items []any) []gridRowView {
	rows := make([]gridRowView, 0, len(items))
	for _, item := range items {
		row := gridRowView{}
		for _, column := range builder.Columns {
			value := valueFromPath(item, column.Name)
			if column.Formatter != nil {
				row.Cells = append(row.Cells, column.Formatter(item, value))
				continue
			}
			row.Cells = append(row.Cells, toHTML(formatValue(value)))
		}
		id := formatValue(valueFromPath(item, "ID"))
		row.ID = id
		// Set row selector state
		if builder.RowSelectorChecked != nil {
			row.Checked = builder.RowSelectorChecked(item)
		}
		if builder.RowSelectorDisabled != nil {
			row.Disabled = builder.RowSelectorDisabled(item)
		}
		if !builder.DisableView {
			row.Actions = append(row.Actions, gridActionView{
				Label:  "View",
				URL:    joinURL(baseURL, id),
				Method: http.MethodGet,
				Class:  gridActionClass(grid.ActionGhost),
			})
		}
		if !builder.DisableEdit {
			row.Actions = append(row.Actions, gridActionView{
				Label:  "Edit",
				URL:    joinURL(baseURL, id, "edit"),
				Method: http.MethodGet,
				Class:  gridActionClass(grid.ActionGhost),
			})
		}
		if !builder.DisableDelete {
			row.Actions = append(row.Actions, gridActionView{
				Label:   "Delete",
				URL:     joinURL(baseURL, id, "delete"),
				Method:  http.MethodPost,
				Confirm: "Delete this record?",
				Class:   gridActionClass(grid.ActionDanger),
			})
		}
		for _, action := range builder.RowActions {
			if action == nil || action.URL == nil {
				continue
			}
			url := strings.TrimSpace(action.URL(item))
			if url == "" {
				continue
			}
			method := strings.ToUpper(strings.TrimSpace(action.Method))
			if method == "" {
				method = http.MethodGet
			}
			row.Actions = append(row.Actions, gridActionView{
				Label:   action.Label,
				URL:     url,
				Method:  method,
				Confirm: action.Confirm,
				Class:   gridActionClass(action.Style),
			})
		}
		rows = append(rows, row)
	}
	return rows
}

func (a *App) buildGridPageActions(baseURL string, builder *grid.Builder, resource Resource) []gridActionView {
	actions := make([]gridActionView, 0, len(builder.PageActions)+1)
	if !builder.DisableCreate && resource.BuildForm != nil {
		actions = append(actions, gridActionView{
			Label:  fallback(builder.CreateLabel, "Create"),
			URL:    joinURL(baseURL, "new"),
			Method: http.MethodGet,
			Class:  gridActionClass(grid.ActionPrimary),
		})
	}
	for _, action := range builder.PageActions {
		if action == nil || strings.TrimSpace(action.URL) == "" {
			continue
		}
		method := strings.ToUpper(strings.TrimSpace(action.Method))
		if method == "" {
			method = http.MethodGet
		}
		actions = append(actions, gridActionView{
			Label:   action.Label,
			URL:     action.URL,
			Method:  method,
			Confirm: action.Confirm,
			Class:   gridActionClass(action.Style),
		})
	}
	return actions
}

func (a *App) buildGridBatchActions(baseURL string, builder *grid.Builder) []gridActionView {
	actions := make([]gridActionView, 0, len(builder.BatchActions)+1)
	// Add default batch delete if not disabled
	if !builder.DisableBatchDelete {
		actions = append(actions, gridActionView{
			Label:   "Delete",
			URL:     joinURL(baseURL, "batch-delete"),
			Method:  http.MethodPost,
			Class:   gridActionClass(grid.ActionDanger),
			Confirm: "Are you sure you want to delete selected records?",
		})
	}
	for _, action := range builder.BatchActions {
		if action == nil {
			continue
		}
		actions = append(actions, gridActionView{
			Label:   action.Label,
			URL:     action.URL,
			Method:  action.Method,
			Confirm: action.Confirm,
			Class:   gridActionClass(action.Style),
		})
	}
	// Add AJAX batch action handlers
	for _, action := range builder.BatchActionHandlers {
		if action == nil {
			continue
		}
		actions = append(actions, gridActionView{
			Label:   action.Label,
			URL:     joinURL(baseURL, "batch-action", url.QueryEscape(action.Label)),
			Method:  action.Method,
			Confirm: action.Confirm,
			Class:   gridActionClass(action.Style),
		})
	}
	return actions
}

func (a *App) buildGridTools(builder *grid.Builder) []gridToolView {
	tools := make([]gridToolView, 0, len(builder.Tools))
	for _, tool := range builder.Tools {
		if tool == nil {
			continue
		}
		tools = append(tools, gridToolView{
			HTML:   tool.Render(),
			Script: tool.Script(),
		})
	}
	return tools
}

func gridActionClass(style grid.ActionStyle) string {
	switch style {
	case grid.ActionPrimary:
		return "btn btn-primary"
	case grid.ActionDanger:
		return "btn btn-danger"
	case grid.ActionGhost:
		return "btn btn-ghost"
	default:
		return "btn btn-ghost"
	}
}

func (a *App) buildGridFilters(builder *grid.Builder, query ListQuery) []gridFilterView {
	filters := make([]gridFilterView, 0, len(builder.Filters))
	for _, filter := range builder.Filters {
		view := gridFilterView{
			Name:     filter.Name,
			Label:    filter.Label,
			Kind:     string(filter.Kind),
			Value:    query.Filters[filter.Name],
			InputKey: "f_" + filter.Name,
		}
		for _, option := range filter.Options {
			view.Options = append(view.Options, gridOptionView{
				Label:    option.Label,
				Value:    option.Value,
				Selected: query.Filters[filter.Name] == option.Value,
			})
		}
		filters = append(filters, view)
	}
	return filters
}

func buildPagination(baseURL string, query ListQuery, total int64, page, perPage int) []paginationLink {
	totalPages := int((total + int64(perPage) - 1) / int64(perPage))
	if totalPages <= 1 {
		return nil
	}
	links := make([]paginationLink, 0, totalPages)
	for i := 1; i <= totalPages; i++ {
		values := url.Values{}
		values.Set("page", strconv.Itoa(i))
		values.Set("per_page", strconv.Itoa(perPage))
		if query.Search != "" {
			values.Set("q", query.Search)
		}
		if query.Sort != "" {
			values.Set("sort", query.Sort)
			values.Set("direction", query.Direction)
		}
		for key, value := range query.Filters {
			if value != "" {
				values.Set("f_"+key, value)
			}
		}
		links = append(links, paginationLink{
			Label:   strconv.Itoa(i),
			URL:     buildURL(baseURL, values),
			Active:  i == page,
			Current: i == page,
		})
	}
	return links
}

func (a *App) buildFormView(resource Resource, builder *form.Builder, record any, csrf, id string) formView {
	return a.buildFormViewState(resource, builder, record, csrf, id, nil, nil, "")
}

func (a *App) buildFormViewState(resource Resource, builder *form.Builder, record any, csrf, id string, submitted Values, fieldErrors map[string]string, message string) formView {
	baseURL := joinURL(a.cfg.Prefix, resource.Path)
	view := formView{
		Title:       fallback(builder.Title, resource.Title),
		Description: fallback(builder.Description, resource.Description),
		Message:     message,
		Action:      baseURL,
		Method:      http.MethodPost,
		Enctype:     "application/x-www-form-urlencoded",
		ShowDelete:  !builder.HideDelete && id != "",
		DeleteLabel: builder.DeleteLabel,
		SubmitLabel: builder.SubmitLabel,
		CSRF:        csrf,
		BackURL:     fallback(builder.CancelBackURL, baseURL),
	}
	if id != "" {
		view.Action = joinURL(baseURL, id)
		view.Method = http.MethodPut
		view.DeleteAction = joinURL(baseURL, id, "delete")
	}
	for _, field := range builder.Fields {
		entry := formFieldView{
			Name:              field.Name,
			SecondName:        field.SecondName,
			Label:             field.Label,
			Type:              string(field.Type),
			Multiple:          field.Multiple,
			Help:              field.Help,
			Required:          field.Required,
			Readonly:          field.Readonly,
			Disabled:          field.Disabled,
			Placeholder:       field.Placeholder,
			SecondPlaceholder: field.SecondPlaceholder,
			Error:             fieldErrors[field.Name],
			Min:               field.MinVal,
			Max:               field.MaxVal,
			Step:              field.StepVal,
			Inline:            field.IsInline,
			CurrencySymbol:    field.CurrencySymbol,
			EditorHeight:      field.EditorHeight,
			KeyPlaceholder:    field.KeyPlaceholder,
			ValuePlaceholder:  field.ValuePlaceholder,
			MaskFormat:        field.MaskFormat,
			SliderMin:         field.SliderMin,
			SliderMax:         field.SliderMax,
			SliderStep:        field.SliderStep,
			SliderPostfix:     field.SliderPostfix,
			AutocompleteURL:   field.AutocompleteURL,
			UploadDir:         field.UploadDir,
			UploadMaxCount:    field.UploadMaxCount,
			IsSortable:        field.IsSortable,
			HtmlContent:       field.HtmlContent,
		}
		if field.Type == form.FieldUpload {
			view.Enctype = "multipart/form-data"
		}
		if record != nil {
			valuePath := field.Name
			if field.ValuePath != "" {
				valuePath = field.ValuePath
			}
			if field.Type == form.FieldMulti {
				entry.Values = valueListFromPath(record, valuePath)
			} else if field.Type == form.FieldRepeater {
				entry.RepeaterFields = buildRepeaterChildViews(field)
				entry.RepeaterRows = repeaterRowsFromValue(valueFromPath(record, valuePath), field)
			} else if field.Type == form.FieldKeyValue {
				entry.KeyValueRows = keyValueRowsFromValue(valueFromPath(record, valuePath), field)
			} else {
				entry.Value = formatInputValue(valueFromPath(record, valuePath), string(field.Type))
				if field.Type == form.FieldSwitch {
					entry.Checked = entry.Value == "true" || entry.Value == "1"
				}
				if field.Type == form.FieldUpload {
					if field.Multiple {
						entry.FileURLs = a.uploadValuesFromRecord(record, field)
					} else {
						entry.FileURL = entry.Value
						entry.IsImage = isImagePath(entry.FileURL)
					}
				}
				if field.Type == form.FieldDateRange && field.SecondName != "" {
					secondValuePath := field.SecondName
					if field.SecondValuePath != "" {
						secondValuePath = field.SecondValuePath
					}
					entry.SecondValue = formatInputValue(valueFromPath(record, secondValuePath), "date")
				}
			}
		}
		if submitted != nil {
			if values, ok := submitted[field.Name]; ok {
				switch field.Type {
				case form.FieldMulti:
					entry.Values = append([]string(nil), values...)
				case form.FieldRepeater:
					entry.RepeaterFields = buildRepeaterChildViews(field)
					entry.RepeaterRows = repeaterRowsFromSubmitted(submitted, field)
				case form.FieldKeyValue:
					entry.KeyValueRows = keyValueRowsFromSubmitted(submitted, field)
				case form.FieldSwitch:
					entry.Checked = len(values) > 0 && (values[0] == "1" || strings.EqualFold(values[0], "true"))
				case form.FieldUpload:
					// Browsers cannot safely repopulate file inputs after failed submit.
				default:
					if len(values) > 0 {
						entry.Value = values[0]
					}
				}
			}
			if field.Type == form.FieldDateRange && field.SecondName != "" {
				if values, ok := submitted[field.SecondName]; ok && len(values) > 0 {
					entry.SecondValue = values[0]
				}
			}
		}
		if field.Type == form.FieldRepeater && len(entry.RepeaterFields) == 0 {
			entry.RepeaterFields = buildRepeaterChildViews(field)
			entry.RepeaterRows = repeaterRowsFromValue(nil, field)
		}
		if field.Type == form.FieldKeyValue && len(entry.KeyValueRows) == 0 {
			entry.KeyValueRows = []keyValueRowView{{FieldName: field.Name, Index: 0, Key: "", Value: ""}}
		}
		for _, option := range field.Options {
			selected := entry.Value == option.Value
			if field.Type == form.FieldMulti {
				selected = false
				for _, value := range entry.Values {
					if value == option.Value {
						selected = true
						break
					}
				}
			}
			entry.Options = append(entry.Options, gridOptionView{
				Label:    option.Label,
				Value:    option.Value,
				Selected: selected,
			})
		}
		view.Fields = append(view.Fields, entry)
	}
	return view
}

func (a *App) buildShowView(resource Resource, builder *show.Builder, record any, id string) showView {
	view := showView{
		Title:       fallback(builder.Title, resource.Title),
		Description: fallback(builder.Description, resource.Description),
		EditURL:     joinURL(a.cfg.Prefix, resource.Path, id, "edit"),
		BackURL:     joinURL(a.cfg.Prefix, resource.Path),
	}
	for _, item := range builder.Items {
		switch item.Type {
		case show.ItemDivider:
			view.Items = append(view.Items, showItemView{Type: "divider", Title: item.Title})
		default:
			value := valueFromPath(record, item.Name)
			rendered := toHTML(formatValue(value))
			if item.Formatter != nil {
				rendered = item.Formatter(record, value)
			}
			view.Items = append(view.Items, showItemView{
				Type:  "field",
				Label: item.Label,
				Value: rendered,
			})
		}
	}
	return view
}

func buildRepeaterChildViews(field *form.Field) []repeaterChildView {
	children := make([]repeaterChildView, 0, len(field.RepeaterFields))
	for _, child := range field.RepeaterFields {
		view := repeaterChildView{
			Name:        child.Name,
			Label:       child.Label,
			Type:        string(child.Type),
			Placeholder: child.Placeholder,
		}
		for _, option := range child.Options {
			view.Options = append(view.Options, gridOptionView{Label: option.Label, Value: option.Value})
		}
		children = append(children, view)
	}
	return children
}

func repeaterRowsFromSubmitted(submitted Values, field *form.Field) []repeaterRowView {
	rowsByIndex := map[int]map[string]string{}
	for _, child := range field.RepeaterFields {
		for index := 0; ; index++ {
			key := fmt.Sprintf("%s.%d.%s", field.Name, index, child.Name)
			values, ok := submitted[key]
			if !ok {
				if index > 20 {
					break
				}
				if _, exists := rowsByIndex[index]; !exists {
					continue
				}
				break
			}
			if _, ok := rowsByIndex[index]; !ok {
				rowsByIndex[index] = map[string]string{}
			}
			if len(values) > 0 {
				rowsByIndex[index][child.Name] = values[0]
			}
		}
	}
	return normalizeRepeaterRows(rowsByIndex, field.RepeaterMinRows)
}

func repeaterRowsFromValue(value any, field *form.Field) []repeaterRowView {
	rowsByIndex := map[int]map[string]string{}
	switch typed := value.(type) {
	case string:
		var decoded []map[string]string
		if strings.TrimSpace(typed) != "" && json.Unmarshal([]byte(typed), &decoded) == nil {
			for i, row := range decoded {
				rowsByIndex[i] = row
			}
		}
	case []map[string]string:
		for i, row := range typed {
			rowsByIndex[i] = row
		}
	default:
		rv := reflect.ValueOf(value)
		for rv.IsValid() && rv.Kind() == reflect.Pointer {
			if rv.IsNil() {
				return normalizeRepeaterRows(rowsByIndex, field.RepeaterMinRows)
			}
			rv = rv.Elem()
		}
		if rv.IsValid() && rv.Kind() == reflect.Slice {
			for i := 0; i < rv.Len(); i++ {
				item := rv.Index(i).Interface()
				row := map[string]string{}
				for _, child := range field.RepeaterFields {
					row[child.Name] = formatValue(valueFromPath(item, child.Name))
				}
				rowsByIndex[i] = row
			}
		}
	}
	return normalizeRepeaterRows(rowsByIndex, field.RepeaterMinRows)
}

func normalizeRepeaterRows(rowsByIndex map[int]map[string]string, minRows int) []repeaterRowView {
	maxRows := minRows
	for index := range rowsByIndex {
		if index+2 > maxRows {
			maxRows = index + 2
		}
	}
	if maxRows < 1 {
		maxRows = 1
	}
	rows := make([]repeaterRowView, 0, maxRows)
	for i := 0; i < maxRows; i++ {
		row := rowsByIndex[i]
		if row == nil {
			row = map[string]string{}
		}
		rows = append(rows, repeaterRowView{Index: i, Values: row})
	}
	return rows
}

// KeyValue helper functions
func keyValueRowsFromValue(value any, field *form.Field) []keyValueRowView {
	rows := []keyValueRowView{}
	if value == nil {
		return rows
	}
	// Parse JSON value
	var data []map[string]string
	switch v := value.(type) {
	case string:
		if v == "" {
			return rows
		}
		if err := json.Unmarshal([]byte(v), &data); err != nil {
			return rows
		}
	case []map[string]string:
		data = v
	case []any:
		for _, item := range v {
			if m, ok := item.(map[string]any); ok {
				row := map[string]string{}
				for key, val := range m {
					row[key] = fmt.Sprint(val)
				}
				data = append(data, row)
			}
		}
	}
	for i, item := range data {
		rows = append(rows, keyValueRowView{
			FieldName:        field.Name,
			Index:            i,
			Key:              item["key"],
			Value:            item["value"],
			KeyPlaceholder:   field.KeyPlaceholder,
			ValuePlaceholder: field.ValuePlaceholder,
		})
	}
	return rows
}

func keyValueRowsFromSubmitted(submitted map[string][]string, field *form.Field) []keyValueRowView {
	rows := []keyValueRowView{}
	// Find all keys with pattern: fieldName[index][key] and fieldName[index][value]
	keyPattern := regexp.MustCompile(`^` + regexp.QuoteMeta(field.Name) + `\[(\d+)\]\[key\]$`)

	indices := make(map[int]bool)
	for key := range submitted {
		if matches := keyPattern.FindStringSubmatch(key); matches != nil {
			if idx, err := strconv.Atoi(matches[1]); err == nil {
				indices[idx] = true
			}
		}
	}

	for idx := range indices {
		keyName := fmt.Sprintf("%s[%d][key]", field.Name, idx)
		valueName := fmt.Sprintf("%s[%d][value]", field.Name, idx)

		keyVal := ""
		if vals, ok := submitted[keyName]; ok && len(vals) > 0 {
			keyVal = vals[0]
		}

		valueVal := ""
		if vals, ok := submitted[valueName]; ok && len(vals) > 0 {
			valueVal = vals[0]
		}

		rows = append(rows, keyValueRowView{
			FieldName:        field.Name,
			Index:            idx,
			Key:              keyVal,
			Value:            valueVal,
			KeyPlaceholder:   field.KeyPlaceholder,
			ValuePlaceholder: field.ValuePlaceholder,
		})
	}

	// Sort by index
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Index < rows[j].Index
	})

	return rows
}

func collectKeyValueValue(formValues map[string][]string, field *form.Field) (string, error) {
	rows := make([]map[string]string, 0)

	keyPattern := regexp.MustCompile(`^` + regexp.QuoteMeta(field.Name) + `\[(\d+)\]\[key\]$`)
	indices := make(map[int]bool)

	for key := range formValues {
		if matches := keyPattern.FindStringSubmatch(key); matches != nil {
			if idx, err := strconv.Atoi(matches[1]); err == nil {
				indices[idx] = true
			}
		}
	}

	for idx := range indices {
		keyName := fmt.Sprintf("%s[%d][key]", field.Name, idx)
		valueName := fmt.Sprintf("%s[%d][value]", field.Name, idx)

		keyVal := ""
		if vals, ok := formValues[keyName]; ok && len(vals) > 0 {
			keyVal = vals[0]
		}

		valueVal := ""
		if vals, ok := formValues[valueName]; ok && len(vals) > 0 {
			valueVal = vals[0]
		}

		// Skip empty rows
		if keyVal != "" || valueVal != "" {
			rows = append(rows, map[string]string{
				"key":   keyVal,
				"value": valueVal,
			})
		}
	}

	if len(rows) == 0 {
		return "", nil
	}

	jsonBytes, err := json.Marshal(rows)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func collectRepeaterValue(formValues map[string][]string, field *form.Field) (string, error) {
	rows := make([]map[string]string, 0)
	limit := 12
	if field.RepeaterMinRows > limit {
		limit = field.RepeaterMinRows
	}
	for index := 0; index < limit; index++ {
		row := map[string]string{}
		nonEmpty := false
		for _, child := range field.RepeaterFields {
			key := fmt.Sprintf("%s.%d.%s", field.Name, index, child.Name)
			values := formValues[key]
			if len(values) == 0 {
				continue
			}
			value := strings.TrimSpace(values[0])
			if value != "" {
				nonEmpty = true
			}
			row[child.Name] = value
		}
		if nonEmpty {
			rows = append(rows, row)
		}
	}
	if len(rows) == 0 {
		return "", nil
	}
	encoded, err := json.Marshal(rows)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

func (a *App) menuView(items []NavigationItem, currentPath string) []menuItemView {
	out := make([]menuItemView, 0, len(items))
	for _, item := range items {
		url := ""
		if item.URI != "" {
			if item.URI == "/" {
				url = normalizePath(a.cfg.Prefix)
			} else if strings.HasPrefix(item.URI, "/") {
				url = normalizePath(item.URI)
			} else {
				url = joinURL(a.cfg.Prefix, item.URI)
			}
		}
		children := a.menuView(item.Children, currentPath)
		prefix := normalizePath(a.cfg.Prefix)
		// Active if: exact URL match, or current path is a sub-path of this item (but not for root/dashboard)
		active := false
		if url != "" {
			if currentPath == url {
				active = true
			} else if url != prefix {
				// For non-root items, also match if current path is a child path
				active = strings.HasPrefix(currentPath, strings.TrimRight(url, "/")+"/")
			}
		}
		// Expanded if this item is active or any child is active/expanded
		expanded := active
		for _, child := range children {
			if child.Active || child.Expanded {
				expanded = true
				break
			}
		}
		out = append(out, menuItemView{
			Title:    item.Title,
			Icon:     item.Icon,
			URL:      url,
			Active:   active,
			Expanded: expanded,
			Children: children,
		})
	}
	return out
}

func cloneQuery(values url.Values) url.Values {
	cloned := url.Values{}
	for key, list := range values {
		for _, value := range list {
			cloned.Add(key, value)
		}
	}
	return cloned
}

func routeAllows(route Route, method string) bool {
	if len(route.Methods) == 0 {
		return true
	}
	for _, allowed := range route.Methods {
		if strings.EqualFold(strings.TrimSpace(allowed), method) {
			return true
		}
	}
	return false
}

func (a *App) handleToolForm(w http.ResponseWriter, r *http.Request, state *sessionState, identity *Identity, formWidget *widgetform.ToolForm) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		// Render form
		renderCtx := formWidget.Render(state.CSRF, a.cfg.Prefix)
		data := map[string]interface{}{
			"Form": renderCtx,
		}
		// Check if async request
		if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if err := a.templates.ExecuteTemplate(w, "toolform", data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		// Render within shell
		a.renderShell(w, r, state, identity, "form", pageData{
			PageTitle: formWidget.Builder().Title,
			Content: toolFormContent{
				Form:  data,
				Title: formWidget.Builder().Title,
			},
		})

	case http.MethodPost:
		// Process form submission
		if err := r.ParseForm(); err != nil {
			a.writeError(w, http.StatusBadRequest, err)
			return
		}
		if r.FormValue("_csrf") != state.CSRF {
			a.writeError(w, http.StatusBadRequest, errors.New("invalid csrf token"))
			return
		}

		resp, err := formWidget.Process(ctx, r.PostForm)
		if err != nil {
			a.writeError(w, http.StatusInternalServerError, err)
			return
		}

		// Check if AJAX request
		if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}

		// HTML response - render form with message
		renderCtx := formWidget.Render(state.CSRF, a.cfg.Prefix)
		data := map[string]interface{}{
			"Form":    renderCtx,
			"Message": resp.Message,
			"Success": resp.Success,
		}

		a.renderShell(w, r, state, identity, "form", pageData{
			PageTitle: formWidget.Builder().Title,
			Message:   resp.Message,
			Content: toolFormContent{
				Form:    data,
				Title:   formWidget.Builder().Title,
				Success: resp.Success,
				Message: resp.Message,
			},
		})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// toolFormContent is the view data for tool form rendering.
type toolFormContent struct {
	Form    map[string]interface{}
	Title   string
	Success bool
	Message string
}
