package models

import "gorm.io/gorm"

type PomodoroSession struct {
	gorm.Model
	TaskTitle string `json:"task_title"`
	Duration  int    `json:"duration"` // in minutes
	Completed bool   `json:"completed"`
}
