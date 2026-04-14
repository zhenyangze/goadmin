package goadmin

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zhenyangze/goadmin/form"
)

func (a *App) uploadBaseURL() string {
	return joinURL(a.cfg.Prefix, a.cfg.UploadPath)
}

func (a *App) serveUploads(w http.ResponseWriter, r *http.Request) bool {
	prefix := a.uploadBaseURL()
	current := normalizePath(r.URL.Path)
	if current != prefix && !strings.HasPrefix(current, prefix+"/") {
		return false
	}
	http.StripPrefix(prefix, http.FileServer(http.Dir(a.cfg.UploadDir))).ServeHTTP(w, r)
	return true
}

func (a *App) saveUpload(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	if a.cfg.SaveUpload != nil {
		return a.cfg.SaveUpload(ctx, file, header)
	}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	token, err := randomToken(12)
	if err != nil {
		return "", err
	}
	relativePath := filepath.Join(time.Now().Format("2006/01/02"), token+ext)
	diskPath := filepath.Join(a.cfg.UploadDir, relativePath)
	if err := os.MkdirAll(filepath.Dir(diskPath), 0o755); err != nil {
		return "", err
	}
	dst, err := os.Create(diskPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}
	return joinURL(a.uploadBaseURL(), filepath.ToSlash(relativePath)), nil
}

func (a *App) validateUpload(field *form.Field, header *multipart.FileHeader) error {
	if field.MaxFileSize > 0 && header.Size > field.MaxFileSize {
		return errors.New("uploaded file exceeds the configured size limit")
	}
	if len(field.AllowedExtensions) == 0 {
		return nil
	}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	for _, allowed := range field.AllowedExtensions {
		normalized := strings.ToLower(strings.TrimSpace(allowed))
		if normalized == "" {
			continue
		}
		if !strings.HasPrefix(normalized, ".") {
			normalized = "." + normalized
		}
		if ext == normalized {
			return nil
		}
	}
	return errors.New("uploaded file type is not allowed")
}

func (a *App) deleteUpload(ctx context.Context, publicPath string) error {
	if strings.TrimSpace(publicPath) == "" {
		return nil
	}
	if a.cfg.DeleteUpload != nil {
		return a.cfg.DeleteUpload(ctx, publicPath)
	}
	base := a.uploadBaseURL()
	if publicPath != base && !strings.HasPrefix(publicPath, base+"/") {
		return nil
	}
	relative := strings.TrimPrefix(publicPath, base)
	relative = strings.TrimPrefix(relative, "/")
	if relative == "" {
		return nil
	}
	return os.Remove(filepath.Join(a.cfg.UploadDir, filepath.FromSlash(relative)))
}
