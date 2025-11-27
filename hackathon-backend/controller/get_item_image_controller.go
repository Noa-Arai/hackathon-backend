package controller

import (
	"net/http"
	"strconv"

	"hackathon-backend/usecase"
)

type GetItemImageController struct {
	Usecase *usecase.GetItemImageUsecase
}

func NewGetItemImageController(u *usecase.GetItemImageUsecase) *GetItemImageController {
	return &GetItemImageController{Usecase: u}
}

func (c *GetItemImageController) Handle(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	indexStr := r.URL.Query().Get("index")
	if indexStr == "" {
		indexStr = "1"
	}

	index, _ := strconv.Atoi(indexStr)

	// DBから画像取得
	data, mime, err := c.Usecase.Execute(id, index)
	if err != nil || len(data) == 0 {
		http.Error(w, "image not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", mime)
	w.Write(data)
}
