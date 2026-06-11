package cmdb

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SyncHandler struct {
	repo   Repository
	apiKey APIKey
}

type syncRequest struct {
	ModelID  string          `json:"model_id"`
	SyncMode string          `json:"sync_mode"`
	Assets   []ImportedAsset `json:"assets"`
}

func NewSyncHandler(repo Repository, apiKey APIKey) *SyncHandler {
	return &SyncHandler{repo: repo, apiKey: apiKey}
}

func (h *SyncHandler) RegisterRoutes(router gin.IRouter) {
	router.POST("/cmdb/assets/sync", h.SyncAssets)
	router.POST("/cmdb/assets/import/preview", h.PreviewCSVImport)
}

func (h *SyncHandler) SyncAssets(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "read request failed"})
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	if !h.authorized(c, body) {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "invalid api signature"})
		return
	}

	var req syncRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}
	if req.SyncMode == "" {
		req.SyncMode = "incremental"
	}
	if req.SyncMode != "incremental" && req.SyncMode != "full" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid sync mode"})
		return
	}

	synced := make([]Asset, 0, len(req.Assets))
	for _, imported := range req.Assets {
		asset := Asset{
			ModelID:    req.ModelID,
			UniqueKey:  imported.UniqueKey,
			Attributes: imported.Attributes,
		}
		saved, err := h.repo.UpsertAsset(c.Request.Context(), asset, "api:"+h.apiKey.KeyID)
		if err != nil {
			writeAssetError(c, err)
			return
		}
		synced = append(synced, saved)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"synced": len(synced), "assets": synced}})
}

func (h *SyncHandler) PreviewCSVImport(c *gin.Context) {
	rows, err := ParseCSVAssets(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": rows})
}

func (h *SyncHandler) authorized(c *gin.Context, body []byte) bool {
	if !h.apiKey.Active() {
		return false
	}
	if c.GetHeader("X-Api-Key") != h.apiKey.KeyID {
		return false
	}
	return h.apiKey.VerifySignature(body, c.GetHeader("X-Signature"))
}
