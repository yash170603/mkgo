package handlers

import (
	"fmt"
	"hospital/internal/models"
	"hospital/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DoctorHandler struct {
	doctorService services.DoctorService
}

func NewDoctorHandler(doctorService services.DoctorService) *DoctorHandler {
	return &DoctorHandler{
		doctorService: doctorService,
	}
}

func (h *DoctorHandler) GetPatients(c *gin.Context) {
	doctorID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id, err := uuid.Parse(doctorID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error parsing doctor ID"})
		return
	}
	patients, err := h.doctorService.GetPatients(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"patients": patients})
}

func (h *DoctorHandler) CreatePrescription(c *gin.Context) {
	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient ID"})
		return
	}

	doctorID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	doctorId, err := uuid.Parse(doctorID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error parsing doctor ID"})
		return
	}

	var prescription models.Prescription
	if err := c.ShouldBindJSON(&prescription); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prescription.DoctorID = doctorId

	if err := h.doctorService.CreatePrescription(patientID, &prescription); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "prescription created successfully", "prescription": prescription})
}

func (h *DoctorHandler) UpdatePrescription(c *gin.Context) {
	patientIDStr := c.Param("patient_id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient ID"})
		return
	}

	doctorID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	doctorId, err := uuid.Parse(doctorID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error parsing doctor ID"})
		return
	}

	var prescription models.Prescription
	if err := c.ShouldBindJSON(&prescription); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prescription.DoctorID = doctorId
	fmt.Println("this is line 106 prescription", prescription)
	if err := h.doctorService.UpdatePrescription(patientID, &prescription); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "prescription updated successfully"})
}

func (h *DoctorHandler) GetAppointments(c *gin.Context) {
	doctorID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	doctorId, err := uuid.Parse(doctorID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error parsing doctor ID"})
		return
	}

	appointments, err := h.doctorService.GetAppointments(doctorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"appointments": appointments})
}

func (h *DoctorHandler) GetAppointmentsByDate(c *gin.Context) {
	doctorID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date parameter is required (format: YYYY-MM-DD)"})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format (expected: YYYY-MM-DD)"})
		return
	}

	doctorId, err := uuid.Parse(doctorID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error parsing doctor ID"})
		return
	}

	appointments, err := h.doctorService.GetAppointmentsByDate(doctorId, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"appointments": appointments})
}
