package dao

import (
	"database/sql"
	"hackathon-backend/model"
	"log"
)

type ItemDAO struct {
	DB *sql.DB
}

func NewItemDAO(db *sql.DB) *ItemDAO {
	return &ItemDAO{DB: db}
}

func (d *ItemDAO) Insert(item *model.Item) error {

	_, err := d.DB.Exec(`
		INSERT INTO items 
		(user_id, title, description, price, image_data, image_type)
		VALUES (?, ?, ?, ?, ?, ?)
	`,
		item.UserID,
		item.Title,
		item.Description,
		item.Price,
		item.ImageData,
		item.ImageType,
	)

	if err != nil {
		log.Printf("INSERT ERROR: %v", err)
	}

	return err
}

func (d *ItemDAO) FindAll() ([]model.Item, error) {

	rows, err := d.DB.Query(`
		SELECT id, user_id, title, description, price, image_type, created_at
		FROM items
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.Item

	for rows.Next() {
		var i model.Item

		// image_data は返さない（一覧が重くなるため）
		err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.Description,
			&i.Price,
			&i.ImageType,
			&i.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// 画像取得エンドポイント用
func (d *ItemDAO) FindImageByID(id string) ([]byte, string, error) {
	row := d.DB.QueryRow(`
		SELECT image_data, image_type
		FROM items
		WHERE id = ?
	`, id)

	var data []byte
	var mime string
	err := row.Scan(&data, &mime)
	return data, mime, err
}
