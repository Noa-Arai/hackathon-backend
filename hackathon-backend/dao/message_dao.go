package dao

import (
	"database/sql"
	"hackathon-backend/model"
	"time"
)

type MessageRepository interface {
	Insert(msg *model.Message) error
	ListChat(userID, partnerID string, itemID int64) ([]model.Message, error)
	ListUserRooms(userID string) ([]model.MessageRoom, error)
}

type MessageDAO struct {
	DB *sql.DB
}

func NewMessageDAO(db *sql.DB) *MessageDAO {
	return &MessageDAO{DB: db}
}

func (d *MessageDAO) Insert(msg *model.Message) error {
	query := `
        INSERT INTO messages (from_user_id, to_user_id, item_id, text)
        VALUES (?, ?, ?, ?)
    `
	_, err := d.DB.Exec(query, msg.FromUserID, msg.ToUserID, msg.ItemID, msg.Text)
	return err
}

func (d *MessageDAO) ListChat(userID, partnerID string, itemID int64) ([]model.Message, error) {
	query := `
        SELECT id, from_user_id, to_user_id, item_id, text, created_at, read_at
        FROM messages
        WHERE item_id = ?
          AND ((from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?))
        ORDER BY created_at ASC
    `
	rows, err := d.DB.Query(query, itemID, userID, partnerID, partnerID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Message
	for rows.Next() {
		var m model.Message
		err := rows.Scan(&m.ID, &m.FromUserID, &m.ToUserID, &m.ItemID, &m.Text, &m.CreatedAt, &m.ReadAt)
		if err != nil {
			return nil, err
		}
		list = append(list, m)
	}

	// --- 受信者側の未読を読む側が開いたら read_at を更新 ---
	update := `
        UPDATE messages
        SET read_at = ?
        WHERE item_id = ?
          AND to_user_id = ?
          AND read_at IS NULL
    `
	d.DB.Exec(update, time.Now(), itemID, userID)

	return list, nil
}

func (d *MessageDAO) ListUserRooms(userID string) ([]model.MessageRoom, error) {
	query := `
    SELECT
        m.item_id,
        CASE WHEN m.from_user_id = ? THEN m.to_user_id ELSE m.from_user_id END AS partner_id,
        u.name AS partner_name,
        u.avatar_type AS partner_avatar,
        m.text AS last_message,
        m.created_at AS last_time,
        (
            SELECT COUNT(*)
            FROM messages
            WHERE item_id = m.item_id
              AND to_user_id = ?
              AND read_at IS NULL
        ) AS unread_count
    FROM messages m
    JOIN users u
      ON u.id = CASE WHEN m.from_user_id = ? THEN m.to_user_id ELSE m.from_user_id END
    WHERE m.from_user_id = ? OR m.to_user_id = ?
    GROUP BY m.item_id, partner_id
    ORDER BY last_time DESC
    `

	rows, err := d.DB.Query(query, userID, userID, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []model.MessageRoom
	for rows.Next() {
		var r model.MessageRoom
		err := rows.Scan(&r.ItemID, &r.PartnerID, &r.PartnerName, &r.PartnerAvatar, &r.LastMessage, &r.LastTime, &r.UnreadCount)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, r)
	}

	return rooms, nil
}
