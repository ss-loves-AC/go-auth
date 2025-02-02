package models

import "time"

type User struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"unique;not null"`
	HashPassword string    `json:"-" gorm:"not null"` // don't include this field in encoding or decoding
	CreatedAt    time.Time `json:"created_at"`
}


