package models

import "gorm.io/gorm"

type Habit struct {
	gorm.Model
	Name           string `json:"name"`
	Frequency      string `json:"frequency"`
	CompletedToday bool   `json:"completed_today"`
	Streak         int    `json:"streak"`
}
