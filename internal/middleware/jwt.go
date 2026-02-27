package middleware

import (
	"net/http"
	"strings"

	"github.com/didikurnia/api-quest/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth returns a middleware that strictly requires a valid Bearer token.
func JWTAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !validateToken(c, cfg) {
			return
		}
		c.Next()
	}
}

// OptionalJWTAuth returns a middleware that validates a token ONLY if one is provided.
// If no Authorization header is present, the request proceeds without auth.
// If a header IS present but invalid, the request is rejected with 401.
func OptionalJWTAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.Next()
			return
		}
		if !validateToken(c, cfg) {
			return
		}
		c.Next()
	}
}

// validateToken parses and validates a Bearer token. Returns false if invalid (response already sent).
func validateToken(c *gin.Context, cfg *config.Config) bool {
	header := c.GetHeader("Authorization")
	if header == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
		return false
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
		return false
	}

	tokenStr := parts[1]
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		c.Set("user", claims)
	}

	return true
}
