package dao

import (
	"database/sql"
	"hackathon-backend/model"
)

type ItemRepository interface {
	Insert(item *model.Item) error
}

type ItemDAO struct {
	DB *sql.DB
}

func NewItemDAO(db *sql.DB) *ItemDAO {
	return &ItemDAO{DB: db}
}

func (d *ItemDAO) Insert(item *model.Item) error {
	_, err := d.DB.Exec(
		"INSERT INTO items (user_id, title, description, price, image_url) VALUES (?, ?, ?, ?, ?)",
		item.UserID, item.Title, item.Description, item.Price, item.ImageURL,
	)
	return err
}
