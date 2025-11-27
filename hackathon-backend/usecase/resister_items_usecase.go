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

func (uc *RegisterItemUsecase) Execute(
	userID, title, desc string, price int,
	img1 []byte, img1Type string,
	img2 []byte, img2Type string,
	img3 []byte, img3Type string,
) error {

	item := &model.Item{
		UserID:      userID,
		Title:       title,
		Description: desc,
		Price:       price,

		Image1Data: img1,
		Image1Type: img1Type,
		Image2Data: img2,
		Image2Type: img2Type,
		Image3Data: img3,
		Image3Type: img3Type,
	}

	return uc.Repo.Insert(item)
}
