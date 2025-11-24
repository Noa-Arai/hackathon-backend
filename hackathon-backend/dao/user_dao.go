package dao

import (
	"database/sql"
	"hackathon-backend/model"
)

type RegisterUserRepository interface {
	Insert(u *model.User) error
	FindByEmail(email string) (*model.User, error)
}

type UserDAO struct {
	DB *sql.DB
}

func NewUserDAO(db *sql.DB) *UserDAO {
	return &UserDAO{DB: db}
}

func (d *UserDAO) Insert(u *model.User) error {
	_, err := d.DB.Exec(
		"INSERT INTO users (id, name, email, password_hash) VALUES (?, ?, ?, ?)",
		u.ID, u.Name, u.Email, u.PasswordHash,
	)
	return err
}

func (d *UserDAO) FindByEmail(email string) (*model.User, error) {
	row := d.DB.QueryRow(
		"SELECT id, name, email, password_hash, created_at FROM users WHERE email = ?",
		email,
	)

	var user model.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
