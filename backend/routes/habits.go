package routes

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rayzox/tickr-backend/database"
	"github.com/rayzox/tickr-backend/models"
)

func RegisterHabitRoutes(r *gin.Engine) {
	r.GET("/habits", GetHabits)
	r.POST("/habits", CreateHabit)
	r.PUT("/habits/:id", UpdateHabit)
	r.DELETE("/habits/:id", DeleteHabit)
}

func GetHabits(c *gin.Context) {
	var habits []models.Habit
	database.DB.Find(&habits)
	c.JSON(http.StatusOK, habits)
}

func CreateHabit(c *gin.Context) {
	var habit models.Habit
	if err := c.ShouldBindJSON(&habit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Create(&habit)
	c.JSON(http.StatusOK, habit)
}

func UpdateHabit(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid habit ID"})
		return
	}

	var habit models.Habit
	if err := database.DB.First(&habit, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
		return
	}

	// Bind JSON and check for errors
	if err := c.ShouldBindJSON(&habit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the updated habit
	if err := database.DB.Save(&habit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update habit"})
		return
	}

	c.JSON(http.StatusOK, habit)
}

func DeleteHabit(c *gin.Context) {
	// Convert id param to uint
	idParam := c.Param("id")
	log.Printf("Received ID param: %s", idParam) // Debug log

	idUint64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		log.Printf("Error parsing ID: %v", err) // Debug log
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid habit ID"})
		return
	}
	id := uint(idUint64)
	log.Printf("Parsed ID: %d", id) // Debug log

	var habit models.Habit
	// Use First to check if it exists
	result := database.DB.First(&habit, id)
	if result.Error != nil {
		log.Printf("Habit not found with ID %d: %v", id, result.Error) // Debug log
		c.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
		return
	}

	// Delete
	deleteResult := database.DB.Delete(&habit)
	if deleteResult.Error != nil {
		log.Printf("Error deleting habit: %v", deleteResult.Error) // Debug log
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete habit"})
		return
	}

	log.Printf("Successfully deleted habit with ID: %d", id) // Debug log
	c.JSON(http.StatusOK, gin.H{"message": "Habit deleted"})
}
