package controller

import (
	"encoding/json"
	"hackathon-backend/usecase"
	"log"
	"net/http"
)

type SearchUserController struct {
	Usecase *usecase.SearchUserUsecase
}

func NewSearchUserController(u *usecase.SearchUserUsecase) *SearchUserController {
	return &SearchUserController{Usecase: u}
}

func (c *SearchUserController) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	name := r.URL.Query().Get("name")
	users, err := c.Usecase.Execute(name)
	if err != nil {
		log.Printf("fail: usecase.Execute, %v\n", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
