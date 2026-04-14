package gormstore

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/zhenyangze/goadmin"
	"gorm.io/gorm"
)

type RelationUpdateStrategy string

const (
	RelationReplace RelationUpdateStrategy = "replace"
	RelationAppend  RelationUpdateStrategy = "append"
	RelationMerge   RelationUpdateStrategy = "merge"
)

// RelationHandler is the extension API used by RepositoryWithRelations.
type RelationHandler[TParent any] interface {
	FormKey() string
	ApplyCreate(context.Context, *gorm.DB, *TParent, goadmin.Values) error
	ApplyUpdate(context.Context, *gorm.DB, *TParent, goadmin.Values) error
	ApplyDelete(context.Context, *gorm.DB, *TParent) error
}

// RepeaterRelationConfig wires one relation-backed repeater field to a child model.
type RepeaterRelationConfig[TChild any] struct {
	FormKey          string
	ForeignKeyField  string
	ForeignKeyColumn string
	SortField        string
	IDField          string
	UpdateStrategy   RelationUpdateStrategy
	DeleteMissing    bool
	HardDelete       bool
}

type RepositoryWithRelations[TParent any] struct {
	*Repository[TParent]
	Relations []RelationHandler[TParent]
}

func NewWithRelations[TParent any](db *gorm.DB) *RepositoryWithRelations[TParent] {
	return &RepositoryWithRelations[TParent]{Repository: New[TParent](db)}
}

func (r *RepositoryWithRelations[TParent]) AddRelation(handler RelationHandler[TParent]) *RepositoryWithRelations[TParent] {
	r.Relations = append(r.Relations, handler)
	return r
}

func (r *RepositoryWithRelations[TParent]) Create(ctx context.Context, values goadmin.Values) error {
	parentValues := cloneWithoutRelationKeys(values, r.Relations)
	var item TParent
	if err := r.assignValues(&item, parentValues); err != nil {
		return err
	}
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := r.runHooks(ctx, tx, &item, parentValues, "", r.Hooks.BeforeCreate); err != nil {
			return err
		}
		if err := tx.Create(&item).Error; err != nil {
			return err
		}
		for _, relation := range r.Relations {
			if err := relation.ApplyCreate(ctx, tx, &item, values); err != nil {
				return err
			}
		}
		return r.runHooks(ctx, tx, &item, parentValues, "", r.Hooks.AfterCreate)
	})
}

func (r *RepositoryWithRelations[TParent]) Update(ctx context.Context, id string, values goadmin.Values) error {
	parentValues := cloneWithoutRelationKeys(values, r.Relations)
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var item TParent
		if err := tx.First(&item, id).Error; err != nil {
			return err
		}
		if err := r.assignValues(&item, parentValues); err != nil {
			return err
		}
		if err := r.runHooks(ctx, tx, &item, parentValues, id, r.Hooks.BeforeUpdate); err != nil {
			return err
		}
		if err := tx.Save(&item).Error; err != nil {
			return err
		}
		for _, relation := range r.Relations {
			if err := relation.ApplyUpdate(ctx, tx, &item, values); err != nil {
				return err
			}
		}
		return r.runHooks(ctx, tx, &item, parentValues, id, r.Hooks.AfterUpdate)
	})
}

func (r *RepositoryWithRelations[TParent]) Delete(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var item TParent
		if err := tx.First(&item, id).Error; err != nil {
			return err
		}
		if err := r.runHooks(ctx, tx, &item, nil, id, r.Hooks.BeforeDelete); err != nil {
			return err
		}
		for _, relation := range r.Relations {
			if err := relation.ApplyDelete(ctx, tx, &item); err != nil {
				return err
			}
		}
		if err := tx.Delete(&item).Error; err != nil {
			return err
		}
		return r.runHooks(ctx, tx, &item, nil, id, r.Hooks.AfterDelete)
	})
}

type RepeaterRelation[TParent any, TChild any] struct {
	Config RepeaterRelationConfig[TChild]
}

func NewRepeaterRelation[TParent any, TChild any](cfg RepeaterRelationConfig[TChild]) *RepeaterRelation[TParent, TChild] {
	if cfg.IDField == "" {
		cfg.IDField = "ID"
	}
	if cfg.UpdateStrategy == "" {
		cfg.UpdateStrategy = RelationReplace
	}
	return &RepeaterRelation[TParent, TChild]{Config: cfg}
}

func NewWithRepeater[TParent any, TChild any](db *gorm.DB, cfg RepeaterRelationConfig[TChild]) *RepositoryWithRelations[TParent] {
	return NewWithRelations[TParent](db).AddRelation(NewRepeaterRelation[TParent, TChild](cfg))
}

func (r *RepeaterRelation[TParent, TChild]) FormKey() string {
	return r.Config.FormKey
}

func (r *RepeaterRelation[TParent, TChild]) ApplyCreate(_ context.Context, tx *gorm.DB, parent *TParent, values goadmin.Values) error {
	rows, err := decodeRepeaterRows(values.First(r.Config.FormKey))
	if err != nil {
		return err
	}
	children, err := r.buildChildren(rows, primaryUint(*parent))
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return tx.Create(&children).Error
	}
	return nil
}

func (r *RepeaterRelation[TParent, TChild]) ApplyUpdate(_ context.Context, tx *gorm.DB, parent *TParent, values goadmin.Values) error {
	if _, ok := values[r.Config.FormKey]; !ok {
		return nil
	}
	rows, err := decodeRepeaterRows(values.First(r.Config.FormKey))
	if err != nil {
		return err
	}
	parentID := primaryUint(*parent)
	existing, err := r.loadChildren(tx, parentID)
	if err != nil {
		return err
	}
	switch r.Config.UpdateStrategy {
	case RelationAppend:
		children, err := r.buildChildren(rows, parentID)
		if err != nil {
			return err
		}
		if len(children) > 0 {
			return tx.Create(&children).Error
		}
		return nil
	case RelationMerge:
		return r.mergeChildren(tx, parentID, existing, rows)
	default:
		if err := r.deleteChildren(tx, existing); err != nil {
			return err
		}
		children, err := r.buildChildren(rows, parentID)
		if err != nil {
			return err
		}
		if len(children) > 0 {
			return tx.Create(&children).Error
		}
		return nil
	}
}

func (r *RepeaterRelation[TParent, TChild]) ApplyDelete(_ context.Context, tx *gorm.DB, parent *TParent) error {
	existing, err := r.loadChildren(tx, primaryUint(*parent))
	if err != nil {
		return err
	}
	return r.deleteChildren(tx, existing)
}

func (r *RepeaterRelation[TParent, TChild]) mergeChildren(tx *gorm.DB, parentID uint, existing []TChild, rows []map[string]string) error {
	byID := map[string]TChild{}
	for _, child := range existing {
		byID[childIDString(child, r.Config.IDField)] = child
	}
	seen := map[string]bool{}
	for index, row := range rows {
		id := strings.TrimSpace(row[r.Config.IDField])
		if child, ok := byID[id]; ok && id != "" {
			if err := r.applyRow(&child, row, parentID, index); err != nil {
				return err
			}
			if err := tx.Save(&child).Error; err != nil {
				return err
			}
			seen[id] = true
			continue
		}
		children, err := r.buildChildren([]map[string]string{row}, parentID)
		if err != nil {
			return err
		}
		if len(children) > 0 {
			if err := tx.Create(&children).Error; err != nil {
				return err
			}
		}
	}
	if r.Config.DeleteMissing {
		for _, child := range existing {
			id := childIDString(child, r.Config.IDField)
			if id != "" && !seen[id] {
				if err := r.deleteOne(tx, &child); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (r *RepeaterRelation[TParent, TChild]) loadChildren(tx *gorm.DB, parentID uint) ([]TChild, error) {
	var children []TChild
	query := tx.Where(r.foreignKeyColumn()+" = ?", parentID)
	if r.Config.SortField != "" {
		query = query.Order(toDBName(r.Config.SortField) + " asc")
	}
	if err := query.Find(&children).Error; err != nil {
		return nil, err
	}
	return children, nil
}

func (r *RepeaterRelation[TParent, TChild]) deleteChildren(tx *gorm.DB, children []TChild) error {
	for i := range children {
		if err := r.deleteOne(tx, &children[i]); err != nil {
			return err
		}
	}
	return nil
}

func (r *RepeaterRelation[TParent, TChild]) deleteOne(tx *gorm.DB, child *TChild) error {
	if r.Config.HardDelete {
		return tx.Unscoped().Delete(child).Error
	}
	return tx.Delete(child).Error
}

func (r *RepeaterRelation[TParent, TChild]) foreignKeyColumn() string {
	if r.Config.ForeignKeyColumn != "" {
		return r.Config.ForeignKeyColumn
	}
	return toDBName(r.Config.ForeignKeyField)
}

func (r *RepeaterRelation[TParent, TChild]) buildChildren(rows []map[string]string, parentID uint) ([]TChild, error) {
	children := make([]TChild, 0, len(rows))
	for index, row := range rows {
		var child TChild
		if err := r.applyRow(&child, row, parentID, index); err != nil {
			return nil, err
		}
		children = append(children, child)
	}
	return children, nil
}

func (r *RepeaterRelation[TParent, TChild]) applyRow(child *TChild, row map[string]string, parentID uint, index int) error {
	root := reflect.ValueOf(child).Elem()
	for key, value := range row {
		if key == r.Config.IDField {
			continue
		}
		if value == "" {
			continue
		}
		if err := setString(root, key, value); err != nil {
			return err
		}
	}
	if err := setUintField(root, r.Config.ForeignKeyField, uint64(parentID)); err != nil {
		return err
	}
	if r.Config.SortField != "" {
		if err := setIntField(root, r.Config.SortField, int64(index+1)); err != nil {
			return err
		}
	}
	return nil
}

func decodeRepeaterRows(raw string) ([]map[string]string, error) {
	if raw == "" {
		return nil, nil
	}
	var rows []map[string]string
	if err := json.Unmarshal([]byte(raw), &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func cloneWithoutRelationKeys[TParent any](values goadmin.Values, handlers []RelationHandler[TParent]) goadmin.Values {
	cloned := goadmin.Values{}
	skip := map[string]bool{}
	for _, handler := range handlers {
		skip[handler.FormKey()] = true
	}
	for name, list := range values {
		if skip[name] {
			continue
		}
		cloned[name] = append([]string(nil), list...)
	}
	return cloned
}

func primaryUint[T any](value T) uint {
	rv := reflect.ValueOf(value)
	for rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return 0
	}
	id := rv.FieldByName("ID")
	if !id.IsValid() {
		return 0
	}
	switch id.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uint(id.Uint())
	default:
		return 0
	}
}

func childIDString[T any](value T, field string) string {
	rv := reflect.ValueOf(value)
	for rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	target := rv.FieldByName(field)
	if !target.IsValid() {
		return ""
	}
	switch target.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.String:
		return fmt.Sprint(target.Interface())
	default:
		return fmt.Sprint(target.Interface())
	}
}

func setUintField(root reflect.Value, name string, value uint64) error {
	field := indirectField(root, name)
	if !field.IsValid() || !field.CanSet() {
		return nil
	}
	if field.Kind() >= reflect.Uint && field.Kind() <= reflect.Uint64 {
		field.SetUint(value)
	}
	return nil
}

func setIntField(root reflect.Value, name string, value int64) error {
	field := indirectField(root, name)
	if !field.IsValid() || !field.CanSet() {
		return nil
	}
	if field.Kind() >= reflect.Int && field.Kind() <= reflect.Int64 {
		field.SetInt(value)
	}
	return nil
}
