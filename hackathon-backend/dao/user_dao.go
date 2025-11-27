package dao

import (
	"database/sql"
	"hackathon-backend/model"
)

type UserDAO struct {
	DB *sql.DB
}

func NewUserDAO(db *sql.DB) *UserDAO {
	return &UserDAO{DB: db}
}

// ---- 新規登録用 Insert ----
func (d *UserDAO) Insert(u *model.User) error {
	_, err := d.DB.Exec(`
		INSERT INTO users (id, name, email, password_hash)
		VALUES (?, ?, ?, ?)
	`,
		u.ID, u.Name, u.Email, u.PasswordHash,
	)
	return err
}

// ---- ログイン用 email 検索 ----
func (d *UserDAO) FindByEmail(email string) (*model.User, error) {
	row := d.DB.QueryRow(`
		SELECT id, name, email, password_hash, bio, birthday
		FROM users
		WHERE email = ?
	`, email)

	var u model.User
	err := row.Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.Bio,
		&u.Birthday,
	)
	if err != nil {
		return nil, err
	}

	u.AvatarURL = "/users/avatar?id=" + u.ID
	return &u, nil
}

// ---- 自分のプロフィール取得 ----
func (d *UserDAO) FindByID(id string) (*model.User, error) {
	row := d.DB.QueryRow(`
		SELECT id, name, email, bio, birthday
		FROM users WHERE id = ?`, id)

	var u model.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Bio, &u.Birthday)
	if err != nil {
		return nil, err
	}

	u.AvatarURL = "/users/avatar?id=" + u.ID
	return &u, nil
}

// ---- プロフィール更新 ----
func (d *UserDAO) UpdateProfile(u *model.User) error {
	_, err := d.DB.Exec(`
		UPDATE users SET name=?, bio=?, birthday=? WHERE id=?
	`,
		u.Name, u.Bio, u.Birthday, u.ID,
	)
	return err
}

// ---- アバター保存 ----
func (d *UserDAO) UpdateAvatar(id string, data []byte, mime string) error {
	_, err := d.DB.Exec(`
		UPDATE users SET avatar_data=?, avatar_type=? WHERE id=?
	`, data, mime, id)
	return err
}

// ---- アバター画像取得 ----
func (d *UserDAO) FindAvatar(id string) ([]byte, string, error) {
	row := d.DB.QueryRow(`
		SELECT avatar_data, avatar_type FROM users WHERE id=?
	`, id)

	var data []byte
	var mime string
	err := row.Scan(&data, &mime)
	return data, mime, err
}
