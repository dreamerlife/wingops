package cmdb

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ModelHandler struct {
	repo Repository
}

func NewModelHandler(repo Repository) *ModelHandler {
	return &ModelHandler{repo: repo}
}

func (h *ModelHandler) RegisterReadRoutes(router gin.IRouter) {
	router.GET("/cmdb/model-groups", h.ListModelGroups)
	router.GET("/cmdb/model-groups/:id/models", h.ListModels)
	router.GET("/cmdb/models/:id", h.GetModel)
}

func (h *ModelHandler) RegisterWriteRoutes(router gin.IRouter) {
	router.POST("/cmdb/model-groups", h.CreateModelGroup)
	router.POST("/cmdb/model-groups/:id/models", h.CreateModel)
	router.PUT("/cmdb/models/:id", h.UpdateModel)
	router.DELETE("/cmdb/models/:id", h.DeleteModel)
}

func (h *ModelHandler) ListModelGroups(c *gin.Context) {
	groups, err := h.repo.ListModelGroups(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "list model groups failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": groups})
}

func (h *ModelHandler) CreateModelGroup(c *gin.Context) {
	var group ModelGroup
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}
	created, err := h.repo.CreateModelGroup(c.Request.Context(), group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "create model group failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": created})
}

func (h *ModelHandler) ListModels(c *gin.Context) {
	models, err := h.repo.ListModels(c.Request.Context(), c.Param("id"))
	if err != nil {
		writeModelError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": models})
}

func (h *ModelHandler) CreateModel(c *gin.Context) {
	var model Model
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}
	model.GroupID = c.Param("id")
	created, err := h.repo.CreateModel(c.Request.Context(), model)
	if err != nil {
		writeModelError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": created})
}

func (h *ModelHandler) GetModel(c *gin.Context) {
	model, err := h.repo.GetModel(c.Request.Context(), c.Param("id"))
	if err != nil {
		writeModelError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": model})
}

func (h *ModelHandler) UpdateModel(c *gin.Context) {
	var model Model
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}
	model.ID = c.Param("id")
	updated, err := h.repo.UpdateModel(c.Request.Context(), model)
	if err != nil {
		writeModelError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": updated})
}

func (h *ModelHandler) DeleteModel(c *gin.Context) {
	if err := h.repo.DeleteModel(c.Request.Context(), c.Param("id")); err != nil {
		writeModelError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"deleted": true}})
}

func writeModelError(c *gin.Context, err error) {
	if errors.Is(err, ErrModelGroupNotFound) || errors.Is(err, ErrModelNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": err.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
}
