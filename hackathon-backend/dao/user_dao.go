package dao

import (
	"database/sql"
	"hackathon-backend/model"
)

type UserRepository interface {
	Insert(u *model.User) error
	FindByEmail(email string) (*model.User, error)
}

type UserDAO struct {
	DB *sql.DB
}

func NewUserDAO(db *sql.DB) *UserDAO {
	return &UserDAO{DB: db}
}

// 新規登録（INSERT）
func (d *UserDAO) Insert(u *model.User) error {
	_, err := d.DB.Exec(`
        INSERT INTO users (id, name, email, password_hash, created_at)
        VALUES (?, ?, ?, ?, NOW())
    `,
		u.ID, u.Name, u.Email, u.PasswordHash,
	)
	return err
}

// emailを使ってユーザーを取得（重複チェックに使う）
func (d *UserDAO) FindByEmail(email string) (*model.User, error) {
	row := d.DB.QueryRow(`
        SELECT id, name, email, password_hash, created_at
        FROM users WHERE email = ?
    `, email)

	var u model.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 見つからなかった → nil を返す
		}
		return nil, err
	}

	return &u, nil
}
