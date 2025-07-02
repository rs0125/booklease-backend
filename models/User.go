package models

import "time"

type User struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UID            string    `gorm:"uniqueIndex"`
	Username       string    `json:"username"`
	PhoneNumber    string    `json:"phone_number"`
	RegistrationNo string    `json:"registration_no"`
	IsAdmin        bool      `json:"is_admin"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
