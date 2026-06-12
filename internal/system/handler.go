package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo Repository
}

type updateConfigRequest struct {
	Value string `json:"value"`
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	h.RegisterReadRoutes(router)
	h.RegisterWriteRoutes(router)
}

func (h *Handler) RegisterReadRoutes(router gin.IRouter) {
	router.GET("/system/configs", h.ListConfigs)
}

func (h *Handler) RegisterWriteRoutes(router gin.IRouter) {
	router.PUT("/system/configs/:key", h.UpdateConfig)
}

func (h *Handler) ListConfigs(c *gin.Context) {
	configs, err := h.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "list configs failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": configs})
}

func (h *Handler) UpdateConfig(c *gin.Context) {
	var req updateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}

	config := Config{Key: c.Param("key"), Value: req.Value}
	if !config.Valid() {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "config key is required"})
		return
	}
	saved, err := h.repo.Save(c.Request.Context(), config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "save config failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": saved})
}
