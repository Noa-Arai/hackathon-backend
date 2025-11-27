package controller

import (
	"hackathon-backend/middleware"
	"hackathon-backend/usecase"
	"io"
	"net/http"
)

type UpdateAvatarController struct {
	Usecase *usecase.UpdateAvatarUsecase
}

func NewUpdateAvatarController(u *usecase.UpdateAvatarUsecase) *UpdateAvatarController {
	return &UpdateAvatarController{Usecase: u}
}

func (c *UpdateAvatarController) Handle(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	// multipart/form-data から "avatar" ファイルを取得
	file, header, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, "avatar file required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "read file error", http.StatusInternalServerError)
		return
	}

	mime := header.Header.Get("Content-Type")

	// Usecase 実行
	if err := c.Usecase.Execute(userID, data, mime); err != nil {
		http.Error(w, "update avatar failed", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"status":"ok"}`))
}
