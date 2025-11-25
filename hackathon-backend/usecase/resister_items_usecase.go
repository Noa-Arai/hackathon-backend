package usecase

import (
	"strconv"

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

	uid, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return err
	}

	item := &model.Item{
		UserID:      uid,
		Title:       title,
		Description: desc,
		Price:       price,
		ImageURL:    img,
	}

	return uc.Repo.Insert(item)
}
