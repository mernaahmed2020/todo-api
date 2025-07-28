package models

import "time"

type Tag struct {
	Tag   string  `gorm:"primaryKey;size:50" json:"tag"`
	Todos []*Todo `gorm:"many2many:todo_tags" json:"-"`
}

type Todo struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title" gorm:"not null"`
	Completed   bool       `json:"completed"`
	Category    string     `json:"category"`
	Priority    string     `json:"priority"`
	CompletedAt *time.Time `json:"completedAt"`
	DueDate     *time.Time `json:"dueDate"`
	Tags        []*Tag     `gorm:"many2many:todo_tags;" json:"tags"`
}
