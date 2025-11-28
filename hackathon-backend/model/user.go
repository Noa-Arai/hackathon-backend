package model

import "errors"

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`

	PasswordHash string `json:"-"`

	Bio      string `json:"bio"`
	Birthday string `json:"birthday"`

	AvatarURL  string `json:"avatar_url"`
	AvatarData []byte `json:"-"`
	AvatarType string `json:"-"`
}

// ============================
// Validate
// ============================
func (u *User) Validate() error {

	if u.Name == "" {
		return errors.New("name is required")
	}

	if u.Email == "" {
		return errors.New("email is required")
	}

	if u.PasswordHash == "" {
		return errors.New("password is required")
	}

	// 追加したい制約があればここに書く
	// 例:
	// if len(u.Name) > 50 { return errors.New("name too long") }

	return nil
}
