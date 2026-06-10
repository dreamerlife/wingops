package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"wingops/internal/auth"
	"wingops/internal/http/middleware"
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
	api := router.Group("/api/v1")
	authHandler.RegisterRoutes(api)

	protected := api.Group("")
	protected.Use(middleware.Auth("dev-secret-change-before-production"))
	authHandler.RegisterUserRoutes(protected.Group("", middleware.RequirePermission("auth.user.read")))
	authHandler.RegisterRoleRoutes(protected.Group("", middleware.RequirePermission("auth.role.read")))

	return router
}
