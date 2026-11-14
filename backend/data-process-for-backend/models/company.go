package models

import (
	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	Name           string `gorm:"type:varchar(255);not null" json:"name"`
	RegistrationNo string `gorm:"type:varchar(100);uniqueIndex" json:"registration_no"`
	Address        string `gorm:"type:text" json:"address"`

	// --- Admin Control ---
	AdminID *uint  `json:"admin_id"`
	Admin   *Admin `gorm:"foreignKey:AdminID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"admin,omitempty"`

	// --- CEO Management (Optional) ---
	CEOID *uint     `json:"ceo_id"`
	CEO   *Employee `gorm:"foreignKey:CEOID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"ceo,omitempty"`

	// --- Employees Relation ---
	// If a company is deleted, its employees are also cleaned up or unlinked
	Employees []Employee `gorm:"foreignKey:CompanyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"employees,omitempty"`
}