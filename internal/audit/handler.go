package audit

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	router.GET("/audit/logs", h.ListLogs)
}

func (h *Handler) ListLogs(c *gin.Context) {
	logs, err := h.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "list audit logs failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": logs})
}
