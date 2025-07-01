package models

import "time"

type User struct {
	ID             uint `gorm:"primaryKey"`
	Username       string
	Password       string
	RegistrationNo string
	IsAdmin        bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
