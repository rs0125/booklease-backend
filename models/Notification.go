package models

import "time"

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	Seen      bool      `json:"seen"`
	CreatedAt time.Time `json:"created_at"`
}
