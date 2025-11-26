package model

import "time"

type Message struct {
	ID         int64     `json:"id"`
	FromUserID string    `json:"from_user_id"`
	ToUserID   string    `json:"to_user_id"`
	ItemID     int64     `json:"item_id"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"created_at"`
}
