package router

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/home-renovation/platform/internal/config"
	"github.com/home-renovation/platform/internal/handler"
	"github.com/home-renovation/platform/internal/middleware"
	"github.com/home-renovation/platform/internal/repository"
	"github.com/home-renovation/platform/internal/service"
)

// Deps 路由装配依赖。
type Deps struct {
	Config     *config.Config
	Logger     *slog.Logger
	UserSvc    service.UserService
	AuditSvc   service.AuditService
	AuditRepo  repository.AuditLogRepository
	ProjectH   *handler.ProjectHandler
	DesignH    *handler.DesignHandler
	MaterialH  *handler.MaterialHandler
	BudgetH    *handler.BudgetHandler
	ConstructionH *handler.ConstructionHandler
	AuditH     *handler.AuditHandler
	UploadH    *handler.UploadHandler
}

// New 构建并配置 Gin 引擎。
func New(deps Deps) *gin.Engine {
	if deps.Config.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(middleware.RequestLoggerMiddleware(deps.Logger))
	engine.Use(middleware.ErrorHandlerMiddleware(deps.Logger))
	engine.Use(middleware.AuditLogMiddleware(deps.AuditRepo))

	engine.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 上传文件静态访问。
	if deps.Config.Upload.Dir != "" {
		engine.Static(deps.Config.Upload.PublicPrefix, deps.Config.Upload.Dir)
	}

	api := engine.Group("/api/v1")
	auth := middleware.AuthMiddleware(deps.UserSvc)

	registerAuthRoutes(api, handler.NewAuthHandler(deps.UserSvc), auth)
	registerProjectRoutes(api, deps.ProjectH, auth, middleware.RBACMiddleware)
	registerDesignRoutes(api, deps.DesignH, auth, middleware.RBACMiddleware)
	registerMaterialRoutes(api, deps.MaterialH, auth, middleware.RBACMiddleware)
	registerBudgetRoutes(api, deps.BudgetH, auth, middleware.RBACMiddleware)
	registerConstructionRoutes(api, deps.ConstructionH, auth, middleware.RBACMiddleware)

	uploadGroup := api.Group("/upload", auth, middleware.UploadMiddleware(deps.Config.Upload))
	uploadGroup.POST("", deps.UploadH.Upload)

	api.GET("/audit-logs", auth, middleware.RBACMiddleware("Admin"), deps.AuditH.List)

	return engine
}
