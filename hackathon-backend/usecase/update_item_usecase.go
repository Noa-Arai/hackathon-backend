package usecase

import "hackathon-backend/model"

type UpdateItemRepository interface {
	Update(item *model.Item) error
}

type UpdateItemUsecase struct {
	Repo UpdateItemRepository
}

func NewUpdateItemUsecase(r UpdateItemRepository) *UpdateItemUsecase {
	return &UpdateItemUsecase{Repo: r}
}

func (uc *UpdateItemUsecase) Execute(
	itemID, userID, title, desc string, price int,
	img1 []byte, img1Type string,
	img2 []byte, img2Type string,
	img3 []byte, img3Type string,
) error {

	item := &model.Item{
		ID:          itemID,
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

	return uc.Repo.Update(item)
}
