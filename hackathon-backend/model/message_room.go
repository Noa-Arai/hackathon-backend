package model

import "time"

type MessageRoom struct {
	ItemID        int64     `json:"item_id"`
	PartnerID     string    `json:"partner_id"`
	PartnerName   string    `json:"partner_name"`
	PartnerAvatar string    `json:"partner_avatar"`
	LastMessage   string    `json:"last_message"`
	LastTime      time.Time `json:"last_time"`
	UnreadCount   int       `json:"unread_count"`
}
