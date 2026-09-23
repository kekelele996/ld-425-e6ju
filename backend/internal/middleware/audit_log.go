package middleware

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/home-renovation/platform/internal/model"
	"github.com/home-renovation/platform/internal/repository"
)

// AuditLogMiddleware 记录写操作日志。
func AuditLogMiddleware(repo repository.AuditLogRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		method := c.Request.Method
		if method == "GET" || method == "HEAD" || method == "OPTIONS" {
			return
		}
		userID, _ := c.Get(string(ContextUserID))
		username, _ := c.Get(string(ContextUsername))
		uid, _ := userID.(uint)
		uname, _ := username.(string)
		if uid == 0 && uname == "" {
			return
		}
		resource := resourceFromPath(c.FullPath())
		entry := &model.AuditLog{
			UserID:   uid,
			Username: uname,
			Action:   method,
			Resource: resource,
			Detail:   c.Request.Method + " " + c.Request.URL.Path,
			IP:       c.ClientIP(),
		}
		_ = repo.Create(entry)
	}
}

func resourceFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	// 路径形如 /api/v1/projects，取资源段。
	for i, part := range parts {
		if part == "v1" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return strconv.Itoa(0)
}
