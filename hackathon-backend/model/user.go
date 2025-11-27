package model

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`

	// 認証用
	PasswordHash string `json:"-"`

	// プロフィール
	Bio      string `json:"bio"`
	Birthday string `json:"birthday"`

	// アバター関連
	AvatarURL  string `json:"avatar_url"`
	AvatarData []byte `json:"-"`
	AvatarType string `json:"-"`
}
