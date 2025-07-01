package api

import (
	"net/http"

	"bookapi/models"
	"bookapi/services"

	"github.com/gin-gonic/gin"
)

func GetFAQ(c *gin.Context) {
	var FAQ []models.FAQ
	if err := services.DB.Find(&FAQ).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch FAQ"})
		return
	}
	c.JSON(http.StatusOK, FAQ)
}
