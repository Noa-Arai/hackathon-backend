package dao

import (
	"database/sql"
	"hackathon-backend/model"
)

type MessageRepository interface {
	Insert(msg *model.Message) error
	FindByItemID(itemID int64) ([]model.Message, error)

	// ★ 追加する
	ListUserRooms(userID string) ([]model.MessageRoom, error)
	MarkAsRead(itemID int64, userID string) error
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

// DM一覧（ユーザーが関わる item ごとの最新メッセージ）
func (d *MessageDAO) ListUserRooms(userID string) ([]model.MessageRoom, error) {
	rows, err := d.DB.Query(`
        SELECT
            m.item_id,
            i.title,
            CASE 
                WHEN m.from_user_id = ? THEN m.to_user_id
                ELSE m.from_user_id
            END AS partner_id,
            m.text AS latest_text,
            m.created_at,
            (
                SELECT COUNT(*) FROM messages 
                WHERE item_id = m.item_id
                  AND to_user_id = ?
                  AND read_at IS NULL
            ) AS unread_count
        FROM messages m
        JOIN items i ON m.item_id = i.id
        WHERE m.from_user_id = ? OR m.to_user_id = ?
        GROUP BY m.item_id, partner_id
        ORDER BY m.created_at DESC
    `, userID, userID, userID, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []model.MessageRoom
	for rows.Next() {
		var r model.MessageRoom
		if err := rows.Scan(
			&r.ItemID,
			&r.ItemTitle,
			&r.PartnerID,
			&r.LatestText,
			&r.UpdatedAt,
			&r.UnreadCount,
		); err != nil {
			return nil, err
		}
		rooms = append(rooms, r)
	}
	return rooms, nil
}

func (d *MessageDAO) MarkAsRead(itemID int64, userID string) error {
	_, err := d.DB.Exec(`
        UPDATE messages 
        SET read_at = NOW() 
        WHERE item_id = ? 
          AND to_user_id = ? 
          AND read_at IS NULL
    `, itemID, userID)
	return err
}
