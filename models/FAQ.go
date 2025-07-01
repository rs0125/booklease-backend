package models

type FAQ struct {
	ID       uint `gorm:"primaryKey"`
	Question string
	Answer   string
}