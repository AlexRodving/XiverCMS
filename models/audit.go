package models

import (
	"time"

	"gorm.io/gorm"
)

// AuditLog represents audit logs for tracking user actions
type AuditLog struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Action      string `json:"action" gorm:"not null"`  // create, update, delete, login, etc.
	Subject     string `json:"subject" gorm:"not null"` // user, content-type, entry, etc.
	SubjectID   *uint  `json:"subjectId"`               // ID of the subject
	Description string `json:"description"`
	IPAddress   string `json:"ipAddress"`
	UserAgent   string `json:"userAgent"`

	UserID *uint `json:"userId"`
	User   *User `json:"user,omitempty" gorm:"foreignKey:UserID"`

	// Additional data stored as JSON
	Metadata JSONB `json:"metadata,omitempty" gorm:"type:jsonb"`
}

// ContentHistory represents version history of content entries
type ContentHistory struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	ContentEntryID uint         `json:"contentEntryId" gorm:"not null"`
	ContentEntry   ContentEntry `json:"contentEntry,omitempty" gorm:"foreignKey:ContentEntryID"`

	// Snapshot of the entry at this point in time
	Data JSONB `json:"data" gorm:"type:jsonb"`

	// Change information
	ChangeType string `json:"changeType"` // created, updated, published, unpublished
	ChangeNote string `json:"changeNote"`

	ChangedByID *uint `json:"changedById"`
	ChangedBy   *User `json:"changedBy,omitempty" gorm:"foreignKey:ChangedByID"`
}
