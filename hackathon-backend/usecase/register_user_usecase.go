package usecase

import (
	"errors"
	"hackathon-backend/model"
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailExists = errors.New("email already used")
	ErrInvalidUser = errors.New("invalid user")
)

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

	existing, err := uc.Repo.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if existing != nil {
		return "", ErrEmailExists
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()

	user := &model.User{
		ID:           id,
		Name:         name,
		Email:        email,
		PasswordHash: string(hashed),
	}

	// ★ Validate 呼び出し
	if err := user.Validate(); err != nil {
		return "", ErrInvalidUser
	}

	return id, uc.Repo.Insert(user)
}
