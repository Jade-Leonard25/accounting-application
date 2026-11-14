package models

import (
	"gorm.io/gorm"
)

type AdminRole string

const (
	RoleSuperAdmin AdminRole = "SUPER_ADMIN"
	RoleAdmin      AdminRole = "ADMIN"
	RoleManager    AdminRole = "MANAGER"
)

type Admin struct {
	gorm.Model
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	FirstName    string    `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName     string    `gorm:"type:varchar(100);not null" json:"last_name"`
	Role         AdminRole `gorm:"type:varchar(50);default:'ADMIN';not null" json:"role"`

	// --- Relations ---
	// An admin owns multiple companies. If an admin is deleted, their companies' AdminID is set to NULL.
	Companies []Company `gorm:"foreignKey:AdminID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"companies,omitempty"`
}