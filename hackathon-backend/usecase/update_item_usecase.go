package usecase

import (
	"hackathon-backend/dao"
	"hackathon-backend/model"
)

type UpdateItemUsecase struct {
	ItemDAO *dao.ItemDAO
}

func NewUpdateItemUsecase(itemDAO *dao.ItemDAO) *UpdateItemUsecase {
	return &UpdateItemUsecase{
		ItemDAO: itemDAO,
	}
}

func (uc *UpdateItemUsecase) Execute(item *model.Item) error {
	return uc.ItemDAO.Update(item)
}
