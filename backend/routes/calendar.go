package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rayzox/tickr-backend/database"
	"github.com/rayzox/tickr-backend/models"
)

func RegisterCalendarRoutes(r *gin.Engine) {
	calendar := r.Group("/calendar")
	{
		calendar.GET("/events", GetCalendarEvents)
		calendar.POST("/events", CreateCalendarEvent)
		calendar.PUT("/events/:id", UpdateCalendarEvent)
		calendar.DELETE("/events/:id", DeleteCalendarEvent)
		calendar.GET("/events/range", GetEventsInRange)
	}
}

func GetCalendarEvents(c *gin.Context) {
	var events []models.Event

	// Use Find with proper preloading
	if err := database.DB.
		Preload("Task").
		Preload("Habit").
		Order("event_date ASC").
		Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

func CreateCalendarEvent(c *gin.Context) {
	var event models.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse the Date string to EventDate if provided
	if event.Date != "" {
		// Try parsing with time format including timezone
		if parsedDate, err := time.Parse(time.RFC3339, event.Date); err == nil {
			event.EventDate = parsedDate
		} else if parsedDate, err := time.Parse("2006-01-02T15:04:05", event.Date); err == nil {
			event.EventDate = parsedDate
		} else if parsedDate, err := time.Parse("2006-01-02", event.Date); err == nil {
			event.EventDate = parsedDate
		}
	}

	if err := database.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	// Preload related data before returning
	database.DB.Preload("Task").Preload("Habit").First(&event, event.ID)

	c.JSON(http.StatusOK, event)
}

func UpdateCalendarEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	var event models.Event
	if err := database.DB.First(&event, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}

	c.JSON(http.StatusOK, event)
}

func DeleteCalendarEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	var event models.Event
	if err := database.DB.First(&event, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if err := database.DB.Delete(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}

func GetEventsInRange(c *gin.Context) {
	startDate := c.Query("start")
	endDate := c.Query("end")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start and end dates are required"})
		return
	}

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
		return
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
		return
	}

	var events []models.Event
	if err := database.DB.Where("event_date BETWEEN ? AND ?", start, end.AddDate(0, 0, 1)).
		Preload("Task").Preload("Habit").Order("event_date ASC").Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
		return
	}

	c.JSON(http.StatusOK, events)
}
