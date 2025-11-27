package usecase

import "hackathon-backend/dao"

type GetItemImageUsecase struct {
	Repo *dao.ItemDAO
}

func NewGetItemImageUsecase(r *dao.ItemDAO) *GetItemImageUsecase {
	return &GetItemImageUsecase{Repo: r}
}

func (uc *GetItemImageUsecase) Execute(id string, index int) ([]byte, string, error) {
	return uc.Repo.FindImage(id, index)
}
