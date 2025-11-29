package usecase

type UpdateAvatarRepository interface {
	UpdateAvatar(id string, data []byte, mime string) error
}

type UpdateAvatarUsecase struct {
	Repo UpdateAvatarRepository
}

func NewUpdateAvatarUsecase(r UpdateAvatarRepository) *UpdateAvatarUsecase {
	return &UpdateAvatarUsecase{Repo: r}
}

func (uc *UpdateAvatarUsecase) Execute(id string, data []byte, mime string) error {
	return uc.Repo.UpdateAvatar(id, data, mime)
}
