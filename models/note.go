package models

import "time"

type Note struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	UserID     uint       `gorm:"index;not null" json:"user_id"`
	NotebookID uint       `gorm:"index;not null" json:"notebook_id"`
	Title      string     `gorm:"not null" json:"title"`
	Content    string     `json:"content"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	Transcript Transcript `gorm:"constraint:OnDelete:CASCADE;"`
}
