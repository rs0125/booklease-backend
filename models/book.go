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

type Book struct {
	ID          uint `gorm:"primaryKey"`
	Title       string
	Author      string
	Description string
	Category    string
	Available   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Note struct {
	ID          uint `gorm:"primaryKey"`
	Title       string
	Subject     string
	Description string
	FilePath    string
	IsPublic    bool
	UploadedBy  uint
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

type Admin struct {
	ID        uint `gorm:"primaryKey"`
	AdminID   uint
	CreatedAt time.Time
}

type FAQ struct {
	ID       uint `gorm:"primaryKey"`
	Question string
	Answer   string
}
