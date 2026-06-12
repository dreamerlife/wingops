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

type userRequest struct {
	Username    string   `json:"username" binding:"required"`
	Password    string   `json:"password"`
	DisplayName string   `json:"display_name" binding:"required"`
	Status      string   `json:"status"`
	RoleNames   []string `json:"role_names"`
}

type roleRequest struct {
	Name            string   `json:"name" binding:"required"`
	DisplayName     string   `json:"display_name" binding:"required"`
	PermissionCodes []string `json:"permission_codes"`
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

func (h *Handler) RegisterUserWriteRoutes(router gin.IRouter) {
	router.POST("/auth/users", h.CreateUser)
	router.PUT("/auth/users/:id", h.UpdateUser)
	router.DELETE("/auth/users/:id", h.DeleteUser)
}

func (h *Handler) RegisterRoleRoutes(router gin.IRouter) {
	router.GET("/auth/roles", h.ListRoles)
	router.GET("/auth/permissions", h.ListPermissions)
}

func (h *Handler) RegisterRoleWriteRoutes(router gin.IRouter) {
	router.POST("/auth/roles", h.CreateRole)
	router.PUT("/auth/roles/:name", h.UpdateRole)
	router.DELETE("/auth/roles/:name", h.DeleteRole)
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
	users, err := h.service.ListUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "list users failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": users})
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req userRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}
	user, err := h.service.CreateUser(c.Request.Context(), User{
		Username:    req.Username,
		DisplayName: req.DisplayName,
		Status:      req.Status,
	}, req.Password, req.RoleNames)
	if err != nil {
		writeAuthManagementError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": user})
}

func (h *Handler) UpdateUser(c *gin.Context) {
	var req userRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}
	roles := make([]Role, 0, len(req.RoleNames))
	for _, name := range req.RoleNames {
		roles = append(roles, Role{Name: name})
	}
	user, err := h.service.UpdateUser(c.Request.Context(), User{
		ID:          c.Param("id"),
		Username:    req.Username,
		DisplayName: req.DisplayName,
		Status:      req.Status,
		Roles:       roles,
	}, req.Password)
	if err != nil {
		writeAuthManagementError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": user})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	if err := h.service.DeleteUser(c.Request.Context(), c.Param("id")); err != nil {
		writeAuthManagementError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"deleted": true}})
}

func (h *Handler) ListRoles(c *gin.Context) {
	roles, err := h.service.ListRoles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "list roles failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": roles})
}

func (h *Handler) CreateRole(c *gin.Context) {
	var req roleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}
	role, err := h.service.CreateRole(c.Request.Context(), roleFromRequest(req))
	if err != nil {
		writeAuthManagementError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": role})
}

func (h *Handler) UpdateRole(c *gin.Context) {
	var req roleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}
	req.Name = c.Param("name")
	role, err := h.service.UpdateRole(c.Request.Context(), roleFromRequest(req))
	if err != nil {
		writeAuthManagementError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": role})
}

func (h *Handler) DeleteRole(c *gin.Context) {
	if err := h.service.DeleteRole(c.Request.Context(), c.Param("name")); err != nil {
		writeAuthManagementError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"deleted": true}})
}

func (h *Handler) ListPermissions(c *gin.Context) {
	permissions, err := h.service.ListPermissions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "list permissions failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": permissions})
}

func roleFromRequest(req roleRequest) Role {
	permissions := make([]Permission, 0, len(req.PermissionCodes))
	for _, code := range req.PermissionCodes {
		permissions = append(permissions, Permission{Code: code})
	}
	return Role{Name: req.Name, DisplayName: req.DisplayName, Permissions: permissions}
}

func writeAuthManagementError(c *gin.Context, err error) {
	if errors.Is(err, ErrUserNotFound) || errors.Is(err, ErrRoleNotFound) || errors.Is(err, ErrPermissionNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": err.Error()})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
}
