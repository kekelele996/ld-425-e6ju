package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/home-renovation/platform/internal/dto"
	"github.com/home-renovation/platform/internal/service"
	"github.com/home-renovation/platform/internal/utils"
)

// BudgetHandler 预算处理器。
type BudgetHandler struct {
	service service.BudgetService
}

// NewBudgetHandler 构造预算处理器。
func NewBudgetHandler(service service.BudgetService) *BudgetHandler {
	return &BudgetHandler{service: service}
}

// Create 创建预算项。
func (h *BudgetHandler) Create(c *gin.Context) {
	var req dto.CreateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	item, err := h.service.Create(&req)
	if err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, toBudgetDTO(item))
}

// Get 获取预算项详情。
func (h *BudgetHandler) Get(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.Error(err)
		return
	}
	item, err := h.service.GetByID(id)
	if err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, toBudgetDTO(item))
}

// List 获取预算列表。
func (h *BudgetHandler) List(c *gin.Context) {
	page, pageSize := parsePage(c)
	if projectIDStr := c.Query("project_id"); projectIDStr != "" {
		id, err := parseUintString(projectIDStr)
		if err != nil {
			c.Error(err)
			return
		}
		items, err := h.service.ListByProjectID(id)
		if err != nil {
			c.Error(err)
			return
		}
		utils.Success(c, toBudgetDTOList(items))
		return
	}
	category := c.Query("category")
	items, total, err := h.service.List(0, category, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, utils.PageResult{List: toBudgetDTOList(items), Total: total, Page: page, PageSize: pageSize})
}

// Update 更新预算项。
func (h *BudgetHandler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.Error(err)
		return
	}
	var req dto.UpdateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	item, err := h.service.Update(id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, toBudgetDTO(item))
}

// Delete 删除预算项。
func (h *BudgetHandler) Delete(c *gin.Context) {
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
