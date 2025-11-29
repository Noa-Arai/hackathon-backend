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
		(user_id, title, description, price,
		 image1_data, image1_type,
		 image2_data, image2_type,
		 image3_data, image3_type)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		item.UserID,
		item.Title, item.Description, item.Price,
		item.Image1Data, item.Image1Type,
		item.Image2Data, item.Image2Type,
		item.Image3Data, item.Image3Type,
	)

	if err != nil {
		log.Printf("INSERT ERROR: %v", err)
	}

	return err
}

func (d *ItemDAO) FindAll() ([]model.Item, error) {

	rows, err := d.DB.Query(`
		SELECT id, user_id, title, description, price, created_at
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

		err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.Description,
			&i.Price,
			&i.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// --- 一覧用URLを組み立てる ---
		i.Image1URL = "/items/image1?id=" + i.ID
		i.Image2URL = "/items/image2?id=" + i.ID
		i.Image3URL = "/items/image3?id=" + i.ID

		items = append(items, i)
	}

	return items, nil
}

func (d *ItemDAO) FindImage(id string, index int) ([]byte, string, error) {
	var query string

	switch index {
	case 1:
		query = `SELECT image1_data, image1_type FROM items WHERE id = ?`
	case 2:
		query = `SELECT image2_data, image2_type FROM items WHERE id = ?`
	case 3:
		query = `SELECT image3_data, image3_type FROM items WHERE id = ?`
	}

	row := d.DB.QueryRow(query, id)

	var data []byte
	var mime string
	err := row.Scan(&data, &mime)

	return data, mime, err
}

func (d *ItemDAO) Update(item *model.Item) error {
	_, err := d.DB.Exec(`
UPDATE items
SET 
  title = ?,
  description = ?,
  price = ?,
  image1_data = ?, image1_type = ?,
  image2_data = ?, image2_type = ?,
  image3_data = ?, image3_type = ?
WHERE id = ? AND user_id = ?
`,
		item.Title,
		item.Description,
		item.Price,
		item.Image1Data, item.Image1Type,
		item.Image2Data, item.Image2Type,
		item.Image3Data, item.Image3Type,
		item.ID, item.UserID,
	)

	return err
}
