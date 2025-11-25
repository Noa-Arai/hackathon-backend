package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"hackathon-backend/usecase"
)

type LoginUserController struct {
	Usecase *usecase.LoginUserUsecase
}

func NewLoginUserController(u *usecase.LoginUserUsecase) *LoginUserController {
	return &LoginUserController{Usecase: u}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *LoginUserController) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("json decode error: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token, err := c.Usecase.Execute(req.Email, req.Password)
	if err != nil {
		log.Printf("login error: %v", err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	res := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
