package auth

import (
	"context"
	"errors"
	"strconv"

	"github.com/zhenyangze/goadmin"
	"gorm.io/gorm"
)

// UserRepository manages built-in users plus their role assignments.
type UserRepository struct {
	DB *gorm.DB
}

// RoleRepository manages built-in roles plus permission/menu assignments.
type RoleRepository struct {
	DB *gorm.DB
}

// PermissionRepository manages built-in permissions.
type PermissionRepository struct {
	DB *gorm.DB
}

// MenuRepository manages built-in menus and exposes a tree.
type MenuRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{DB: db}
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{DB: db}
}

func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{DB: db}
}

func (r *UserRepository) List(ctx context.Context, query goadmin.ListQuery) (goadmin.ListResult, error) {
	db := r.DB.WithContext(ctx).Model(&User{}).Preload("Roles")
	if query.Search != "" {
		like := "%" + query.Search + "%"
		db = db.Where("username LIKE ? OR name LIKE ?", like, like)
	}
	if username := query.Filters["Username"]; username != "" {
		db = db.Where("username = ?", username)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return goadmin.ListResult{}, err
	}
	page, perPage := normalizePage(query)
	var users []User
	if err := db.Order("id desc").Offset((page - 1) * perPage).Limit(perPage).Find(&users).Error; err != nil {
		return goadmin.ListResult{}, err
	}
	items := make([]any, 0, len(users))
	for i := range users {
		items = append(items, users[i])
	}
	return goadmin.ListResult{Items: items, Total: total, Page: page, PerPage: perPage}, nil
}

func (r *UserRepository) Get(ctx context.Context, id string) (any, error) {
	var user User
	if err := r.DB.WithContext(ctx).Preload("Roles").First(&user, id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Create(ctx context.Context, values goadmin.Values) error {
	password, err := HashPassword(values.First("Password"))
	if err != nil {
		return err
	}
	user := User{
		Username: values.First("Username"),
		Name:     values.First("Name"),
		Password: password,
	}
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		return replaceRoles(tx, &user, values.All("RoleIDs"))
	})
}

func (r *UserRepository) Update(ctx context.Context, id string, values goadmin.Values) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user User
		if err := tx.Preload("Roles").First(&user, id).Error; err != nil {
			return err
		}
		user.Username = values.First("Username")
		user.Name = values.First("Name")
		if password := values.First("Password"); password != "" {
			hashed, err := HashPassword(password)
			if err != nil {
				return err
			}
			user.Password = hashed
		}
		if err := tx.Save(&user).Error; err != nil {
			return err
		}
		return replaceRoles(tx, &user, values.All("RoleIDs"))
	})
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	var total int64
	if err := r.DB.WithContext(ctx).Model(&User{}).Count(&total).Error; err != nil {
		return err
	}
	if total <= 1 {
		return errors.New("cannot delete the last admin user")
	}
	return r.DB.WithContext(ctx).Delete(&User{}, id).Error
}

func (r *RoleRepository) List(ctx context.Context, query goadmin.ListQuery) (goadmin.ListResult, error) {
	db := r.DB.WithContext(ctx).Model(&Role{}).Preload("Permissions").Preload("Menus")
	if query.Search != "" {
		like := "%" + query.Search + "%"
		db = db.Where("name LIKE ? OR slug LIKE ?", like, like)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return goadmin.ListResult{}, err
	}
	page, perPage := normalizePage(query)
	var roles []Role
	if err := db.Order("id asc").Offset((page - 1) * perPage).Limit(perPage).Find(&roles).Error; err != nil {
		return goadmin.ListResult{}, err
	}
	items := make([]any, 0, len(roles))
	for i := range roles {
		items = append(items, roles[i])
	}
	return goadmin.ListResult{Items: items, Total: total, Page: page, PerPage: perPage}, nil
}

func (r *RoleRepository) Get(ctx context.Context, id string) (any, error) {
	var role Role
	if err := r.DB.WithContext(ctx).Preload("Permissions").Preload("Menus").First(&role, id).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) Create(ctx context.Context, values goadmin.Values) error {
	role := Role{Name: values.First("Name"), Slug: values.First("Slug")}
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&role).Error; err != nil {
			return err
		}
		if err := replacePermissions(tx, &role, values.All("PermissionIDs")); err != nil {
			return err
		}
		return replaceMenus(tx, &role, values.All("MenuIDs"))
	})
}

func (r *RoleRepository) Update(ctx context.Context, id string, values goadmin.Values) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var role Role
		if err := tx.Preload("Permissions").Preload("Menus").First(&role, id).Error; err != nil {
			return err
		}
		role.Name = values.First("Name")
		role.Slug = values.First("Slug")
		if err := tx.Save(&role).Error; err != nil {
			return err
		}
		if err := replacePermissions(tx, &role, values.All("PermissionIDs")); err != nil {
			return err
		}
		return replaceMenus(tx, &role, values.All("MenuIDs"))
	})
}

func (r *RoleRepository) Delete(ctx context.Context, id string) error {
	var role Role
	if err := r.DB.WithContext(ctx).First(&role, id).Error; err != nil {
		return err
	}
	if role.Slug == AdministratorRole {
		return errors.New("cannot delete the administrator role")
	}
	return r.DB.WithContext(ctx).Delete(&Role{}, id).Error
}

func (r *PermissionRepository) List(ctx context.Context, query goadmin.ListQuery) (goadmin.ListResult, error) {
	db := r.DB.WithContext(ctx).Model(&Permission{})
	if query.Search != "" {
		like := "%" + query.Search + "%"
		db = db.Where("name LIKE ? OR slug LIKE ? OR path LIKE ?", like, like, like)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return goadmin.ListResult{}, err
	}
	page, perPage := normalizePage(query)
	var permissions []Permission
	if err := db.Order("id asc").Offset((page - 1) * perPage).Limit(perPage).Find(&permissions).Error; err != nil {
		return goadmin.ListResult{}, err
	}
	items := make([]any, 0, len(permissions))
	for i := range permissions {
		items = append(items, permissions[i])
	}
	return goadmin.ListResult{Items: items, Total: total, Page: page, PerPage: perPage}, nil
}

func (r *PermissionRepository) Get(ctx context.Context, id string) (any, error) {
	var permission Permission
	if err := r.DB.WithContext(ctx).First(&permission, id).Error; err != nil {
		return nil, err
	}
	return permission, nil
}

func (r *PermissionRepository) Create(ctx context.Context, values goadmin.Values) error {
	permission := Permission{Name: values.First("Name"), Slug: values.First("Slug"), Path: values.First("Path")}
	return r.DB.WithContext(ctx).Create(&permission).Error
}

func (r *PermissionRepository) Update(ctx context.Context, id string, values goadmin.Values) error {
	var permission Permission
	if err := r.DB.WithContext(ctx).First(&permission, id).Error; err != nil {
		return err
	}
	permission.Name = values.First("Name")
	permission.Slug = values.First("Slug")
	permission.Path = values.First("Path")
	return r.DB.WithContext(ctx).Save(&permission).Error
}

func (r *PermissionRepository) Delete(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Delete(&Permission{}, id).Error
}

func (r *MenuRepository) List(ctx context.Context, query goadmin.ListQuery) (goadmin.ListResult, error) {
	db := r.DB.WithContext(ctx).Model(&Menu{})
	if query.Search != "" {
		like := "%" + query.Search + "%"
		db = db.Where("title LIKE ? OR uri LIKE ? OR permission_slug LIKE ?", like, like, like)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return goadmin.ListResult{}, err
	}
	page, perPage := normalizePage(query)
	var menus []Menu
	if err := db.Order("parent_id asc, `order` asc, id asc").Offset((page - 1) * perPage).Limit(perPage).Find(&menus).Error; err != nil {
		return goadmin.ListResult{}, err
	}
	items := make([]any, 0, len(menus))
	for i := range menus {
		items = append(items, menus[i])
	}
	return goadmin.ListResult{Items: items, Total: total, Page: page, PerPage: perPage}, nil
}

func (r *MenuRepository) Get(ctx context.Context, id string) (any, error) {
	var menu Menu
	if err := r.DB.WithContext(ctx).First(&menu, id).Error; err != nil {
		return nil, err
	}
	return menu, nil
}

func (r *MenuRepository) Create(ctx context.Context, values goadmin.Values) error {
	menu := Menu{
		ParentID:       parseUint(values.First("ParentID")),
		Order:          parseInt(values.First("Order")),
		Title:          values.First("Title"),
		Icon:           values.First("Icon"),
		URI:            values.First("URI"),
		PermissionSlug: values.First("PermissionSlug"),
	}
	return r.DB.WithContext(ctx).Create(&menu).Error
}

func (r *MenuRepository) Update(ctx context.Context, id string, values goadmin.Values) error {
	var menu Menu
	if err := r.DB.WithContext(ctx).First(&menu, id).Error; err != nil {
		return err
	}
	menu.ParentID = parseUint(values.First("ParentID"))
	menu.Order = parseInt(values.First("Order"))
	menu.Title = values.First("Title")
	menu.Icon = values.First("Icon")
	menu.URI = values.First("URI")
	menu.PermissionSlug = values.First("PermissionSlug")
	return r.DB.WithContext(ctx).Save(&menu).Error
}

func (r *MenuRepository) Delete(ctx context.Context, id string) error {
	var count int64
	if err := r.DB.WithContext(ctx).Model(&Menu{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("delete child menus before removing this menu")
	}
	return r.DB.WithContext(ctx).Delete(&Menu{}, id).Error
}

func (r *MenuRepository) Tree(ctx context.Context) ([]goadmin.TreeNode, error) {
	var menus []Menu
	if err := r.DB.WithContext(ctx).Order("parent_id asc, `order` asc, id asc").Find(&menus).Error; err != nil {
		return nil, err
	}
	nodesByParent := map[uint][]goadmin.TreeNode{}
	for _, menu := range menus {
		nodesByParent[menu.ParentID] = append(nodesByParent[menu.ParentID], goadmin.TreeNode{
			ID:          toString(menu.ID),
			ParentID:    toString(menu.ParentID),
			Title:       menu.Title,
			Description: menu.URI,
		})
	}
	return attachMenuChildren(nodesByParent, 0), nil
}

func replaceRoles(tx *gorm.DB, user *User, ids []string) error {
	var roles []Role
	if len(ids) > 0 {
		if err := tx.Where("id IN ?", ids).Find(&roles).Error; err != nil {
			return err
		}
	}
	return tx.Model(user).Association("Roles").Replace(&roles)
}

func replacePermissions(tx *gorm.DB, role *Role, ids []string) error {
	var permissions []Permission
	if len(ids) > 0 {
		if err := tx.Where("id IN ?", ids).Find(&permissions).Error; err != nil {
			return err
		}
	}
	return tx.Model(role).Association("Permissions").Replace(&permissions)
}

func replaceMenus(tx *gorm.DB, role *Role, ids []string) error {
	var menus []Menu
	if len(ids) > 0 {
		if err := tx.Where("id IN ?", ids).Find(&menus).Error; err != nil {
			return err
		}
	}
	return tx.Model(role).Association("Menus").Replace(&menus)
}

func attachMenuChildren(itemsByParent map[uint][]goadmin.TreeNode, parentID uint) []goadmin.TreeNode {
	nodes := itemsByParent[parentID]
	for i := range nodes {
		nodes[i].Children = attachMenuChildren(itemsByParent, parseUint(nodes[i].ID))
	}
	return nodes
}

func normalizePage(query goadmin.ListQuery) (int, int) {
	page := query.Page
	perPage := query.PerPage
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	return page, perPage
}

func parseInt(raw string) int {
	if raw == "" {
		return 0
	}
	value, _ := strconv.Atoi(raw)
	return value
}

func parseUint(raw string) uint {
	if raw == "" {
		return 0
	}
	value, _ := strconv.ParseUint(raw, 10, 64)
	return uint(value)
}

func toString[T ~uint | ~uint64 | ~uint32 | ~uint16 | ~uint8](value T) string {
	return strconv.FormatUint(uint64(value), 10)
}
