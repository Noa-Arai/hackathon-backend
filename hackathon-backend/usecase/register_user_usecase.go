package usecase

import (
	"hackathon-backend/model"
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// UserRepository インターフェース（daoが実装する）
type UserRepository interface {
	Insert(u *model.User) error
}

type RegisterUserUsecase struct {
	Repo UserRepository
}

func NewRegisterUserUsecase(repo UserRepository) *RegisterUserUsecase {
	return &RegisterUserUsecase{Repo: repo}
}

// Execute はバリデーション→ID生成→保存を行う
func (uc *RegisterUserUsecase) Execute(name string, age int) (string, error) {
	user := &model.User{Name: name, Age: age}
	if !user.Validate() {
		return "", ErrInvalidUser
	}

	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	user.ID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()

	if err := uc.Repo.Insert(user); err != nil {
		return "", err
	}
	return user.ID, nil
}
