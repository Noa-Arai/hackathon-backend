package controller

import (
	"encoding/json"
	"net/http"

	"hackathon-backend/usecase"
)

type AIController struct {
	Usecase *usecase.AIUsecase
}

func NewAIController(uc *usecase.AIUsecase) *AIController {
	return &AIController{Usecase: uc}
}

// ------------------------------
// 共通: JSONレスポンスヘルパー
// ------------------------------
func jsonResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// ------------------------------
// 商品説明生成 API
// ------------------------------
type DescribeRequest struct {
	Title string `json:"title"`
}

func (c *AIController) Describe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonResponse(w, map[string]string{"error": "Method Not Allowed"}, http.StatusMethodNotAllowed)
		return
	}

	var req DescribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, map[string]string{"error": "Invalid request"}, http.StatusBadRequest)
		return
	}

	text, err := c.Usecase.GenerateDescription(req.Title)
	if err != nil {
		jsonResponse(w, map[string]string{"error": "AI generation failed"}, http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]string{
		"description": text,
	}, http.StatusOK)
}

// ------------------------------
// QA API
// ------------------------------
type AskRequest struct {
	Question string `json:"question"`
}

func (c *AIController) Ask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonResponse(w, map[string]string{"error": "Method Not Allowed"}, http.StatusMethodNotAllowed)
		return
	}

	var req AskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, map[string]string{"error": "Invalid request"}, http.StatusBadRequest)
		return
	}

	answer, err := c.Usecase.AnswerQuestion(req.Question)
	if err != nil {
		jsonResponse(w, map[string]string{"error": "AI answer failed"}, http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]string{
		"answer": answer,
	}, http.StatusOK)
}

// 感情検索のエンドポイント
func (c *AIController) HandleEmotionSearch(w http.ResponseWriter, r *http.Request) {
	// クエリパラメータから感情を取得 (?emotion=イライラ)
	emotion := r.URL.Query().Get("emotion")
	if emotion == "" {
		http.Error(w, "Emotion is required", http.StatusBadRequest)
		return
	}

	// UseCaseを呼ぶ
	keywords, err := c.Usecase.SuggestKeywords(emotion)
	if err != nil {
		http.Error(w, "AI Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// JSONで返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"emotion":  emotion,
		"keywords": keywords,
	})
}
