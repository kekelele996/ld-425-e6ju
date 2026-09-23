package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/home-renovation/platform/internal/dto"
	"github.com/home-renovation/platform/internal/service"
	"github.com/home-renovation/platform/internal/utils"
)

// ProjectHandler 项目处理器。
type ProjectHandler struct {
	service service.ProjectService
}

// NewProjectHandler 构造项目处理器。
func NewProjectHandler(service service.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

// Create 创建项目。
func (h *ProjectHandler) Create(c *gin.Context) {
	var req dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	project, err := h.service.Create(&req)
	if err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, toProjectDTO(project))
}

// Get 获取项目详情。
func (h *ProjectHandler) Get(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.Error(err)
		return
	}
	project, err := h.service.GetByID(id)
	if err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, toProjectDTO(project))
}

// List 获取项目列表。
func (h *ProjectHandler) List(c *gin.Context) {
	page, pageSize := parsePage(c)
	status := c.Query("status")
	keyword := c.Query("keyword")
	projects, total, err := h.service.List(status, keyword, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, utils.PageResult{List: toProjectDTOList(projects), Total: total, Page: page, PageSize: pageSize})
}

// Update 更新项目。
func (h *ProjectHandler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.Error(err)
		return
	}
	var req dto.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	project, err := h.service.Update(id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, toProjectDTO(project))
}

// Delete 删除项目。
func (h *ProjectHandler) Delete(c *gin.Context) {
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

// UpdateStatus 项目状态流转。
func (h *ProjectHandler) UpdateStatus(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.Error(err)
		return
	}
	var req dto.UpdateProjectStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	project, err := h.service.UpdateStatus(id, req.Status)
	if err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, toProjectDTO(project))
}
