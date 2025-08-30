// models/pomodoro.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type PomodoroSession struct {
	gorm.Model
	Phase       string    `json:"phase"`    // 'work', 'short', 'long'
	Duration    int       `json:"duration"` // duration in minutes
	CompletedAt time.Time `json:"completed_at"`
}

type PomodoroStats struct {
	CompletedSessions int `json:"completedSessions"`
	TotalMinutes      int `json:"totalMinutes"`
	TodaySessions     int `json:"todaySessions"`
}
