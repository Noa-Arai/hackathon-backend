package controller

import (
	"hackathon-backend/usecase"
	"net/http"
)

type GetItemImageController struct {
	Usecase *usecase.GetItemImageUsecase
}

func NewGetItemImageController(uc *usecase.GetItemImageUsecase) *GetItemImageController {
	return &GetItemImageController{Usecase: uc}
}

func (c *GetItemImageController) Handle(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	data, mime, err := c.Usecase.Execute(id)
	if err != nil || data == nil {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	// MIME を設定
	w.Header().Set("Content-Type", mime)

	// 画像データを返す
	w.Write(data)
}
