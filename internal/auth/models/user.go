package models

import "time"

type User struct {
	DeletedAt time.Time `json:"deleted_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	Role      string    `json:"role"`
	Password  string    `json:"password"`
	Username  string    `json:"username"`
	ID        int64     `json:"id"`
}
