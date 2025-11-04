package models

import (
	"time"

	"gorm.io/gorm"
)

// ContentRelation represents a relationship between content entries
type ContentRelation struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Source entry
	SourceContentTypeUID string `json:"sourceContentTypeUid" gorm:"not null;index"`
	SourceEntryID        uint   `json:"sourceEntryId" gorm:"not null;index"`
	SourceFieldName      string `json:"sourceFieldName" gorm:"not null"` // Field name in schema

	// Target entry
	TargetContentTypeUID string `json:"targetContentTypeUid" gorm:"not null;index"`
	TargetEntryID        uint   `json:"targetEntryId" gorm:"not null;index"`

	// Relation type: oneToOne, oneToMany, manyToOne, manyToMany
	RelationType string `json:"relationType" gorm:"not null;default:manyToOne"`

	// Order for oneToMany and manyToMany
	Order int `json:"order" gorm:"default:0"`
}

// ContentField defines a field in content type schema
type ContentField struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"` // string, text, number, boolean, date, relation, media, component
	Required    bool        `json:"required"`
	Unique      bool        `json:"unique"`
	Default     interface{} `json:"default,omitempty"`
	Description string      `json:"description,omitempty"`

	// For relation type
	RelationType      string `json:"relationType,omitempty"`      // oneToOne, oneToMany, manyToOne, manyToMany
	TargetContentType string `json:"targetContentType,omitempty"` // UID of target content type

	// For media type
	Multiple bool `json:"multiple,omitempty"` // Single file or multiple files

	// For component type
	ComponentType string `json:"componentType,omitempty"` // UID of component type

	// Validation
	MinLength *int     `json:"minLength,omitempty"`
	MaxLength *int     `json:"maxLength,omitempty"`
	Min       *float64 `json:"min,omitempty"`
	Max       *float64 `json:"max,omitempty"`
	Pattern   string   `json:"pattern,omitempty"` // Regex pattern

	// Options for enum/select
	Options []map[string]interface{} `json:"options,omitempty"`
}

// ComponentType represents reusable component schemas
type ComponentType struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	UID         string         `json:"uid" gorm:"uniqueIndex;not null"`
	DisplayName string         `json:"displayName" gorm:"not null"`
	Description string         `json:"description"`
	Category    string         `json:"category"` // For grouping components
	Fields      []ContentField `json:"fields" gorm:"type:jsonb"`
}
