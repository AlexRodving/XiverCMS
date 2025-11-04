package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// JSONB type for PostgreSQL JSON fields
type JSONB map[string]interface{}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONB)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, j)
}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// ContentType defines the structure of content models
type ContentType struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	UID         string `json:"uid" gorm:"uniqueIndex;not null"`
	Kind        string `json:"kind" gorm:"default:collectionType"` // collectionType, singleType
	DisplayName string `json:"displayName" gorm:"not null"`
	Description string `json:"description"`
	IsVisible   bool   `json:"isVisible" gorm:"default:true"`
	AccessType  string `json:"accessType" gorm:"default:public"` // public, authenticated, moderator, admin

	// Schema definition stored as JSON
	Schema JSONB `json:"schema" gorm:"type:jsonb"`

	Entries []ContentEntry `json:"entries,omitempty" gorm:"foreignKey:ContentTypeID"`
}

// ContentEntry represents an entry of a content type
type ContentEntry struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	PublishedAt *time.Time     `json:"publishedAt"`

	ContentTypeID uint        `json:"contentTypeId" gorm:"not null"`
	ContentType   ContentType `json:"contentType,omitempty" gorm:"foreignKey:ContentTypeID"`

	// Dynamic fields stored as JSON
	Data JSONB `json:"data" gorm:"type:jsonb"`

	// Status: draft, published
	Status string `json:"status" gorm:"default:draft"`

	// Created by user
	CreatedByID *uint `json:"createdById"`
	CreatedBy   *User `json:"createdBy,omitempty" gorm:"foreignKey:CreatedByID"`

	UpdatedByID *uint `json:"updatedById"`
	UpdatedBy   *User `json:"updatedBy,omitempty" gorm:"foreignKey:UpdatedByID"`
}

// MediaFile represents uploaded media files
type MediaFile struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Name        string `json:"name" gorm:"not null"`
	Alternative string `json:"alternative"`
	Caption     string `json:"caption"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Format      string `json:"format"` // jpg, png, pdf, etc.
	Mime        string `json:"mime" gorm:"not null"`
	Size        int64  `json:"size"` // size in bytes
	URL         string `json:"url" gorm:"not null"`
	Path        string `json:"path" gorm:"not null"`

	Provider string `json:"provider" gorm:"default:local"` // local, aws-s3, cloudinary, etc.

	CreatedByID *uint `json:"createdById"`
	CreatedBy   *User `json:"createdBy,omitempty" gorm:"foreignKey:CreatedByID"`
}
