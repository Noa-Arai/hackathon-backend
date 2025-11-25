package usecase

import (
	"hackathon-backend/dao"
	"hackathon-backend/model"
)

type PurchaseUsecase struct {
	Repo dao.PurchaseRepository
}

func NewPurchaseUsecase(repo dao.PurchaseRepository) *PurchaseUsecase {
	return &PurchaseUsecase{Repo: repo}
}

func (uc *PurchaseUsecase) Execute(itemID int64, buyerID string) error {

	purchase := &model.Purchase{
		ItemID:  itemID,
		BuyerID: buyerID,
	}

	return uc.Repo.Insert(purchase)
}
