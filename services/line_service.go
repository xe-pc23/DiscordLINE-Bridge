package services

import (
	"line-discord-bridge/config"
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

type LineService struct {
	client     *messaging_api.MessagingApiAPI
	lastUserID string // 最後にメッセージを送ってきたユーザーのID
}

var LineServiceInstance *LineService

// InitLineService initializes LINE bot client
func InitLineService() error {
	client, err := messaging_api.NewMessagingApiAPI(
		config.AppConfig.LineChannelAccessToken,
	)
	if err != nil {
		return err
	}

	LineServiceInstance = &LineService{
		client: client,
	}

	log.Println("LINE Service initialized")
	return nil
}

// SendMessage sends a text message to LINE user (reply)
func (s *LineService) SendMessage(replyToken string, message string) error {
	_, err := s.client.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: replyToken,
			Messages: []messaging_api.MessageInterface{
				messaging_api.TextMessage{
					Text: message,
				},
			},
		},
	)

	if err != nil {
		log.Printf("Failed to send LINE message: %v", err)
		return err
	}

	log.Printf("Sent LINE message: %s", message)
	return nil
}

// PushMessage sends a text message to LINE user (push)
func (s *LineService) PushMessage(userID string, message string) error {
	_, err := s.client.PushMessage(
		&messaging_api.PushMessageRequest{
			To: userID,
			Messages: []messaging_api.MessageInterface{
				messaging_api.TextMessage{
					Text: message,
				},
			},
		},
		"",
	)

	if err != nil {
		log.Printf("Failed to push LINE message: %v", err)
		return err
	}

	log.Printf("Pushed LINE message to %s: %s", userID, message)
	return nil
}

// SetLastUserID sets the last user ID who sent a message
func (s *LineService) SetLastUserID(userID string) {
	s.lastUserID = userID
	log.Printf("Set last LINE user ID: %s", userID)
}

// GetLastUserID gets the last user ID
func (s *LineService) GetLastUserID() string {
	return s.lastUserID
}

// ParseRequest parses LINE webhook request
func (s *LineService) ParseRequest(channelSecret string, req *http.Request) (*webhook.CallbackRequest, error) {
	return webhook.ParseRequest(channelSecret, req)
}
