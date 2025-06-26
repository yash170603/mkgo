package services

import (
	"errors"
	"hospital/internal/config"
	"hospital/internal/database"
	"hospital/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReceptionistServiceInterface defines the contract for receptionist operations
type ReceptionistServiceInterface interface {
	// Patient operations
	CreatePatient(patient *models.Patient) error
	GetPatients(page, limit int) ([]models.Patient, int64, error)
	GetPatient(patientID uuid.UUID) (*models.Patient, error)
	UpdatePatient(patientID uuid.UUID, patient *models.Patient) error
	DeletePatient(patientID uuid.UUID) error

	GetDoctor() (models.User, error)

	// Appointment operations
	CreateAppointment(appointment *models.Appointment) error
	GetAppointments(page, limit int) ([]models.Appointment, int64, error)
	GetAppointment(appointmentID uuid.UUID) (*models.Appointment, error)
	UpdateAppointment(patientID uuid.UUID, appointmentID uuid.UUID, parsedTime time.Time, status string, notes string) (*models.Appointment, error)
	DeleteAppointment(appointmentID uuid.UUID) error
	GetAllAppointments() ([]models.Appointment, error)
}

// ReceptionistService implements ReceptionistServiceInterface
type ReceptionistService struct {
	db  *database.DB
	cfg config.Config
}

// NewReceptionistService creates a new receptionist service instance
func NewReceptionistService(db *database.DB, cfg config.Config) ReceptionistServiceInterface {
	return &ReceptionistService{
		db:  db,
		cfg: cfg,
	}
}

// Patient Operations

func (s *ReceptionistService) CreatePatient(patient *models.Patient) error {
	if patient.Name == "" || patient.Email == "" || patient.Phone == "" || patient.Address == "" {
		return errors.New("name, email, phone, and address are required fields")
	}

	// Check if patient with email already exists, meaning error should be nil,
	var existingPatient models.Patient
	if err := s.db.Conn.Where("email = ?", patient.Email).First(&existingPatient).Error; err == nil {
		return errors.New("patient with this email already exists")
	}

	if err := s.db.Conn.Create(patient).Error; err != nil {
		return err
	}

	return nil
}

func (s *ReceptionistService) GetDoctor() (models.User, error) {
	var doctor models.User
	if err := s.db.Conn.Where("role = 'doctor'").First(&doctor).Error; err != nil {
		return models.User{}, err
	}
	return doctor, nil
}

func (s *ReceptionistService) GetPatients(page, limit int) ([]models.Patient, int64, error) {
	var patients []models.Patient
	var total int64

	// Count total records
	if err := s.db.Conn.Model(&models.Patient{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get patients with pagination and preload doctor information
	if err := s.db.Conn.Preload("User").
		Offset(offset).
		Limit(limit).
		Find(&patients).Error; err != nil {
		return nil, 0, err
	}

	return patients, total, nil
}

func (s *ReceptionistService) GetPatient(patientID uuid.UUID) (*models.Patient, error) {
	var patient models.Patient

	if err := s.db.Conn.Preload("User").
		Where("id = ?", patientID).
		First(&patient).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("patient not found")
		}
		return nil, err
	}

	return &patient, nil
}

func (s *ReceptionistService) UpdatePatient(patientID uuid.UUID, patient *models.Patient) error {
	// Check if patient exists
	var existingPatient models.Patient
	if err := s.db.Conn.Where("id = ?", patientID).First(&existingPatient).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("patient not found")
		}
		return err
	}

	// Check if email is being changed and if new email already exists
	if patient.Email != existingPatient.Email {
		var emailCheck models.Patient
		if err := s.db.Conn.Where("email = ? AND id != ?", patient.Email, patientID).First(&emailCheck).Error; err == nil {
			return errors.New("patient with this email already exists")
		}
	}

	// Update patient
	patient.ID = patientID // Ensure ID doesn't change

	patient.UserID = existingPatient.UserID
	if err := s.db.Conn.Save(patient).Error; err != nil {
		return err
	}

	return nil
}

func (s *ReceptionistService) DeletePatient(patientID uuid.UUID) error {
	result := s.db.Conn.Delete(&models.Patient{}, patientID)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("patient not found")
	}

	return nil
}

// Appointment Operations

func (s *ReceptionistService) CreateAppointment(appointment *models.Appointment) error {
	// Validate required fields
	if appointment.PatientID == uuid.Nil || appointment.DoctorID == uuid.Nil {
		return errors.New("patient_id and doctor_id are required")
	}

	// Check if patient exists
	var patient models.Patient
	if err := s.db.Conn.Where("id = ?", appointment.PatientID).First(&patient).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("patient not found")
		}
		return err
	}

	// Check if doctor exists and has correct role
	var doctor models.User
	if err := s.db.Conn.Where("id = ? AND role = ?", appointment.DoctorID, "doctor").First(&doctor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("doctor not found")
		}
		return err
	}

	// Check for conflicting appointments (same doctor, same time)
	var existingAppointment models.Appointment
	if err := s.db.Conn.Where("doctor_id = ? AND appointment_date = ? AND status != ?",
		appointment.DoctorID, appointment.AppointmentDate, "cancelled").
		First(&existingAppointment).Error; err == nil {
		return errors.New("doctor already has an appointment at this time")
	}

	if err := s.db.Conn.Create(appointment).Error; err != nil {
		return err
	}

	return nil
}

func (s *ReceptionistService) GetAppointments(page, limit int) ([]models.Appointment, int64, error) {
	var appointments []models.Appointment
	var total int64

	// Count total records
	if err := s.db.Conn.Model(&models.Appointment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get appointments with pagination and preload related data
	if err := s.db.Conn.Preload("Patient").
		Preload("Doctor").
		Offset(offset).
		Limit(limit).
		Order("appointment_date DESC").
		Find(&appointments).Error; err != nil {
		return nil, 0, err
	}

	return appointments, total, nil
}

func (s *ReceptionistService) GetAppointment(appointmentID uuid.UUID) (*models.Appointment, error) {
	var appointment models.Appointment

	if err := s.db.Conn.Preload("Patient").
		Preload("Doctor").
		Where("id = ?", appointmentID).
		First(&appointment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("appointment not found")
		}
		return nil, err
	}

	return &appointment, nil
}

// func (s *ReceptionistService) UpdateAppointment(appointmentID uuid.UUID, appointment *models.Appointment) error {
// 	// Check if appointment exists
// 	var existingAppointment models.Appointment
// 	if err := s.db.Conn.Where("id = ?", appointmentID).First(&existingAppointment).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return errors.New("appointment not found")
// 		}
// 		return err
// 	}

// 	// If doctor or date is being changed, check for conflicts
// 	if appointment.DoctorID != existingAppointment.DoctorID ||
// 		appointment.AppointmentDate != existingAppointment.AppointmentDate {

// 		var conflictingAppointment models.Appointment
// 		if err := s.db.Conn.Where("doctor_id = ? AND appointment_date = ? AND id != ? AND status != ?",
// 			appointment.DoctorID, appointment.AppointmentDate, appointmentID, "cancelled").
// 			First(&conflictingAppointment).Error; err == nil {
// 			return errors.New("doctor already has an appointment at this time")
// 		}
// 	}

// 	// Update appointment
// 	appointment.ID = appointmentID // Ensure ID doesn't change
// 	if err := s.db.Conn.Save(appointment).Error; err != nil {
// 		return err
// 	}

//		return nil
//	}
func (s *ReceptionistService) UpdateAppointment(patientID, appointmentID uuid.UUID, date time.Time, status, notes string) (*models.Appointment, error) {
	var existing models.Appointment
	if err := s.db.Conn.Where("id = ? AND patient_id = ?", appointmentID, patientID).First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("appointment not found")
		}
		return nil, err
	}

	// Check for doctor conflict if date changed
	if !existing.AppointmentDate.Equal(date) {
		var conflict models.Appointment
		if err := s.db.Conn.Where("doctor_id = ? AND appointment_date = ? AND id != ? AND status != ?",
			existing.DoctorID, date, appointmentID, "cancelled").First(&conflict).Error; err == nil {
			return nil, errors.New("doctor already has an appointment at this time")
		}
	}

	// Update fields
	existing.AppointmentDate = date
	existing.Status = status
	existing.Notes = notes

	if err := s.db.Conn.Save(&existing).Error; err != nil {
		return nil, err
	}

	return &existing, nil
}

func (s *ReceptionistService) GetAllAppointments() ([]models.Appointment, error) {
	var appointments []models.Appointment
	if err := s.db.Conn.Find(&appointments).Error; err != nil {
		return nil, err
	}
	return appointments, nil
}

func (s *ReceptionistService) DeleteAppointment(appointmentID uuid.UUID) error {
	result := s.db.Conn.Delete(&models.Appointment{}, appointmentID)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("appointment not found")
	}

	return nil
}
