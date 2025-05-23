package models

import "time"

type Category struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
	Name      string    `json:"name"`
	ID        int64     `json:"id"`
}

type CategoryStats struct {
	CategoryName  string `json:"category_name"`
	CategoryID    int64  `json:"category_id"`
	ProductCount  int    `json:"product_count"`
	TotalQuantity int    `json:"total_quantity"`
	TotalValue    int    `json:"total_value"`
}
