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
        (user_id, title, description, price, category, is_lucky_bag,
        image1_data, image1_type,
        image2_data, image2_type,
        image3_data, image3_type)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `,
		item.UserID,
		item.Title, item.Description, item.Price, item.Category, item.IsLuckyBag,
		item.Image1Data, item.Image1Type,
		item.Image2Data, item.Image2Type,
		item.Image3Data, item.Image3Type,
	)

	if err != nil {
		log.Printf("INSERT ERROR: %v", err)
	}

	return err
}

// dao/item_dao.go の FindAll 関数をこれに書き換えてください

func (d *ItemDAO) FindAll() ([]model.Item, error) {
	// 🔥 修正: 画像データの中身ではなく「存在するかどうか(0か1か)」を取得する
	rows, err := d.DB.Query(`
		SELECT 
			id, user_id, title, description, price, category, is_lucky_bag, created_at,
			(image1_data IS NOT NULL AND LENGTH(image1_data) > 0) as has_img1,
			(image2_data IS NOT NULL AND LENGTH(image2_data) > 0) as has_img2,
			(image3_data IS NOT NULL AND LENGTH(image3_data) > 0) as has_img3
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
		var hasImg1, hasImg2, hasImg3 bool // 画像があるかどうかのフラグ

		err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.Description,
			&i.Price,
			&i.Category,
			&i.IsLuckyBag,
			&i.CreatedAt,
			&hasImg1, // 🔥 追加
			&hasImg2, // 🔥 追加
			&hasImg3, // 🔥 追加
		)
		if err != nil {
			return nil, err
		}

		// 🔥 修正: 画像がある場合だけURLを設定する
		if hasImg1 {
			i.Image1URL = "/items/image1?id=" + i.ID
		}
		if hasImg2 {
			i.Image2URL = "/items/image2?id=" + i.ID
		}
		if hasImg3 {
			i.Image3URL = "/items/image3?id=" + i.ID
		}

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
	// 基本の更新クエリ
	query := `
		UPDATE items
		SET 
			title = ?,
			description = ?,
			price = ?,
			category = ?,
			is_lucky_bag = ?
	`
	args := []interface{}{
		item.Title,
		item.Description,
		item.Price,
		item.Category,
		item.IsLuckyBag,
	}

	// 🔥 画像1がある場合のみ更新対象にする
	if len(item.Image1Data) > 0 {
		query += `, image1_data = ?, image1_type = ?`
		args = append(args, item.Image1Data, item.Image1Type)
	}
	// 🔥 画像2がある場合のみ更新対象にする
	if len(item.Image2Data) > 0 {
		query += `, image2_data = ?, image2_type = ?`
		args = append(args, item.Image2Data, item.Image2Type)
	}
	// 🔥 画像3がある場合のみ更新対象にする
	if len(item.Image3Data) > 0 {
		query += `, image3_data = ?, image3_type = ?`
		args = append(args, item.Image3Data, item.Image3Type)
	}

	// 最後にWHERE句をつける
	query += ` WHERE id = ? AND user_id = ?`
	args = append(args, item.ID, item.UserID)

	_, err := d.DB.Exec(query, args...)
	return err
}
