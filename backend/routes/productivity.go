package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rayzox/tickr-backend/database"
	"github.com/rayzox/tickr-backend/models"
)

func RegisterProductivityRoutes(r *gin.Engine) {
	productivity := r.Group("/productivity")
	{
		productivity.GET("/dashboard", GetDashboardData)
		productivity.POST("/focus/start", StartFocusSession)
		productivity.PUT("/focus/:id/complete", CompleteFocusSession)
		productivity.GET("/focus/active", GetActiveFocusSession)
		productivity.POST("/schedule/task", ScheduleTask)
		productivity.POST("/schedule/habit", ScheduleHabit)
		productivity.GET("/analytics/weekly", GetWeeklyAnalytics)
	}
}

func GetDashboardData(c *gin.Context) {
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.AddDate(0, 0, 1)

	var todayTasks []models.Task
	var todayHabits []models.Habit
	var todayEvents []models.Event
	var todayPomodoros []models.PomodoroSession

	// Get today's data
	database.DB.Where("due_date = ?", today.Format("2006-01-02")).Find(&todayTasks)
	database.DB.Find(&todayHabits)
	database.DB.Where("event_date >= ? AND event_date < ?", today, tomorrow).Find(&todayEvents)
	database.DB.Where("completed_at >= ? AND completed_at < ?", today, tomorrow).Find(&todayPomodoros)

	// Get active focus session
	var activeFocus models.FocusSession
	database.DB.Where("end_time IS NULL").Preload("Task").First(&activeFocus)

	response := gin.H{
		"today_tasks":     todayTasks,
		"today_habits":    todayHabits,
		"today_events":    todayEvents,
		"today_pomodoros": todayPomodoros,
		"active_focus":    activeFocus,
		"date":            today,
	}

	c.JSON(http.StatusOK, response)
}

func StartFocusSession(c *gin.Context) {
	var request struct {
		TaskID          uint `json:"task_id" binding:"required"`
		PlannedDuration int  `json:"planned_duration"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if task exists
	var task models.Task
	if err := database.DB.First(&task, request.TaskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Check for active sessions
	var activeSession models.FocusSession
	if err := database.DB.Where("end_time IS NULL").First(&activeSession).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Another focus session is already active"})
		return
	}

	session := models.FocusSession{
		TaskID:          request.TaskID,
		StartTime:       time.Now(),
		PlannedDuration: request.PlannedDuration,
	}

	if err := database.DB.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start focus session"})
		return
	}

	database.DB.Preload("Task").First(&session, session.ID)
	c.JSON(http.StatusOK, session)
}

func CompleteFocusSession(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	var request struct {
		Notes      string `json:"notes"`
		Completed  bool   `json:"completed"`
		PomodoroID *uint  `json:"pomodoro_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var session models.FocusSession
	if err := database.DB.First(&session, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Focus session not found"})
		return
	}

	now := time.Now()
	session.EndTime = &now
	session.ActualDuration = int(now.Sub(session.StartTime).Minutes())
	session.Notes = request.Notes
	session.Completed = request.Completed

	if err := database.DB.Save(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete session"})
		return
	}

	// Update task pomodoro count if productive
	if request.Completed && request.PomodoroID != nil {
		var task models.Task
		if err := database.DB.First(&task, session.TaskID).Error; err == nil {
			task.CompletedPomodoros++
			database.DB.Save(&task)

			database.DB.Model(&models.PomodoroSession{}).Where("id = ?", *request.PomodoroID).
				Update("task_id", task.ID)
		}
	}

	c.JSON(http.StatusOK, session)
}

func GetActiveFocusSession(c *gin.Context) {
	var session models.FocusSession
	if err := database.DB.Where("end_time IS NULL").Preload("Task").First(&session).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"active_session": nil})
		return
	}
	c.JSON(http.StatusOK, session)
}

func ScheduleTask(c *gin.Context) {
	var request struct {
		TaskID    uint      `json:"task_id" binding:"required"`
		EventDate time.Time `json:"event_date" binding:"required"`
		Duration  int       `json:"duration"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var task models.Task
	if err := database.DB.First(&task, request.TaskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	event := models.Event{
		Title:       "Work on: " + task.Title,
		Description: task.Description,
		EventDate:   request.EventDate,
		Duration:    request.Duration,
		Priority:    task.Priority,
		EventType:   "task",
		TaskID:      &task.ID,
	}

	if err := database.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to schedule task"})
		return
	}

	task.CalendarEventID = &event.ID
	database.DB.Save(&task)

	c.JSON(http.StatusOK, event)
}

func ScheduleHabit(c *gin.Context) {
	var request struct {
		HabitID   uint      `json:"habit_id" binding:"required"`
		StartDate time.Time `json:"start_date" binding:"required"`
		Days      int       `json:"days"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var habit models.Habit
	if err := database.DB.First(&habit, request.HabitID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
		return
	}

	var events []models.Event
	currentDate := request.StartDate

	for i := 0; i < request.Days; i++ {
		shouldSchedule := false

		switch habit.Frequency {
		case "daily":
			shouldSchedule = true
		case "weekly":
			shouldSchedule = currentDate.Weekday() == request.StartDate.Weekday()
		case "monthly":
			shouldSchedule = currentDate.Day() == request.StartDate.Day()
		}

		if shouldSchedule {
			targetTime := "09:00"
			if habit.TargetTime != "" {
				targetTime = habit.TargetTime
			}

			eventDateTime, _ := time.Parse("2006-01-02 15:04",
				currentDate.Format("2006-01-02")+" "+targetTime)

			event := models.Event{
				Title:       "🎯 " + habit.Name,
				Description: "Habit reminder",
				EventDate:   eventDateTime,
				Duration:    30,
				Priority:    "medium",
				EventType:   "habit",
				HabitID:     &habit.ID,
			}
			events = append(events, event)
		}

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	if len(events) > 0 {
		if err := database.DB.Create(&events).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to schedule habits"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"scheduled_events": len(events), "events": events})
}

func GetWeeklyAnalytics(c *gin.Context) {
	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())

	var weeklyStats []gin.H

	for i := 0; i < 7; i++ {
		day := startOfWeek.AddDate(0, 0, i)
		nextDay := day.AddDate(0, 0, 1)

		var completedTasks, totalTasks int64
		var completedHabits, totalHabits int64
		var pomodoroSessions int64
		var totalPomodoroMinutes int64

		database.DB.Model(&models.Task{}).Where("completed = ? AND updated_at >= ? AND updated_at < ?", true, day, nextDay).Count(&completedTasks)
		database.DB.Model(&models.Task{}).Where("due_date = ?", day.Format("2006-01-02")).Count(&totalTasks)

		database.DB.Model(&models.Habit{}).Where("completed_today = ? AND updated_at >= ? AND updated_at < ?", true, day, nextDay).Count(&completedHabits)
		database.DB.Model(&models.Habit{}).Count(&totalHabits)

		database.DB.Model(&models.PomodoroSession{}).Where("completed_at >= ? AND completed_at < ?", day, nextDay).Count(&pomodoroSessions)
		database.DB.Model(&models.PomodoroSession{}).Where("completed_at >= ? AND completed_at < ? AND phase = ?", day, nextDay, "work").
			Select("COALESCE(SUM(duration), 0)").Scan(&totalPomodoroMinutes)

		weeklyStats = append(weeklyStats, gin.H{
			"date":               day.Format("2006-01-02"),
			"day_name":           day.Format("Monday"),
			"completed_tasks":    completedTasks,
			"total_tasks":        totalTasks,
			"completed_habits":   completedHabits,
			"total_habits":       totalHabits,
			"pomodoro_sessions":  pomodoroSessions,
			"productive_minutes": totalPomodoroMinutes,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"week_start":  startOfWeek.Format("2006-01-02"),
		"daily_stats": weeklyStats,
	})
}
