package usecase

import "hackathon-backend/model"

type UpdateProfileRepository interface {
	UpdateProfile(user *model.User) error
}

type UpdateProfileUsecase struct {
	Repo UpdateProfileRepository
}

func NewUpdateProfileUsecase(r UpdateProfileRepository) *UpdateProfileUsecase {
	return &UpdateProfileUsecase{Repo: r}
}

func (uc *UpdateProfileUsecase) Execute(userID, name, bio, birthday string) error {
	u := &model.User{
		ID:       userID,
		Name:     name,
		Bio:      bio,
		Birthday: birthday,
	}

	return uc.Repo.UpdateProfile(u)
}
