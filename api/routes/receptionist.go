package routes

import (
	"hospital/api/handlers"
	"hospital/api/middleware"
	"hospital/internal/config"
	"hospital/internal/database"
	"hospital/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterReceptionist(apiGroup *gin.RouterGroup, cfg config.Config, db *database.DB) {
	// Create service interface - this returns the interface, not concrete type
	receptionistService := services.NewReceptionistService(db, cfg)
	receptionistHandler := handlers.NewReceptionistHandler(receptionistService)

	authGroup := apiGroup.Group("/receptionist")
	authGroup.Use(middleware.AuthMiddleware(cfg))
	authGroup.Use(middleware.RoleMiddleware("receptionist"))

	// Patient routes
	authGroup.POST("/patients", receptionistHandler.CreatePatient)
	authGroup.GET("/patients", receptionistHandler.GetPatients)
	authGroup.GET("/patients/:patient_id", receptionistHandler.GetPatient)
	authGroup.PUT("/patients/:patient_id", receptionistHandler.UpdatePatient)
	authGroup.DELETE("/patients/:patient_id", receptionistHandler.DeletePatient)

	// Appointment routes -
	authGroup.POST("/patients/:patient_id/appointments", receptionistHandler.CreateAppointment)                //done
	authGroup.GET("/patients/:patient_id/appointments", receptionistHandler.GetAppointments)                   //done
	authGroup.GET("/patients/:patient_id/appointments/:appointment_id", receptionistHandler.GetAppointment)    //done but yk u goota manually select the appointment_id from client , not id but yk how itll be handled
	authGroup.PUT("/patients/:patient_id/appointments/:appointment_id", receptionistHandler.UpdateAppointment) //done
	authGroup.DELETE("/patients/:patient_id/appointments/:appointment_id", receptionistHandler.DeleteAppointment)
	// New endpoint to fetch all appointments
	authGroup.GET("/appointments", receptionistHandler.GetAllAppointments)
}
