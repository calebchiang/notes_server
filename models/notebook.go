package models

import "time"

type Notebook struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"index;not null"`
	Title     string `gorm:"not null"`
	Color     string
	CreatedAt time.Time
	UpdatedAt time.Time

	Notes []Note `gorm:"constraint:OnDelete:CASCADE;"`
}
