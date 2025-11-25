package model

import "time"

type Purchase struct {
	ID          int64     `json:"id"`
	ItemID      int64     `json:"item_id"`
	BuyerID     string    `json:"buyer_id"`
	PurchasedAt time.Time `json:"purchased_at"`
}
