package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	router.POST("/auth/login", h.Login)
}

func (h *Handler) RegisterUserRoutes(router gin.IRouter) {
	router.GET("/auth/users", h.ListUsers)
}

func (h *Handler) RegisterRoleRoutes(router gin.IRouter) {
	router.GET("/auth/roles", h.ListRoles)
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}

	token, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "invalid username or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "login failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": token})
}

func (h *Handler) ListUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []gin.H{
		{
			"username":     "admin",
			"display_name": "管理员",
			"status":       "active",
		},
	}})
}

func (h *Handler) ListRoles(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []gin.H{
		{"name": "system_admin", "display_name": "系统管理员"},
		{"name": "ops_admin", "display_name": "运维管理员"},
		{"name": "ops_operator", "display_name": "运维操作员"},
		{"name": "readonly", "display_name": "只读用户"},
	}})
}
