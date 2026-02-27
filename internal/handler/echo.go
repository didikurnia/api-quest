package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Echo handles POST /echo — Level 2.
// It echoes back the exact JSON body received, byte-for-byte.
func Echo(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil || len(body) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	c.Data(http.StatusOK, "application/json", body)
}
