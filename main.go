package main

import (
	"line-discord-bridge/config"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 設定読み込み
	config.LoadConfig()

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

	// サーバー起動
	port := config.AppConfig.Port
	log.Printf("Server starting on port %s", port)
	//起動しなかったときのエラー処理
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
