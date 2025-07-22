package models

import "time"

type Todo struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title" gorm:"not null"`
	Completed   bool       `json:"completed"`
	Category    string     `json:"category"`
	Priority    string     `json:"priority"`
	CompletedAt *time.Time `json:"completedAt"`
	DueDate     *time.Time `json:"dueDate"`
}
