package models

import "time"

type Notebook struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Title     string    `gorm:"not null" json:"title"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Notes []Note `gorm:"constraint:OnDelete:CASCADE;" json:"notes"`
}
