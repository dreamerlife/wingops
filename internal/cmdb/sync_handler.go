package cmdb

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SyncHandler struct {
	repo Repository
}

type syncRequest struct {
	ModelID  string          `json:"model_id"`
	SyncMode string          `json:"sync_mode"`
	Assets   []ImportedAsset `json:"assets"`
}

func NewSyncHandler(repo Repository, apiKey APIKey) *SyncHandler {
	return &SyncHandler{repo: repo}
}

func (h *SyncHandler) RegisterRoutes(router gin.IRouter) {
	router.POST("/cmdb/assets/sync", h.SyncAssets)
	router.POST("/cmdb/assets/import/preview", h.PreviewCSVImport)
}

func (h *SyncHandler) RegisterAPIKeyReadRoutes(router gin.IRouter) {
	router.GET("/cmdb/api-keys", h.ListAPIKeys)
}

func (h *SyncHandler) RegisterAPIKeyWriteRoutes(router gin.IRouter) {
	router.POST("/cmdb/api-keys", h.CreateAPIKey)
	router.DELETE("/cmdb/api-keys/:id", h.RevokeAPIKey)
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
	apiKeyID, _ := c.Get("api_key_id")

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
		saved, err := h.repo.UpsertAsset(c.Request.Context(), asset, "api:"+stringValue(apiKeyID))
		if err != nil {
			writeAssetError(c, err)
			return
		}
		synced = append(synced, saved)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"synced": len(synced), "assets": synced}})
}

func (h *SyncHandler) ListAPIKeys(c *gin.Context) {
	keys, err := h.repo.ListAPIKeys(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "list api keys failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": keys})
}

func (h *SyncHandler) CreateAPIKey(c *gin.Context) {
	var key APIKey
	if err := c.ShouldBindJSON(&key); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}
	created, err := h.repo.CreateAPIKey(c.Request.Context(), key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "create api key failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": created})
}

func (h *SyncHandler) RevokeAPIKey(c *gin.Context) {
	if err := h.repo.RevokeAPIKey(c.Request.Context(), c.Param("id")); err != nil {
		if err == ErrAPIKeyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "revoke api key failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"revoked": true}})
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
	key, err := h.repo.GetAPIKeyByKeyID(c.Request.Context(), c.GetHeader("X-Api-Key"))
	if err != nil || !key.Active() {
		return false
	}
	c.Set("api_key_id", key.KeyID)
	return key.VerifySignature(body, c.GetHeader("X-Signature"))
}

func stringValue(value any) string {
	text, ok := value.(string)
	if !ok {
		return ""
	}
	return text
}
