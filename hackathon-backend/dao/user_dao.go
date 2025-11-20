package dao

import (
	"database/sql"
	"hackathon-backend/model"
)

type UserRepository interface {
	Insert(u *model.User) error
	FindByName(name string) ([]model.User, error)
}

type UserDAO struct {
	DB *sql.DB
}

func NewUserDAO(db *sql.DB) *UserDAO {
	return &UserDAO{DB: db}
}

func (d *UserDAO) Insert(u *model.User) error {
	tx, err := d.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec("INSERT INTO user (id, name, age) VALUES (?, ?, ?)", u.ID, u.Name, u.Age)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (d *UserDAO) FindByName(name string) ([]model.User, error) {
	rows, err := d.DB.Query("SELECT id, name, age FROM user WHERE name = ?", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []model.User{}
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Age); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
