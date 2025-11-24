package usecase

import (
	"errors"
	"hackathon-backend/model"

	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

var ErrEmailAlreadyUsed = errors.New("email already used")

// UserRepository は DAO が実装するインターフェース
type UserRepository interface {
	Insert(u *model.User) error
	FindByEmail(email string) (*model.User, error)
}

type RegisterUserUsecase struct {
	Repo UserRepository
}

func NewRegisterUserUsecase(repo UserRepository) *RegisterUserUsecase {
	return &RegisterUserUsecase{Repo: repo}
}

// Execute: name, email, password を受け取る本物の /signup 処理
func (uc *RegisterUserUsecase) Execute(name, email, password string) (*model.User, error) {

	// 1. email 重複チェック
	existing, err := uc.Repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailAlreadyUsed
	}

	// 2. password ハッシュ化
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 3. User 構造体を作成
	user := &model.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashed),
	}

	// 4. ULID 生成（今のコードはとても良い）
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	user.ID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()

	// 5. DB 保存
	if err := uc.Repo.Insert(user); err != nil {
		return nil, err
	}

	return user, nil
}
