package controller

import (
	"encoding/json"
	"hackathon-backend/middleware"
	"hackathon-backend/usecase"
	"net/http"
)

type ProfileReq struct {
	Name     string `json:"name"`
	Bio      string `json:"bio"`
	Birthday string `json:"birthday"`
}

type UpdateProfileController struct {
	Usecase *usecase.UpdateProfileUsecase
}

func NewUpdateProfileController(u *usecase.UpdateProfileUsecase) *UpdateProfileController {
	return &UpdateProfileController{Usecase: u}
}

func (c *UpdateProfileController) Handle(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	var req ProfileReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := c.Usecase.Execute(userID, req.Name, req.Bio, req.Birthday); err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"status":"ok"}`))
}
