package controller

import (
	"hackathon-backend/middleware"
	"hackathon-backend/usecase"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

type RegisterItemController struct {
	Usecase *usecase.RegisterItemUsecase
}

func NewRegisterItemController(u *usecase.RegisterItemUsecase) *RegisterItemController {
	return &RegisterItemController{Usecase: u}
}

func readFile(f multipart.File) []byte {
	if f == nil {
		return nil
	}
	b, _ := io.ReadAll(f)
	return b
}

func (c *RegisterItemController) Handle(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middleware.UserIDKey).(string)

	if err := r.ParseMultipartForm(30 << 20); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	price, _ := strconv.Atoi(r.FormValue("price"))

	category := r.FormValue("category")
	isLuckyBag := r.FormValue("is_lucky_bag") == "true"

	// --- 画像3つ取得 ---
	file1, header1, _ := r.FormFile("file1")
	file2, header2, _ := r.FormFile("file2")
	file3, header3, _ := r.FormFile("file3")

	img1 := readFile(file1)
	img2 := readFile(file2)
	img3 := readFile(file3)

	t1 := ""
	t2 := ""
	t3 := ""

	if header1 != nil {
		t1 = header1.Header.Get("Content-Type")
	}
	if header2 != nil {
		t2 = header2.Header.Get("Content-Type")
	}
	if header3 != nil {
		t3 = header3.Header.Get("Content-Type")
	}

	err := c.Usecase.Execute(
		userID,
		title,
		description,
		price,
		category,
		isLuckyBag,
		img1, t1,
		img2, t2,
		img3, t3,
	)

	if err != nil {
		http.Error(w, "Insert failed", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"status":"ok"}`))
}
