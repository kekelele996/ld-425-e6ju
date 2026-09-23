package router

import (
	"github.com/gin-gonic/gin"
	"github.com/home-renovation/platform/internal/constants"
	"github.com/home-renovation/platform/internal/handler"
)

// registerMaterialRoutes 注册材料路由。
func registerMaterialRoutes(group *gin.RouterGroup, h *handler.MaterialHandler, auth gin.HandlerFunc, rbac func(...string) gin.HandlerFunc) {
	group.GET("/materials", auth, h.List)
	group.GET("/materials/:id", auth, h.Get)
	group.POST("/materials", auth, rbac(constants.RoleAdmin, constants.RoleDesigner, constants.RoleProjectManager), h.Create)
	group.PUT("/materials/:id", auth, rbac(constants.RoleAdmin, constants.RoleDesigner, constants.RoleContractor, constants.RoleProjectManager), h.Update)
	group.DELETE("/materials/:id", auth, rbac(constants.RoleAdmin), h.Delete)
	group.PUT("/materials/:id/status", auth, rbac(constants.RoleAdmin, constants.RoleDesigner, constants.RoleContractor, constants.RoleProjectManager), h.UpdateStatus)
}
