package controller

import (
	"encoding/json"
	"net/http"

	"hackathon-backend/usecase"
)

type GetItemsController struct {
	Usecase      *usecase.GetItemsUsecase
	ImageUsecase *usecase.GetItemImageUsecase
}

func NewGetItemsController(
	u *usecase.GetItemsUsecase,
	img *usecase.GetItemImageUsecase,
) *GetItemsController {
	return &GetItemsController{
		Usecase:      u,
		ImageUsecase: img,
	}
}

// ===== 商品一覧 GET =====
func (c *GetItemsController) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	items, err := c.Usecase.Execute()
	if err != nil {
		http.Error(w, "Failed to fetch items", http.StatusInternalServerError)
		return
	}

	// 一覧に3枚の画像エンドポイントを付与
	for i := range items {
		id := items[i].ID
		items[i].Image1URL = "/items/image1?id=" + id
		items[i].Image2URL = "/items/image2?id=" + id
		items[i].Image3URL = "/items/image3?id=" + id
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// ===== 画像取得 GET =====
func (c *GetItemsController) HandleImage(index int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Missing id", 400)
			return
		}

		data, mime, err := c.ImageUsecase.Execute(id, index)
		if err != nil || len(data) == 0 {
			http.Error(w, "Not Found", 404)
			return
		}

		w.Header().Set("Content-Type", mime)
		w.Write(data)
	}
}
