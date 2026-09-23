package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/home-renovation/platform/internal/service"
	"github.com/home-renovation/platform/internal/utils"
)

// AuditHandler 操作日志处理器。
type AuditHandler struct {
	service service.AuditService
}

// NewAuditHandler 构造操作日志处理器。
func NewAuditHandler(service service.AuditService) *AuditHandler {
	return &AuditHandler{service: service}
}

// List 获取操作日志列表。
func (h *AuditHandler) List(c *gin.Context) {
	page, pageSize := parsePage(c)
	logs, total, err := h.service.List(page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, utils.PageResult{List: logs, Total: total, Page: page, PageSize: pageSize})
}
