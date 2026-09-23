package router

import (
	"github.com/gin-gonic/gin"
	"github.com/home-renovation/platform/internal/constants"
	"github.com/home-renovation/platform/internal/handler"
)

// registerBudgetRoutes 注册预算路由。
func registerBudgetRoutes(group *gin.RouterGroup, h *handler.BudgetHandler, auth gin.HandlerFunc, rbac func(...string) gin.HandlerFunc) {
	group.GET("/budgets", auth, h.List)
	group.GET("/budgets/:id", auth, h.Get)
	group.POST("/budgets", auth, rbac(constants.RoleAdmin, constants.RoleProjectManager), h.Create)
	group.PUT("/budgets/:id", auth, rbac(constants.RoleAdmin, constants.RoleProjectManager), h.Update)
	group.DELETE("/budgets/:id", auth, rbac(constants.RoleAdmin), h.Delete)
}
