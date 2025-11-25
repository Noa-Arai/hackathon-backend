package model

import "time"

type Item struct {
	ID          int64     `json:"id"`
	UserID      string    `json:"user_id"` // ← 修正
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
}
