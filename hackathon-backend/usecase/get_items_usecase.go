package usecase

import "hackathon-backend/model"

type ItemRepository interface {
	FindAll() ([]model.Item, error)
}

type GetItemsUsecase struct {
	Repo ItemRepository
}

func NewGetItemsUsecase(repo ItemRepository) *GetItemsUsecase {
	return &GetItemsUsecase{Repo: repo}
}

func (uc *GetItemsUsecase) Execute() ([]model.Item, error) {
	return uc.Repo.FindAll()
}
