package model

import "time"

type Item struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int    `json:"price"`

	// ★ 追加（BLOB保存用）
	ImageData []byte `json:"-"`          // フロントに直接返さない
	ImageType string `json:"image_type"` // "image/png" など

	CreatedAt time.Time `json:"created_at"`
}
