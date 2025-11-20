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
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (c *RegisterUserController) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var req UserReqForHTTPPost
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("fail: json.Decode, %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := c.Usecase.Execute(req.Name, req.Age)
	if err != nil {
		log.Printf("fail: usecase.Execute, %v\n", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	res := map[string]string{"id": id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
