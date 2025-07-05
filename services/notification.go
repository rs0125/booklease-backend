package services

import (
	"bookapi/models"
	"time"
)

func CreateNotification(userID uint, notifType, message string) error {
	notif := models.Notification{
		UserID:    userID,
		Type:      notifType,
		Message:   message,
		Seen:      false,
		CreatedAt: time.Now(),
	}
	return DB.Create(&notif).Error
}
