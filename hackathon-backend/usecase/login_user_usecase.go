package usecase

import (
	"errors"
	"hackathon-backend/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ★ ここに定義するのが正解！
var ErrInvalidLogin = errors.New("invalid email or password")

type LoginUserRepository interface {
	FindByEmail(email string) (*model.User, error)
}

type LoginUserUsecase struct {
	Repo      LoginUserRepository
	JWTSecret string
}

func NewLoginUserUsecase(repo LoginUserRepository, secret string) *LoginUserUsecase {
	return &LoginUserUsecase{Repo: repo, JWTSecret: secret}
}

func (uc *LoginUserUsecase) Execute(email, password string) (string, error) {

	user, err := uc.Repo.FindByEmail(email)
	if err != nil || user == nil {
		return "", ErrInvalidLogin
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidLogin
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.JWTSecret))
}
