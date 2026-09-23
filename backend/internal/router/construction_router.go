package router

import (
	"github.com/gin-gonic/gin"
	"github.com/home-renovation/platform/internal/constants"
	"github.com/home-renovation/platform/internal/handler"
)

// registerConstructionRoutes 注册施工路由。
func registerConstructionRoutes(group *gin.RouterGroup, h *handler.ConstructionHandler, auth gin.HandlerFunc, rbac func(...string) gin.HandlerFunc) {
	group.GET("/constructions", auth, h.List)
	group.GET("/constructions/:id", auth, h.Get)
	group.POST("/constructions", auth, rbac(constants.RoleAdmin, constants.RoleProjectManager), h.Create)
	group.PUT("/constructions/:id", auth, rbac(constants.RoleAdmin, constants.RoleProjectManager, constants.RoleContractor), h.Update)
	group.DELETE("/constructions/:id", auth, rbac(constants.RoleAdmin), h.Delete)
	group.PUT("/constructions/:id/status", auth, rbac(constants.RoleAdmin, constants.RoleProjectManager, constants.RoleContractor), h.UpdateStatus)
	group.PUT("/constructions/:id/accept", auth, rbac(constants.RoleAdmin, constants.RoleProjectManager, constants.RoleContractor), h.Accept)
}
