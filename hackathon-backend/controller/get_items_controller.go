package controller

import (
	"encoding/json"
	"net/http"

	"hackathon-backend/usecase"
)

type GetItemsController struct {
	Usecase *usecase.GetItemsUsecase
}

func NewGetItemsController(u *usecase.GetItemsUsecase) *GetItemsController {
	return &GetItemsController{Usecase: u}
}

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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
