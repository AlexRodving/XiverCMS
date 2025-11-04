package models

import (
	"time"

	"gorm.io/gorm"
)

// APIToken represents API tokens for programmatic access
type APIToken struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Name        string     `json:"name" gorm:"not null"`
	Token       string     `json:"token" gorm:"uniqueIndex;not null"`
	Description string     `json:"description"`
	Type        string     `json:"type" gorm:"default:read-only"` // read-only, full-access
	LastUsedAt  *time.Time `json:"lastUsedAt"`
	ExpiresAt   *time.Time `json:"expiresAt"`

	CreatedByID *uint `json:"createdById"`
	CreatedBy   *User `json:"createdBy,omitempty" gorm:"foreignKey:CreatedByID"`
}
