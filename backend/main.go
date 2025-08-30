package main

import (
	"github.com/rayzox/tickr-backend/database"
	"github.com/rayzox/tickr-backend/models"
	"github.com/rayzox/tickr-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	// connect DB & migrate
	database.ConnectDatabase()
	database.DB.AutoMigrate(&models.Task{})

	// register routes
	routes.RegisterTaskRoutes(r)
	routes.RegisterHabitRoutes(r)

	r.Run(":8080")
}
