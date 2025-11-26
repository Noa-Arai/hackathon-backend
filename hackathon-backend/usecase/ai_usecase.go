package usecase

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/genai"
)

type AIUsecase struct{}

func NewAIUsecase() *AIUsecase {
	return &AIUsecase{}
}

func (uc *AIUsecase) GenerateDescription(title string) (string, error) {
	apiKey := os.Getenv("GOOGLE_GENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GOOGLE_GENAI_API_KEY is not set")
	}

	// v0.3.0 は ClientConfig を渡す
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return "", err
	}
	// v0.3.0 には Close() が存在しない

	prompt := fmt.Sprintf(
		"商品名: %s\nこの商品の魅力が自然に伝わる200文字以内の商品説明文を作成してください。",
		title,
	)

	// genai.Text(prompt) は []*genai.Content を返す
	resp, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash-exp",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", err
	}

	text, err := resp.Text()
	if err != nil {
		return "", err
	}

	return text, nil
}

func (uc *AIUsecase) AnswerQuestion(question string) (string, error) {
	apiKey := os.Getenv("GOOGLE_GENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GOOGLE_GENAI_API_KEY is not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf("質問：%s\n分かりやすく簡潔に答えてください。", question)

	resp, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash-exp",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", err
	}

	text, err := resp.Text()
	if err != nil {
		return "", err
	}

	return text, nil
}
