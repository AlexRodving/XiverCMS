package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Email        string `json:"email" gorm:"uniqueIndex;not null"`
	Username     string `json:"username" gorm:"uniqueIndex;not null"`
	Password     string `json:"-" gorm:"not null"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	IsActive     bool   `json:"isActive" gorm:"default:true"`
	IsSuperAdmin bool   `json:"isSuperAdmin" gorm:"default:false"`

	Roles []Role `json:"roles,omitempty" gorm:"many2many:user_roles;"`
}

type Role struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Name        string `json:"name" gorm:"uniqueIndex;not null"`
	Description string `json:"description"`
	Type        string `json:"type" gorm:"default:custom"` // public, custom

	Users       []User       `json:"users,omitempty" gorm:"many2many:user_roles;"`
	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
}

type Permission struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Action     string `json:"action" gorm:"not null"`                // create, read, update, delete, publish
	Subject    string `json:"subject" gorm:"not null"`               // all, user, role, content-type, etc.
	Properties string `json:"properties,omitempty" gorm:"type:json"` // JSON field for additional properties
	Conditions string `json:"conditions,omitempty" gorm:"type:json"` // JSON field for conditions

	Roles []Role `json:"roles,omitempty" gorm:"many2many:role_permissions;"`
}
