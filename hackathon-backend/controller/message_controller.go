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
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	fromUserID, _ := r.Context().Value(middleware.UserIDKey).(string)

	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.ToUserID == "" || req.Text == "" {
		http.Error(w, "to_user_id and text required", http.StatusBadRequest)
		return
	}

	if err := c.Usecase.SendMessage(fromUserID, req.ToUserID, req.ItemID, req.Text); err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"status":"sent"}`))
}

// GET /messages?item_id=1
func (c *MessageController) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	itemIDStr := r.URL.Query().Get("item_id")
	if itemIDStr == "" {
		http.Error(w, "item_id required", http.StatusBadRequest)
		return
	}

	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid item_id", http.StatusBadRequest)
		return
	}

	messages, err := c.Usecase.GetMessages(itemID)
	if err != nil {
		http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// GET /messages/list  ← DMルーム一覧
func (c *MessageController) ListRooms(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	rooms, err := c.Usecase.ListUserRooms(userID)
	if err != nil {
		http.Error(w, "failed to load message rooms", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms)
}

// POST /messages/read
func (c *MessageController) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	itemIDStr := r.URL.Query().Get("item_id")
	if itemIDStr == "" {
		http.Error(w, "item_id required", http.StatusBadRequest)
		return
	}

	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid item_id", http.StatusBadRequest)
		return
	}

	if err := c.Usecase.MarkAsRead(itemID, userID); err != nil {
		http.Error(w, "failed to update read status", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}
