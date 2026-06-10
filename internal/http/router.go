package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"wingops/internal/auth"
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
	authHandler.RegisterRoutes(router.Group("/api/v1"))

	return router
}
