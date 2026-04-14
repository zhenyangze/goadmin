package gormstore

import (
	"context"
	"testing"

	"github.com/zhenyangze/goadmin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type hookParent struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Items []hookChild
}

type hookChild struct {
	ID           uint `gorm:"primaryKey"`
	HookParentID uint
	Label        string
	Sort         int
}

func TestRepositoryWithRelationsHooksAndMultipleRepeaters(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:hooks-test?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&hookParent{}, &hookChild{}); err != nil {
		t.Fatal(err)
	}

	repo := NewWithRelations[hookParent](db).
		AddRelation(NewRepeaterRelation[hookParent, hookChild](RepeaterRelationConfig[hookChild]{
			FormKey:          "Items",
			ForeignKeyField:  "HookParentID",
			ForeignKeyColumn: "hook_parent_id",
			SortField:        "Sort",
			IDField:          "ID",
			UpdateStrategy:   RelationMerge,
			DeleteMissing:    false,
		}))

	var calls []string
	repo.Hooks.BeforeCreate = append(repo.Hooks.BeforeCreate, func(_ context.Context, ctx HookContext[hookParent]) error {
		calls = append(calls, "before_create")
		ctx.Item.Name = "created-by-hook"
		return nil
	})
	repo.Hooks.AfterCreate = append(repo.Hooks.AfterCreate, func(_ context.Context, ctx HookContext[hookParent]) error {
		calls = append(calls, "after_create")
		if ctx.Item.ID == 0 {
			t.Fatal("expected ID after create hook")
		}
		return nil
	})
	repo.Hooks.BeforeUpdate = append(repo.Hooks.BeforeUpdate, func(_ context.Context, ctx HookContext[hookParent]) error {
		calls = append(calls, "before_update")
		ctx.Item.Name = ctx.Item.Name + "-hooked"
		return nil
	})
	repo.Hooks.AfterDelete = append(repo.Hooks.AfterDelete, func(_ context.Context, _ HookContext[hookParent]) error {
		calls = append(calls, "after_delete")
		return nil
	})

	createValues := goadmin.Values{
		"Name":  {"plain"},
		"Items": {`[{"Label":"one"},{"Label":"two"}]`},
	}
	if err := repo.Create(context.Background(), createValues); err != nil {
		t.Fatal(err)
	}

	var parent hookParent
	if err := db.Preload("Items").First(&parent).Error; err != nil {
		t.Fatal(err)
	}
	if parent.Name != "created-by-hook" {
		t.Fatalf("expected create hook mutation, got %q", parent.Name)
	}
	if len(parent.Items) != 2 {
		t.Fatalf("expected 2 child rows, got %d", len(parent.Items))
	}

	updateValues := goadmin.Values{
		"Name":  {"updated"},
		"Items": {`[{"ID":"1","Label":"one-updated"}]`},
	}
	if err := repo.Update(context.Background(), "1", updateValues); err != nil {
		t.Fatal(err)
	}
	if err := db.Preload("Items", func(tx *gorm.DB) *gorm.DB { return tx.Order("sort asc, id asc") }).First(&parent, 1).Error; err != nil {
		t.Fatal(err)
	}
	if parent.Name != "updated-hooked" {
		t.Fatalf("expected update hook mutation, got %q", parent.Name)
	}
	if len(parent.Items) != 2 || parent.Items[0].Label != "one-updated" || parent.Items[1].Label != "two" {
		t.Fatalf("expected merge relation behavior, got %+v", parent.Items)
	}

	if err := repo.Delete(context.Background(), "1"); err != nil {
		t.Fatal(err)
	}
	if len(calls) != 4 || calls[0] != "before_create" || calls[1] != "after_create" || calls[2] != "before_update" || calls[3] != "after_delete" {
		t.Fatalf("unexpected hook calls: %+v", calls)
	}
}
