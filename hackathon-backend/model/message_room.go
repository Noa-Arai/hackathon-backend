package model

import "time"

type MessageRoom struct {
	ItemID      int64     `json:"item_id"`
	ItemTitle   string    `json:"item_title"`
	PartnerID   string    `json:"partner_id"`
	LatestText  string    `json:"latest_text"`
	UpdatedAt   time.Time `json:"updated_at"`
	UnreadCount int       `json:"unread_count"`
}
