package usecase

import "hackathon-backend/model"

type SearchRepository interface {
	FindByName(name string) ([]model.User, error)
}

type SearchUserUsecase struct {
	Repo SearchRepository
}

func NewSearchUserUsecase(repo SearchRepository) *SearchUserUsecase {
	return &SearchUserUsecase{Repo: repo}
}

func (uc *SearchUserUsecase) Execute(name string) ([]model.User, error) {
	if name == "" {
		return nil, ErrInvalidName
	}
	return uc.Repo.FindByName(name)
}
