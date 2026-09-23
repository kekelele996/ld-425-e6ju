package middleware

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/home-renovation/platform/internal/config"
	"github.com/home-renovation/platform/internal/utils"
)

// UploadMiddleware 文件上传中间件，校验大小与扩展名并保存文件。
func UploadMiddleware(cfg config.UploadConfig) gin.HandlerFunc {
	allowedExts := map[string]struct{}{
		".jpg": {}, ".jpeg": {}, ".png": {}, ".webp": {}, ".pdf": {}, ".dwg": {}, ".zip": {},
	}
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			utils.AbortError(c, http.StatusBadRequest, 40000, "file field is required")
			return
		}
		if file.Size > cfg.MaxSizeMB*1024*1024 {
			utils.AbortError(c, http.StatusBadRequest, 40000, fmt.Sprintf("file exceeds %dMB limit", cfg.MaxSizeMB))
			return
		}
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if _, ok := allowedExts[ext]; !ok {
			utils.AbortError(c, http.StatusBadRequest, 40000, "unsupported file extension")
			return
		}
		url, err := utils.SaveUploadedFile(cfg.Dir, cfg.PublicPrefix, file)
		if err != nil {
			c.Error(fmt.Errorf("save uploaded file: %w", err))
			return
		}
		c.Set("upload_url", url)
		c.Next()
	}
}
