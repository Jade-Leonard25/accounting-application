package models

import (
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	ID                   string  `gorm:"type:varchar(50);primaryKey" json:"id"`
	FirstName            string  `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName             string  `gorm:"type:varchar(100);not null" json:"last_name"`
	Email                string  `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	ImageURL             *string `gorm:"type:text" json:"image"`
	UniqueIdentification string  `gorm:"type:varchar(20)" json:"uniqeIdentification"`

	// --- Address ---
	Address              string  `gorm:"type:text" json:"address"`

	// --- Employment & Schedule Metadata ---
	EmploymentStatus string    `gorm:"type:varchar(50);default:'Full-time'" json:"employment_status"` // Full-time, Part-time, Contractor
	PayFrequency     string    `gorm:"type:varchar(30);default:'Semi-monthly'" json:"pay_frequency"`  // Monthly, Bi-weekly, Semi-monthly
	HireDate         time.Time `json:"hire_date"`

	// --- Compensation (All values in Cents) ---
	BaseSalaryCents int `gorm:"not null" json:"base_salary_cents"`
	HourlyRateCents int `gorm:"default:0" json:"hourly_rate_cents"` // For hourly or overtime calculation

	// --- Allowances (All values in Cents) ---
	HousingAllowanceCents   int `gorm:"default:0" json:"housing_allowance_cents"`
	TransportAllowanceCents int `gorm:"default:0" json:"transport_allowance_cents"`
	MealAllowanceCents      int `gorm:"default:0" json:"meal_allowance_cents"`

	// --- Tax & Statutory Identification ---
	TaxID            string `gorm:"type:varchar(50)" json:"tax_id"`
	SocialSecurityNo string `gorm:"type:varchar(50)" json:"social_security_no"`
	TaxFilingStatus  string `gorm:"type:varchar(30);default:'Single'" json:"tax_filing_status"` // Single, Married, Head of Household

	// --- Benefits & Deductions (All values in Cents) ---
	HealthInsuranceCents int `gorm:"default:0" json:"health_insurance_cents"`
	RetirementFundCents  int `gorm:"default:0" json:"retirement_fund_cents"`  // 401k, Pension, or Provident fund
	OtherDeductionsCents int `gorm:"default:0" json:"other_deductions_cents"` // Company loans, cash advances, etc.
}