package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequirePermission(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/write", func(c *gin.Context) {
		c.Set("permissions", map[string]struct{}{"cmdb.asset.read": {}})
	}, RequirePermission("cmdb.asset.write"), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	router.GET("/read", func(c *gin.Context) {
		c.Set("permissions", map[string]struct{}{"cmdb.asset.read": {}})
	}, RequirePermission("cmdb.asset.read"), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	writeRec := httptest.NewRecorder()
	router.ServeHTTP(writeRec, httptest.NewRequest(http.MethodGet, "/write", nil))
	if writeRec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for missing write permission, got %d", writeRec.Code)
	}

	readRec := httptest.NewRecorder()
	router.ServeHTTP(readRec, httptest.NewRequest(http.MethodGet, "/read", nil))
	if readRec.Code != http.StatusOK {
		t.Fatalf("expected 200 for read permission, got %d", readRec.Code)
	}
}
