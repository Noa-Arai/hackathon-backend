package dao

import (
	"database/sql"
	"hackathon-backend/model"
	"log"
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

	if err != nil {
		log.Printf("INSERT ERROR: %v", err) // ★これを追加
	}

	return err
}

func (d *ItemDAO) FindAll() ([]model.Item, error) {
	rows, err := d.DB.Query("SELECT id, user_id, title, description, price, image_url, created_at FROM items ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.Item{}
	for rows.Next() {
		var item model.Item
		if err := rows.Scan(
			&item.ID, &item.UserID, &item.Title, &item.Description,
			&item.Price, &item.ImageURL, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
