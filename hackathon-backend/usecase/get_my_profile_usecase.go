package usecase

import "hackathon-backend/model"

type GetMyProfileRepository interface {
	FindByID(id string) (*model.User, error)
}

type GetMyProfileUsecase struct {
	Repo GetMyProfileRepository
}

func NewGetMyProfileUsecase(r GetMyProfileRepository) *GetMyProfileUsecase {
	return &GetMyProfileUsecase{Repo: r}
}

func (uc *GetMyProfileUsecase) Execute(userID string) (*model.User, error) {
	return uc.Repo.FindByID(userID)
}
