package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping handles GET /ping — Level 1.
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
