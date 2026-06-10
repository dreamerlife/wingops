package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"wingops/internal/audit"
)

func Audit(repo audit.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		actorID, _ := c.Get("actor_id")
		_ = repo.Append(c.Request.Context(), audit.Log{
			ActorID:    stringValue(actorID),
			Method:     c.Request.Method,
			Path:       c.Request.URL.Path,
			StatusCode: c.Writer.Status(),
			Resource:   resourceFromPath(c.Request.URL.Path),
			CreatedAt:  time.Now(),
		})
	}
}

func resourceFromPath(path string) string {
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")
	if len(parts) >= 3 {
		return parts[2]
	}
	return ""
}

func stringValue(value any) string {
	text, ok := value.(string)
	if !ok {
		return ""
	}
	return text
}
