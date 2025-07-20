package models

import "time"

type Todo struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Completed   bool       `json:"completed"`
	Category    string     `json:"category"`
	Priority    string     `json:"priority"`
	CompletedAt *time.Time `json:"completedAt"`
	DueDate     *time.Time `json:"dueDate"`
}
