package routes

import (
	"hospital/api/handlers"
	"hospital/api/middleware"
	"hospital/internal/config"
	"hospital/internal/database"
	"hospital/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterDoctor(apiGroup *gin.RouterGroup, cfg config.Config, db *database.DB) {
	doctorService := services.NewDoctorService(db, cfg)
	doctorHandler := handlers.NewDoctorHandler(doctorService)

	authGroup := apiGroup.Group("/doctor")
	authGroup.Use(middleware.AuthMiddleware(cfg))
	authGroup.Use(middleware.RoleMiddleware("doctor"))

	// Patient routes
	authGroup.GET("/patients", doctorHandler.GetPatients) //done

	// Prescription routes
	authGroup.POST("/prescriptions/:patient_id", doctorHandler.CreatePrescription) //done
	authGroup.PUT("/prescriptions/:patient_id", doctorHandler.UpdatePrescription)  //done

	// Appointment routes
	authGroup.GET("/appointments", doctorHandler.GetAppointments)               //done
	authGroup.GET("/appointments/by-date", doctorHandler.GetAppointmentsByDate) //done
}
