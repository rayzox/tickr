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
	database.DB.AutoMigrate(&models.Habit{})
	database.DB.AutoMigrate(&models.Event{}) // Your existing event model - now enhanced
	database.DB.AutoMigrate(&models.PomodoroSession{})
	database.DB.AutoMigrate(&models.FocusSession{}) // New focus model

	// register routes
	routes.RegisterTaskRoutes(r)
	routes.RegisterHabitRoutes(r)
	routes.RegisterPomodoroRoutes(r)
	routes.RegisterCalendarRoutes(r)     // New calendar routes
	routes.RegisterProductivityRoutes(r) // New productivity routes

	r.Run(":8080")
}
