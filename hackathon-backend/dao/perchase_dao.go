package dao

import (
	"database/sql"
	"hackathon-backend/model"
)

type PurchaseRepository interface {
	Insert(p *model.Purchase) error
}

type PurchaseDAO struct {
	DB *sql.DB
}

func NewPurchaseDAO(db *sql.DB) *PurchaseDAO {
	return &PurchaseDAO{DB: db}
}

func (d *PurchaseDAO) Insert(p *model.Purchase) error {
	_, err := d.DB.Exec(
		"INSERT INTO purchases (item_id, buyer_id) VALUES (?, ?)",
		p.ItemID, p.BuyerID,
	)
	return err
}
