package usecase

type ImageRepository interface {
	FindImageByID(id string) ([]byte, string, error)
}

type GetItemImageUsecase struct {
	Repo ImageRepository
}

func NewGetItemImageUsecase(r ImageRepository) *GetItemImageUsecase {
	return &GetItemImageUsecase{Repo: r}
}

func (uc *GetItemImageUsecase) Execute(id string) ([]byte, string, error) {
	return uc.Repo.FindImageByID(id)
}
