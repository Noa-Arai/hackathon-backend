package controller

import (
	"hackathon-backend/dao"
	"net/http"
)

type GetAvatarController struct {
	DAO *dao.UserDAO
}

func NewGetAvatarController(d *dao.UserDAO) *GetAvatarController {
	return &GetAvatarController{DAO: d}
}

func (c *GetAvatarController) Handle(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	data, mime, err := c.DAO.FindAvatar(id)
	if err != nil || data == nil {
		http.Error(w, "No avatar", 404)
		return
	}

	w.Header().Set("Content-Type", mime)
	w.Write(data)
}
