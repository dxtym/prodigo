package models

import "time"

type Product struct {
	CreatedAt  time.Time `json:"created_at"`
	Title      string    `json:"title"`
	UpdatedAt  time.Time `json:"updated_at"`
	Image      string    `json:"image"`
	DeletedAt  time.Time `json:"deleted_at"`
	Status     string    `json:"status"`
	ID         int64     `json:"id"`
	CategoryID int       `json:"category_id"`
	Price      int       `json:"price"`
	Quantity   int       `json:"quantity"`
}

type ProductFilterSearch struct {
	Search       string
	CategoryName string
	Status       string
	PriceMin     int
	PriceMax     int
}
