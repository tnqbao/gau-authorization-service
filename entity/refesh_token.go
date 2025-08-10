package entity

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Token     string    `gorm:"uniqueIndex"`
	DeviceID  string    `gorm:"uniqueIndex"`
	ExpiresAt time.Time `gorm:"index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UserID    uuid.UUID `gorm:"index;not null"`
}
