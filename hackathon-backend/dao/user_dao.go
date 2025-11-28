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

// -----------------------
// Find user by email (Login 用)
// -----------------------
func (d *UserDAO) FindByEmail(email string) (*model.User, error) {

	query := `
        SELECT
            id,
            name,
            email,
            password_hash,
            created_at,     -- ★追加
            avatar_data,
            avatar_type,
            bio,
            birthday
        FROM users
        WHERE email = ?
        LIMIT 1
    `

	row := d.DB.QueryRow(query, email)

	var createdAt string // ★ 受け取り専用の捨て変数

	var u model.User
	err := row.Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&createdAt, // ★ created_at をここで捨てる
		&u.AvatarData,
		&u.AvatarType,
		&u.Bio,
		&u.Birthday,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	u.AvatarURL = "/users/avatar?id=" + u.ID
	return &u, nil
}

// -----------------------
// Find by ID（プロフィール / アバター取得）
// -----------------------
func (d *UserDAO) FindByID(id string) (*model.User, error) {
	row := d.DB.QueryRow(`
		SELECT id, name, email, password_hash, bio, birthday, avatar_data, avatar_type
		FROM users
		WHERE id = ?
	`, id)

	var u model.User
	err := row.Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.Bio,
		&u.Birthday,
		&u.AvatarData,
		&u.AvatarType,
	)
	if err != nil {
		return nil, err
	}

	u.AvatarURL = "/users/avatar?id=" + u.ID
	return &u, nil
}

// -----------------------
// Create user
// -----------------------
func (d *UserDAO) Insert(u *model.User) error {
	_, err := d.DB.Exec(`
		INSERT INTO users (id, name, email, password_hash)
		VALUES (?, ?, ?, ?)
	`,
		u.ID, u.Name, u.Email, u.PasswordHash,
	)
	return err
}

// -----------------------
// Update Profile
// -----------------------
func (d *UserDAO) UpdateProfile(u *model.User) error {
	_, err := d.DB.Exec(`
		UPDATE users
		SET name=?, bio=?, birthday=?
		WHERE id=?
	`,
		u.Name, u.Bio, u.Birthday, u.ID,
	)
	return err
}

// -----------------------
// Save Avatar
// -----------------------
func (d *UserDAO) UpdateAvatar(id string, data []byte, mime string) error {
	_, err := d.DB.Exec(`
		UPDATE users
		SET avatar_data=?, avatar_type=?
		WHERE id=?
	`, data, mime, id)
	return err
}

// -----------------------
// Get Avatar
// -----------------------
func (d *UserDAO) FindAvatar(id string) ([]byte, string, error) {
	row := d.DB.QueryRow(`
		SELECT avatar_data, avatar_type
		FROM users WHERE id=?
	`, id)

	var data []byte
	var mime string
	err := row.Scan(&data, &mime)
	return data, mime, err
}
