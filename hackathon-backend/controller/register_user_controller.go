package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"hackathon-backend/usecase"
)

type RegisterUserController struct {
	Usecase *usecase.RegisterUserUsecase
}

func NewRegisterUserController(u *usecase.RegisterUserUsecase) *RegisterUserController {
	return &RegisterUserController{Usecase: u}
}

type UserReqForHTTPPost struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *RegisterUserController) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UserReqForHTTPPost
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("json decode error: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := c.Usecase.Execute(req.Name, req.Email, req.Password)
	if err != nil {
		log.Printf("usecase error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := map[string]string{"id": id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
