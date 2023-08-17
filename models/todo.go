package models

import (
	"time"
)

// Todo - todo model
type Todo struct {
	ID        *int64    `json:"id"`
	UserID    *int64    `json:"user_id"`
	Title     *string   `json:"title"`
	Completed *bool     `json:"completed"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
