package models

import (
	"time"

	"github.com/google/uuid"
)

type Appointment struct {
	ID              uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	PatientID       uuid.UUID `gorm:"not null" json:"patient_id"`
	DoctorID        uuid.UUID `gorm:"not null" json:"doctor_id"`
	AppointmentDate time.Time `gorm:"not null" json:"appointment_date"`
	Status          string    `gorm:"type:text CHECK (status IN ('scheduled','completed','cancelled'));default:'scheduled'" json:"status"`
	Notes           string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Patient Patient `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	Doctor  User    `gorm:"foreignKey:DoctorID" json:"doctor,omitempty"`
}
