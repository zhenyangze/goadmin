package goadmin

import (
	"context"
	"mime/multipart"

	"github.com/zhenyangze/goadmin/theme"
)

// UploadHandler stores an uploaded file and returns the persisted public path/URL.
type UploadHandler func(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error)
type DeleteUploadHandler func(ctx context.Context, publicPath string) error

// Config defines the reusable admin app settings.
type Config struct {
	AppName       string
	Title         string
	Prefix        string
	SessionSecret string
	SessionCookie string
	UploadDir     string
	UploadPath    string
	SaveUpload    UploadHandler
	DeleteUpload  DeleteUploadHandler
	Theme         theme.Theme
}

// WithDefaults normalizes configuration values.
func (c Config) WithDefaults() Config {
	if c.AppName == "" {
		c.AppName = "Go Admin"
	}
	if c.Title == "" {
		c.Title = c.AppName
	}
	if c.Prefix == "" {
		c.Prefix = "/admin"
	}
	if c.SessionCookie == "" {
		c.SessionCookie = "goadmin_session"
	}
	if c.UploadDir == "" {
		c.UploadDir = "tmp/goadmin/uploads"
	}
	if c.UploadPath == "" {
		c.UploadPath = "uploads"
	}
	if c.Theme.Accent == "" {
		c.Theme = theme.Default()
	}
	return c
}
