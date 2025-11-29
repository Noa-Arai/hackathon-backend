package controller

import (
	"encoding/json"
	"hackathon-backend/usecase"
	"net/http"
	"strconv"
)

type MessageController struct {
	UC *usecase.MessageUsecase
}

func NewMessageController(uc *usecase.MessageUsecase) *MessageController {
	return &MessageController{UC: uc}
}

func (c *MessageController) Send(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var body struct {
		PartnerID string `json:"partner_id"`
		ItemID    int64  `json:"item_id"`
		Text      string `json:"text"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	err := c.UC.SendMessage(userID, body.PartnerID, body.ItemID, body.Text)
	if err != nil {
		http.Error(w, "Failed to send", 500)
		return
	}

	w.Write([]byte(`{"status":"ok"}`))
}

func (c *MessageController) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	itemID, _ := strconv.ParseInt(r.URL.Query().Get("item_id"), 10, 64)
	partnerID := r.URL.Query().Get("partner_id")

	list, err := c.UC.ListChat(userID, partnerID, itemID)
	if err != nil {
		http.Error(w, "Failed", 500)
		return
	}

	json.NewEncoder(w).Encode(list)
}

func (c *MessageController) ListRooms(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	rooms, err := c.UC.ListUserRooms(userID)
	if err != nil {
		http.Error(w, "Failed", 500)
		return
	}

	json.NewEncoder(w).Encode(rooms)
}
