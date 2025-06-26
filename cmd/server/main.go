package main

import (
	"hospital/api"
	"hospital/internal/config"
	"hospital/internal/database"
	"hospital/internal/seeder"
)

func main() {

	cfg := config.New()
	db := database.Connect(cfg.DatabaseConfig)
	seeder.SeedUsers(db.Conn)

	api := api.New(db, cfg)
	api.Run(cfg.APIConfig.Port)

}
