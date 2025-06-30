package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HelloHandler(c *gin.Context) {
	c.String(http.StatusOK, "✅ Authenticated! Welcome to the protected route.")
}
