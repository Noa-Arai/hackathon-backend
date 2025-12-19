package usecase

import (
	"context"
	"fmt"
	"os"
	"strings"

	"google.golang.org/genai"
)

type AIUsecase struct{}

func NewAIUsecase() *AIUsecase {
	return &AIUsecase{}
}

func (uc *AIUsecase) SuggestKeywords(emotion string) ([]string, error) {
	apiKey := os.Getenv("GOOGLE_GENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GOOGLE_GENAI_API_KEY is not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, err
	}

	// カンマ区切りで返させるのがポイント
	prompt := fmt.Sprintf(`
        ユーザーの今の気分は「%s」です。
        この気分に寄り添う、または解消するためのフリマアプリの「商品検索キーワード」を3つ提案してください。
        
        ルール:
        - 日本語で答えること
        - キーワードをカンマ(,)区切りで並べるだけ（余計な文章は禁止）
        - 例: 入浴剤,チョコレート,アロマキャンドル
    `, emotion)

	contents := []*genai.Content{
		{
			Parts: []*genai.Part{
				{Text: prompt},
			},
		},
	}

	resp, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash-exp",
		contents,
		nil,
	)
	if err != nil {
		return nil, err
	}

	text, err := resp.Text()
	if err != nil {
		return nil, err
	}

	// 文字列 "入浴剤, チョコ" を ["入浴剤", "チョコ"] に分割して整形
	rawKeywords := strings.Split(text, ",")
	var keywords []string
	for _, k := range rawKeywords {
		cleaned := strings.TrimSpace(k)
		if cleaned != "" {
			keywords = append(keywords, cleaned)
		}
	}

	return keywords, nil
}

// 商品説明生成
func (uc *AIUsecase) GenerateDescription(title string) (string, error) {
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

	prompt := fmt.Sprintf(
		"商品名: %s\nこの商品の魅力が自然に伝わる200文字以内の商品説明文を作成してください。",
		title,
	)

	// === v0.3.0 用に正しく変換 ===
	contents := []*genai.Content{
		{
			Parts: []*genai.Part{
				{Text: prompt},
			},
		},
	}

	resp, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash-exp",
		contents,
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

// QA（質問回答）
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

	contents := []*genai.Content{
		{
			Parts: []*genai.Part{
				{Text: prompt},
			},
		},
	}

	resp, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash-exp",
		contents,
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
