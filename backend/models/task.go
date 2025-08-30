package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	DueDate   string `json:"due_date"`
	Priority  string `json:"priority"`
}
