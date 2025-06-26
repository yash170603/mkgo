package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name         string    `gorm:"not null" json:"name"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	Role         string    `gorm:"type:text CHECK (role IN ('receptionist','doctor'));not null" json:"role"`
	PasswordHash string    `gorm:"not null" json:"-"`

	// Relationships
	Patients      []Patient      `gorm:"foreignKey:UserID" json:"patients,omitempty"`
	Prescriptions []Prescription `gorm:"foreignKey:DoctorID" json:"prescriptions,omitempty"`
	Appointments  []Appointment  `gorm:"foreignKey:DoctorID" json:"appointments,omitempty"`
}
