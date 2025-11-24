package model

import "time"

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // ← "-" は「JSON に出さない」という意味
	CreatedAt    time.Time `json:"created_at"`
}

// 今は最低限。必要なら後でバリデーション追加も可能。
func (u *User) Validate() bool {
	if u.Name == "" || len(u.Name) > 50 {
		return false
	}
	if u.Email == "" {
		return false
	}
	if u.PasswordHash == "" {
		return false
	}
	return true
}
