package models

import "time"

type Note struct {
	ID         uint   `gorm:"primaryKey"`
	UserID     uint   `gorm:"index;not null"`
	NotebookID uint   `gorm:"index;not null"`
	Title      string `gorm:"not null"`
	Content    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
