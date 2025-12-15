package model

import "time"

type Item struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Category    string `json:"category"`

	// --- 保存用（BLOB） ---
	Image1Data []byte `json:"-"`
	Image1Type string `json:"-"`
	Image2Data []byte `json:"-"`
	Image2Type string `json:"-"`
	Image3Data []byte `json:"-"`
	Image3Type string `json:"-"`

	// --- フロント用URL ---
	Image1URL string `json:"image1_url"`
	Image2URL string `json:"image2_url"`
	Image3URL string `json:"image3_url"`

	CreatedAt time.Time `json:"created_at"`
}
