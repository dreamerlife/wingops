package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "authorization required"})
			return
		}

		tokenValue, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || tokenValue == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "invalid authorization header"})
			return
		}

		token, err := jwt.Parse(tokenValue, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenUnverifiable
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "invalid token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("permissions", permissionSetFromClaims(claims["permissions"]))
			if actorID, err := claims.GetSubject(); err == nil {
				c.Set("actor_id", actorID)
			}
		}

		c.Next()
	}
}

func RequirePermission(code string) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, ok := c.Get("permissions")
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": 403, "message": "permission denied"})
			return
		}
		permissions, ok := value.(map[string]struct{})
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": 403, "message": "permission denied"})
			return
		}
		if _, allowed := permissions[code]; !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": 403, "message": "permission denied"})
			return
		}

		c.Next()
	}
}

func permissionSetFromClaims(value any) map[string]struct{} {
	permissions := make(map[string]struct{})
	items, ok := value.([]any)
	if !ok {
		return permissions
	}
	for _, item := range items {
		code, ok := item.(string)
		if ok && code != "" {
			permissions[code] = struct{}{}
		}
	}
	return permissions
}
