package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	DueDate   time.Time `json:"due_date"` // Changed from string to time.Time
	Priority  string    `json:"priority"`

	// New integration fields
	Description        string            `json:"description"`
	DueDateParsed      *time.Time        `json:"due_date_parsed"`
	CalendarEventID    *uint             `json:"calendar_event_id"`
	CalendarEvent      *Event            `json:"calendar_event,omitempty" gorm:"foreignKey:CalendarEventID"`
	PomodoroSessions   []PomodoroSession `json:"pomodoro_sessions,omitempty" gorm:"foreignKey:TaskID"`
	EstimatedPomodoros int               `json:"estimated_pomodoros"`
	CompletedPomodoros int               `json:"completed_pomodoros"`
}
