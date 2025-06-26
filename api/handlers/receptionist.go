package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"hospital/internal/models"
	"hospital/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReceptionistHandler struct {
	receptionistService services.ReceptionistServiceInterface
}

func NewReceptionistHandler(receptionistService services.ReceptionistServiceInterface) *ReceptionistHandler {
	return &ReceptionistHandler{
		receptionistService: receptionistService,
	}
}

// Patient Handlers

func (h *ReceptionistHandler) CreatePatient(c *gin.Context) {
	var patient models.Patient

	if err := c.ShouldBindJSON(&patient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	doctor, err := h.receptionistService.GetDoctor()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid or non-existent doctor",
		})
		return
	}
	fmt.Println("this is the doctor ", doctor)

	patient.UserID = doctor.ID
	fmt.Println(patient.UserID)

	if err := h.receptionistService.CreatePatient(&patient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error at creating patient": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Patient created successfully",
		"patient": patient,
	})
}

func (h *ReceptionistHandler) GetPatients(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	patients, total, err := h.receptionistService.GetPatients(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve patients",
		})
		return
	}

	// Calculate pagination info
	totalPages := (int(total) + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"patients": patients,
		"pagination": gin.H{
			"current_page": page,
			"total_pages":  totalPages,
			"total_count":  total,
			"per_page":     limit,
		},
	})
}

func (h *ReceptionistHandler) GetPatient(c *gin.Context) {
	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid patient ID format",
		})
		return
	}

	patient, err := h.receptionistService.GetPatient(patientID)
	if err != nil {
		if err.Error() == "patient not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Patient not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve patient",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"patient": patient,
	})
}

func (h *ReceptionistHandler) UpdatePatient(c *gin.Context) {
	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid patient ID format",
		})
		return
	}

	var patient models.Patient
	if err := c.ShouldBindJSON(&patient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := h.receptionistService.UpdatePatient(patientID, &patient); err != nil {
		if err.Error() == "patient not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Patient not found",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Patient updated successfully",
		"patient": patient,
	})
}

func (h *ReceptionistHandler) DeletePatient(c *gin.Context) {
	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid patient ID format",
		})
		return
	}

	if err := h.receptionistService.DeletePatient(patientID); err != nil {
		if err.Error() == "patient not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Patient not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete patient",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Patient deleted successfully",
	})
}

// Appointment Handlers

// func (h *ReceptionistHandler) CreateAppointment(c *gin.Context) {
// 	var appointment models.Appointment

// 	if err := c.ShouldBindJSON(&appointment); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error":   "Invalid request body",
// 			"details": err.Error(),
// 		})
// 		return
// 	}

// 	// Parse date and time in format "02/01/2006 15:04"
// 	appointment.AppointmentDate, err := time.Parse("02/01/2006 15:04", appointment.AppointmentDate)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error":   "Invalid appointment date format",
// 			"details": "Please use format: DD/MM/YYYY HH:MM",
// 		})
// 		return
// 	}

// 	// Validate that appointment is not in the past
// 	if appointment.AppointmentDate.Before(time.Now()) {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error":   "Invalid appointment date",
// 			"details": "Appointment date must be in the future",
// 		})
// 		return
// 	}

// 	patientId := c.Param("patient_id")
// 	patientIds, err := uuid.Parse(patientId)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid patient ID format",
// 		})
// 		return
// 	}

// 	appointment.PatientID = patientIds
// 	// fetch patient first from above api, then extract doctors detaisl from that to attach doctor in appointment struct
// 	patient, err := h.receptionistService.GetPatient(patientIds)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error":   "Failed to fetch patient details",
// 			"details": err.Error(),
// 		})
// 		return
// 	}

// 	appointment.DoctorID, err = uuid.Parse(patient.UserID.String())
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid doctor ID format",
// 		})
// 		return
// 	}

// 	if err := h.receptionistService.CreateAppointment(&appointment); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error":   "Failed to create appointment",
// 			"details": err.Error(),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{
// 		"message":     "Appointment created successfully",
// 		"appointment": appointment,
// 	})
// }

func (h *ReceptionistHandler) CreateAppointment(c *gin.Context) {
	type AppointmentInput struct {
		AppointmentDate string `json:"appointment_date"` // e.g. "03/07/2026 12:00"
		Status          string `json:"status"`
		Notes           string `json:"notes"`
	}

	var input AppointmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Parse date and time from string
	parsedTime, err := time.Parse("02/01/2006 15:04", input.AppointmentDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid appointment date format",
			"details": "Please use format: DD/MM/YYYY HH:MM",
		})
		return
	}

	// Validate that appointment is not in the past
	if parsedTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid appointment date",
			"details": "Appointment date must be in the future",
		})
		return
	}

	// Get patient ID from URL
	patientId := c.Param("patient_id")
	patientUUID, err := uuid.Parse(patientId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID format"})
		return
	}

	// Fetch patient details to get doctor
	patient, err := h.receptionistService.GetPatient(patientUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to fetch patient details",
			"details": err.Error(),
		})
		return
	}

	// Prepare appointment
	appointment := models.Appointment{
		PatientID:       patientUUID,
		DoctorID:        patient.UserID,
		AppointmentDate: parsedTime,
		Status:          input.Status,
		Notes:           input.Notes,
	}

	// Save appointment
	if err := h.receptionistService.CreateAppointment(&appointment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to create appointment",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Appointment created successfully",
		"appointment": appointment,
	})
}

func (h *ReceptionistHandler) GetAppointments(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	appointments, total, err := h.receptionistService.GetAppointments(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve appointments",
		})
		return
	}

	// Calculate pagination info
	totalPages := (int(total) + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"appointments": appointments,
		"pagination": gin.H{
			"current_page": page,
			"total_pages":  totalPages,
			"total_count":  total,
			"per_page":     limit,
		},
	})
}

func (h *ReceptionistHandler) GetAppointment(c *gin.Context) {
	appointmentIDStr := c.Param("appointment_id")
	appointmentID, err := uuid.Parse(appointmentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid appointment ID format",
		})
		return
	}

	appointment, err := h.receptionistService.GetAppointment(appointmentID)
	if err != nil {
		if err.Error() == "appointment not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Appointment not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve appointment",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"appointment": appointment,
	})
}

// func (h *ReceptionistHandler) UpdateAppointment(c *gin.Context) {
// 	appointmentIDStr := c.Param("appointment_id")
// 	appointmentID, err := uuid.Parse(appointmentIDStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid appointment ID format",
// 		})
// 		return
// 	}

// 	var appointment models.Appointment
// 	if err := c.ShouldBindJSON(&appointment); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error":   "Invalid request body",
// 			"details": err.Error(),
// 		})
// 		return
// 	}

// 	if err := h.receptionistService.UpdateAppointment(appointmentID, &appointment); err != nil {
// 		if err.Error() == "appointment not found" {
// 			c.JSON(http.StatusNotFound, gin.H{
// 				"error": "Appointment not found",
// 			})
// 			return
// 		}
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

//		c.JSON(http.StatusOK, gin.H{
//			"message":     "Appointment updated successfully",
//			"appointment": appointment,
//		})
//	}
func (h *ReceptionistHandler) UpdateAppointment(c *gin.Context) {
	appointmentIDStr := c.Param("appointment_id")
	appointmentID, err := uuid.Parse(appointmentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID format"})
		return
	}

	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID format"})
		return
	}

	type AppointmentUpdateInput struct {
		AppointmentDate string `json:"appointment_date"`
		Status          string `json:"status"`
		Notes           string `json:"notes"`
	}

	var input AppointmentUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Parse and validate appointment date
	parsedTime, err := time.Parse("02/01/2006 15:04", input.AppointmentDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid appointment date format",
			"details": "Please use format: DD/MM/YYYY HH:MM",
		})
		return
	}
	if parsedTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Appointment date must be in the future"})
		return
	}

	// Call service to update appointment safely
	updatedAppointment, err := h.receptionistService.UpdateAppointment(patientID, appointmentID, parsedTime, input.Status, input.Notes)
	if err != nil {
		if err.Error() == "appointment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Appointment updated successfully",
		"appointment": updatedAppointment,
	})
}

func (h *ReceptionistHandler) DeleteAppointment(c *gin.Context) {
	appointmentIDStr := c.Param("appointment_id")
	appointmentID, err := uuid.Parse(appointmentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid appointment ID format",
		})
		return
	}

	if err := h.receptionistService.DeleteAppointment(appointmentID); err != nil {
		if err.Error() == "appointment not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Appointment not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete appointment",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Appointment deleted successfully",
	})
}
