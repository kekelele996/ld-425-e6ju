package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/home-renovation/platform/internal/utils"
)

// UploadHandler 文件上传处理器。
type UploadHandler struct{}

// NewUploadHandler 构造文件上传处理器。
func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

// Upload 上传文件，URL 由 UploadMiddleware 写入上下文。
func (h *UploadHandler) Upload(c *gin.Context) {
	url, _ := c.Get("upload_url")
	utils.Success(c, gin.H{"url": url})
}
