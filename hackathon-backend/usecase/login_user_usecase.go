package usecase

import (
	"errors"
	"hackathon-backend/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidLogin = errors.New("invalid email or password")
)

type LoginUserRepository interface {
	FindByEmail(email string) (*model.User, error)
}

type LoginUserUsecase struct {
	Repo      LoginUserRepository
	JWTSecret string
}

func NewLoginUserUsecase(repo LoginUserRepository, secret string) *LoginUserUsecase {
	return &LoginUserUsecase{
		Repo:      repo,
		JWTSecret: secret,
	}
}

func (uc *LoginUserUsecase) Execute(email, password string) (string, error) {

	user, err := uc.Repo.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrInvalidLogin
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidLogin
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(uc.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
