package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/home-renovation/platform/internal/dto"
	"github.com/home-renovation/platform/internal/service"
	"github.com/home-renovation/platform/internal/utils"
)

// ConstructionHandler 施工节点处理器。
type ConstructionHandler struct {
	service service.ConstructionService
}

// NewConstructionHandler 构造施工处理器。
func NewConstructionHandler(service service.ConstructionService) *ConstructionHandler {
	return &ConstructionHandler{service: service}
}

// Create 创建施工节点。
func (h *ConstructionHandler) Create(c *gin.Context) {
	var req dto.CreateConstructionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	node, err := h.service.Create(&req)
	if err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, toConstructionDTO(node))
}

// Get 获取施工节点详情。
func (h *ConstructionHandler) Get(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.Error(err)
		return
	}
	node, err := h.service.GetByID(id)
	if err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, toConstructionDTO(node))
}

// List 获取施工节点列表。
func (h *ConstructionHandler) List(c *gin.Context) {
	page, pageSize := parsePage(c)
	if projectIDStr := c.Query("project_id"); projectIDStr != "" {
		id, err := parseUintString(projectIDStr)
		if err != nil {
			c.Error(err)
			return
		}
		nodes, err := h.service.ListByProjectID(id)
		if err != nil {
			c.Error(err)
			return
		}
		utils.Success(c, toConstructionDTOList(nodes))
		return
	}
	status := c.Query("status")
	nodes, total, err := h.service.List(0, status, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, utils.PageResult{List: toConstructionDTOList(nodes), Total: total, Page: page, PageSize: pageSize})
}

// Update 更新施工节点。
func (h *ConstructionHandler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.Error(err)
		return
	}
	var req dto.UpdateConstructionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	node, err := h.service.Update(id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, toConstructionDTO(node))
}

// Delete 删除施工节点。
func (h *ConstructionHandler) Delete(c *gin.Context) {
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

// UpdateStatus 施工状态流转。
func (h *ConstructionHandler) UpdateStatus(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.Error(err)
		return
	}
	var req dto.UpdateConstructionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	node, err := h.service.UpdateStatus(id, req.Status)
	if err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, toConstructionDTO(node))
}

// Accept 施工验收。
func (h *ConstructionHandler) Accept(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.Error(err)
		return
	}
	var req dto.AcceptConstructionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	node, err := h.service.Accept(id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, toConstructionDTO(node))
}
