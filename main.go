package main

import (
	"line-discord-bridge/config"
	"line-discord-bridge/handlers"
	"line-discord-bridge/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 設定読み込み
	config.LoadConfig()

	// LINE Service初期化
	if err := services.InitLineService(); err != nil {
		log.Fatalf("Failed to initialize LINE service: %v", err)
	}

	// Discord Service初期化
	if err := services.InitDiscordService(); err != nil {
		log.Fatalf("Failed to initialize Discord service: %v", err)
	}
	defer services.DiscordServiceInstance.Close()

	// Gemini Service初期化
	if config.AppConfig.GeminiAPIKey != "" {
		if err := services.InitGeminiService(); err != nil {
			log.Printf("Failed to initialize Gemini service: %v", err)
		} else {
			defer services.GeminiServiceInstance.Close()
		}
	} else {
		log.Println("Gemini API Key not found, skipping Gemini service initialization")
	}

	// Ginのルーター作成
	router := gin.Default()

	// ヘルスチェックエンドポイント
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "LINE-Discord Bridge Server",
		})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	// LINE Webhook エンドポイント
	// 複数のパターンに対応
	router.POST("/webhook/line", handlers.LineWebhookHandler)
	router.POST("/webhook", handlers.LineWebhookHandler)
	router.POST("/callback", handlers.LineWebhookHandler)

	// サーバー起動
	port := config.AppConfig.Port
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
