package models

import "time"

// separate table for both refresh and revoked token
type RefreshToken struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	JTI       string    `json:"jti" gorm:"unique;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"index"`
}

type RevokedToken struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	JTI       string    `json:"jti" gorm:"unique;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"index"`
}
