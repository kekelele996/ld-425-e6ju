package router

import (
	"github.com/gin-gonic/gin"
	"github.com/home-renovation/platform/internal/constants"
	"github.com/home-renovation/platform/internal/handler"
)

// registerProjectRoutes 注册项目路由。
func registerProjectRoutes(group *gin.RouterGroup, h *handler.ProjectHandler, auth gin.HandlerFunc, rbac func(...string) gin.HandlerFunc) {
	group.GET("/projects", auth, h.List)
	group.GET("/projects/:id", auth, h.Get)
	group.POST("/projects", auth, rbac(constants.RoleAdmin, constants.RoleProjectManager), h.Create)
	group.PUT("/projects/:id", auth, rbac(constants.RoleAdmin, constants.RoleProjectManager), h.Update)
	group.DELETE("/projects/:id", auth, rbac(constants.RoleAdmin), h.Delete)
	group.PUT("/projects/:id/status", auth, rbac(constants.RoleAdmin, constants.RoleProjectManager), h.UpdateStatus)
}
