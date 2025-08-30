package models

import (
	"time"

	"gorm.io/gorm"
)

type Habit struct {
	gorm.Model
	Name           string `json:"name"`
	Frequency      string `json:"frequency"`
	CompletedToday bool   `json:"completed_today"`
	Streak         int    `json:"streak"`

	// New integration fields
	Color           string     `json:"color"`
	LastCompletedAt *time.Time `json:"last_completed_at"`
	TargetTime      string     `json:"target_time"`
	CalendarEvents  []Event    `json:"calendar_events,omitempty" gorm:"foreignKey:HabitID"`
}
