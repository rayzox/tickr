package models

import (
	"time"

	"gorm.io/gorm"
)

type PomodoroSession struct {
	gorm.Model
	Phase       string    `json:"phase"`
	Duration    int       `json:"duration"`
	CompletedAt time.Time `json:"completed_at"`

	// New integration fields
	TaskID     *uint  `json:"task_id"`
	Task       *Task  `json:"task,omitempty" gorm:"foreignKey:TaskID"`
	Notes      string `json:"notes"`
	Productive bool   `json:"productive"`
}

type PomodoroStats struct {
	CompletedSessions int `json:"completedSessions"`
	TotalMinutes      int `json:"totalMinutes"`
	TodaySessions     int `json:"todaySessions"`
}
