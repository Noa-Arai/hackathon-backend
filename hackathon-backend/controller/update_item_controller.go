package controller

import (
	"io"
	"net/http"
	"strconv"

	"hackathon-backend/usecase"
)

type UpdateItemController struct {
	Usecase *usecase.UpdateItemUsecase
}

func NewUpdateItemController(u *usecase.UpdateItemUsecase) *UpdateItemController {
	return &UpdateItemController{Usecase: u}
}

func (c *UpdateItemController) Handle(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}

	userID := r.Context().Value("user_id").(string)

	itemID := r.FormValue("item_id")
	title := r.FormValue("title")
	desc := r.FormValue("description")
	price, _ := strconv.Atoi(r.FormValue("price"))

	var imgData [3][]byte
	var imgType [3]string

	files := r.MultipartForm.File["images"]

	for i := 0; i < 3; i++ {
		if i < len(files) {
			f, _ := files[i].Open()
			bin, _ := io.ReadAll(f)
			f.Close()

			imgData[i] = bin
			imgType[i] = files[i].Header.Get("Content-Type")
		}
	}

	err := c.Usecase.Execute(
		itemID, userID, title, desc, price,
		imgData[0], imgType[0],
		imgData[1], imgType[1],
		imgData[2], imgType[2],
	)
	if err != nil {
		http.Error(w, "update failed", 500)
		return
	}

	w.Write([]byte(`{"ok":true}`))
}
