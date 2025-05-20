package models

import "time"

type Category struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
	Name      string    `json:"name"`
	ID        int64     `json:"id"`
}
