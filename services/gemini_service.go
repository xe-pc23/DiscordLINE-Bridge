package services

import (
	"context"
	"encoding/json"
	"fmt"
	"line-discord-bridge/config"
	"log"
	"sync"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type ChatMessage struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

type AdviceResponse struct {
	ShouldAdvise     bool   `json:"should_advise"`
	AdviceForDiscord string `json:"advice_for_discord"`
	AdviceForLine    string `json:"advice_for_line"`
}

type GeminiService struct {
	client  *genai.Client
	model   *genai.GenerativeModel
	history []ChatMessage
	mu      sync.Mutex
}

var GeminiServiceInstance *GeminiService

func InitGeminiService() error {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(config.AppConfig.GeminiAPIKey))
	if err != nil {
		return fmt.Errorf("failed to create gemini client: %v", err)
	}

	model := client.GenerativeModel("gemini-2.0-flash-lite")
	model.ResponseMIMEType = "application/json"

	GeminiServiceInstance = &GeminiService{
		client:  client,
		model:   model,
		history: make([]ChatMessage, 0),
	}

	log.Println("Gemini Service initialized")
	return nil
}

func (s *GeminiService) Close() {
	if s.client != nil {
		s.client.Close()
	}
}

func (s *GeminiService) AddMessage(sender, content string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.history = append(s.history, ChatMessage{
		Sender:  sender,
		Content: content,
	})

	// 履歴は最新20件のみ保持
	if len(s.history) > 20 {
		s.history = s.history[len(s.history)-20:]
	}
}

func (s *GeminiService) AnalyzeChat(ctx context.Context) (*AdviceResponse, error) {
	s.mu.Lock()
	historyJSON, _ := json.Marshal(s.history)
	s.mu.Unlock()

	prompt := fmt.Sprintf(`
あなたはLINEとDiscordの会話を仲介するAIアシスタントだよ。
以下の会話履歴を分析して、特に直近のLINEユーザーからの返答内容について、Discordユーザー向けにアドバイスしてね。
LINEユーザーへのアドバイスは不要だよ。
アドバイスは本当に必要な場合のみでOK。挨拶や単純な会話、相槌のみの場合は不要だよ。
会話が停滞している時、誤解が生じそうな時、より良い表現がある時などにアドバイスして。
口調はタメ口で、短く簡潔にお願い！

会話履歴:
%s

以下のJSON形式で出力してね:
{
  "should_advise": boolean,
  "advice_for_discord": "string (LINEユーザーへの返信についてのアドバイスなど)",
  "advice_for_line": ""
}
`, string(historyJSON))

	resp, err := s.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %v", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no content generated")
	}

	var advice AdviceResponse
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			if err := json.Unmarshal([]byte(txt), &advice); err != nil {
				log.Printf("Failed to unmarshal gemini response: %v, text: %s", err, txt)
				continue
			}
			return &advice, nil
		}
	}

	return nil, fmt.Errorf("no text part in response")
}
