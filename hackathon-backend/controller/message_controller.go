package controller

import (
	"encoding/json"
	"hackathon-backend/middleware"
	"hackathon-backend/usecase"
	"net/http"
	"strconv"
)

type SendMessageRequest struct {
	ToUserID string `json:"to_user_id"`
	ItemID   int64  `json:"item_id"`
	Text     string `json:"text"`
}

type MessageController struct {
	Usecase *usecase.MessageUsecase
}

func NewMessageController(uc *usecase.MessageUsecase) *MessageController {
	return &MessageController{Usecase: uc}
}

// POST /messages
func (c *MessageController) Send(w http.ResponseWriter, r *http.Request) {
	fromUserID, _ := r.Context().Value(middleware.UserIDKey).(string)

	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.ToUserID == "" || req.Text == "" {
		http.Error(w, "to_user_id and text required", http.StatusBadRequest)
		return
	}

	err := c.Usecase.SendMessage(fromUserID, req.ToUserID, req.ItemID, req.Text)
	if err != nil {
		http.Error(w, "send failed", 500)
		return
	}

	w.Write([]byte(`{"status":"sent"}`))
}

// GET /messages?item_id=xxx
func (c *MessageController) List(w http.ResponseWriter, r *http.Request) {
	itemIDStr := r.URL.Query().Get("item_id")

	itemID, _ := strconv.ParseInt(itemIDStr, 10, 64)
	msgs, err := c.Usecase.GetMessages(itemID)
	if err != nil {
		http.Error(w, "load failed", 500)
		return
	}

	json.NewEncoder(w).Encode(msgs)
}

// GET /messages/rooms
func (c *MessageController) ListRooms(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	rooms, err := c.Usecase.ListUserRooms(userID)
	if err != nil {
		http.Error(w, "list failed", 500)
		return
	}
	json.NewEncoder(w).Encode(rooms)
}

// GET /messages/list?item_id=xxx&partner_id=xxx
func (c *MessageController) ListChat(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	itemIDStr := r.URL.Query().Get("item_id")
	partnerID := r.URL.Query().Get("partner_id")

	itemID, _ := strconv.ParseInt(itemIDStr, 10, 64)

	msgs, err := c.Usecase.ListChat(userID, partnerID, itemID)
	if err != nil {
		http.Error(w, "chat load failed", 500)
		return
	}

	json.NewEncoder(w).Encode(msgs)
}

// POST /messages/read?item_id=xxx
func (c *MessageController) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	itemIDStr := r.URL.Query().Get("item_id")
	itemID, _ := strconv.ParseInt(itemIDStr, 10, 64)

	err := c.Usecase.MarkAsRead(itemID, userID)
	if err != nil {
		http.Error(w, "mark failed", 500)
		return
	}

	w.Write([]byte(`{"status":"ok"}`))
}
