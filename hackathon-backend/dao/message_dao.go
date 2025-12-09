package dao

import (
	"database/sql"
	"fmt"
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
	// ✅ 修正1: NULL許容のカラム(read_at)を正しく扱う
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
		// ReadAt はポインタ型 (*time.Time) なので、NULLなら nil が入る（これでOK）
		err := rows.Scan(&m.ID, &m.FromUserID, &m.ToUserID, &m.ItemID, &m.Text, &m.CreatedAt, &m.ReadAt)
		if err != nil {
			return nil, err
		}
		list = append(list, m)
	}

	// ✅ 既読更新処理 (変更なし)
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

// 🔥 ここが修正の本丸です
func (d *MessageDAO) ListUserRooms(userID string) ([]model.MessageRoom, error) {
	// ✅ 修正2: GROUP BY を廃止し、最新順に全件取得する
	// COALESCE(u.avatar_type, '') でNULLエラーを回避
	query := `
    SELECT
        m.item_id,
        m.from_user_id,
        m.to_user_id,
        m.text,
        m.created_at,
        u.name,
        COALESCE(u.avatar_type, ''), 
        (
            SELECT COUNT(*)
            FROM messages sub
            WHERE sub.item_id = m.item_id
              AND sub.to_user_id = ? 
              AND sub.read_at IS NULL
        ) AS unread_count
    FROM messages m
    JOIN users u
      ON u.id = CASE 
          WHEN m.from_user_id = ? THEN m.to_user_id 
          ELSE m.from_user_id 
      END
    WHERE m.from_user_id = ? OR m.to_user_id = ?
    ORDER BY m.created_at DESC
    `

	// 引数の順番に注意:
	// 1. unread_count用 (userID)
	// 2. JOIN用 (userID)
	// 3. WHERE from (userID)
	// 4. WHERE to (userID)
	rows, err := d.DB.Query(query, userID, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []model.MessageRoom

	// ✅ 修正3: Go側で重複排除（同じ部屋ならスキップ）する
	// キー: "item_id-partner_id"
	seen := make(map[string]bool)

	for rows.Next() {
		var itemID int64
		var fromID, toID, text, partnerName, partnerAvatar string
		var createdAt time.Time
		var unreadCount int

		err := rows.Scan(&itemID, &fromID, &toID, &text, &createdAt, &partnerName, &partnerAvatar, &unreadCount)
		if err != nil {
			return nil, err
		}

		// 相手のIDを特定
		partnerID := toID
		if fromID != userID {
			partnerID = fromID
		}

		// ユニークキーを作成 (商品ID + 相手ID)
		key := fmt.Sprintf("%d-%s", itemID, partnerID)

		// すでに見た部屋ならスキップ（ORDER BY DESCなので、最初に来たのが最新メッセージ）
		if seen[key] {
			continue
		}
		seen[key] = true

		// リストに追加
		rooms = append(rooms, model.MessageRoom{
			ItemID:        itemID,
			PartnerID:     partnerID,
			PartnerName:   partnerName,
			PartnerAvatar: partnerAvatar,
			LastMessage:   text,
			LastTime:      createdAt,
			UnreadCount:   unreadCount,
		})
	}

	return rooms, nil
}
