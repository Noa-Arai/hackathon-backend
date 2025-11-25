package controller

import (
	"encoding/json"
	"hackathon-backend/middleware"
	"hackathon-backend/usecase"
	"net/http"
)

type PurchaseRequest struct {
	ItemID int64 `json:"item_id"`
}

type PurchaseController struct {
	Usecase *usecase.PurchaseUsecase
}

func NewPurchaseController(uc *usecase.PurchaseUsecase) *PurchaseController {
	return &PurchaseController{Usecase: uc}
}

func (c *PurchaseController) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	buyerID := r.Context().Value(middleware.UserIDKey).(string)

	var req PurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := c.Usecase.Execute(req.ItemID, buyerID); err != nil {
		http.Error(w, "Failed to purchase", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"purchased"}`))
}
