// routes/pomodoro_routes.go
package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rayzox/tickr-backend/database"
	"github.com/rayzox/tickr-backend/models"
)

func RegisterPomodoroRoutes(r *gin.Engine) {
	pomodoro := r.Group("/pomodoro")
	{
		pomodoro.POST("/sessions", CreatePomodoroSession)
		pomodoro.GET("/stats", GetPomodoroStats)
		pomodoro.GET("/sessions", GetPomodoroSessions)
		pomodoro.DELETE("/sessions", ClearPomodoroSessions)
	}
}

func CreatePomodoroSession(c *gin.Context) {
	var session models.PomodoroSession
	if err := c.ShouldBindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set completed_at to now if not provided
	if session.CompletedAt.IsZero() {
		session.CompletedAt = time.Now()
	}

	if err := database.DB.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, session)
}

func GetPomodoroStats(c *gin.Context) {
	var totalSessions int64
	var totalMinutes int64
	var todaySessions int64

	// Get total completed sessions
	database.DB.Model(&models.PomodoroSession{}).Count(&totalSessions)

	// Get total minutes
	database.DB.Model(&models.PomodoroSession{}).Select("COALESCE(SUM(duration), 0)").Scan(&totalMinutes)

	// Get today's sessions (from midnight today)
	today := time.Now().Truncate(24 * time.Hour)
	database.DB.Model(&models.PomodoroSession{}).
		Where("completed_at >= ?", today).
		Count(&todaySessions)

	stats := models.PomodoroStats{
		CompletedSessions: int(totalSessions),
		TotalMinutes:      int(totalMinutes),
		TodaySessions:     int(todaySessions),
	}

	c.JSON(http.StatusOK, stats)
}

func GetPomodoroSessions(c *gin.Context) {
	var sessions []models.PomodoroSession

	// Get recent sessions (last 50)
	database.DB.Order("completed_at DESC").Limit(50).Find(&sessions)

	c.JSON(http.StatusOK, sessions)
}

func ClearPomodoroSessions(c *gin.Context) {
	if err := database.DB.Where("1 = 1").Delete(&models.PomodoroSession{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All sessions cleared"})
}
