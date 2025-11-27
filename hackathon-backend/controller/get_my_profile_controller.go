package controller

import (
	"encoding/json"
	"hackathon-backend/middleware"
	"hackathon-backend/usecase"
	"net/http"
)

type GetMyProfileController struct {
	Usecase *usecase.GetMyProfileUsecase
}

func NewGetMyProfileController(u *usecase.GetMyProfileUsecase) *GetMyProfileController {
	return &GetMyProfileController{Usecase: u}
}

func (c *GetMyProfileController) Handle(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	res, err := c.Usecase.Execute(userID)
	if err != nil {
		http.Error(w, "Failed to get profile", 500)
		return
	}

	json.NewEncoder(w).Encode(res)
}
