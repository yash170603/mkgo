package models

import (
	"time"

	"github.com/google/uuid"
)

type Patient struct {
	ID                   uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name                 string         `gorm:"not null" json:"name"`
	Email                string         `gorm:"uniqueIndex;not null" json:"email"`
	Phone                string         `gorm:"not null" json:"phone"`
	Address              string         `gorm:"not null" json:"address"`
	CreatedAt            time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	PatientPrescriptions []Prescription `gorm:"foreignKey:PatientID" json:"patient_prescriptions,omitempty"`
	PatientAppointments  []Appointment  `gorm:"foreignKey:PatientID" json:"patient_appointments,omitempty"`
	UserID               uuid.UUID      `gorm:"not null" json:"user_id"`                            // Foreign key to the doctor treating the patient
	User                 User           `gorm:"foreignKey:UserID" json:"doctor_assigned,omitempty"` // Relationship to User model
}
