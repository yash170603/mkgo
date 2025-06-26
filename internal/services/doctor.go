package services

import (
	"errors"
	"fmt"
	"hospital/internal/config"
	"hospital/internal/database"
	"hospital/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DoctorService interface {
	GetPatients(doctorID uuid.UUID) ([]models.Patient, error)
	CreatePrescription(patientID uuid.UUID, prescription *models.Prescription) error
	UpdatePrescription(patientID uuid.UUID, prescription *models.Prescription) error
	GetAppointments(doctorID uuid.UUID) ([]models.Appointment, error)
	GetAppointmentsByDate(doctorID uuid.UUID, date time.Time) ([]models.Appointment, error)
}

type doctorService struct {
	db  *database.DB
	cfg config.Config
}

func NewDoctorService(db *database.DB, cfg config.Config) DoctorService {
	return &doctorService{
		db:  db,
		cfg: cfg,
	}
}

func (s *doctorService) GetPatients(doctorID uuid.UUID) ([]models.Patient, error) {
	var patients []models.Patient
	err := s.db.Conn.Where("user_id = ?", doctorID).Find(&patients).Error
	if err != nil {
		return nil, err
	}
	return patients, nil
}

func (s *doctorService) CreatePrescription(patientID uuid.UUID, prescription *models.Prescription) error {
	// Verify patient exists
	var patient models.Patient
	if err := s.db.Conn.First(&patient, "id = ?", patientID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("patient not found")
		}
		return err
	}
	prescription.PatientID = patientID
	return s.db.Conn.Create(prescription).Error
}

// func (s *doctorService) UpdatePrescription(patientID uuid.UUID, prescription *models.Prescription) error {
// 	// Verify patient exists
// 	var patient models.Patient
// 	if err := s.db.Conn.First(&patient, "id = ?", patientID).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return errors.New("patient not found")
// 		}
// 		return err
// 	}

// 	return s.db.Conn.Model(prescription).Where("patient_id = ?", patientID).Updates(prescription).Error
// }

func (s *doctorService) UpdatePrescription(patientID uuid.UUID, prescription *models.Prescription) error {
	// Verify patient exists
	var patient models.Patient
	if err := s.db.Conn.First(&patient, "id = ?", patientID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("patient not found")
		}
		return err
	}
	fmt.Println("this is the patient at service line 79", patient)
	// Update prescription directly by conditions without using Model()
	result := s.db.Conn.Model(&models.Prescription{}).
		Where("patient_id = ? AND doctor_id = ?", patientID, prescription.DoctorID).
		Updates(map[string]interface{}{
			"medication":   prescription.Medication,
			"dosage":       prescription.Dosage,
			"instructions": prescription.Instructions,
		})
	fmt.Println("this is the result at service line 88", result)
	if result.Error != nil {
		return result.Error
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return errors.New("prescription not found for this patient and doctor")
	}

	return nil
}

func (s *doctorService) GetAppointments(doctorID uuid.UUID) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := s.db.Conn.Preload("Patient").Where("doctor_id = ?", doctorID).Find(&appointments).Error
	if err != nil {
		return nil, err
	}
	return appointments, nil
}

func (s *doctorService) GetAppointmentsByDate(doctorID uuid.UUID, date time.Time) ([]models.Appointment, error) {
	var appointments []models.Appointment
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := s.db.Conn.Preload("Patient").Where("doctor_id = ? AND appointment_date >= ? AND appointment_date < ?",
		doctorID, startOfDay, endOfDay).Find(&appointments).Error
	if err != nil {
		return nil, err
	}
	return appointments, nil
}
