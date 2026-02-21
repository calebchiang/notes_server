package models

import "time"

type Transcript struct {
	ID             uint   `gorm:"primaryKey"`
	NoteID         uint   `gorm:"uniqueIndex;not null"`
	TranscriptJSON string `gorm:"type:longtext;not null"`
	Source         string `gorm:"not null"`
	SourceID       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
