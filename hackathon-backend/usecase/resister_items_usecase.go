package usecase

import (
	"hackathon-backend/model"
)

type RegisterItemRepository interface {
	Insert(item *model.Item) error
}

type RegisterItemUsecase struct {
	Repo RegisterItemRepository
}

func NewRegisterItemUsecase(r RegisterItemRepository) *RegisterItemUsecase {
	return &RegisterItemUsecase{Repo: r}
}

func (uc *RegisterItemUsecase) Execute(userID string, title, desc string, price int, img string) error {

	item := &model.Item{
		UserID:      userID, // ← ULID をそのまま保存
		Title:       title,
		Description: desc,
		Price:       price,
		ImageURL:    img,
	}

	return uc.Repo.Insert(item)
}
