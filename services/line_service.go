package services

import (
	"line-discord-bridge/config"
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

type LineService struct {
	client *messaging_api.MessagingApiAPI
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

// SendMessage sends a text message to LINE user
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

// ParseRequest parses LINE webhook request
func (s *LineService) ParseRequest(channelSecret string, req *http.Request) (*webhook.CallbackRequest, error) {
	return webhook.ParseRequest(channelSecret, req)
}
