package controller

import (
	"encoding/json"
	"hackathon-backend/middleware"
	"hackathon-backend/usecase"
	"net/http"
)

type RegisterItemController struct {
	Usecase *usecase.RegisterItemUsecase
}

func NewRegisterItemController(u *usecase.RegisterItemUsecase) *RegisterItemController {
	return &RegisterItemController{Usecase: u}
}

type ItemReqBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	ImageURL    string `json:"image_url"`
}

func (c *RegisterItemController) Handle(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	var body ItemReqBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := c.Usecase.Execute(userID, body.Title, body.Description, body.Price, body.ImageURL)
	if err != nil {
		http.Error(w, "Insert failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status":"ok"}`))
}
