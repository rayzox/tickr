package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"` // Keep your existing ISO format

	// New integration fields
	EventDate time.Time `json:"event_date"` // Parsed date for queries
	Priority  string    `json:"priority"`   // 'low', 'medium', 'high'
	AllDay    bool      `json:"all_day"`
	Duration  int       `json:"duration"` // Duration in minutes

	// Relationships
	EventType string `json:"event_type"` // 'task', 'habit', 'pomodoro', 'custom'
	TaskID    *uint  `json:"task_id"`
	HabitID   *uint  `json:"habit_id"`
	Task      *Task  `json:"task,omitempty" gorm:"foreignKey:TaskID"`
	Habit     *Habit `json:"habit,omitempty" gorm:"foreignKey:HabitID"`
}
