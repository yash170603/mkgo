package database

import (
	"fmt"
	"hospital/internal/config"
	"hospital/internal/models"
	"log"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

type DB struct {
	Conn *gorm.DB
}

func Connect(cfg config.DatabaseConfig) *DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		cfg.Host, cfg.User, cfg.Password, cfg.DB, cfg.Port)
	var err error
	db := &DB{}
	db.Conn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	db.Conn.AutoMigrate(&models.User{}, &models.Patient{}, &models.Appointment{}, &models.Prescription{}) //  User and Patient models are migrated
	log.Println("Connected to database successfully")
	return db
}
