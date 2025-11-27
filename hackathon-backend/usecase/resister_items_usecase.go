package usecase

import "hackathon-backend/model"

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
	userID string,
	title string,
	desc string,
	price int,
	imageData []byte,
	imageType string,
) error {

	item := &model.Item{
		UserID:      userID,
		Title:       title,
		Description: desc,
		Price:       price,
		ImageData:   imageData,
		ImageType:   imageType,
	}

	return uc.Repo.Insert(item)
}
