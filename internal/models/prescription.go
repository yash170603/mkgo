package models

import (
	"time"

	"github.com/google/uuid"
)

type Prescription struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	PatientID    uuid.UUID `gorm:"not null" json:"patient_id"`
	DoctorID     uuid.UUID `gorm:"not null" json:"doctor_id"`
	Medication   string    `gorm:"not null" json:"medication"`
	Dosage       string    `gorm:"not null" json:"dosage"`
	Instructions string    `gorm:"type:text" json:"instructions"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Patient Patient `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	Doctor  User    `gorm:"foreignKey:DoctorID" json:"doctor,omitempty"`
}
