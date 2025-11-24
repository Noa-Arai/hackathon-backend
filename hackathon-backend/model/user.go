package model

type User struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	CreatedAt    string `json:"created_at"`
}

func (u *User) Validate() bool {
	if u.Name == "" || len(u.Name) > 50 {
		return false
	}
	if u.Email == "" || len(u.Email) > 255 {
		return false
	}
	if u.PasswordHash == "" {
		return false
	}
	return true
}
