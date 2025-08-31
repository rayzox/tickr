package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/rayzox/tickr-backend/database"
	"github.com/rayzox/tickr-backend/models"

	"github.com/gin-gonic/gin"
)

func RegisterTaskRoutes(r *gin.Engine) {
	r.GET("/tasks", GetTasks)
	r.POST("/tasks", CreateTask)
	r.PUT("/tasks/:id", UpdateTask)
	r.DELETE("/tasks/:id", DeleteTask)
}

func GetTasks(c *gin.Context) {
	var tasks []models.Task
	database.DB.Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

// tasks.go - Update CreateTask function
func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if due_date is provided as a string and parse it
	var requestData map[string]interface{}
	if err := c.ShouldBindJSON(&requestData); err == nil {
		if dueDateStr, ok := requestData["due_date"].(string); ok && dueDateStr != "" {
			// Try parsing in different formats
			if parsedDate, err := time.Parse(time.RFC3339, dueDateStr); err == nil {
				task.DueDate = parsedDate
			} else if parsedDate, err := time.Parse("2006-01-02T15:04:05", dueDateStr); err == nil {
				task.DueDate = parsedDate
			} else if parsedDate, err := time.Parse("2006-01-02", dueDateStr); err == nil {
				// Set default time to 9:00 AM if only date is provided
				task.DueDate = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 9, 0, 0, 0, parsedDate.Location())
			}
		}
	}

	if err := database.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func UpdateTask(c *gin.Context) {
	var task models.Task
	if err := database.DB.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.BindJSON(&task)
	database.DB.Save(&task)
	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam) // Use Atoi instead of ParseUint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var task models.Task
	if err := database.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if err := database.DB.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
