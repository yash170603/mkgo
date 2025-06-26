package routes

import (
	"hospital/api/handlers"
	"hospital/internal/config"
	"hospital/internal/database"
	"hospital/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterAuth(apiGroup *gin.RouterGroup, cfg config.Config, db *database.DB) {
	authService := services.NewAuthService(db, cfg)
	authHandler := handlers.NewHandler(authService)

	authGroup := apiGroup.Group("/auth")
	authGroup.POST("/login", authHandler.LoginHandler)
	authGroup.POST("/logout", authHandler.LogoutHandler)
}
