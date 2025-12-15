package controller

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"hackathon-backend/middleware"
	"hackathon-backend/model"
	"hackathon-backend/usecase"
)

type UpdateItemController struct {
	Usecase *usecase.UpdateItemUsecase
}

func NewUpdateItemController(uc *usecase.UpdateItemUsecase) *UpdateItemController {
	return &UpdateItemController{Usecase: uc}
}

func (c *UpdateItemController) Handle(w http.ResponseWriter, r *http.Request) {

	// =======================
	// JWT ユーザー取得
	// =======================
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// =======================
	// POST 以外拒否
	// =======================
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// =======================
	// multipart パース
	// =======================
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		http.Error(w, "Invalid multipart form", http.StatusBadRequest)
		return
	}

	// =======================
	// フォーム値
	// =======================
	itemID := r.FormValue("id")
	title := r.FormValue("title")
	description := r.FormValue("description")
	category := r.FormValue("category")
	priceStr := r.FormValue("price")

	price, _ := strconv.Atoi(priceStr)

	// =======================
	// Model を作成
	// =======================
	item := &model.Item{
		ID:          itemID, // ID は string のままで OK
		UserID:      userID,
		Title:       title,
		Description: description,
		Category:    category,
		Price:       price,
	}

	// =======================
	// image1
	// =======================
	file1, header1, _ := r.FormFile("image1")
	if file1 != nil {
		defer file1.Close()
		data, _ := io.ReadAll(file1)
		item.Image1Data = data
		item.Image1Type = header1.Header.Get("Content-Type")
	}

	// =======================
	// image2
	// =======================
	file2, header2, _ := r.FormFile("image2")
	if file2 != nil {
		defer file2.Close()
		data, _ := io.ReadAll(file2)
		item.Image2Data = data
		item.Image2Type = header2.Header.Get("Content-Type")
	}

	// =======================
	// image3
	// =======================
	file3, header3, _ := r.FormFile("image3")
	if file3 != nil {
		defer file3.Close()
		data, _ := io.ReadAll(file3)
		item.Image3Data = data
		item.Image3Type = header3.Header.Get("Content-Type")
	}

	// =======================
	// 更新実行
	// =======================
	if err := c.Usecase.Execute(item); err != nil {
		fmt.Println("UPDATE ERROR:", err)
		http.Error(w, "Failed to update item", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"status":"ok"}`))
}
