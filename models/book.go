package models

import "time"

type Book struct {
	ID          uint `gorm:"primaryKey"`
	Title       string
	Type        string
	Subject     string
	Author      string
	Description string
	Category    string
	Available   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Rental struct {
	ID         uint `gorm:"primaryKey"`
	UserID     uint
	BookID     *uint
	NotesID    *uint
	RentedFrom time.Time
	DueDate    time.Time
	IsReturned bool
}

type Wishlist struct {
	ID      uint `gorm:"primaryKey"`
	UserID  uint
	BookID  uint
	AddedAt time.Time
}
