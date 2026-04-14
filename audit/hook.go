// Package audit provides repository hooks for automatic audit logging.
package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/zhenyangze/goadmin"
	"gorm.io/gorm"
)

// AuditLog represents an audit log entry
type AuditLog struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"index"` // 操作用户ID
	Actor      string    `gorm:"size:120"` // 操作用户名
	Action     string    `gorm:"size:40;index"` // create, update, delete
	Resource   string    `gorm:"size:120;index"` // 资源名称
	ResourceID string    `gorm:"size:120;index"` // 资源ID
	IP         string    `gorm:"size:64"`
	UserAgent  string    `gorm:"size:500"`
	Level      string    `gorm:"size:40;index"` // info, warning, danger
	Detail     string    `gorm:"size:2000"` // 操作详情(JSON)
	CreatedAt  time.Time `gorm:"index"`
}

// TableName returns the table name
func (AuditLog) TableName() string {
	return "admin_audit_logs"
}

// IdentityFunc retrieves identity from context
type IdentityFunc func(ctx context.Context) *goadmin.Identity

// RequestInfoFunc retrieves request info from context
type RequestInfoFunc func(ctx context.Context) (ip, userAgent string)

// Hook implements goadmin.RepositoryHook for automatic audit logging
type Hook struct {
	DB              *gorm.DB
	ResourceName    string
	GetIdentity     IdentityFunc
	GetRequestInfo  RequestInfoFunc
	LogCreate       bool // 是否记录创建操作
	LogUpdate       bool // 是否记录更新操作
	LogDelete       bool // 是否记录删除操作
	MaxDetailLength int  // 详情最大长度
}

// NewHook creates a new audit hook
func NewHook(db *gorm.DB, resourceName string) *Hook {
	return &Hook{
		DB:              db,
		ResourceName:    resourceName,
		LogCreate:       true,
		LogUpdate:       true,
		LogDelete:       true,
		MaxDetailLength: 2000,
		GetIdentity:     defaultGetIdentity,
		GetRequestInfo:  defaultGetRequestInfo,
	}
}

// SetIdentityFunc sets a custom identity retrieval function
func (h *Hook) SetIdentityFunc(fn IdentityFunc) *Hook {
	h.GetIdentity = fn
	return h
}

// SetRequestInfoFunc sets a custom request info retrieval function
func (h *Hook) SetRequestInfoFunc(fn RequestInfoFunc) *Hook {
	h.GetRequestInfo = fn
	return h
}

// BeforeCreate implements RepositoryHook
func (h *Hook) BeforeCreate(ctx context.Context, values goadmin.Values) error {
	return nil
}

// AfterCreate implements RepositoryHook
func (h *Hook) AfterCreate(ctx context.Context, id string, values goadmin.Values) error {
	if !h.LogCreate {
		return nil
	}
	return h.log(ctx, "create", id, values, nil)
}

// BeforeUpdate implements RepositoryHook
func (h *Hook) BeforeUpdate(ctx context.Context, id string, values goadmin.Values) error {
	return nil
}

// AfterUpdate implements RepositoryHook
func (h *Hook) AfterUpdate(ctx context.Context, id string, values goadmin.Values, old any) error {
	if !h.LogUpdate {
		return nil
	}
	return h.log(ctx, "update", id, values, old)
}

// BeforeDelete implements RepositoryHook
func (h *Hook) BeforeDelete(ctx context.Context, id string) error {
	return nil
}

// AfterDelete implements RepositoryHook
func (h *Hook) AfterDelete(ctx context.Context, id string, deleted any) error {
	if !h.LogDelete {
		return nil
	}
	return h.log(ctx, "delete", id, nil, deleted)
}

// log creates an audit log entry
func (h *Hook) log(ctx context.Context, action, resourceID string, values goadmin.Values, old any) error {
	identity := h.GetIdentity(ctx)
	if identity == nil {
		// 未登录用户操作，可选记录为 system 或跳过
		identity = &goadmin.Identity{
			ID:       0,
			Username: "system",
		}
	}

	ip, userAgent := h.GetRequestInfo(ctx)

	detail := h.buildDetail(action, values, old)
	if len(detail) > h.MaxDetailLength {
		detail = detail[:h.MaxDetailLength] + "..."
	}

	log := AuditLog{
		UserID:     identity.ID,
		Actor:      identity.DisplayName,
		Action:     action,
		Resource:   h.ResourceName,
		ResourceID: resourceID,
		IP:         ip,
		UserAgent:  userAgent,
		Level:      h.getLevel(action),
		Detail:     detail,
		CreatedAt:  time.Now(),
	}

	// 使用后台 context 避免事务影响
	return h.DB.Create(&log).Error
}

// buildDetail builds the detail JSON
func (h *Hook) buildDetail(action string, values goadmin.Values, old any) string {
	detail := map[string]any{
		"action": action,
	}

	if values != nil && len(values) > 0 {
		// 过滤敏感字段
		filtered := make(goadmin.Values)
		for k, v := range values {
			if isSensitiveField(k) {
				filtered[k] = []string{"[REDACTED]"}
			} else {
				filtered[k] = v
			}
		}
		detail["values"] = filtered
	}

	if old != nil {
		// 尝试序列化旧值（脱敏）
		oldMap := sanitizeForLog(old)
		detail["old"] = oldMap
	}

	jsonBytes, err := json.Marshal(detail)
	if err != nil {
		return fmt.Sprintf("{\"action\":\"%s\",\"error\":\"marshal failed\"}", action)
	}
	return string(jsonBytes)
}

// getLevel determines log level based on action
func (h *Hook) getLevel(action string) string {
	switch action {
	case "delete":
		return "danger"
	case "update":
		return "warning"
	default:
		return "info"
	}
}

// isSensitiveField checks if field contains sensitive data
func isSensitiveField(name string) bool {
	sensitive := []string{"password", "passwd", "pwd", "secret", "token", "api_key", "apikey", "credential"}
	lower := strings.ToLower(name)
	for _, s := range sensitive {
		if strings.Contains(lower, s) {
			return true
		}
	}
	return false
}

// sanitizeForLog removes sensitive data from log
func sanitizeForLog(v any) any {
	// 使用 JSON 序列化再反序列化来清理
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil
	}

	var m map[string]any
	if err := json.Unmarshal(bytes, &m); err != nil {
		return v
	}

	// 移除敏感字段
	for k := range m {
		if isSensitiveField(k) {
			m[k] = "[REDACTED]"
		}
	}

	return m
}

// defaultGetIdentity retrieves identity from context
func defaultGetIdentity(ctx context.Context) *goadmin.Identity {
	if v := ctx.Value("identity"); v != nil {
		if id, ok := v.(*goadmin.Identity); ok {
			return id
		}
	}
	return nil
}

// defaultGetRequestInfo retrieves request info from context
func defaultGetRequestInfo(ctx context.Context) (string, string) {
	if v := ctx.Value("request_info"); v != nil {
		if info, ok := v.(RequestInfo); ok {
			return info.IP, info.UserAgent
		}
	}
	return "", ""
}

// RequestInfo holds request information for audit logging
type RequestInfo struct {
	IP        string
	UserAgent string
}
