package controller

import (
	"encoding/json"
	"fmt"
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
	userIDVal := r.Context().Value("userID")
	if userIDVal == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDVal.(string)

	var body struct {
		PartnerID string `json:"partner_id"`
		ItemID    int64  `json:"item_id"`
		Text      string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if body.Text == "" {
		http.Error(w, "Text is requires", http.StatusBadRequest)
		return
	}
	err := c.UC.SendMessage(userID, body.PartnerID, body.ItemID, body.Text)
	if err != nil {

		fmt.Printf("❌ SendMessage Erroe: %v\n", err)

		http.Error(w, "Failed to send", http.StatusInternalServerError)
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
