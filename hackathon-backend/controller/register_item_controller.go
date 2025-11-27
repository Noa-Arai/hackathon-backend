package controller

import (
	"hackathon-backend/middleware"
	"hackathon-backend/usecase"
	"io"
	"net/http"
	"strconv"
)

type RegisterItemController struct {
	Usecase *usecase.RegisterItemUsecase
}

func NewRegisterItemController(u *usecase.RegisterItemUsecase) *RegisterItemController {
	return &RegisterItemController{Usecase: u}
}

func (c *RegisterItemController) Handle(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middleware.UserIDKey).(string)

	if err := r.ParseMultipartForm(20 << 20); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	priceStr := r.FormValue("price")
	price, _ := strconv.Atoi(priceStr)

	file, header, err := r.FormFile("file")

	var img []byte
	var imgType string

	if err == nil {
		defer file.Close()
		img, _ = io.ReadAll(file)
		imgType = header.Header.Get("Content-Type")
	}

	if err := c.Usecase.Execute(
		userID,
		title,
		description,
		price,
		img,
		imgType,
	); err != nil {
		http.Error(w, "Insert failed", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"status":"ok"}`))
}
