package api

import (
	"net/http"
	"strconv"

	"bookapi/models"
	"bookapi/services"

	"github.com/gin-gonic/gin"
)

func Notifications(c *gin.Context) {
	uid := c.GetString("uid")

	var user models.User
	if err := services.DB.Where("uid = ?", uid).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found " + err.Error()})
		return
	}

	var notifs []models.Notification
	if err := services.DB.
		Where("user_id = ?", user.ID).
		Order("created_at DESC").
		Find(&notifs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifs)
}

func DeleteNotification(c *gin.Context) {
	uid := c.GetString("uid")

	// Parse notification ID
	notifIDStr := c.Param("id")
	notifID, err := strconv.ParseUint(notifIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID" + err.Error()})
		return
	}

	// Find user
	var user models.User
	if err := services.DB.Where("uid = ?", uid).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found" + err.Error()})
		return
	}

	// Find notification
	var notif models.Notification
	if err := services.DB.First(&notif, notifID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found" + err.Error()})
		return
	}

	// Only allow user to delete their own notifications
	if notif.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this notification" + err.Error()})
		return
	}

	// Delete it
	if err := services.DB.Delete(&notif).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete notification" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted"})
}

func MarkNotificationSeen(c *gin.Context) {
	uid := c.GetString("uid")
	notificationID := c.Param("id")

	var user models.User
	if err := services.DB.Where("uid = ?", uid).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var notification models.Notification
	if err := services.DB.Where("id = ? AND user_id = ?", notificationID, user.ID).First(&notification).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}

	notification.Seen = true
	if err := services.DB.Save(&notification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as seen"})
}

func DeleteAllNotifications(c *gin.Context) {
	uid := c.GetString("uid")

	var user models.User
	if err := services.DB.Where("uid = ?", uid).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if err := services.DB.Where("user_id = ?", user.ID).Delete(&models.Notification{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete notifications" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All notifications deleted"})
}
