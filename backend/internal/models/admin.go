package models

import "time"

type Admin struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null;size:255"`
	PasswordHash string    `json:"-" gorm:"not null"`
	Name         string    `json:"name" gorm:"size:255"`
	CreatedAt    time.Time `json:"created_at"`
}
