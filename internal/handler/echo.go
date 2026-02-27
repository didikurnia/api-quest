package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Echo handles POST /echo — Level 2.
// It echoes back whatever JSON body is sent.
func Echo(c *gin.Context) {
	var body map[string]any
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	c.JSON(http.StatusOK, body)
}
