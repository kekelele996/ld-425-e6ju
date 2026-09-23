package router

import (
	"github.com/gin-gonic/gin"
	"github.com/home-renovation/platform/internal/constants"
	"github.com/home-renovation/platform/internal/handler"
)

// registerDesignRoutes 注册设计路由。
func registerDesignRoutes(group *gin.RouterGroup, h *handler.DesignHandler, auth gin.HandlerFunc, rbac func(...string) gin.HandlerFunc) {
	group.GET("/designs", auth, h.List)
	group.GET("/designs/:id", auth, h.Get)
	group.POST("/designs", auth, rbac(constants.RoleAdmin, constants.RoleDesigner, constants.RoleProjectManager), h.Create)
	group.PUT("/designs/:id", auth, rbac(constants.RoleAdmin, constants.RoleDesigner, constants.RoleProjectManager), h.Update)
	group.DELETE("/designs/:id", auth, rbac(constants.RoleAdmin), h.Delete)
	group.PUT("/designs/:id/submit", auth, rbac(constants.RoleAdmin, constants.RoleDesigner, constants.RoleProjectManager), h.Submit)
	group.PUT("/designs/:id/review", auth, rbac(constants.RoleAdmin, constants.RoleOwner), h.Review)
}
