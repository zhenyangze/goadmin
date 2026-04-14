package auth

import (
	"context"
	"errors"
	"time"

	"github.com/zhenyangze/goadmin"
	"gorm.io/gorm"
)

// LoginLog 记录用户登录/登出/失败行为
type LoginLog struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    *uint     // 用户ID（登录失败时可能为空）
	Username  string    `gorm:"size:120;index"` // 尝试登录的用户名
	IP        string    `gorm:"size:64"`
	UserAgent string    `gorm:"size:500"`
	Action    string    `gorm:"size:40;index"` // login, logout, failed, locked
	Reason    string    `gorm:"size:255"`      // 失败原因
	CreatedAt time.Time `gorm:"index"`
}

// TableName 返回表名
func (LoginLog) TableName() string {
	return "admin_login_logs"
}

// LoginLogger 定义登录日志记录接口
type LoginLogger interface {
	// LogLogin 记录登录行为
	LogLogin(ctx context.Context, userID *uint, username, ip, userAgent, action, reason string) error
	// GetRecentFailures 获取最近登录失败次数
	GetRecentFailures(ctx context.Context, username string, since time.Time) (int64, error)
}

// LoginLogRepository 提供登录日志的存储访问
type LoginLogRepository struct {
	DB *gorm.DB
}

// NewLoginLogRepository 创建登录日志仓库
func NewLoginLogRepository(db *gorm.DB) *LoginLogRepository {
	return &LoginLogRepository{DB: db}
}

// LogLogin 记录登录行为
func (r *LoginLogRepository) LogLogin(ctx context.Context, userID *uint, username, ip, userAgent, action, reason string) error {
	log := LoginLog{
		UserID:    userID,
		Username:  username,
		IP:        ip,
		UserAgent: userAgent,
		Action:    action,
		Reason:    reason,
		CreatedAt: time.Now(),
	}
	return r.DB.WithContext(ctx).Create(&log).Error
}

// GetRecentFailures 获取指定用户名在最近时间段内的失败次数
func (r *LoginLogRepository) GetRecentFailures(ctx context.Context, username string, since time.Time) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).
		Model(&LoginLog{}).
		Where("username = ? AND action = ? AND created_at >= ?", username, "failed", since).
		Count(&count).Error
	return count, err
}

// List 列出登录日志（用于管理界面）
func (r *LoginLogRepository) List(ctx context.Context, query goadmin.ListQuery) (goadmin.ListResult, error) {
	db := r.DB.WithContext(ctx).Model(&LoginLog{})

	if query.Search != "" {
		like := "%" + query.Search + "%"
		db = db.Where("username LIKE ? OR ip LIKE ? OR user_agent LIKE ?", like, like, like)
	}

	if action := query.Filters["Action"]; action != "" {
		db = db.Where("action = ?", action)
	}

	if username := query.Filters["Username"]; username != "" {
		db = db.Where("username = ?", username)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return goadmin.ListResult{}, err
	}

	page, perPage := normalizePage(query)
	var logs []LoginLog
	if err := db.Order("created_at desc, id desc").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&logs).Error; err != nil {
		return goadmin.ListResult{}, err
	}

	items := make([]any, 0, len(logs))
	for i := range logs {
		items = append(items, logs[i])
	}

	return goadmin.ListResult{
		Items:   items,
		Total:   total,
		Page:    page,
		PerPage: perPage,
	}, nil
}

// Get 获取单条登录日志
func (r *LoginLogRepository) Get(ctx context.Context, id string) (any, error) {
	var log LoginLog
	if err := r.DB.WithContext(ctx).First(&log, id).Error; err != nil {
		return nil, err
	}
	return log, nil
}

// Create 创建登录日志（管理界面不需要手动创建）
func (r *LoginLogRepository) Create(ctx context.Context, values goadmin.Values) error {
	return errors.New("login logs are created automatically")
}

// Update 更新登录日志（不允许）
func (r *LoginLogRepository) Update(ctx context.Context, id string, values goadmin.Values) error {
	return errors.New("login logs are read-only")
}

// Delete 删除登录日志
func (r *LoginLogRepository) Delete(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Delete(&LoginLog{}, id).Error
}
