package usecase

import (
	"hackathon-backend/model"
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

// Repository interface
type RegisterUserRepository interface {
	Insert(u *model.User) error
	FindByEmail(email string) (*model.User, error)
}

type RegisterUserUsecase struct {
	Repo RegisterUserRepository
}

func NewRegisterUserUsecase(repo RegisterUserRepository) *RegisterUserUsecase {
	return &RegisterUserUsecase{Repo: repo}
}

func (uc *RegisterUserUsecase) Execute(name, email, password string) (string, error) {

	// email 重複チェック
	existing, err := uc.Repo.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if existing != nil {
		return "", ErrEmailExists
	}

	// パスワードハッシュ化
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// ID 生成（ULID）
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()

	// モデル作成
	user := &model.User{
		ID:           id,
		Name:         name,
		Email:        email,
		PasswordHash: string(hashed),
	}

	// 簡単なバリデーション
	if user.Name == "" || user.Email == "" || password == "" {
		return "", ErrInvalidUser
	}

	// DB保存
	if err := uc.Repo.Insert(user); err != nil {
		return "", err
	}

	return id, nil
}
