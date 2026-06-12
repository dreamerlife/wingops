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

type Dependencies struct {
	AuthRepository   auth.Repository
	AuditRepository  audit.Repository
	SystemRepository system.Repository
	CMDBRepository   cmdb.Repository
	JWTSecret        string
	TokenTTL         time.Duration
}

func NewRouter() *gin.Engine {
	repo, err := auth.NewDevelopmentRepository()
	if err != nil {
		panic(err)
	}
	return NewRouterWithDependencies(Dependencies{
		AuthRepository:   repo,
		AuditRepository:  audit.NewMemoryRepository(),
		SystemRepository: system.NewMemoryRepository(system.Config{Key: "platform.name", Value: "WingOps"}),
		CMDBRepository:   cmdb.NewMemoryRepository(),
		JWTSecret:        "dev-secret-change-before-production",
		TokenTTL:         time.Hour,
	})
}

func NewRouterWithDependencies(deps Dependencies) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	if deps.JWTSecret == "" {
		deps.JWTSecret = "dev-secret-change-before-production"
	}
	if deps.TokenTTL <= 0 {
		deps.TokenTTL = time.Hour
	}
	authService := auth.NewService(deps.AuthRepository, deps.JWTSecret, deps.TokenTTL)
	authHandler := auth.NewHandler(authService)
	auditRepo := deps.AuditRepository
	auditHandler := audit.NewHandler(auditRepo)
	systemHandler := system.NewHandler(deps.SystemRepository)
	cmdbRepo := deps.CMDBRepository
	cmdbModelHandler := cmdb.NewModelHandler(cmdbRepo)
	cmdbAssetHandler := cmdb.NewAssetHandler(cmdbRepo)
	cmdbSyncHandler := cmdb.NewSyncHandler(cmdbRepo, cmdb.NewDevelopmentAPIKey())
	api := router.Group("/api/v1")
	authHandler.RegisterRoutes(api)

	protected := api.Group("")
	protected.Use(middleware.Auth(deps.JWTSecret))
	protected.Use(middleware.Audit(auditRepo))
	authHandler.RegisterUserRoutes(protected.Group("", middleware.RequirePermission("auth.user.read")))
	authHandler.RegisterUserWriteRoutes(protected.Group("", middleware.RequirePermission("auth.user.write")))
	authHandler.RegisterRoleRoutes(protected.Group("", middleware.RequirePermission("auth.role.read")))
	authHandler.RegisterRoleWriteRoutes(protected.Group("", middleware.RequirePermission("auth.role.write")))
	auditHandler.RegisterRoutes(protected.Group("", middleware.RequirePermission("audit.log.read")))
	systemHandler.RegisterReadRoutes(protected.Group("", middleware.RequirePermission("system.config.read")))
	systemHandler.RegisterWriteRoutes(protected.Group("", middleware.RequirePermission("system.config.write")))
	cmdbModelHandler.RegisterReadRoutes(protected.Group("", middleware.RequirePermission("cmdb.model.read")))
	cmdbModelHandler.RegisterWriteRoutes(protected.Group("", middleware.RequirePermission("cmdb.model.write")))
	cmdbAssetHandler.RegisterReadRoutes(protected.Group("", middleware.RequirePermission("cmdb.asset.read")))
	cmdbAssetHandler.RegisterWriteRoutes(protected.Group("", middleware.RequirePermission("cmdb.asset.write")))
	cmdbSyncHandler.RegisterAPIKeyReadRoutes(protected.Group("", middleware.RequirePermission("cmdb.apikey.read")))
	cmdbSyncHandler.RegisterAPIKeyWriteRoutes(protected.Group("", middleware.RequirePermission("cmdb.apikey.write")))
	cmdbSyncHandler.RegisterRoutes(api)

	return router
}
