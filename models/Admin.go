package models

import "time"

type Admin struct {
	ID        uint `gorm:"primaryKey"`
	AdminID   uint
	CreatedAt time.Time
}