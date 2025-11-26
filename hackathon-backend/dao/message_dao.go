package dao

import (
	"database/sql"
	"hackathon-backend/model"
)

type MessageRepository interface {
	Insert(msg *model.Message) error
	FindByItemID(itemID int64) ([]model.Message, error)
}

type MessageDAO struct {
	DB *sql.DB
}

func NewMessageDAO(db *sql.DB) *MessageDAO {
	return &MessageDAO{DB: db}
}

func (d *MessageDAO) Insert(msg *model.Message) error {
	_, err := d.DB.Exec(
		"INSERT INTO messages (from_user_id, to_user_id, item_id, text) VALUES (?, ?, ?, ?)",
		msg.FromUserID, msg.ToUserID, msg.ItemID, msg.Text,
	)
	return err
}

func (d *MessageDAO) FindByItemID(itemID int64) ([]model.Message, error) {
	rows, err := d.DB.Query(
		"SELECT id, from_user_id, to_user_id, item_id, text, created_at FROM messages WHERE item_id = ? ORDER BY created_at ASC",
		itemID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []model.Message

	for rows.Next() {
		var msg model.Message
		if err := rows.Scan(
			&msg.ID,
			&msg.FromUserID,
			&msg.ToUserID,
			&msg.ItemID,
			&msg.Text,
			&msg.CreatedAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
