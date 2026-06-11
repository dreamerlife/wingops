package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"wingops/internal/audit"
	"wingops/internal/auth"
	"wingops/internal/cmdb"
	"wingops/internal/http/middleware"
	"wingops/internal/system"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	repo, err := auth.NewDevelopmentRepository()
	if err != nil {
		panic(err)
	}
	authService := auth.NewService(repo, "dev-secret-change-before-production", time.Hour)
	authHandler := auth.NewHandler(authService)
	auditRepo := audit.NewMemoryRepository()
	auditHandler := audit.NewHandler(auditRepo)
	systemRepo := system.NewMemoryRepository(system.Config{Key: "platform.name", Value: "WingOps"})
	systemHandler := system.NewHandler(systemRepo)
	cmdbRepo := cmdb.NewMemoryRepository()
	cmdbModelHandler := cmdb.NewModelHandler(cmdbRepo)
	cmdbAssetHandler := cmdb.NewAssetHandler(cmdbRepo)
	api := router.Group("/api/v1")
	authHandler.RegisterRoutes(api)

	protected := api.Group("")
	protected.Use(middleware.Auth("dev-secret-change-before-production"))
	protected.Use(middleware.Audit(auditRepo))
	authHandler.RegisterUserRoutes(protected.Group("", middleware.RequirePermission("auth.user.read")))
	authHandler.RegisterRoleRoutes(protected.Group("", middleware.RequirePermission("auth.role.read")))
	auditHandler.RegisterRoutes(protected.Group("", middleware.RequirePermission("audit.log.read")))
	systemHandler.RegisterReadRoutes(protected.Group("", middleware.RequirePermission("system.config.read")))
	systemHandler.RegisterWriteRoutes(protected.Group("", middleware.RequirePermission("system.config.write")))
	cmdbModelHandler.RegisterReadRoutes(protected.Group("", middleware.RequirePermission("cmdb.model.read")))
	cmdbModelHandler.RegisterWriteRoutes(protected.Group("", middleware.RequirePermission("cmdb.model.write")))
	cmdbAssetHandler.RegisterReadRoutes(protected.Group("", middleware.RequirePermission("cmdb.asset.read")))
	cmdbAssetHandler.RegisterWriteRoutes(protected.Group("", middleware.RequirePermission("cmdb.asset.write")))

	return router
}
