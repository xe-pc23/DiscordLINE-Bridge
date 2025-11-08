package handlers

import (
	"line-discord-bridge/config"
	"line-discord-bridge/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

// LineWebhookHandler handles LINE webhook events
func LineWebhookHandler(c *gin.Context) {
	// Webhookリクエストをパース
	cb, err := services.LineServiceInstance.ParseRequest(
		config.AppConfig.LineChannelSecret,
		c.Request,
	)
	if err != nil {
		log.Printf("Failed to parse LINE request: %v", err)
		if err == webhook.ErrInvalidSignature {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	// 各イベントを処理
	for _, event := range cb.Events {
		log.Printf("LINE Event type: %s", event.GetType())

		// メッセージイベントのみ処理
		switch e := event.(type) {
		case webhook.MessageEvent:
			handleLineMessage(e)
		default:
			log.Printf("Unsupported event type: %T", event)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func handleLineMessage(event webhook.MessageEvent) {
	// ユーザーIDを保存（Discordからの返信用）
	if event.Source != nil {
		switch source := event.Source.(type) {
		case webhook.UserSource:
			services.LineServiceInstance.SetLastUserID(source.UserId)
		}
	}

	// メッセージタイプを確認
	switch message := event.Message.(type) {
	case webhook.TextMessageContent:
		log.Printf("Received LINE message: %s", message.Text)

		// Discordに転送
		if services.DiscordServiceInstance != nil {
			err := services.DiscordServiceInstance.SendDM(
				config.AppConfig.DiscordUserID,
				message.Text,
			)
			if err != nil {
				log.Printf("Failed to forward to Discord: %v", err)
			}
		}

	default:
		log.Printf("Unsupported message type: %T", message)
	}
}
