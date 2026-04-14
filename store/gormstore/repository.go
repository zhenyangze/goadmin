package gormstore

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/zhenyangze/goadmin"
	"gorm.io/gorm"
)

// TreeConfig enables tree rendering on a generic GORM model.
type TreeConfig struct {
	IDField          string
	ParentField      string
	TitleField       string
	DescriptionField string
	OrderField       string
}

// Repository provides generic CRUD and tree access for a GORM model.
type Repository[T any] struct {
	DB           *gorm.DB
	SearchFields []string
	FilterFields []string
	Preloads     []string
	DefaultOrder string
	Mutators     map[string]func(string) (any, error)
	TreeConfig   *TreeConfig
	Hooks        Hooks[T]
	StdHooks     []goadmin.RepositoryHook // Standard hooks interface
	resourceName string
}

// AddHook implements goadmin.HookableRepository.
func (r *Repository[T]) AddHook(hook goadmin.RepositoryHook) {
	r.StdHooks = append(r.StdHooks, hook)
}

// SetResourceName implements goadmin.HookableRepository.
func (r *Repository[T]) SetResourceName(name string) {
	r.resourceName = name
}

// runStdHooks executes standard RepositoryHook methods.
func (r *Repository[T]) runStdHooks(ctx context.Context, method func(hook goadmin.RepositoryHook) error) error {
	for _, hook := range r.StdHooks {
		if err := method(hook); err != nil {
			return err
		}
	}
	return nil
}

type HookFunc[T any] func(context.Context, HookContext[T]) error

type HookContext[T any] struct {
	Tx     *gorm.DB
	Item   *T
	Values goadmin.Values
	ID     string
}

type Hooks[T any] struct {
	BeforeCreate []HookFunc[T]
	AfterCreate  []HookFunc[T]
	BeforeUpdate []HookFunc[T]
	AfterUpdate  []HookFunc[T]
	BeforeDelete []HookFunc[T]
	AfterDelete  []HookFunc[T]
}

// New creates a generic GORM-backed repository.
func New[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{
		DB:       db,
		Mutators: map[string]func(string) (any, error){},
	}
}

// List implements goadmin.Repository.
func (r *Repository[T]) List(ctx context.Context, query goadmin.ListQuery) (goadmin.ListResult, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PerPage < 1 {
		query.PerPage = 10
	}

	db := r.DB.WithContext(ctx).Model(new(T))
	for _, preload := range r.Preloads {
		db = db.Preload(preload)
	}

	if query.Search != "" && len(r.SearchFields) > 0 {
		like := "%" + query.Search + "%"
		db = db.Where(r.searchClause(), repeatArgs(like, len(r.SearchFields))...)
	}

	allowedFilters := make(map[string]bool, len(r.FilterFields))
	for _, field := range r.FilterFields {
		allowedFilters[field] = true
	}
	for field, value := range query.Filters {
		if value == "" {
			continue
		}
		if len(allowedFilters) > 0 && !allowedFilters[field] {
			continue
		}
		db = db.Where(toDBName(field)+" = ?", value)
	}

	order := r.DefaultOrder
	if query.Sort != "" {
		dir := "asc"
		if strings.EqualFold(query.Direction, "desc") {
			dir = "desc"
		}
		order = fmt.Sprintf("%s %s", toDBName(query.Sort), dir)
	}
	if order != "" {
		db = db.Order(order)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return goadmin.ListResult{}, err
	}

	var items []T
	if err := db.Offset((query.Page - 1) * query.PerPage).Limit(query.PerPage).Find(&items).Error; err != nil {
		return goadmin.ListResult{}, err
	}

	records := make([]any, 0, len(items))
	for i := range items {
		records = append(records, items[i])
	}

	return goadmin.ListResult{
		Items:   records,
		Total:   total,
		Page:    query.Page,
		PerPage: query.PerPage,
	}, nil
}

// Get fetches one record by primary key.
func (r *Repository[T]) Get(ctx context.Context, id string) (any, error) {
	var item T
	db := r.DB.WithContext(ctx)
	for _, preload := range r.Preloads {
		db = db.Preload(preload)
	}
	if err := db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return item, nil
}

// Create inserts a record using form values.
func (r *Repository[T]) Create(ctx context.Context, values goadmin.Values) error {
	var item T
	if err := r.assignValues(&item, values); err != nil {
		return err
	}

	// Run standard BeforeCreate hooks
	if err := r.runStdHooks(ctx, func(hook goadmin.RepositoryHook) error {
		return hook.BeforeCreate(ctx, values)
	}); err != nil {
		return err
	}

	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := r.runHooks(ctx, tx, &item, values, "", r.Hooks.BeforeCreate); err != nil {
			return err
		}
		if err := tx.Create(&item).Error; err != nil {
			return err
		}
		if err := r.runHooks(ctx, tx, &item, values, "", r.Hooks.AfterCreate); err != nil {
			return err
		}

		// Run standard AfterCreate hooks (with the ID)
		id := fmt.Sprintf("%v", reflect.ValueOf(item).FieldByName("ID").Interface())
		return r.runStdHooks(ctx, func(hook goadmin.RepositoryHook) error {
			return hook.AfterCreate(ctx, id, values)
		})
	})
}

// Update mutates an existing record.
func (r *Repository[T]) Update(ctx context.Context, id string, values goadmin.Values) error {
	var item T
	if err := r.DB.WithContext(ctx).First(&item, id).Error; err != nil {
		return err
	}

	// Keep a copy of old item for hooks
	oldItem := item

	if err := r.assignValues(&item, values); err != nil {
		return err
	}

	// Run standard BeforeUpdate hooks
	if err := r.runStdHooks(ctx, func(hook goadmin.RepositoryHook) error {
		return hook.BeforeUpdate(ctx, id, values)
	}); err != nil {
		return err
	}

	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := r.runHooks(ctx, tx, &item, values, id, r.Hooks.BeforeUpdate); err != nil {
			return err
		}
		if err := tx.Save(&item).Error; err != nil {
			return err
		}
		if err := r.runHooks(ctx, tx, &item, values, id, r.Hooks.AfterUpdate); err != nil {
			return err
		}

		// Run standard AfterUpdate hooks
		return r.runStdHooks(ctx, func(hook goadmin.RepositoryHook) error {
			return hook.AfterUpdate(ctx, id, values, oldItem)
		})
	})
}

// Delete removes a record.
func (r *Repository[T]) Delete(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var item T
		if err := tx.First(&item, id).Error; err != nil {
			return err
		}

		// Run standard BeforeDelete hooks
		if err := r.runStdHooks(ctx, func(hook goadmin.RepositoryHook) error {
			return hook.BeforeDelete(ctx, id)
		}); err != nil {
			return err
		}

		if err := r.runHooks(ctx, tx, &item, nil, id, r.Hooks.BeforeDelete); err != nil {
			return err
		}
		if err := tx.Delete(&item).Error; err != nil {
			return err
		}
		if err := r.runHooks(ctx, tx, &item, nil, id, r.Hooks.AfterDelete); err != nil {
			return err
		}

		// Run standard AfterDelete hooks
		return r.runStdHooks(ctx, func(hook goadmin.RepositoryHook) error {
			return hook.AfterDelete(ctx, id, item)
		})
	})
}

// Tree renders the model as a nested hierarchy when configured.
func (r *Repository[T]) Tree(ctx context.Context) ([]goadmin.TreeNode, error) {
	if r.TreeConfig == nil {
		return nil, errors.New("tree is not configured")
	}

	db := r.DB.WithContext(ctx).Model(new(T))
	if r.TreeConfig.OrderField != "" {
		db = db.Order(r.TreeConfig.OrderField + " asc")
	}

	var items []T
	if err := db.Find(&items).Error; err != nil {
		return nil, err
	}

	nodesByParent := map[string][]goadmin.TreeNode{}
	for i := range items {
		record := reflect.ValueOf(items[i])
		nodesByParent[valueAsString(record, r.TreeConfig.ParentField)] = append(
			nodesByParent[valueAsString(record, r.TreeConfig.ParentField)],
			goadmin.TreeNode{
				ID:          valueAsString(record, r.TreeConfig.IDField),
				ParentID:    valueAsString(record, r.TreeConfig.ParentField),
				Title:       valueAsString(record, r.TreeConfig.TitleField),
				Description: valueAsString(record, r.TreeConfig.DescriptionField),
			},
		)
	}

	return attachTreeChildren(nodesByParent, "0"), nil
}

func (r *Repository[T]) searchClause() string {
	parts := make([]string, 0, len(r.SearchFields))
	for _, field := range r.SearchFields {
		parts = append(parts, toDBName(field)+" LIKE ?")
	}
	return strings.Join(parts, " OR ")
}

func (r *Repository[T]) assignValues(target *T, values goadmin.Values) error {
	root := reflect.ValueOf(target).Elem()
	for name, list := range values {
		if mutator, ok := r.Mutators[name]; ok {
			raw := ""
			if len(list) > 0 {
				raw = list[0]
			}
			if raw == "" {
				continue
			}
			value, err := mutator(raw)
			if err != nil {
				return err
			}
			if err := setField(root, name, reflect.ValueOf(value)); err != nil {
				return err
			}
			continue
		}
		if len(list) > 1 {
			field := indirectField(root, name)
			if field.IsValid() && field.CanSet() && field.Kind() == reflect.String {
				field.SetString(strings.Join(list, ","))
				continue
			}
		}
		raw := ""
		if len(list) > 0 {
			raw = list[0]
		}
		if raw == "" {
			continue
		}
		if err := setString(root, name, raw); err != nil {
			return err
		}
	}
	return nil
}

func setString(root reflect.Value, name, raw string) error {
	field := indirectField(root, name)
	if !field.IsValid() || !field.CanSet() {
		return fmt.Errorf("unknown field %q", name)
	}
	switch field.Kind() {
	case reflect.String:
		field.SetString(raw)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Type().PkgPath() == "time" && field.Type().Name() == "Time" {
			ts, err := time.Parse("2006-01-02 15:04", raw)
			if err != nil {
				ts, err = time.Parse("2006-01-02T15:04", raw)
			}
			if err != nil {
				ts, err = time.Parse("2006-01-02", raw)
			}
			if err != nil {
				ts, err = time.Parse(time.RFC3339, raw)
			}
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(ts))
			return nil
		}
		value, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(value)
	case reflect.Bool:
		field.SetBool(raw == "1" || strings.EqualFold(raw, "true") || strings.EqualFold(raw, "on"))
	case reflect.Float32, reflect.Float64:
		value, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return err
		}
		field.SetFloat(value)
	default:
		if field.Type().PkgPath() == "time" && field.Type().Name() == "Time" {
			ts, err := time.Parse("2006-01-02 15:04", raw)
			if err != nil {
				ts, err = time.Parse("2006-01-02T15:04", raw)
			}
			if err != nil {
				ts, err = time.Parse("2006-01-02", raw)
			}
			if err != nil {
				ts, err = time.Parse(time.RFC3339, raw)
			}
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(ts))
			return nil
		}
		return fmt.Errorf("unsupported field type for %q", name)
	}
	return nil
}

func toDBName(field string) string {
	if field == "" {
		return ""
	}
	var out []rune
	for i, r := range field {
		if unicode.IsUpper(r) {
			if i > 0 {
				out = append(out, '_')
			}
			out = append(out, unicode.ToLower(r))
			continue
		}
		out = append(out, r)
	}
	return string(out)
}

func setField(root reflect.Value, name string, value reflect.Value) error {
	field := indirectField(root, name)
	if !field.IsValid() || !field.CanSet() {
		return fmt.Errorf("unknown field %q", name)
	}
	if value.Type().AssignableTo(field.Type()) {
		field.Set(value)
		return nil
	}
	if value.Type().ConvertibleTo(field.Type()) {
		field.Set(value.Convert(field.Type()))
		return nil
	}
	return fmt.Errorf("cannot assign %s to %s", value.Type(), field.Type())
}

func indirectField(root reflect.Value, name string) reflect.Value {
	field := root
	if field.Kind() == reflect.Pointer {
		field = field.Elem()
	}
	return field.FieldByName(name)
}

func valueAsString(value reflect.Value, field string) string {
	target := indirectField(value, field)
	if !target.IsValid() {
		return ""
	}
	return fmt.Sprint(target.Interface())
}

func attachTreeChildren(itemsByParent map[string][]goadmin.TreeNode, parentID string) []goadmin.TreeNode {
	nodes := itemsByParent[parentID]
	for i := range nodes {
		nodes[i].Children = attachTreeChildren(itemsByParent, nodes[i].ID)
	}
	return nodes
}

func repeatArgs(arg string, count int) []any {
	values := make([]any, 0, count)
	for i := 0; i < count; i++ {
		values = append(values, arg)
	}
	return values
}

func (r *Repository[T]) runHooks(ctx context.Context, tx *gorm.DB, item *T, values goadmin.Values, id string, hooks []HookFunc[T]) error {
	for _, hook := range hooks {
		if hook == nil {
			continue
		}
		if err := hook(ctx, HookContext[T]{
			Tx:     tx,
			Item:   item,
			Values: values,
			ID:     id,
		}); err != nil {
			return err
		}
	}
	return nil
}
