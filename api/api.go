package api

import (
	"fmt"
	"hospital/api/routes"
	"hospital/internal/config"
	"hospital/internal/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Api struct {
	App *gin.Engine
}

func New(db *database.DB, cfg config.Config) *Api {
	r := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000", "http://localhost:8081",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}
	r.Use(cors.New(corsConfig))

	apiGroup := r.Group("/api")
	routes.RegisterAuth(apiGroup, cfg, db)
	routes.RegisterDoctor(apiGroup, cfg, db)
	routes.RegisterReceptionist(apiGroup, cfg, db)

	return &Api{App: r}
}

func (a *Api) Run(port int) error {
	return a.App.Run(fmt.Sprintf(":%d", port))
}
