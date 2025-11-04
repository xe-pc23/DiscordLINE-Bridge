package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                   string
	LineChannelSecret      string
	LineChannelAccessToken string
	DiscordBotToken        string
	DiscordUserID          string
}

var AppConfig *Config

// LoadConfig loads environment variables
func LoadConfig() {
	// .envファイルを読み込み（存在しない場合はスキップ）
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	AppConfig = &Config{
		Port:                   getEnv("PORT", "8080"),
		LineChannelSecret:      getEnv("LINE_CHANNEL_SECRET", ""),
		LineChannelAccessToken: getEnv("LINE_CHANNEL_ACCESS_TOKEN", ""),
		DiscordBotToken:        getEnv("DISCORD_BOT_TOKEN", ""),
		DiscordUserID:          getEnv("DISCORD_USER_ID", ""),
	}

}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
