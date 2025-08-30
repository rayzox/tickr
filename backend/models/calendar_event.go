package models

import (
	"time"

	"gorm.io/gorm"
)

type CalendarEvent struct {
	gorm.Model
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	EventDate   time.Time `json:"event_date" binding:"required"`
	Priority    string    `json:"priority"`
	AllDay      bool      `json:"all_day"`
	Duration    int       `json:"duration"` // Duration in minutes

	// Relationships
	EventType string `json:"event_type"` // 'task', 'habit', 'pomodoro', 'custom'
	TaskID    *uint  `json:"task_id"`
	HabitID   *uint  `json:"habit_id"`
	Task      *Task  `json:"task,omitempty" gorm:"foreignKey:TaskID"`
	Habit     *Habit `json:"habit,omitempty" gorm:"foreignKey:HabitID"`
}

type FocusSession struct {
	gorm.Model
	TaskID          uint       `json:"task_id"`
	Task            Task       `json:"task" gorm:"foreignKey:TaskID"`
	StartTime       time.Time  `json:"start_time"`
	EndTime         *time.Time `json:"end_time"`
	PlannedDuration int        `json:"planned_duration"`
	ActualDuration  int        `json:"actual_duration"`
	Notes           string     `json:"notes"`
	Completed       bool       `json:"completed"`
}
