package handler

import (
	"net/http"
	"time"

	"github.com/didikurnia/api-quest/internal/config"
	"github.com/didikurnia/api-quest/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	cfg *config.Config
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

// Token handles POST /auth/token — Level 5.
// Issues a JWT for valid credentials (admin/password).
func (h *AuthHandler) Token(c *gin.Context) {
	var req model.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	// Hardcoded credentials as expected by the challenge
	if req.Username != "admin" || req.Password != "password" {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "invalid credentials"})
		return
	}

	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   req.Username,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(h.cfg.JWTExpiry)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, model.TokenResponse{Token: signed})
}
