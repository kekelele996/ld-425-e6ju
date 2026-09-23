package router

import (
	"github.com/gin-gonic/gin"
	"github.com/home-renovation/platform/internal/handler"
)

// registerAuthRoutes 注册认证路由。
func registerAuthRoutes(group *gin.RouterGroup, authHandler *handler.AuthHandler, authMiddleware gin.HandlerFunc) {
	group.POST("/auth/login", authHandler.Login)
	group.GET("/auth/me", authMiddleware, authHandler.Me)
}
