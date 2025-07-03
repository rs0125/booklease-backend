package models

import "time"

type Book struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Title       string `json:"title"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
	Category    string `json:"category"`
	FilePath    string `json:"file_path"`
	Available   bool   `json:"available"`
	UploadedBy  uint   `json:"uploaded_by"` // foreign key (user ID)

	// GORM association — joins Book to the uploading user
	Uploader User `gorm:"foreignKey:UploadedBy" json:"uploader,omitempty"`

	Type      string    `json:"type"`
	ValidTill time.Time `json:"valid_till"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Rental struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `json:"user_id"`
	BookID      *uint     `json:"book_id"`
	OwnerID     *uint     `json:"owner_id"`
	Description string    `json:"description"`
	RentedFrom  time.Time `json:"rented_from"`
	DueDate     time.Time `json:"due_date"`
	IsReturned  bool      `json:"is_returned"`
}

type Wishlist struct {
	ID uint `gorm:"primaryKey" json:"id"`

	UserID uint `json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	BookID uint `json:"book_id"`
	Book   Book `gorm:"foreignKey:BookID" json:"book,omitempty"`

	CreatedAt time.Time `json:"created_at"`

	// ValidTill time.Time `json:"valid_till"`
	// Available   bool   `json:"available"` //these can be accessed through referencing book foreign key

	// ID      uint `gorm:"primaryKey"`
	// UserID  uint
	// BookID  uint
	// AddedAt time.Time
}
