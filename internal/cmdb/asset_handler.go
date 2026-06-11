package cmdb

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AssetHandler struct {
	repo Repository
}

func NewAssetHandler(repo Repository) *AssetHandler {
	return &AssetHandler{repo: repo}
}

func (h *AssetHandler) RegisterReadRoutes(router gin.IRouter) {
	router.GET("/cmdb/assets", h.ListAssets)
	router.GET("/cmdb/assets/:id", h.GetAsset)
	router.GET("/cmdb/assets/:id/change-logs", h.ListChangeLogs)
}

func (h *AssetHandler) RegisterWriteRoutes(router gin.IRouter) {
	router.POST("/cmdb/assets", h.CreateAsset)
	router.PUT("/cmdb/assets/:id", h.UpdateAsset)
	router.DELETE("/cmdb/assets/:id", h.DeleteAsset)
}

func (h *AssetHandler) ListAssets(c *gin.Context) {
	assets, err := h.repo.ListAssets(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "list assets failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": assets})
}

func (h *AssetHandler) CreateAsset(c *gin.Context) {
	var asset Asset
	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}
	created, err := h.repo.CreateAsset(c.Request.Context(), asset, actorID(c))
	if err != nil {
		writeAssetError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": created})
}

func (h *AssetHandler) GetAsset(c *gin.Context) {
	asset, err := h.repo.GetAsset(c.Request.Context(), c.Param("id"))
	if err != nil {
		writeAssetError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": asset})
}

func (h *AssetHandler) UpdateAsset(c *gin.Context) {
	var asset Asset
	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}
	asset.ID = c.Param("id")
	updated, err := h.repo.UpdateAsset(c.Request.Context(), asset, actorID(c))
	if err != nil {
		writeAssetError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": updated})
}

func (h *AssetHandler) DeleteAsset(c *gin.Context) {
	if err := h.repo.DeleteAsset(c.Request.Context(), c.Param("id")); err != nil {
		writeAssetError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"deleted": true}})
}

func (h *AssetHandler) ListChangeLogs(c *gin.Context) {
	logs, err := h.repo.ListAssetChangeLogs(c.Request.Context(), c.Param("id"))
	if err != nil {
		writeAssetError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": logs})
}

func writeAssetError(c *gin.Context, err error) {
	if errors.Is(err, ErrModelNotFound) || errors.Is(err, ErrAssetNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": err.Error()})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
}

func actorID(c *gin.Context) string {
	value, ok := c.Get("actor_id")
	if !ok {
		return ""
	}
	actor, ok := value.(string)
	if !ok {
		return ""
	}
	return actor
}
