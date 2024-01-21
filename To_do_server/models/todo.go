package models

import "time"

type Todo struct {
	ID          string
	Title       string
	Description string
	IsCompleted bool
	CreateAt    time.Time
}
