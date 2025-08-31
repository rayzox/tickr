package models

import (
	"time"

	"gorm.io/gorm"
)

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
