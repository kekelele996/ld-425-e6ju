package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/home-renovation/platform/internal/dto"
	"github.com/home-renovation/platform/internal/middleware"
	"github.com/home-renovation/platform/internal/service"
	"github.com/home-renovation/platform/internal/utils"
)

// DesignHandler 设计阶段处理器。
type DesignHandler struct {
	service service.DesignService
}

// NewDesignHandler 构造设计处理器。
func NewDesignHandler(service service.DesignService) *DesignHandler {
	return &DesignHandler{service: service}
}

// Create 创建设计阶段。
func (h *DesignHandler) Create(c *gin.Context) {
	var req dto.CreateDesignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	phase, err := h.service.Create(&req)
	if err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, toDesignDTO(phase))
}

// Get 获取设计阶段详情。
func (h *DesignHandler) Get(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.Error(err)
		return
	}
	phase, err := h.service.GetByID(id)
	if err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, toDesignDTO(phase))
}

// List 获取设计阶段列表。
func (h *DesignHandler) List(c *gin.Context) {
	page, pageSize := parsePage(c)
	if projectIDStr := c.Query("project_id"); projectIDStr != "" {
		id, err := parseUintString(projectIDStr)
		if err != nil {
			c.Error(err)
			return
		}
		phases, err := h.service.ListByProjectID(id)
		if err != nil {
			c.Error(err)
			return
		}
		utils.Success(c, toDesignDTOList(phases))
		return
	}
	phases, total, err := h.service.List(page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, utils.PageResult{List: toDesignDTOList(phases), Total: total, Page: page, PageSize: pageSize})
}

// Update 更新设计阶段。
func (h *DesignHandler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.Error(err)
		return
	}
	var req dto.UpdateDesignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	phase, err := h.service.Update(id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, toDesignDTO(phase))
}

// Delete 删除设计阶段。
func (h *DesignHandler) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.Error(err)
		return
	}
	if err := h.service.Delete(id); err != nil {
		c.Error(err)
		return
	}
	utils.SuccessMessage(c, "deleted", nil)
}

// Submit 提交设计审核。
func (h *DesignHandler) Submit(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.Error(err)
		return
	}
	var req dto.SubmitDesignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	phase, err := h.service.Submit(id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, toDesignDTO(phase))
}

// Review 业主审核设计。
func (h *DesignHandler) Review(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.Error(err)
		return
	}
	var req dto.ReviewDesignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	reviewerID, _ := c.Get(string(middleware.ContextUserID))
	phase, err := h.service.Review(id, reviewerID.(uint), &req)
	if err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, toDesignDTO(phase))
}
